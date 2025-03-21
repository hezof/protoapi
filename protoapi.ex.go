package protoapi

import (
	"fmt"
	"github.com/hezof/framework"
	"github.com/hezof/protojson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**********************************
 * 扩展Error实现StatusResult接口
 **********************************/

func (e *Error) Error() string {
	// 不能使用ToJson()会在EncodeField过程形成死循环
	return framework.ToJson(e)
}

func (e *Error) DecodeField(r *protojson.JsonDecoder, f string) {
	switch f {
	case profile.ResultCodeField:
		protojson.DecodeUint32(r, &e.Code)
	case profile.ResultNameField:
		protojson.DecodeString(r, &e.Name)
	case profile.ResultMessageField:
		protojson.DecodeString(r, &e.Message)
	}
}

func (e *Error) EncodeField(w *protojson.JsonEncoder) {
	protojson.EncodeUint32_WithEmpty(w, profile.ResultCodeField, e.Code)
	protojson.EncodeString_OmitEmpty(w, profile.ResultNameField, e.Name)
	protojson.EncodeString_OmitEmpty(w, profile.ResultMessageField, e.Message)
}

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (e *Error) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Status<<framework.ErrorCodeBits|e.Code), e.Message)
}

func (e *Error) SetStatus(status uint32) {
	e.Status = status
}

func (e *Error) SetName(name string) {
	e.Name = name
}

func (e *Error) SetMessage(message string) {
	if len(e.Details) > 0 {
		e.Message = fmt.Sprintf(message, framework.AnySlice(e.Details)...)
	} else {
		e.Message = message
	}
}

var _ framework.Error = (*Error)(nil)
