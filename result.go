package protoapi

import (
	"fmt"
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
	_CodeBits   = 22
	_CodeMask   = 1<<_CodeBits - 1
	_StatusBits = 9
	_StatusMask = 1<<_StatusBits - 1
)

// StatusResult 带状态码的结果
type StatusResult interface {
	error
	FieldCodec
	GRPCStatus() *status.Status
	GetCode() uint32
	SetStatus(status uint32)
	GetStatus() uint32
	SetName(name string)
	GetName() string
	SetMessage(message string)
	GetMessage() string
}

// SimpleResult 带状态的结果. 必须注意status与code的约定取值范围!
type SimpleResult struct {
	Status  uint32   // 状态代码(http).
	Code    uint32   // 错误代码. 0表示成功
	Name    string   // 错误名称. OK表示成功
	Message string   // 错误消息.
	Details []string `json:"-"` // 错误参数.
	Data    any      `json:"-"` // 结果数据
}

func (sr *SimpleResult) Error() string {
	// 不能使用ToJson()会在EncodeField过程形成死循环
	bs, _ := MarshalJSON(sr)
	return UnsafeString(bs)
}
func (sr *SimpleResult) DecodeField(r *JsonDecoder, f string) {
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

func (sr *SimpleResult) EncodeField(w *JsonEncoder) {
	EncodeUint32_WithEmpty(w, profile.ResultCodeField, sr.Code)
	EncodeString_OmitEmpty(w, profile.ResultNameField, sr.Name)
	EncodeString_OmitEmpty(w, profile.ResultMessageField, sr.Message)
	EncodeAny_OmitEmpty(w, profile.ResultDataField, sr.Data)
}

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (sr *SimpleResult) GRPCStatus() *status.Status {
	return status.New(codes.Code(sr.Status<<_CodeBits|sr.Code), sr.Message)
}

func (sr *SimpleResult) GetCode() uint32 {
	return sr.Code
}

func (sr *SimpleResult) GetStatus() uint32 {
	return sr.Status
}

func (sr *SimpleResult) SetStatus(status uint32) {
	sr.Status = status
}

func (sr *SimpleResult) GetName() string {
	return sr.Name
}

func (sr *SimpleResult) SetName(name string) {
	sr.Name = name
}

func (sr *SimpleResult) GetMessage() string {
	return sr.Message
}

func (sr *SimpleResult) SetMessage(message string) {
	if len(sr.Details) > 0 {
		message = fmt.Sprintf(message, sr.Details)
	}
	sr.Message = message
}

var _ StatusResult = (*SimpleResult)(nil)

// StatusError 创建StatusResult错误实例. 必须注意status与code的取值范围:
// - Status 取值范围(0,1024)
// - Code 取值范围(0,4194304)
func StatusError(status uint32, code uint32, message string, details ...string) StatusResult {

	status &= _StatusMask
	code &= _CodeMask

	if len(details) > 0 {
		message = fmt.Sprintf(message, as(details)...)
	}
	return &SimpleResult{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// StatusErrorFrom 定义统一的error转换为*Result规则
func StatusErrorFrom(err error) StatusResult {

	// 内部错误
	if result, ok := err.(StatusResult); ok {
		return result
	}
	// grpc错误
	if sta, ok := status.FromError(err); ok {
		// Grpc Code转换为StatusResult的Status/Code
		return &SimpleResult{
			Status:  uint32(sta.Code()) >> _CodeBits,
			Code:    uint32(sta.Code()) & _CodeMask,
			Message: sta.Message(),
		}
	}

	// 其他错误
	return &SimpleResult{
		Status:  profile.DefaultErrorStatus,
		Code:    uint32(codes.Unknown),
		Message: err.Error(),
	}
}
