package protoapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
StatusResult统一结果与错误的数据结构, 并实现与Grpc Error的转换.
因为Grpc Error Status只有Code字段, 约定StatusResult的Status/Code分别存储在高10位与低22位!

约定StatusResult Code取值范围:
- [0,17)         表示保留错误码! grpc内置错误码, 参考codes._maxCode
- [17,4194304)   表示业务错误码! 最大值(2^22 - 1)! 因为Grpc Code的前10位用于表示StatusResult Status

约定StatusResult Status取值范围:
- (0,1023]
*/
const (
	_CodeBits   = 22
	_CodeMask   = 2 << (_CodeBits - 1)
	_StatusBits = 10
	_StatusMask = 2 << (_StatusBits - 1)
)

// 对于protoapi.required的全局默认设置,调用者可在server启动前修改设置!

type StatusResult struct {
	status  uint32 // 状态代码(http).
	code    uint32 // 错误代码. 0表示成功
	name    string // 错误名称. OK表示成功
	message string // 错误消息.
	details []any  // 错误参数.
	data    any    // 结果数据
}

func (sr *StatusResult) GetStatus() uint32 {
	return sr.status
}

func (sr *StatusResult) GetCode() uint32 {
	return sr.code
}

func (sr *StatusResult) GetName() string {
	return sr.name
}

func (sr *StatusResult) GetMessage() string {
	return sr.message
}

func (sr *StatusResult) DecodeField(r *JsonDecoder, f string) {
	switch f {
	case profile.ResultCodeField:
		DecodeUint32(r, &sr.code)
	case profile.ResultNameField:
		DecodeString(r, &sr.name)
	case profile.ResultMessageField:
		DecodeString(r, &sr.message)
	case profile.ResultDataField:
		DecodeAny(r, &sr.data)
	}
}

func (sr *StatusResult) EncodeField(w *JsonEncoder) {
	EncodeUint32_WithEmpty(w, profile.ResultCodeField, sr.code)
	EncodeString_OmitEmpty(w, profile.ResultNameField, sr.name)
	EncodeString_OmitEmpty(w, profile.ResultMessageField, sr.message)
	EncodeAny_OmitEmpty(w, profile.ResultDataField, sr.data)
}

var _ FieldCodec = (*StatusResult)(nil)

func (sr *StatusResult) Error() string {
	return ToJson(sr)
}

var _ error = (*StatusResult)(nil) // 断言接口

// GRPCStatus 支持status.FromError, 实现StatusResult到grpc Status的转换
func (sr *StatusResult) GRPCStatus() *status.Status {
	return status.New(codes.Code(sr.status<<_CodeBits|sr.code), sr.message)
}

// StatusError 创建StatusResult错误实例. 要求:
// - status 取值范围(0,1023]
// - code 取值范围(0,)
func StatusError(status uint32, code uint32, message string, details ...interface{}) *StatusResult {

	status &= _StatusMask
	code &= _CodeMask

	if len(details) > 0 {
		message = Sprintf(message, details...)
	}
	return &StatusResult{
		status:  status,
		code:    code,
		message: message,
		details: details,
	}
}

// StatusErrorFrom 定义统一的error转换为*Result规则
func StatusErrorFrom(err error) *StatusResult {
	// 分类处理错误
	if result, ok := err.(*StatusResult); ok {
		return result
	}
	if result, ok := err.(*Error); ok {
		return &StatusResult{
			status:  result.Status,
			code:    result.Code,
			name:    result.Name,
			message: result.Message,
		}
	}
	if sta, ok := status.FromError(err); ok {
		// Grpc Code转换为StatusResult的Status/Code
		return &StatusResult{
			status:  uint32(sta.Code()) >> _CodeBits,
			code:    uint32(sta.Code()) & _CodeMask,
			message: sta.Message(),
		}
	}
	// unknown error
	return &StatusResult{
		status:  profile.DefaultErrorStatus,
		code:    uint32(codes.Unknown),
		message: err.Error(),
	}
}
