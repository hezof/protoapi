package protoapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (e *Error) Error() string {
	return ToJson(e)
}

func (e *Error) SetStatus(status uint32) {
	e.Status = status
}

func (e *Error) SetName(name string) {
	e.Name = name
}

func (e *Error) SetMessage(message string) {
	if len(e.De)
}

func (e *Error) DecodeField(r *JsonDecoder, f string) {
	switch f {
	case profile.ResultCodeField:
		DecodeUint32(r, &e.Code)
	case profile.ResultNameField:
		DecodeString(r, &e.Name)
	case profile.ResultMessageField:
		DecodeString(r, &e.Message)
	}
}

func (e *Error) EncodeField(w *JsonEncoder) {
	EncodeUint32_WithEmpty(w, profile.ResultCodeField, e.Code)
	EncodeString_OmitEmpty(w, profile.ResultNameField, e.Name)
	EncodeString_OmitEmpty(w, profile.ResultMessageField, e.Message)
}

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (e *Error) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Status<<_CodeBits|e.Code), e.Message)
}

var _ StatusResult = (*Error)(nil)
