package protoapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HandleFunc 处理逻辑函数
type HandleFunc func(ctx *Context)

// Context 处理上下文.必须注意:Context IS NON-THREAD-SAFE!!!
type Context struct {
	stream                              // (不用清理)支持grpc.ServerStream的基类,内部包含对Context的引用. 必须nested struct!
	mux            *mux                 // (不用清理)创建Context的mux实例(不能重置)
	result         StatusResult         // (不用清理)仅仅用于wrap apply result避免反复创建临时Result! 不作其他用途!
	params         Params               // (不用清理)Params, 每次调用node.getValue()都必须重置
	skippedNodes   []skippedNode        // (不用清理)配合gin-tree使用
	cipath         []byte               // (不用清理)case-insensitive lookup path
	cibuff         [4]byte              // (不用清理)case-insensitive lookup buffer
	handle         int                  // (不用清理)初始值为-1!
	resources      map[uint32]*resource // i18n资源
	attribute      map[any]any          // 处理属性
	query          url.Values           // 请求query
	panic          any                  // 处理异常
	ResponseWriter *proxyResponseWriter // 响应ResponseWriter. 不建议nested struct!
	Request        *http.Request        // 请求Request
	Handler        *Handler             // 处理Handler
}

var _ context.Context = (*Context)(nil)

// clean 负责清理外部引用. 确保pool安全!
func (c *Context) clean() *Context {
	c.resources = nil
	c.attribute = nil
	c.query = nil
	c.panic = nil
	c.ResponseWriter.ResponseWriter = nil
	c.Request = nil
	c.Handler = nil

	return c
}

// 使用前重置上下文
func (c *Context) serveNext(w http.ResponseWriter, r *http.Request, h *Handler) {
	c.ResponseWriter.ResponseWriter = w
	c.ResponseWriter.statusCode = 0
	c.Request = r
	c.Handler = h
	c.handle = -1 // 初始值为-1!
	c.Next()
}

func (c *Context) panicNext(p interface{}, h *Handler) {
	// 保留上下文信息,重置handler为panic handler继续处理.
	c.panic = p
	c.Handler = h
	c.handle = -1 // 初始值为-1!
	c.Next()
}

func (c *Context) Next() {
	c.handle++
	for c.handle < len(c.Handler.HandleChain) {
		c.Handler.HandleChain[c.handle](c)
		c.handle++
	}
}

func (c *Context) Aborted() bool {
	return c.handle >= HandleChainCapacity
}

func (c *Context) Abort() {
	c.handle = HandleChainCapacity
}

func (c *Context) AbortWithStatus(statusCode int) {
	c.handle = HandleChainCapacity
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *Context) ParamValue(key string) string {
	v, _ := c.params.Get(key)
	return v
}

/*************************************
 实现context.Context
*************************************/

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Request.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.Request.Context().Done()
}

func (c *Context) Err() error {
	return c.Request.Context().Err()
}

func (c *Context) Value(key any) any {
	if ret, ok := c.attribute[key]; ok {
		return ret
	}
	return c.Request.Context().Value(key)
}

/*************************************
 自定义工具方法
*************************************/

// ReadBody 读取后可能导致body流指针指向最后导致后面无法读取!
func (c *Context) ReadBody() ([]byte, error) {
	if buff, ok := c.Request.Body.(*BuffBody); ok {
		return buff.data, nil
	} else {
		return io.ReadAll(c.Request.Body) // Transfer-Encoding: chunked的情况!
	}
}

func (c *Context) CopyBody() ([]byte, error) {
	if buff, ok := c.Request.Body.(*BuffBody); ok {
		return buff.data, nil
	} else {
		data, err := io.ReadAll(c.Request.Body)
		if err == nil {
			c.Request.Body = &BuffBody{data: data}
		}
		return data, err
	}
}

var (
	jsonContentType     = []string{"application/json; charset=utf-8"}
	htmlContentType     = []string{"text/html; charset=utf-8"}
	plainContentType    = []string{"text/plain; charset=utf-8"}
	streamContentType   = []string{"application/octet-stream"}
	eventsContentType   = []string{"text/event-stream"}
	eventsCacheControl  = []string{"no-cache"}
	keepAliveConnection = []string{"keep-alive"}
	closeConnection     = []string{"close"}
)

func (c *Context) resource(code uint32) *resource {
	if c.resources == nil {
		var lang string
		if vs, ok := c.Request.Header["Accept-Language"]; ok {
			lang = vs[0]
		}
		c.resources = fastGetResMapByAcceptLanguage(lang)
	}
	return c.resources[code]
}

func (c *Context) WriteErrorResult(out io.Writer, err *StatusResult) error {

	if c.mux.closed != 0 {
		c.ResponseWriter.Header()["Connection"] = closeConnection
	}

	c.ResponseWriter.Header()["Content-Type"] = jsonContentType
	c.ResponseWriter.WriteHeader(int(err.Status))

	// 国际化错误消息(延后初始化)
	if rs := c.resource(err.Code); rs != nil {
		// 覆盖status
		if rs.Status > 0 {
			err.Status = rs.Status
		}
		// 支持参数格式化
		if len(err.Details) == 0 {
			err.Message = rs.Message
		} else {
			err.Message = Sprintf(rs.Message, err.Details...)
		}
	}

	enc := GetEncoder(out)
	defer PutEncoder(enc)

	err.EncodeJSON(enc)
	return enc.Close()
}

func (c *Context) WriteApplyResult(out io.Writer, val any) error {
	if c.mux.closed != 0 {
		c.ResponseWriter.Header()["Connection"] = closeConnection
	}
	switch c.Handler.Meta.Result {
	case Http_simple:
		c.ResponseWriter.Header()["Content-Type"] = jsonContentType
		c.ResponseWriter.WriteHeader(int(c.Handler.Meta.Status))

		enc := GetEncoder(out)
		defer PutEncoder(enc)

		c.result.Code = c.Handler.Meta.Status
		c.result.Data = val
		c.result.EncodeJSON(enc)

		return enc.Close()

	case Http_unwrap:
		c.ResponseWriter.Header()["Content-Type"] = jsonContentType
		c.ResponseWriter.WriteHeader(int(c.Handler.Meta.Status))

		if jc, ok := val.(JsonCodec); ok {
			enc := GetEncoder(out)
			defer PutEncoder(enc)

			jc.EncodeJSON(enc)
			return enc.Close()
		} else {
			enc := json.NewEncoder(out)
			enc.SetEscapeHTML(false)
			return enc.Encode(val)
		}

	case Http_events:
		// TODO:
	}
	return fmt.Errorf("unsupport http result: %v", c.Handler.Meta.Result)
}

func (c *Context) WritePlain(status int, data string) error {
	return c.WritePlainBytes(status, UnsafeBytes(data))
}

func (c *Context) WritePlainBytes(status int, data []byte) error {
	if c.mux.closed != 0 {
		c.ResponseWriter.Header()["Connection"] = closeConnection
	}
	c.ResponseWriter.Header()["Content-Type"] = plainContentType
	c.ResponseWriter.WriteHeader(status)
	_, err := c.ResponseWriter.Write(data)
	return err
}

func (c *Context) WriteHtml(status int, data string) error {
	return c.WriteHtmlBytes(status, UnsafeBytes(data))
}

func (c *Context) WriteHtmlBytes(status int, data []byte) error {
	if c.mux.closed != 0 {
		c.ResponseWriter.Header()["Connection"] = closeConnection
	}
	c.ResponseWriter.Header()["Content-Type"] = htmlContentType
	c.ResponseWriter.WriteHeader(status)
	_, err := c.ResponseWriter.Write(data)
	return err
}

/*************************************
 辅助数据结构
*************************************/

// BuffBody 一次性读写Body, 配合ReadBody()一块使用
type BuffBody struct {
	head int
	data []byte
}

func (b *BuffBody) Read(p []byte) (int, error) {
	if b.head >= len(b.data) {
		return 0, io.EOF
	}
	n := copy(p, b.data[b.head:])
	b.head += n
	return n, nil
}

func (b *BuffBody) Close() error {
	b.head = len(b.data)
	return nil
}

var _ io.ReadCloser = (*BuffBody)(nil)

// proxyResponseWriter 截获status用于拦截器(length已经去掉)
type proxyResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *proxyResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

/*************************************
 与gin.Context兼容的方法
*************************************/
