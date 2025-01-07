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
func (ctx *Context) clean() *Context {
	ctx.resources = nil
	ctx.attribute = nil
	ctx.query = nil
	ctx.panic = nil
	ctx.ResponseWriter.ResponseWriter = nil
	ctx.Request = nil
	ctx.Handler = nil

	return ctx
}

// 使用前重置上下文
func (ctx *Context) serveNext(w http.ResponseWriter, r *http.Request, h *Handler) {
	ctx.ResponseWriter.ResponseWriter = w
	ctx.ResponseWriter.statusCode = 0
	ctx.Request = r
	ctx.Handler = h
	ctx.handle = -1 // 初始值为-1!
	ctx.Next()
}

func (ctx *Context) panicNext(p interface{}, h *Handler) {
	// 保留上下文信息,重置handler为panic handler继续处理.
	ctx.panic = p
	ctx.Handler = h
	ctx.handle = -1 // 初始值为-1!
	ctx.Next()
}

func (ctx *Context) Next() {
	ctx.handle++
	for ctx.handle < len(ctx.Handler.HandleChain) {
		ctx.Handler.HandleChain[ctx.handle](ctx)
		ctx.handle++
	}
}

func (ctx *Context) Aborted() bool {
	return ctx.handle >= HandleChainCapacity
}

func (ctx *Context) Abort() {
	ctx.handle = HandleChainCapacity
}

func (ctx *Context) AbortWithStatus(statusCode int) {
	ctx.handle = HandleChainCapacity
	ctx.ResponseWriter.WriteHeader(statusCode)
}

func (ctx *Context) ParamValue(key string) string {
	v, _ := ctx.params.Get(key)
	return v
}

/*************************************
 实现context.Context
*************************************/

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.Request.Context().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.Request.Context().Done()
}

func (ctx *Context) Err() error {
	return ctx.Request.Context().Err()
}

func (ctx *Context) Value(key any) any {
	if ret, ok := ctx.attribute[key]; ok {
		return ret
	}
	return ctx.Request.Context().Value(key)
}

/*************************************
 输出结果
*************************************/

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

func (ctx *Context) resource(code uint32) *resource {
	if ctx.resources == nil {
		var lang string
		if vs, ok := ctx.Request.Header["Accept-Language"]; ok {
			lang = vs[0]
		}
		ctx.resources = fastGetResMapByAcceptLanguage(lang)
	}
	return ctx.resources[code]
}

func (ctx *Context) WriteErrorResult(result *StatusResult) error {
	// graceful关闭期间断开keepalive连接
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}
	// 设置内容类型
	ctx.ResponseWriter.Header()["Content-Type"] = jsonContentType
	// 写出状态与结果
	ctx.ResponseWriter.WriteStatus(result.Status)
	return ctx.writeErrorResult(ctx.ResponseWriter.ResponseWriter, result)
}

func (ctx *Context) writeErrorResult(out io.Writer, result *StatusResult) error {
	// 国际化错误消息(延后初始化)
	if rs := ctx.resource(result.Code); rs != nil {
		// 覆盖status
		if rs.Status > 0 {
			result.Status = rs.Status
		}
		// 支持参数格式化
		if len(result.Details) == 0 {
			result.Message = rs.Message
		} else {
			result.Message = Sprintf(rs.Message, result.Details...)
		}
	}

	w := GetEncoder(out)
	defer PutEncoder(w)

	result.EncodeJSON(w)
	return w.Close()
}

func (ctx *Context) WriteApplyResult(val any) error {
	// graceful关闭期间断开keepalive连接
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}
	// 设置内容类型
	ctx.ResponseWriter.Header()["Content-Type"] = jsonContentType
	// 写出状态与结果
	ctx.ResponseWriter.WriteStatus(ctx.Handler.Meta.Status)
	return ctx.writeApplyResult(ctx.ResponseWriter.ResponseWriter, val)
}

func (ctx *Context) writeApplyResult(out io.Writer, val any) error {
	switch ctx.Handler.Meta.Result {
	case Http_simple:

		w := GetEncoder(ctx.ResponseWriter.ResponseWriter)
		defer PutEncoder(w)

		ctx.result.Code = ctx.Handler.Meta.Status
		ctx.result.Data = val
		ctx.result.EncodeJSON(w)

		return w.Close()

	case Http_unwrap:

		if jc, ok := val.(JsonCodec); ok {
			w := GetEncoder(ctx.ResponseWriter.ResponseWriter)
			defer PutEncoder(w)

			jc.EncodeJSON(w)
			return w.Close()
		} else {
			w := json.NewEncoder(ctx.ResponseWriter.ResponseWriter)
			w.SetEscapeHTML(false)
			return w.Encode(val)
		}

	}
	return fmt.Errorf("unsupport result catalog: %v", ctx.Handler.Meta.Result)
}

func (ctx *Context) WritePlain(status int, data string) error {
	return ctx.WritePlainBytes(status, UnsafeBytes(data))
}

func (ctx *Context) WritePlainBytes(status int, data []byte) error {
	// graceful关闭期间断开keepalive连接
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}
	// 设置内容类型
	ctx.ResponseWriter.Header()["Content-Type"] = plainContentType
	// 写出状态与结果
	ctx.ResponseWriter.WriteHeader(status)
	_, err := ctx.ResponseWriter.Write(data)
	return err
}

func (ctx *Context) WriteHtml(status int, data string) error {
	return ctx.WriteHtmlBytes(status, UnsafeBytes(data))
}

func (ctx *Context) WriteHtmlBytes(status int, data []byte) error {
	// graceful关闭期间断开keepalive连接
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}
	// 设置内容类型
	ctx.ResponseWriter.Header()["Content-Type"] = htmlContentType
	// 写出状态与结果
	ctx.ResponseWriter.WriteHeader(status)
	_, err := ctx.ResponseWriter.Write(data)
	return err
}

/*************************************
 请求内容
*************************************/

// ReadBody 读取后可能导致body流指针指向最后导致后面无法读取!
func (ctx *Context) ReadBody() ([]byte, error) {
	if buff, ok := ctx.Request.Body.(*BuffBody); ok {
		return buff.data, nil
	} else {
		return io.ReadAll(ctx.Request.Body) // Transfer-Encoding: chunked的情况!
	}
}

func (ctx *Context) CopyBody() ([]byte, error) {
	if buff, ok := ctx.Request.Body.(*BuffBody); ok {
		return buff.data, nil
	} else {
		data, err := io.ReadAll(ctx.Request.Body)
		if err == nil {
			ctx.Request.Body = &BuffBody{data: data}
		}
		return data, err
	}
}

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
	rw.ResponseWriter.WriteHeader(rw.statusCode)
}

func (rw *proxyResponseWriter) WriteStatus(status uint32) {
	rw.WriteHeader(int(status))
}

/*************************************
 与gin.Context兼容的方法
*************************************/
