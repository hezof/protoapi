package protoapi

import (
	"io"
	"net/http"
)

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
