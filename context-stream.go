package protoapi

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/*************************************
 实现grpc.ServerStream
*************************************/

func (ctx *Context) SendHeader(md metadata.MD) error {
	h := ctx.ResponseWriter.Header()
	for k, v := range md {
		h[k] = v
	}
	return nil
}

func (ctx *Context) SetHeader(md metadata.MD) error {
	panic("unsupported operation")
}

func (ctx *Context) SetTrailer(md metadata.MD) {
	panic("unsupported operation")
}

func (ctx *Context) SendMsg(m interface{}) (err error) {
	if ctx.streamWriter == nil {
		// 来自restful请求不会设置streamWriter,在此进行初始化!
		if ctx.mux.closed != 0 {
			ctx.ResponseWriter.Header()["Connection"] = closeConnection
		}
		ctx.ResponseWriter.Header()["Content-Type"] = jsonContentType
		ctx.ResponseWriter.WriteStatus(ctx.Handler.Status)
		ctx.streamWriter = ctx.ResponseWriter.ResponseWriter
	}
	return ctx.writeApplyResult(ctx.streamWriter, m)
}

// RecvMsg 为与server.bootstrapStreamInterceptor()逻辑一致, streaming场景不校验MessageValidator
func (ctx *Context) RecvMsg(m interface{}) error {
	if ctx.streamReader == nil {
		// 来自restful请求不会设置streamReader,在此进行初始化!
		ctx.streamReader = ctx.Request.Body
	}
	// 解码请求体. 但验证需要
	return DecodeRequest(ctx.streamReader, m)
}

func (ctx *Context) Context() context.Context {
	return ctx
}

var _ grpc.ServerStream = (*Context)(nil)

type ServerStreamContext[Req any, Res any] struct {
	grpc.ServerStream
}

func (ctx *ServerStreamContext[Req, Res]) Recv() (*Req, error) {
	req := new(Req)
	err := ctx.ServerStream.RecvMsg(req)
	return req, err
}

func (ctx *ServerStreamContext[Req, Res]) Send(res *Res) error {
	return ctx.ServerStream.SendMsg(res)
}

func (ctx *ServerStreamContext[Req, Res]) SendAndClose(res *Res) error {
	return ctx.ServerStream.SendMsg(res)
}

func StreamContext[Req, Res any](ctx *Context) *ServerStreamContext[Req, Res] {
	return &ServerStreamContext[Req, Res]{grpc.ServerStream(ctx)}
}
