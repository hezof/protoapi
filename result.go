package protoapi

import (
	"fmt"
	"github.com/hezof/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
StatusResult统一结果与错误的数据结构, 并实现与Grpc Error的转换.
因为Grpc Error Status只有Code字段, 约定StatusResult的Status/Code分别存储在高10位与低22位!

约定StatusResult Code取值范围:
- [0,17)         表示保留错误码! grpc内置错误码, 参考codes._maxCode
- [17,4194393)   表示业务错误码! 最大值(2^22 - 1)! 因为Grpc Code的前10位用于表示StatusResult Status

约定StatusResult Status取值范围:
- (0,511]
*/
const (
	_CodeBits   = core.CodeBits
	_CodeMask   = core.CodeMask
	_StatusBits = core.StatusBits
	_StatusMask = core.StatusMask
)

// StatusResult 带状态码的结果
type StatusResult interface {
	core.Error
	SetStatus(status uint32)
	SetName(name string)
	SetMessage(message string)
	GRPCStatus() *status.Status
}

// StatusResultModel 带状态的结果. 必须注意status与code的约定取值范围!
type StatusResultModel struct {
	Status  uint32   // 状态代码(http).
	Code    uint32   // 错误代码. 0表示成功
	Name    string   // 错误名称. OK表示成功
	Message string   // 错误消息.
	Details []string `json:"-"` // 错误参数.
	Data    any      `json:"-"` // 结果数据
}

func (sr *StatusResultModel) Error() string {
	// 不能使用ToJson()会在EncodeField过程形成死循环
	return core.ToJson(sr)
}

func (sr *StatusResultModel) GetCode() uint32 {
	return sr.Code
}

func (sr *StatusResultModel) GetStatus() uint32 {
	return sr.Status
}

func (sr *StatusResultModel) SetStatus(status uint32) {
	sr.Status = status
}

func (sr *StatusResultModel) GetName() string {
	return sr.Name
}

func (sr *StatusResultModel) SetName(name string) {
	sr.Name = name
}

func (sr *StatusResultModel) GetMessage() string {
	return sr.Message
}

func (sr *StatusResultModel) SetMessage(message string) {
	if len(sr.Details) > 0 {
		message = fmt.Sprintf(message, sr.Details)
	}
	sr.Message = message
}

func (sr *StatusResultModel) GetDetails() []string {
	return sr.Details
}

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (sr *StatusResultModel) GRPCStatus() *status.Status {
	return status.New(codes.Code(sr.Status<<_CodeBits|sr.Code), sr.Message)
}

var _ StatusResult = (*StatusResultModel)(nil)

func (sr *StatusResultModel) DecodeField(r *JsonDecoder, f string) {
	switch f {
	case profile.ResultCodeField:
		DecodeUint32(r, &sr.Code)
	case profile.ResultNameField:
		DecodeString(r, &sr.Name)
	case profile.ResultMessageField:
		DecodeString(r, &sr.Message)
	case profile.ResultDataField:
		DecodeAny(r, sr.Data)
	}
}

func (sr *StatusResultModel) EncodeField(w *JsonEncoder) {
	EncodeUint32_WithEmpty(w, profile.ResultCodeField, sr.Code)
	EncodeString_OmitEmpty(w, profile.ResultNameField, sr.Name)
	EncodeString_OmitEmpty(w, profile.ResultMessageField, sr.Message)
	EncodeAny_OmitEmpty(w, profile.ResultDataField, sr.Data)
}

var _ FieldCodec = (*StatusResultModel)(nil)

// StatusError 创建StatusResult错误实例. 必须注意status与code的取值范围:
// - Status 取值范围(0,1024)
// - Code 取值范围(0,4194304)
func StatusError(status uint32, code uint32, message string, details ...string) StatusResult {

	status &= _StatusMask
	code &= _CodeMask

	if len(details) > 0 {
		message = fmt.Sprintf(message, core.AnySlice(details)...)
	}
	return &StatusResultModel{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// StatusErrorFrom 定义统一的error转换为*Result规则
func StatusErrorFrom(err error) StatusResult {

	// 内部错误
	if val, ok := err.(StatusResult); ok {
		return val
	}
	if val, ok := err.(core.Error); ok {
		return &StatusResultModel{
			Status:  val.GetStatus(),
			Code:    val.GetCode(),
			Message: val.GetMessage(),
			Details: val.GetDetails(),
		}
	}
	// grpc错误
	if sta, ok := status.FromError(err); ok {
		// Grpc Code转换为StatusResult的Status/Code
		return &StatusResultModel{
			Status:  uint32(sta.Code()) >> _CodeBits,
			Code:    uint32(sta.Code()) & _CodeMask,
			Message: sta.Message(),
		}
	}

	// 其他错误
	return &StatusResultModel{
		Status:  profile.DefaultErrorStatus,
		Code:    uint32(codes.Unknown),
		Message: err.Error(),
	}
}
