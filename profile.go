package protoapi

import (
	"net/http"
	"time"
)

type Profile struct {
	ResultCodeField            string        // code前缀, 默认: `"code":`, 0表示成功
	ResultNameField            string        // name前缀, 默认: `"name":`, OK表示成功
	ResultDataField            string        // data前缀, 默认: `"data":`.
	ResultMessageField         string        // message前缀, 默认: `"message":`
	DecoderBufferSize          int           // 默认8K
	EncoderBufferSize          int           // 默认8K
	MaximumNestingDepth        int           // limit maximum depth of nesting, as allowed by https://tools.ietf.org/html/rfc7159#section-9
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
	ResultCodeField:            `code`,
	ResultNameField:            `name`,
	ResultDataField:            `data`,
	ResultMessageField:         `message`,
	DecoderBufferSize:          8 * 1024,
	EncoderBufferSize:          8 * 1024,
	MaximumNestingDepth:        128,
	HandleChainCapacity:        512,
	InsensitiveCapacity:        256,
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

func InitProfile(p Profile) {
	profile.ResultCodeField = NvlS(p.ResultCodeField, profile.ResultCodeField)
	profile.ResultNameField = NvlS(p.ResultNameField, profile.ResultNameField)
	profile.ResultDataField = NvlS(p.ResultDataField, profile.ResultDataField)
	profile.ResultMessageField = NvlS(p.ResultMessageField, profile.ResultMessageField)
	profile.DecoderBufferSize = NvlI(p.DecoderBufferSize, profile.DecoderBufferSize)
	profile.EncoderBufferSize = NvlI(p.EncoderBufferSize, profile.EncoderBufferSize)
	profile.MaximumNestingDepth = NvlI(p.MaximumNestingDepth, profile.MaximumNestingDepth)
	profile.HandleChainCapacity = NvlI(p.HandleChainCapacity, profile.HandleChainCapacity)
	profile.InsensitiveCapacity = NvlI(p.InsensitiveCapacity, profile.InsensitiveCapacity)
	profile.HttpFormMaxMemory = NvlI(p.HttpFormMaxMemory, profile.HttpFormMaxMemory)
	profile.HttpBodyMaxBytes = NvlI(p.HttpBodyMaxBytes, profile.HttpBodyMaxBytes)
}
