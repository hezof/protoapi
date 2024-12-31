package protoapi

import (
	"context"
	"google.golang.org/grpc/metadata"
)

/*************************************
 实现grpc.streaming
*************************************/

type streaming struct {
	c *Context // 流式上下文
}

func (s *streaming) SendHeader(md metadata.MD) error {
	hd := s.c.ResponseWriter.Header()
	for k, v := range md {
		hd[k] = v
	}
	return nil
}

func (s *streaming) SetHeader(md metadata.MD) error {
	return s.SendHeader(md)
}

func (s *streaming) SetTrailer(md metadata.MD) {
	// non-supported
}

func (s *streaming) SendMsg(m interface{}) (err error) {
	if s.c.streamEncoder == nil {
		if s.c.mux.closed != 0 {
			s.c.ResponseWriter.Header()["Connection"] = closeConnection
		}
		s.c.ResponseWriter.Header()["Content-Type"] = jsonContentType
		s.c.ResponseWriter.WriteHeader(s.c.Handler.Status)
		s.c.streamEncoder = s.c.Handler.NewEncoder(s.c.ResponseWriter.ResponseWriter)
	}
	if s.c.Handler.Unwrap {
		err = s.c.streamEncoder.Encode(m)
	} else {
		s.c.wrapper.Data = m
		err = s.c.streamEncoder.Encode(&s.c.wrapper)
		s.c.wrapper.Data = nil // 及时清理
	}
	return
}

func (s *streaming) RecvMsg(m interface{}) error {
	if s.c.streamDecoder == nil {
		s.c.streamDecoder = NewDecoder(s.c.Request.Body)
	}
	return s.c.streamDecoder.Decode(m)
}

func (s *streaming) Context() context.Context {
	return s.c
}
