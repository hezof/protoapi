package protoapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/hezof/core"
	"io"
	"mime/multipart"
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
	streamReader   io.Reader            // (需要清理)流式读入器, 同时也是流式环境的初始化标识!
	streamWriter   io.Writer            // (需要清理)流式写出器, 同时也是流式环境的初始化标识!
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
	ctx.streamReader = nil
	ctx.streamWriter = nil
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

// resource 请求i18n资源
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

// WriteErrorResult 用于restful写出错误结果
func (ctx *Context) WriteErrorResult(result core.Error) error {
	// graceful关闭期间断开keepalive连接
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}
	// 设置内容类型
	ctx.ResponseWriter.Header()["Content-Type"] = jsonContentType
	// 写出状态与结果
	ctx.ResponseWriter.WriteStatus(result.GetStatus())
	return ctx.writeErrorResult(ctx.ResponseWriter.ResponseWriter, result)
}

// writeErrorResult 用于websocket写出错误结果
func (ctx *Context) writeErrorResult(out io.Writer, result core.Error) error {
	// 国际化错误消息(延后初始化)
	if hasResMap {
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
	}
	return EncodeResponse(out, result)
}

// WriteApplyResult 用于restful写出请求结果
func (ctx *Context) WriteApplyResult(val any) error {
	// graceful关闭期间断开keepalive连接
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}
	// 设置内容类型
	ctx.ResponseWriter.Header()["Content-Type"] = jsonContentType
	// 写出状态与结果(status不为0)
	ctx.ResponseWriter.WriteStatus(ctx.Handler.Status)
	return ctx.writeApplyResult(ctx.ResponseWriter.ResponseWriter, val)
}

// writeApplyResult 用于websocket写出请求结果
func (ctx *Context) writeApplyResult(out io.Writer, val any) error {
	switch ctx.Handler.Result {
	case Result_normal:

		ctx.result.Code = 0
		ctx.result.Data = val

		err := EncodeResponse(out, &ctx.result)

		// 及时清理避免引用
		ctx.result.Data = nil

		return err

	case Result_unwrap:

		return EncodeResponse(out, val)

	default:
		return fmt.Errorf("invalid result type: %v", ctx.Handler.Setting.Meta.Http.Result)
	}

}

// ReadBody 读取后可能导致body流指针指向最后导致后面无法读取!
func (ctx *Context) ReadBody() ([]byte, error) {
	if buff, ok := ctx.Request.Body.(*core.BuffBody); ok {
		return buff.Data, nil
	} else {
		return io.ReadAll(ctx.Request.Body) // Transfer-Encoding: chunked的情况!
	}
}

func (ctx *Context) CopyBody() ([]byte, error) {
	if buff, ok := ctx.Request.Body.(*core.BuffBody); ok {
		buff.Reset()
		return buff.Data, nil
	} else {
		data, err := io.ReadAll(ctx.Request.Body)
		if err == nil {
			ctx.Request.Body = &core.BuffBody{Data: data}
		}
		return data, err
	}
}

func (ctx *Context) ReadJson(val any) error {
	// BuffBody允许重复读
	if buff, ok := ctx.Request.Body.(*core.BuffBody); ok {
		buff.Reset()
	}
	return DecodeRequest(ctx.Request.Body, val)
}

func (ctx *Context) CopyJson(val any) error {
	// BuffBody允许重复读
	if _, err := ctx.CopyBody(); err != nil {
		return err
	}
	return DecodeRequest(ctx.Request.Body, val)
}

func (ctx *Context) WriteJson(status uint32, val any) error {
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}
	ctx.ResponseWriter.Header()["Content-Type"] = jsonContentType
	ctx.ResponseWriter.WriteStatus(status)
	return EncodeResponse(ctx.ResponseWriter.ResponseWriter, val)
}
func (ctx *Context) WritePlain(status int, data string) error {
	return ctx.WritePlainBytes(status, core.UnsafeBytes(data))
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
	return ctx.WriteHtmlBytes(status, core.UnsafeBytes(data))
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
 辅助变量与数据结构
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

var _ http.ResponseWriter = (*proxyResponseWriter)(nil)

/*************************************
 辅助函数与上传/下载
*************************************/

var ErrInvalidContextToReadFile = errors.New("invalid context to read file")
var ErrInvalidContextToWriteFile = errors.New("invalid context to write file")

// ReadFile 配合body=file使用! 专门用于HTTP上传数据! 无法用于GRPC场景!
func ReadFile(ctx context.Context, key string) (multipart.File, *multipart.FileHeader, error) {
	if wc, ok := ctx.(*Context); ok {
		return wc.Request.FormFile(key)
	}
	return nil, nil, ErrInvalidContextToReadFile
}

// WriteFile 专门用于HTTP上传文件! 无法用于GRPC场景!
// 使用WriteFile()/WriterStream()的Service Method实现必须确保返回rsp为nil(即无法用于grpc调用)!
func WriteFile(ctx context.Context, name string, data io.Reader) error {
	if wc, ok := ctx.(*Context); ok {
		if wc.mux.closed != 0 {
			wc.ResponseWriter.Header()["Connection"] = closeConnection
		}
		wc.ResponseWriter.Header()["Content-Type"] = streamContentType
		wc.ResponseWriter.Header()["Content-Disposition"] = []string{"attachment;filename=" + url.QueryEscape(name)} // 防止中文乱码
		wc.ResponseWriter.WriteStatus(wc.Handler.Status)
		if _, err := io.Copy(wc.ResponseWriter.ResponseWriter, data); err != nil {
			return err
		}
		return nil
	}
	return ErrInvalidContextToWriteFile
}

var ErrInvalidContextToReadStream = errors.New("invalid context to read stream")
var ErrInvalidContextToWriteStream = errors.New("invalid context to write stream")

// ReadStream 配合body=stream使用! 专门用于HTTP上传数据! 无法用于GRPC场景!
func ReadStream(ctx context.Context) (io.ReadCloser, error) {
	if wc, ok := ctx.(*Context); ok {
		return wc.Request.Body, nil
	}
	return nil, ErrInvalidContextToReadStream
}

// WriteStream 专门用于HTTP上传文件! 无法用于GRPC场景!
// 使用DownFile()/WriterStream()的Service Method实现必须确保返回rsp为nil(即无法用于grpc调用)!
func WriteStream(ctx context.Context, headers ...http.Header) (io.Writer, error) {
	if wc, ok := ctx.(*Context); ok {
		if wc.mux.closed != 0 {
			wc.ResponseWriter.Header()["Connection"] = closeConnection
		}
		for _, h := range headers {
			for k, vs := range h {
				wc.ResponseWriter.Header()[k] = vs
			}
		}
		wc.ResponseWriter.WriteStatus(wc.Handler.Status)
		return wc.ResponseWriter.ResponseWriter, nil
	}
	return nil, ErrInvalidContextToWriteStream
}
