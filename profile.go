package protoapi

import (
	"net/http"
)

const DefMaxMem = 32 << 20 // 32M

type Profile struct {
	DefaultApplyStatus           uint32 // 默认成功状态码
	DefaultErrorStatus           uint32 // 默认错误状态码
	DefaultBadRequestErrorCode   uint32 // 默认参数解析错误状态码
	DefaultBadRequestErrorStatus uint32 // 默认参数解析错误代码
}

var profile = Profile{
	DefaultApplyStatus:           uint32(http.StatusOK),
	DefaultErrorStatus:           uint32(http.StatusForbidden),
	DefaultBadRequestErrorCode:   uint32(http.StatusBadRequest),
	DefaultBadRequestErrorStatus: uint32(http.StatusBadRequest),
}

func InitProfile(ops ...func(p *Profile)) {
	for _, op := range ops {
		op(&profile)
	}
}
