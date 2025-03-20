package protoapi

import (
	"fmt"
	"github.com/hezof/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**********************************
 * 扩展Error实现StatusResult接口
 **********************************/

func (e *Error) Error() string {
	// 不能使用ToJson()会在EncodeField过程形成死循环
	return base.ToJson(e)
}

func (e *Error) DecodeField(r *base.JsonDecoder, f string) {
	switch f {
	case profile.ResultCodeField:
		base.DecodeUint32(r, &e.Code)
	case profile.ResultNameField:
		base.DecodeString(r, &e.Name)
	case profile.ResultMessageField:
		base.DecodeString(r, &e.Message)
	}
}

func (e *Error) EncodeField(w *base.JsonEncoder) {
	base.EncodeUint32_WithEmpty(w, profile.ResultCodeField, e.Code)
	base.EncodeString_OmitEmpty(w, profile.ResultNameField, e.Name)
	base.EncodeString_OmitEmpty(w, profile.ResultMessageField, e.Message)
}

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (e *Error) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Status<<base.ErrorCodeBits|e.Code), e.Message)
}

func (e *Error) SetStatus(status uint32) {
	e.Status = status
}

func (e *Error) SetName(name string) {
	e.Name = name
}

func (e *Error) SetMessage(message string) {
	if len(e.Details) > 0 {
		e.Message = fmt.Sprintf(message, base.AnySlice(e.Details)...)
	} else {
		e.Message = message
	}
}

var _ base.Error = (*Error)(nil)
