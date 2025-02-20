package protoapi

import (
	"net/http"
	"time"
)

const MAX_MEM = 32 << 20 // 32M

type Profile struct {
	ResultCodeField              string        // code前缀, 默认: `"Code":`, 0表示成功
	ResultNameField              string        // name前缀, 默认: `"Name":`, OK表示成功
	ResultDataField              string        // data前缀, 默认: `"Data":`.
	ResultMessageField           string        // message前缀, 默认: `"Message":`
	DecoderBufferSize            int           // 默认8K
	EncoderBufferSize            int           // 默认8K
	HttpFormMaxMemory            int64         // 32 MB,同gin及多数web框架.
	HttpBodyMaxBytes             int64         // 32 MB,默认请求体的字节数. 注意: 请求体不是响应体, 后者没有限制!
	HttpKeepAlive                time.Duration // 3分钟
	GrpcKeepAlive                time.Duration // 5分钟
	GrpcKeepAlivePolicy          time.Duration // 5分钟
	DefaultApplyStatus           uint32        // 默认成功状态码
	DefaultErrorStatus           uint32        // 默认错误状态码
	DefaultBadRequestErrorCode   uint32        // 默认参数解析错误状态码
	DefaultBadRequestErrorStatus uint32        // 默认参数解析错误代码
}

var profile = Profile{
	ResultCodeField:              `Code`,
	ResultNameField:              `Name`,
	ResultDataField:              `Data`,
	ResultMessageField:           `Message`,
	DecoderBufferSize:            8 * 1024,
	EncoderBufferSize:            8 * 1024,
	HttpFormMaxMemory:            MAX_MEM,
	HttpBodyMaxBytes:             MAX_MEM,
	HttpKeepAlive:                3 * time.Minute,
	GrpcKeepAlive:                5 * time.Minute,
	GrpcKeepAlivePolicy:          5 * time.Minute,
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
