package protoapi

import (
	"fmt"
	"github.com/hezof/core"
	"github.com/hezof/protojson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatusResult 带状态的结果. 必须注意status与code的约定取值范围!
type StatusResult core.StatusResult

func (sr *StatusResult) Error() string {
	// 不能使用ToJson()会在EncodeField过程形成死循环
	return core.ToJson(sr)
}

func (sr *StatusResult) GetCode() uint32 {
	return sr.Code
}

func (sr *StatusResult) GetStatus() uint32 {
	return sr.Status
}

func (sr *StatusResult) SetStatus(status uint32) {
	sr.Status = status
}

func (sr *StatusResult) GetName() string {
	return sr.Name
}

func (sr *StatusResult) SetName(name string) {
	sr.Name = name
}

func (sr *StatusResult) GetMessage() string {
	return sr.Message
}

func (sr *StatusResult) SetMessage(message string) {
	if len(sr.Details) > 0 {
		message = fmt.Sprintf(message, sr.Details)
	}
	sr.Message = message
}

func (sr *StatusResult) GetDetails() []string {
	return sr.Details
}

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (sr *StatusResult) GRPCStatus() *status.Status {
	/*
		StatusResult统一结果与错误的数据结构, 并实现与Grpc Error的转换.
		因为Grpc Error Status只有Code字段, 约定StatusResult的Status/Code分别存储在高10位与低22位!

		约定StatusResult Code取值范围:
		- [0,17)         表示保留错误码! grpc内置错误码, 参考codes._maxCode
		- [17,4194393)   表示业务错误码! 最大值(2^22 - 1)! 因为Grpc Code的前10位用于表示StatusResult Status

		约定StatusResult Status取值范围:
		- (0,511]
	*/
	return status.New(codes.Code(sr.Status<<core.ErrorCodeBits|sr.Code), sr.Message)
}

var _ core.Error = (*StatusResult)(nil)

func (sr *StatusResult) DecodeField(r *protojson.JsonDecoder, f string) {
	switch f {
	case core.ResultCodeField:
		protojson.DecodeUint32(r, &sr.Code)
	case core.ResultNameField:
		protojson.DecodeString(r, &sr.Name)
	case core.ResultMessageField:
		protojson.DecodeString(r, &sr.Message)
	case core.ResultDataField:
		protojson.DecodeAny(r, sr.Data)
	}
}

func (sr *StatusResult) EncodeField(w *protojson.JsonEncoder) {
	protojson.EncodeUint32_WithEmpty(w, core.ResultCodeField, sr.Code)
	protojson.EncodeString_OmitEmpty(w, core.ResultNameField, sr.Name)
	protojson.EncodeString_OmitEmpty(w, core.ResultMessageField, sr.Message)
	protojson.EncodeAny_OmitEmpty(w, core.ResultDataField, sr.Data)
}

var _ protojson.FieldCodec = (*StatusResult)(nil)

// StatusErrorFrom 定义统一的error转换为*Result规则
func StatusErrorFrom(err error) *StatusResult {

	// 内部错误
	if val, ok := err.(*StatusResult); ok {
		return val
	}

	// 错误框架
	if val, ok := err.(*core.StatusResult); ok {
		return (*StatusResult)(val)
	}

	// 外部错误
	if val, ok := err.(core.Error); ok {
		return &StatusResult{
			Status:  val.GetStatus(),
			Code:    val.GetCode(),
			Message: val.GetMessage(),
			Details: val.GetDetails(),
		}
	}

	// grpc错误
	if sta, ok := status.FromError(err); ok {
		/*
			StatusResult统一结果与错误的数据结构, 并实现与Grpc Error的转换.
			因为Grpc Error Status只有Code字段, 约定StatusResult的Status/Code分别存储在高10位与低22位!

			约定StatusResult Code取值范围:
			- [0,17)         表示保留错误码! grpc内置错误码, 参考codes._maxCode
			- [17,4194393)   表示业务错误码! 最大值(2^22 - 1)! 因为Grpc Code的前10位用于表示StatusResult Status

			约定StatusResult Status取值范围:
			- (0,511]
		*/
		return &StatusResult{
			Status:  uint32(sta.Code()) >> core.ErrorCodeBits,
			Code:    uint32(sta.Code()) & core.ErrorCodeMask,
			Message: sta.Message(),
		}
	}

	// 其他错误
	return &StatusResult{
		Status:  profile.DefaultErrorStatus,
		Code:    uint32(codes.Unknown),
		Message: err.Error(),
	}
}
