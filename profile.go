package protoapi

import (
	"net/http"
	"time"
)

type Profile struct {
	ResultCodeField            string        // code前缀, 默认: `"Code":`, 0表示成功
	ResultNameField            string        // name前缀, 默认: `"Name":`, OK表示成功
	ResultDataField            string        // data前缀, 默认: `"Data":`.
	ResultMessageField         string        // message前缀, 默认: `"Message":`
	DecoderBufferSize          int           // 默认8K
	EncoderBufferSize          int           // 默认8K
	HttpFormMaxMemory          int64         // 32 MB,同gin及多数web框架.
	HttpBodyMaxBytes           int64         // 32 MB,默认请求体的字节数. 注意: 请求体不是响应体, 后者没有限制!
	HttpKeepAlive              time.Duration // 3分钟
	GrpcKeepAlive              time.Duration // 5分钟
	GrpcKeepAlivePolicy        time.Duration // 5分钟
	DefaultApplyStatus         uint32
	DefaultErrorStatus         uint32
	DefaultDecodeErrorCode     uint32
	DefaultDecodeErrorStatus   uint32
	DefaultRequiredErrorStatus uint32
	DefaultRequiredErrorCode   uint32
	DefaultValidateErrorStatus uint32
	DefaultValidateErrorCode   uint32
}

var profile = Profile{
	ResultCodeField:            `Code`,
	ResultNameField:            `Name`,
	ResultDataField:            `Data`,
	ResultMessageField:         `Message`,
	DecoderBufferSize:          8 * 1024,
	EncoderBufferSize:          8 * 1024,
	HttpFormMaxMemory:          32 << 20,
	HttpBodyMaxBytes:           32 << 20,
	HttpKeepAlive:              3 * time.Minute,
	GrpcKeepAlive:              5 * time.Minute,
	GrpcKeepAlivePolicy:        5 * time.Minute,
	DefaultApplyStatus:         uint32(http.StatusOK),
	DefaultErrorStatus:         uint32(http.StatusForbidden),
	DefaultDecodeErrorCode:     uint32(http.StatusBadRequest),
	DefaultDecodeErrorStatus:   uint32(http.StatusBadRequest),
	DefaultRequiredErrorStatus: uint32(http.StatusBadRequest),
	DefaultRequiredErrorCode:   uint32(http.StatusBadRequest),
	DefaultValidateErrorStatus: uint32(http.StatusBadRequest),
	DefaultValidateErrorCode:   uint32(http.StatusBadRequest),
}

func InitProfile(ops ...func(p *Profile)) {
	for _, op := range ops {
		op(&profile)
	}
}
