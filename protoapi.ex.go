package protoapi

import (
	"fmt"
	"github.com/hezof/core"
	"github.com/hezof/protojson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**********************************
 * 扩展Error实现StatusResult接口
 **********************************/

func (e *Error) Error() string {
	// 不能使用ToJson()会在EncodeField过程形成死循环
	return core.ToJson(e)
}

func (e *Error) DecodeField(r *protojson.JsonDecoder, f string) {
	switch f {
	case core.ResultCodeField:
		protojson.DecodeUint32(r, &e.Code)
	case core.ResultNameField:
		protojson.DecodeString(r, &e.Name)
	case core.ResultMessageField:
		protojson.DecodeString(r, &e.Message)
	}
}

func (e *Error) EncodeField(w *protojson.JsonEncoder) {
	protojson.EncodeUint32_WithEmpty(w, core.ResultCodeField, e.Code)
	protojson.EncodeString_OmitEmpty(w, core.ResultNameField, e.Name)
	protojson.EncodeString_OmitEmpty(w, core.ResultMessageField, e.Message)
}

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (e *Error) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Status<<core.ErrorCodeBits|e.Code), e.Message)
}

func (e *Error) SetStatus(status uint32) {
	e.Status = status
}

func (e *Error) SetName(name string) {
	e.Name = name
}

func (e *Error) SetMessage(message string) {
	if len(e.Details) > 0 {
		e.Message = fmt.Sprintf(message, core.AnySlice(e.Details)...)
	} else {
		e.Message = message
	}
}

var _ core.Error = (*Error)(nil)
