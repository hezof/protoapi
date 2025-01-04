package protoapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
默认code错误码. 兼容http/grpc的code范围是[0,math.MaxInt32]
统一约定:
- [0,99]              表示保留错误码! (业务/扩展切勿占用)
- [100,999]           表示请求错误码! (与http status code一致)
- [1000,9999]         表示系统错误码!
- [10000,2147483647]  表示业务错误码!
*/

// 对于protoapi.required的全局默认设置,调用者可在server启动前修改设置!

type StatusResult struct {
	Status  uint32 `json:"-"`                 // 状态代码(http).
	Code    uint32 `json:"code"`              // 错误代码. 0表示成功
	Name    string `json:"name,omitempty"`    // 错误名称. OK表示成功
	Message string `json:"message,omitempty"` // 错误消息.
	Details []any  `json:"-"`                 // 错误参数.
	Data    any    `json:"data,omitempty"`    // 结果数据
}

func (sr *StatusResult) DecodeJSON(r *JsonDecoder) {
	panic("unsupported operation")
}

func (sr *StatusResult) EncodeJSON(w *JsonEncoder) {
	EncodeMessage(w, sr, func(w *JsonEncoder, r *StatusResult) {
		EncodeUint32_WithEmpty(w, profile.ResultCodeField, r.Code)
		EncodeString_OmitEmpty(w, profile.ResultNameField, r.Name)
		EncodeString_OmitEmpty(w, profile.ResultMessageField, r.Message)
		EncodeAny_OmitEmpty(w, profile.ResultDataField, r.Data)
	})
}

var _ JsonCodec = (*StatusResult)(nil)

func (sr *StatusResult) Error() string {
	return ToJson(sr)
}

var _ error = (*StatusResult)(nil) // 断言接口

// GRPCStatus 支持status.FromError
func (sr *StatusResult) GRPCStatus() *status.Status {
	return status.New(codes.Code(sr.Code), sr.Message)
}

func StatusError(status uint32, code uint32, message string, args ...interface{}) *StatusResult {
	if len(args) > 0 {
		message = Sprintf(message, args...)
	}
	return &StatusResult{
		Code:    uint32(code),
		Message: message,
		Status:  uint32(status),
		Details: args,
	}
}

// FromError 定义统一的error转换为*Result规则
func FromError(err error, defaultStatus uint32) *StatusResult {
	// 分类处理错误
	if result, ok := err.(*StatusResult); ok {
		result.Status = NvlI(result.Status, defaultStatus)
		// basic error
		return result
	}
	if result, ok := err.(*Error); ok {
		return &StatusResult{
			Status:  NvlI(result.Status, defaultStatus),
			Code:    result.Code,
			Name:    result.Name,
			Message: result.Message,
		}
	}
	if sta, ok := status.FromError(err); ok {
		// grpc error
		return &StatusResult{
			Status:  defaultStatus,
			Code:    uint32(sta.Code()),
			Message: sta.Message(),
		}
	}
	// unknown error
	return &StatusResult{
		Status:  defaultStatus,
		Code:    uint32(codes.Unknown),
		Message: err.Error(),
	}
}
