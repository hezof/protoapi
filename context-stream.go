package protoapi

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/*************************************
 实现grpc.ServerStream
*************************************/

type stream struct {
	c *Context // 流式上下文
}

func (s *stream) SendHeader(md metadata.MD) error {
	h := s.c.ResponseWriter.Header()
	for k, v := range md {
		h[k] = v
	}
	return nil
}

func (s *stream) SetHeader(md metadata.MD) error {
	panic("unsupported operation")
}

func (s *stream) SetTrailer(md metadata.MD) {
	panic("unsupported operation")
}

func (s *stream) SendMsg(m interface{}) (err error) {
	if s.c.streamEncoder == nil {
		if s.c.mux.closed != 0 {
			s.c.ResponseWriter.Header()["Connection"] = closeConnection
		}
		s.c.ResponseWriter.Header()["Content-Type"] = jsonContentType
		s.c.ResponseWriter.WriteHeader(int(s.c.Handler.Meta.Status))
		s.c.streamEncoder = GetEncoder(s.c.ResponseWriter.ResponseWriter)
	}
	
}

func (s *stream) RecvMsg(m interface{}) error {
	if s.c.streamDecoder == nil {
		s.c.streamDecoder = NewDecoder(s.c.Request.Body)
	}
	return s.c.streamDecoder.Decode(m)
}

func (s *stream) Context() context.Context {
	return s.c
}

var _ grpc.ServerStream = (*stream)(nil)

type Test struct {
	stream
}

var _ grpc.ServerStream = (*Test)(nil)
