package protoapi

type Profile struct {
	ResultCodeField             string // code前缀, 默认: `"code":`, 0表示成功
	ResultNameField             string // name前缀, 默认: `"name":`, OK表示成功
	ResultDataField             string // data前缀, 默认: `"data":`.
	ResultMessageField          string // message前缀, 默认: `"message":`
	DecoderBufferSize           int    // 默认8K
	EncoderBufferSize           int    // 默认8K
	MaximumNestingDepth         int    // limit maximum depth of nesting, as allowed by https://tools.ietf.org/html/rfc7159#section-9
	MinimumBufferLength         int    // limit minimum length of buffer
	MaximumErrorLength          int    // limit maximum length of error
	HandleChainCapacity         int    // HandleChain容许的最大长度上限
	InsensitiveCapacity         int    // 用于忽略大小写的查找过程
	ParseMultipartFormMaxMemory int64  // 32 MB,同gin及多数web框架.
	LimitRequestBodyMaxBytes    int64  // 32 MB,默认请求体的字节数. 注意: 请求体不是响应体, 后者没有限制!
}

var profile = Profile{
	ResultCodeField:             `code`,
	ResultNameField:             `name`,
	ResultDataField:             `data`,
	ResultMessageField:          `message`,
	DecoderBufferSize:           8 * 1024,
	EncoderBufferSize:           8 * 1024,
	MaximumNestingDepth:         128,
	MinimumBufferLength:         1024,
	MaximumErrorLength:          13,
	HandleChainCapacity:         512,
	InsensitiveCapacity:         256,
	ParseMultipartFormMaxMemory: 32 << 20,
	LimitRequestBodyMaxBytes:    32 << 20,
}

func InitProfile(p Profile) {
	profile.ResultCodeField = NvlS(p.ResultCodeField, profile.ResultCodeField)
	profile.ResultNameField = NvlS(p.ResultNameField, profile.ResultNameField)
	profile.ResultDataField = NvlS(p.ResultDataField, profile.ResultDataField)
	profile.ResultMessageField = NvlS(p.ResultMessageField, profile.ResultMessageField)
	profile.DecoderBufferSize = NvlI(p.DecoderBufferSize, profile.DecoderBufferSize)
	profile.EncoderBufferSize = NvlI(p.EncoderBufferSize, profile.EncoderBufferSize)
	profile.MaximumNestingDepth = NvlI(p.MaximumNestingDepth, profile.MaximumNestingDepth)
	profile.MinimumBufferLength = NvlI(p.MinimumBufferLength, profile.MinimumBufferLength)
	profile.MaximumErrorLength = NvlI(p.MaximumErrorLength, profile.MaximumErrorLength)
	profile.HandleChainCapacity = NvlI(p.HandleChainCapacity, profile.HandleChainCapacity)
	profile.InsensitiveCapacity = NvlI(p.InsensitiveCapacity, profile.InsensitiveCapacity)
	profile.ParseMultipartFormMaxMemory = NvlI(p.ParseMultipartFormMaxMemory, profile.ParseMultipartFormMaxMemory)
	profile.LimitRequestBodyMaxBytes = NvlI(p.LimitRequestBodyMaxBytes, profile.LimitRequestBodyMaxBytes)
}
