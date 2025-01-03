package protoapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"net/http"
	"sync"
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
var (
	DefaultApplyStatus         = uint32(http.StatusOK)
	DefaultErrorStatus         = uint32(http.StatusForbidden)
	DefaultDecodeErrorCode     = uint32(http.StatusBadRequest)
	DefaultDecodeErrorStatus   = uint32(http.StatusBadRequest)
	DefaultRequiredErrorStatus = uint32(http.StatusBadRequest)
	DefaultRequiredErrorCode   = uint32(http.StatusBadRequest)
	DefaultValidateErrorStatus = uint32(http.StatusBadRequest)
	DefaultValidateErrorCode   = uint32(http.StatusBadRequest)
)

type StatusResult struct {
	Status  uint32         `json:"-"`                 // 状态代码(http).
	Code    uint32         `json:"code"`              // 错误代码. 0表示成功
	Name    string         `json:"name,omitempty"`    // 错误名称. OK表示成功
	Message string         `json:"message,omitempty"` // 错误消息.
	Details []interface{}  `json:"-"`                 // 错误参数.
	Data    MessageEncoder `json:"data,omitempty"`    // 结果数据
}

func (r *StatusResult) ProtoReflect() protoreflect.Message {
	return r.Data.(protoreflect.Message)
}

func (r *StatusResult) EncodeJSON(w *JsonEncoder) error {
	EncodeMessage(w, r, func(w *JsonEncoder, r *StatusResult) {
		EncodeUint32_WithEmpty(w, profile.ResultCodeField, r.Code)
		EncodeString_OmitEmpty(w, profile.ResultNameField, r.Name)
		EncodeString_OmitEmpty(w, profile.ResultMessageField, r.Message)
		EncodeMessageEncoder(w, profile.ResultDataField, r.Data)
	})
	_, err := w.Close()
	return err
}

var _ MessageEncoder = (*StatusResult)(nil)

func (r *StatusResult) Error() string {
	return ToJson(r)
}

var _ error = (*StatusResult)(nil) // 断言接口

func (r *StatusResult) SetStatus(status int) {
	r.Status = uint32(status)
}

// GRPCStatus 支持status.FromError
func (r *StatusResult) GRPCStatus() *status.Status {
	return status.New(codes.Code(r.Code), r.Message)
}

func StatusError(status uint32, code uint32, message string, args ...interface{}) *StatusResult {
	code2status.Store(code, status) // 附加功能: grpc返回的error没有status.
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

// 附加功能: grpc返回的error没有status.
var code2status sync.Map // {uint32: uint32}

// Status 返回StatusError()StatusResult()定义的status
func Status(code uint32) int {
	if st, ok := code2status.Load(code); ok {
		return int(st.(uint32))
	}
	return 0
}

// StatusErrorFrom 定义统一的error转换为*Result规则
// - err: 错误参数, 不能为空, 否则nil panic
func StatusErrorFrom(err error, defaultStatus uint32) *StatusResult {
	// 分类处理错误
	if result, ok := err.(*StatusResult); ok {
		if result.Status == 0 {
			if vl, ok := code2status.Load(result.Code); ok {
				result.Status = vl.(uint32)
			} else {
				result.Status = defaultStatus
			}
		}
		// basic error
		return result
	} else if sta, ok := status.FromError(err); ok {
		// grpc error
		result = &StatusResult{
			Code:    uint32(sta.Code()),
			Message: sta.Message(),
		}
		if vl, ok := code2status.Load(result.Code); ok {
			result.Status = vl.(uint32)
		} else {
			result.Status = defaultStatus
		}
		return result
	} else {
		// unknown error
		result = &StatusResult{
			Status:  defaultStatus,
			Code:    uint32(codes.Unknown),
			Message: err.Error(),
		}
		return result
	}
}
