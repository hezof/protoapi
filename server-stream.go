package protoapi

import "google.golang.org/grpc"

type ServerStreamExtend struct {
	grpc.ServerStream
	*MethodSetting
}

func (ss *ServerStreamExtend) RecvMsg(m any) error {
	err := ss.ServerStream.RecvMsg(m)
	if err != nil {
		return err
	}
	if mv, ok := m.(MessageValidator); ok {
		return mv.Validate(ss.MethodSetting, ss.ServerStream.Context())
	}
	return nil
}

var _ grpc.ServerStream = (*ServerStreamExtend)(nil)
