package protoapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Context 处理上下文.必须注意:Context IS NON-THREAD-SAFE!!!
type Context struct {
	mux            *mux                 // (不用清理)创建Context的mux实例(不能重置)
	result         StatusResult         // (不用清理)仅仅用于wrap apply result避免反复创建临时Result! 不作其他用途!
	params         Params               // (不用清理)Params, 每次调用node.getValue()都必须重置
	skippedNodes   []skippedNode        // (不用清理)配合gin-tree使用
	cipath         []byte               // (不用清理)case-insensitive lookup path
	cibuff         [4]byte              // (不用清理)case-insensitive lookup buffer
	handle         int                  // (不用清理)初始值为-1!
	ResponseWriter proxyResponseWriter  // (需要清理)响应ResponseWriter. 不建议nested struct!
	Request        *http.Request        // (需要清理)请求Request
	Handler        *Handler             // (需要清理)处理Handler
	resources      map[uint32]*resource // (需要清理)i18n资源
	attribute      map[any]any          // (需要清理)处理属性
	query          url.Values           // (需要清理)请求query
	panic          any                  // (需要清理)处理异常
}

// clean 负责清理外部引用. 确保pool安全!
func (ctx *Context) clean() *Context {
	ctx.ResponseWriter.ResponseWriter = nil
	ctx.Request = nil
	ctx.Handler = nil
	ctx.resources = nil
	ctx.attribute = nil
	ctx.query = nil
	ctx.panic = nil
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
	ctx.Handler = h
	ctx.handle = -1 // 初始值为-1!
	ctx.panic = p
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

var _ context.Context = (*Context)(nil)

/*************************************
 输出结果
*************************************/

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
	ctx.ResponseWriter.WriteStatus(result.status)
	return ctx.writeErrorResult(ctx.ResponseWriter.ResponseWriter, result)
}

func (ctx *Context) writeErrorResult(out io.Writer, result *StatusResult) error {
	// 国际化错误消息(延后初始化)
	if rs := ctx.resource(result.GetCode()); rs != nil {
		// 覆盖status
		if rs.Status > 0 {
			result.SetStatus(rs.Status)
		}
		if rs.Name != `` {
			result.SetName(rs.Name)
		}
		if rs.Message != `` {
			result.SetMessage(rs.Message)
		}
	}

	w := GetEncoder(out)
	defer PutEncoder(w)

	EncodeMessage(w, result)
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
	ctx.ResponseWriter.WriteStatus(ctx.Handler.Setting.Meta.Http.Status)
	return ctx.writeApplyResult(ctx.ResponseWriter.ResponseWriter, val)
}

func (ctx *Context) writeApplyResult(out io.Writer, val any) error {
	switch ctx.Handler.Setting.Meta.Http.Result {
	case Http_simple:

		w := GetEncoder(ctx.ResponseWriter.ResponseWriter)
		defer PutEncoder(w)

		ctx.result.Code = ctx.Handler.Setting.Meta.Http.Status
		ctx.result.Data = val
		EncodeMessage(w, &ctx.result)

		return w.Close()

	case Http_unwrap:

		w := GetEncoder(ctx.ResponseWriter.ResponseWriter)
		defer PutEncoder(w)

		EncodeAny(w, val)
		return w.Close()
	}
	return fmt.Errorf("unsupport result catalog: %v", ctx.Handler.Setting.Meta.Http.Result)
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
