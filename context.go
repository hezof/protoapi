package protoapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HandleFunc 处理逻辑函数
type HandleFunc func(ctx *Context)

// Context 处理上下文.必须注意:Context IS NON-THREAD-SAFE!!!
type Context struct {
	streaming                                  // 支持grpc.ServerStream的基类,内部包含对Context的引用. 必须nested struct!
	mux            *mux                        // (不用清理)创建Context的mux实例(不能重置)
	params         Params                      // (不用清理)Params, 每次调用node.getValue()都必须重置
	skippedNodes   []skippedNode               // (不用清理)配合gin-tree使用
	cipath         []byte                      // (不用清理)case-insensitive lookup path
	cibuff         [4]byte                     // (不用清理)case-insensitive lookup buffer
	handle         int                         // (不用清理)初始值为-1!
	query          url.Values                  // 请求query
	attrs          map[interface{}]interface{} // 处理属性
	panic          interface{}                 // 处理异常
	ResponseWriter *proxyResponseWriter        // 响应ResponseWriter. 不建议nested struct!
	Request        *http.Request               // 请求Request
	Handler        *Handler                    // 处理Handler
}

var _ context.Context = (*Context)(nil)

// clean 负责清理外部引用. 确保pool安全!
func (c *Context) clean() *Context {
	c.ResponseWriter.ResponseWriter = nil
	c.Request = nil
	c.Handler = nil
	c.query = nil
	c.attrs = nil
	c.panic = nil
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
	return c.handle >= profile.HandleChainCapacity
}

func (c *Context) Abort() {
	c.handle = profile.HandleChainCapacity
}

func (c *Context) AbortWithStatus(statusCode int) {
	c.handle = profile.HandleChainCapacity
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

func (c *Context) Value(key interface{}) interface{} {
	if ret, ok := c.attrs[key]; ok {
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

/*************************************
 与gin.Context兼容的方法
*************************************/

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
 与Handler相关的方法.
*************************************/

var (
	jsonContentType     = []string{"application/json; charset=utf-8"}
	htmlContentType     = []string{"text/html; charset=utf-8"}
	plainContentType    = []string{"text/plain; charset=utf-8"}
	streamContentType   = []string{"application/octet-streaming"}
	eventsContentType   = []string{"text/event-stream"}
	eventsCacheControl  = []string{"no-cache"}
	keepAliveConnection = []string{"keep-alive"}
	closeConnection     = []string{"close"}
)

// WriteApplyResult 写出请求结果
func (c *Context) WriteApplyResult(rsp MessageEncoder) error {
	if c.mux.closed != 0 {
		c.ResponseWriter.Header()["Connection"] = closeConnection
	}
	c.ResponseWriter.Header()["Content-Type"] = jsonContentType
	c.ResponseWriter.WriteHeader(int(c.Handler.Meta.Http.Status))
	if c.Handler.Meta.Http.Result == SimpleResult {

	} else {

	}

	return nil
}

// WriteErrorResult 写出错误结果
func (c *Context) WriteErrorResult(err error) error {
	// 统一转换错误为StatusError
	result := StatusErrorFrom(err, DefaultErrorStatus)
	// 打印未知道错误作为预警信息!
	if result.Code == CodeUnknown {
		// 打印未知错误(code==2)用于日志预警
		log.Error("unknown error: uri=%v, err=%v", c.Request.RequestURI, err)
	}
	if lenResMap > 0 {
		c.i18nErrorResult(result, fastGetResMapByAcceptLanguage(c.getAcceptLanguage()))
	}
	return c.WriteJson(result.status, result)
}

func (c *Context) WriteJson(status int, obj interface{}) error {
	if c.mux.closed != 0 {
		c.ResponseWriter.Header()["Connection"] = closeConnection
	}
	c.ResponseWriter.Header()["Content-Type"] = jsonContentType
	c.ResponseWriter.WriteHeader(status)
	if me, ok := obj.(MessageEncoder); ok {
		enc := GetEncoder(c.ResponseWriter.ResponseWriter)
		defer PutEncoder(enc)
		return me.EncodeJSON(enc)
	}
	enc := json.NewEncoder(c.ResponseWriter.ResponseWriter)
	enc.SetEscapeHTML(false)
	return enc.Encode(obj)
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
