package protoapi

import (
	"encoding/json"
	"github.com/hezof/protoapi/internal/protojson"
	"io"
	"sync"
)

type (
	JsonCodec   = protojson.JsonCodec
	JsonDecoder = protojson.JsonDecoder
	JsonEncoder = protojson.JsonEncoder
)

var (
	NewJsonDecoder         = protojson.NewJsonDecoder
	NewJsonBuffer          = protojson.NewJsonBuffer
	DecodeBool             = protojson.DecodeBool
	DecodeInt32            = protojson.DecodeInt32
	DecodeInt64            = protojson.DecodeInt64
	DecodeUint32           = protojson.DecodeUint32
	DecodeUint64           = protojson.DecodeUint64
	DecodeFloat            = protojson.DecodeFloat
	DecodeDouble           = protojson.DecodeDouble
	DecodeString           = protojson.DecodeString
	DecodeBytes            = protojson.DecodeBytes
	DecodeEnumName         = protojson.DecodeEnumName
	DecodeEnum             = protojson.DecodeEnum
	DecodeMessage          = protojson.DecodeMessage
	DecodeBoolOptional     = protojson.DecodeBoolOptional
	DecodeInt32Optional    = protojson.DecodeInt32Optional
	DecodeInt64Optional    = protojson.DecodeInt64Optional
	DecodeUint32Optional   = protojson.DecodeUint32Optional
	DecodeUint64Optional   = protojson.DecodeUint64Optional
	DecodeFloat32Optional  = protojson.DecodeFloat32Optional
	DecodeFloat64Optional  = protojson.DecodeFloat64Optional
	DecodeStringOptional   = protojson.DecodeStringOptional
	DecodeBytesOptional    = protojson.DecodeBytesOptional
	DecodeEnumNameOptional = protojson.DecodeEnumNameOptional
	DecodeEnumOptional     = protojson.DecodeEnumOptional
	DecodeMessageOptional  = protojson.DecodeMessageOptional
	DecodeBoolRepeated     = protojson.DecodeBoolRepeated
	DecodeInt32Repeated    = protojson.DecodeInt32Repeated
	DecodeInt64Repeated    = protojson.DecodeInt64Repeated
	DecodeUint32Repeated   = protojson.DecodeUint32Repeated
	DecodeUint64Repeated   = protojson.DecodeUint64Repeated
	DecodeFloat32Repeated  = protojson.DecodeFloat32Repeated
	DecodeFloat64Repeated  = protojson.DecodeFloat64Repeated
	DecodeStringRepeated   = protojson.DecodeStringRepeated
	DecodeBytesRepeated    = protojson.DecodeBytesRepeated
	DecodeEnumNameRepeated = protojson.DecodeEnumNameRepeated
	DecodeEnumRepeated     = protojson.DecodeEnumRepeated
	DecodeMessageRepeated  = protojson.DecodeMessageRepeated
	DecodeBoolMap          = protojson.DecodeBoolMap
	DecodeInt32Map         = protojson.DecodeInt32Map
	DecodeInt64Map         = protojson.DecodeInt64Map
	DecodeUint32Map        = protojson.DecodeUint32Map
	DecodeUint64Map        = protojson.DecodeUint64Map
	DecodeFloat32Map       = protojson.DecodeFloat32Map
	DecodeDoubleMap        = protojson.DecodeDoubleMap
	DecodeStringMap        = protojson.DecodeStringMap
	DecodeBytesMap         = protojson.DecodeBytesMap
	DecodeEnumNameMap      = protojson.DecodeEnumNameMap
	DecodeEnumMap          = protojson.DecodeEnumMap
	DecodeMessageMap       = protojson.DecodeMessageMap
)

var NewJsonEncoder func(out io.Writer, size int) *JsonEncoder = protojson.NewJsonEncoder
var (
	EncodeBool                 = protojson.EncodeBool
	EncodeBool_OmitEmpty       = protojson.EncodeBool_OmitEmpty
	EncodeBool_WithEmpty       = protojson.EncodeBool_WithEmpty
	EncodeBool_ConvEmpty       = protojson.EncodeBool_ConvEmpty
	EncodeInt32                = protojson.EncodeInt32
	EncodeInt32_OmitEmpty      = protojson.EncodeInt32_OmitEmpty
	EncodeInt32_WithEmpty      = protojson.EncodeInt32_WithEmpty
	EncodeInt32_ConvEmpty      = protojson.EncodeInt32_ConvEmpty
	EncodeInt64                = protojson.EncodeInt64
	EncodeInt64_OmitEmpty      = protojson.EncodeInt64_OmitEmpty
	EncodeInt64_WithEmpty      = protojson.EncodeInt64_WithEmpty
	EncodeInt64_ConvEmpty      = protojson.EncodeInt64_ConvEmpty
	EncodeUint32               = protojson.EncodeUint32
	EncodeUint32_OmitEmpty     = protojson.EncodeUint32_OmitEmpty
	EncodeUint32_WithEmpty     = protojson.EncodeUint32_WithEmpty
	EncodeUint32_ConvEmpty     = protojson.EncodeUint32_ConvEmpty
	EncodeUint64               = protojson.EncodeUint64
	EncodeUint64_OmitEmpty     = protojson.EncodeUint64_OmitEmpty
	EncodeUint64_WithEmpty     = protojson.EncodeUint64_WithEmpty
	EncodeUint64_ConvEmpty     = protojson.EncodeUint64_ConvEmpty
	EncodeFloat                = protojson.EncodeFloat
	EncodeFloat_OmitEmpty      = protojson.EncodeFloat_OmitEmpty
	EncodeFloat_WithEmpty      = protojson.EncodeFloat_WithEmpty
	EncodeFloat_ConvEmpty      = protojson.EncodeFloat_ConvEmpty
	EncodeDouble               = protojson.EncodeDouble
	EncodeDouble_OmitEmpty     = protojson.EncodeDouble_OmitEmpty
	EncodeDouble_WithEmpty     = protojson.EncodeDouble_WithEmpty
	EncodeDouble_ConvEmpty     = protojson.EncodeDouble_ConvEmpty
	EncodeString               = protojson.EncodeString
	EncodeString_OmitEmpty     = protojson.EncodeString_OmitEmpty
	EncodeString_WithEmpty     = protojson.EncodeString_WithEmpty
	EncodeString_ConvEmpty     = protojson.EncodeString_ConvEmpty
	EncodeStringHtml           = protojson.EncodeStringHtml
	EncodeStringHtml_OmitEmpty = protojson.EncodeStringHtml_OmitEmpty
	EncodeStringHtml_WithEmpty = protojson.EncodeStringHtml_WithEmpty
	EncodeStringHtml_ConvEmpty = protojson.EncodeStringHtml_ConvEmpty
	EncodeBytes                = protojson.EncodeBytes
	EncodeBytes_OmitEmpty      = protojson.EncodeBytes_OmitEmpty
	EncodeBytes_WithEmpty      = protojson.EncodeBytes_WithEmpty
	EncodeBytes_ConvEmpty      = protojson.EncodeBytes_ConvEmpty
	EncodeEnumName             = protojson.EncodeEnumName
	EncodeEnumName_OmitEmpty   = protojson.EncodeEnumName_OmitEmpty
	EncodeEnumName_WithEmpty   = protojson.EncodeEnumName_WithEmpty
	EncodeEnumName_ConvEmpty   = protojson.EncodeEnumName_ConvEmpty
	EncodeEnum                 = protojson.EncodeEnum
	EncodeEnum_OmitEmpty       = protojson.EncodeEnum_OmitEmpty
	EncodeEnum_WithEmpty       = protojson.EncodeEnum_WithEmpty
	EncodeEnum_ConvEmpty       = protojson.EncodeEnum_ConvEmpty
	EncodeMessage              = protojson.EncodeMessage
	EncodeMessage_OmitEmpty    = protojson.EncodeMessage_OmitEmpty
	EncodeMessage_WithEmpty    = protojson.EncodeMessage_WithEmpty
	EncodeMessage_ConvEmpty    = protojson.EncodeMessage_ConvEmpty
)
var (
	EncodeBoolOptional                 = protojson.EncodeBoolOptional
	EncodeBoolOptional_OmitEmpty       = protojson.EncodeBoolOptional_OmitEmpty
	EncodeBoolOptional_WithEmpty       = protojson.EncodeBoolOptional_WithEmpty
	EncodeBoolOptional_ConvEmpty       = protojson.EncodeBoolOptional_ConvEmpty
	EncodeInt32Optional                = protojson.EncodeInt32Optional
	EncodeInt32Optional_OmitEmpty      = protojson.EncodeInt32Optional_OmitEmpty
	EncodeInt32Optional_WithEmpty      = protojson.EncodeInt32Optional_WithEmpty
	EncodeInt32Optional_ConvEmpty      = protojson.EncodeInt32Optional_ConvEmpty
	EncodeInt64Optional                = protojson.EncodeInt64Optional
	EncodeInt64Optional_OmitEmpty      = protojson.EncodeInt64Optional_OmitEmpty
	EncodeInt64Optional_WithEmpty      = protojson.EncodeInt64Optional_WithEmpty
	EncodeInt64Optional_ConvEmpty      = protojson.EncodeInt64Optional_ConvEmpty
	EncodeUint32Optional               = protojson.EncodeUint32Optional
	EncodeUint32Optional_OmitEmpty     = protojson.EncodeUint32Optional_OmitEmpty
	EncodeUint32Optional_WithEmpty     = protojson.EncodeUint32Optional_WithEmpty
	EncodeUint32Optional_ConvEmpty     = protojson.EncodeUint32Optional_ConvEmpty
	EncodeUint64Optional               = protojson.EncodeUint64Optional
	EncodeUint64Optional_OmitEmpty     = protojson.EncodeUint64Optional_OmitEmpty
	EncodeUint64Optional_WithEmpty     = protojson.EncodeUint64Optional_WithEmpty
	EncodeUint64Optional_ConvEmpty     = protojson.EncodeUint64Optional_ConvEmpty
	EncodeFloatOptional                = protojson.EncodeFloatOptional
	EncodeFloatOptional_OmitEmpty      = protojson.EncodeFloatOptional_OmitEmpty
	EncodeFloatOptional_WithEmpty      = protojson.EncodeFloatOptional_WithEmpty
	EncodeFloatOptional_ConvEmpty      = protojson.EncodeFloatOptional_ConvEmpty
	EncodeDoubleOptional               = protojson.EncodeDoubleOptional
	EncodeDoubleOptional_OmitEmpty     = protojson.EncodeDoubleOptional_OmitEmpty
	EncodeDoubleOptional_WithEmpty     = protojson.EncodeDoubleOptional_WithEmpty
	EncodeDoubleOptional_ConvEmpty     = protojson.EncodeDoubleOptional_ConvEmpty
	EncodeStringOptional               = protojson.EncodeStringOptional
	EncodeStringOptional_OmitEmpty     = protojson.EncodeStringOptional_OmitEmpty
	EncodeStringOptional_WithEmpty     = protojson.EncodeStringOptional_WithEmpty
	EncodeStringOptional_ConvEmpty     = protojson.EncodeStringOptional_ConvEmpty
	EncodeStringHtmlOptional           = protojson.EncodeStringHtmlOptional
	EncodeStringHtmlOptional_OmitEmpty = protojson.EncodeStringHtmlOptional_OmitEmpty
	EncodeStringHtmlOptional_WithEmpty = protojson.EncodeStringHtmlOptional_WithEmpty
	EncodeStringHtmlOptional_ConvEmpty = protojson.EncodeStringHtmlOptional_ConvEmpty
	EncodeBytesOptional                = protojson.EncodeBytesOptional
	EncodeBytesOptional_OmitEmpty      = protojson.EncodeBytesOptional_OmitEmpty
	EncodeBytesOptional_WithEmpty      = protojson.EncodeBytesOptional_WithEmpty
	EncodeBytesOptional_ConvEmpty      = protojson.EncodeBytesOptional_ConvEmpty
	EncodeEnumNameOptional             = protojson.EncodeEnumNameOptional
	EncodeEnumNameOptional_OmitEmpty   = protojson.EncodeEnumNameOptional_OmitEmpty
	EncodeEnumNameOptional_WithEmpty   = protojson.EncodeEnumNameOptional_WithEmpty
	EncodeEnumNameOptional_ConvEmpty   = protojson.EncodeEnumNameOptional_ConvEmpty
	EncodeEnumOptional                 = protojson.EncodeEnumOptional
	EncodeEnumOptional_OmitEmpty       = protojson.EncodeEnumOptional_OmitEmpty
	EncodeEnumOptional_WithEmpty       = protojson.EncodeEnumOptional_WithEmpty
	EncodeEnumOptional_ConvEmpty       = protojson.EncodeEnumOptional_ConvEmpty
	EncodeMessageOptional              = protojson.EncodeMessageOptional
	EncodeMessageOptional_OmitEmpty    = protojson.EncodeMessageOptional_OmitEmpty
	EncodeMessageOptional_WithEmpty    = protojson.EncodeMessageOptional_WithEmpty
	EncodeMessageOptional_ConvEmpty    = protojson.EncodeMessageOptional_ConvEmpty
)
var (
	EncodeBoolRepeated                 = protojson.EncodeBoolRepeated
	EncodeBoolRepeated_OmitEmpty       = protojson.EncodeBoolRepeated_OmitEmpty
	EncodeBoolRepeated_WithEmpty       = protojson.EncodeBoolRepeated_WithEmpty
	EncodeBoolRepeated_ConvEmpty       = protojson.EncodeBoolRepeated_ConvEmpty
	EncodeInt32Repeated                = protojson.EncodeInt32Repeated
	EncodeInt32Repeated_OmitEmpty      = protojson.EncodeInt32Repeated_OmitEmpty
	EncodeInt32Repeated_WithEmpty      = protojson.EncodeInt32Repeated_WithEmpty
	EncodeInt32Repeated_ConvEmpty      = protojson.EncodeInt32Repeated_ConvEmpty
	EncodeInt64Repeated                = protojson.EncodeInt64Repeated
	EncodeInt64Repeated_OmitEmpty      = protojson.EncodeInt64Repeated_OmitEmpty
	EncodeInt64Repeated_WithEmpty      = protojson.EncodeInt64Repeated_WithEmpty
	EncodeInt64Repeated_ConvEmpty      = protojson.EncodeInt64Repeated_ConvEmpty
	EncodeUint32Repeated               = protojson.EncodeUint32Repeated
	EncodeUint32Repeated_OmitEmpty     = protojson.EncodeUint32Repeated_OmitEmpty
	EncodeUint32Repeated_WithEmpty     = protojson.EncodeUint32Repeated_WithEmpty
	EncodeUint32Repeated_ConvEmpty     = protojson.EncodeUint32Repeated_ConvEmpty
	EncodeUint64Repeated               = protojson.EncodeUint64Repeated
	EncodeUint64Repeated_OmitEmpty     = protojson.EncodeUint64Repeated_OmitEmpty
	EncodeUint64Repeated_WithEmpty     = protojson.EncodeUint64Repeated_WithEmpty
	EncodeUint64Repeated_ConvEmpty     = protojson.EncodeUint64Repeated_ConvEmpty
	EncodeFloatRepeated                = protojson.EncodeFloatRepeated
	EncodeFloatRepeated_OmitEmpty      = protojson.EncodeFloatRepeated_OmitEmpty
	EncodeFloatRepeated_WithEmpty      = protojson.EncodeFloatRepeated_WithEmpty
	EncodeFloatRepeated_ConvEmpty      = protojson.EncodeFloatRepeated_ConvEmpty
	EncodeDoubleRepeated               = protojson.EncodeDoubleRepeated
	EncodeDoubleRepeated_OmitEmpty     = protojson.EncodeDoubleRepeated_OmitEmpty
	EncodeDoubleRepeated_WithEmpty     = protojson.EncodeDoubleRepeated_WithEmpty
	EncodeDoubleRepeated_ConvEmpty     = protojson.EncodeDoubleRepeated_ConvEmpty
	EncodeStringRepeated               = protojson.EncodeStringRepeated
	EncodeStringRepeated_OmitEmpty     = protojson.EncodeStringRepeated_OmitEmpty
	EncodeStringRepeated_WithEmpty     = protojson.EncodeStringRepeated_WithEmpty
	EncodeStringRepeated_ConvEmpty     = protojson.EncodeStringRepeated_ConvEmpty
	EncodeStringHtmlRepeated           = protojson.EncodeStringHtmlRepeated
	EncodeStringHtmlRepeated_OmitEmpty = protojson.EncodeStringHtmlRepeated_OmitEmpty
	EncodeStringHtmlRepeated_WithEmpty = protojson.EncodeStringHtmlRepeated_WithEmpty
	EncodeStringHtmlRepeated_ConvEmpty = protojson.EncodeStringHtmlRepeated_ConvEmpty
	EncodeBytesRepeated                = protojson.EncodeBytesRepeated
	EncodeBytesRepeated_OmitEmpty      = protojson.EncodeBytesRepeated_OmitEmpty
	EncodeBytesRepeated_WithEmpty      = protojson.EncodeBytesRepeated_WithEmpty
	EncodeBytesRepeated_ConvEmpty      = protojson.EncodeBytesRepeated_ConvEmpty
	EncodeEnumNameRepeated             = protojson.EncodeEnumNameRepeated
	EncodeEnumNameRepeated_OmitEmpty   = protojson.EncodeEnumNameRepeated_OmitEmpty
	EncodeEnumNameRepeated_WithEmpty   = protojson.EncodeEnumNameRepeated_WithEmpty
	EncodeEnumNameRepeated_ConvEmpty   = protojson.EncodeEnumNameRepeated_ConvEmpty
	EncodeEnumRepeated                 = protojson.EncodeEnumRepeated
	EncodeEnumRepeated_OmitEmpty       = protojson.EncodeEnumRepeated_OmitEmpty
	EncodeEnumRepeated_WithEmpty       = protojson.EncodeEnumRepeated_WithEmpty
	EncodeEnumRepeated_ConvEmpty       = protojson.EncodeEnumRepeated_ConvEmpty
	EncodeMessageRepeated              = protojson.EncodeMessageRepeated
	EncodeMessageRepeated_OmitEmpty    = protojson.EncodeMessageRepeated_OmitEmpty
	EncodeMessageRepeated_WithEmpty    = protojson.EncodeMessageRepeated_WithEmpty
	EncodeMessageRepeated_ConvEmpty    = protojson.EncodeMessageRepeated_ConvEmpty
)
var (
	EncodeBoolMap                 = protojson.EncodeBoolMap
	EncodeBoolMap_OmitEmpty       = protojson.EncodeBoolMap_OmitEmpty
	EncodeBoolMap_WithEmpty       = protojson.EncodeBoolMap_WithEmpty
	EncodeBoolMap_ConvEmpty       = protojson.EncodeBoolMap_ConvEmpty
	EncodeInt32Map                = protojson.EncodeInt32Map
	EncodeInt32Map_OmitEmpty      = protojson.EncodeInt32Map_OmitEmpty
	EncodeInt32Map_WithEmpty      = protojson.EncodeInt32Map_WithEmpty
	EncodeInt32Map_ConvEmpty      = protojson.EncodeInt32Map_ConvEmpty
	EncodeInt64Map                = protojson.EncodeInt64Map
	EncodeInt64Map_OmitEmpty      = protojson.EncodeInt64Map_OmitEmpty
	EncodeInt64Map_WithEmpty      = protojson.EncodeInt64Map_WithEmpty
	EncodeInt64Map_ConvEmpty      = protojson.EncodeInt64Map_ConvEmpty
	EncodeUint32Map               = protojson.EncodeUint32Map
	EncodeUint32Map_OmitEmpty     = protojson.EncodeUint32Map_OmitEmpty
	EncodeUint32Map_WithEmpty     = protojson.EncodeUint32Map_WithEmpty
	EncodeUint32Map_ConvEmpty     = protojson.EncodeUint32Map_ConvEmpty
	EncodeUint64Map               = protojson.EncodeUint64Map
	EncodeUint64Map_OmitEmpty     = protojson.EncodeUint64Map_OmitEmpty
	EncodeUint64Map_WithEmpty     = protojson.EncodeUint64Map_WithEmpty
	EncodeUint64Map_ConvEmpty     = protojson.EncodeUint64Map_ConvEmpty
	EncodeFloatMap                = protojson.EncodeFloatMap
	EncodeFloatMap_OmitEmpty      = protojson.EncodeFloatMap_OmitEmpty
	EncodeFloatMap_WithEmpty      = protojson.EncodeFloatMap_WithEmpty
	EncodeFloatMap_ConvEmpty      = protojson.EncodeFloatMap_ConvEmpty
	EncodeDoubleMap               = protojson.EncodeDoubleMap
	EncodeDoubleMap_OmitEmpty     = protojson.EncodeDoubleMap_OmitEmpty
	EncodeDoubleMap_WithEmpty     = protojson.EncodeDoubleMap_WithEmpty
	EncodeDoubleMap_ConvEmpty     = protojson.EncodeDoubleMap_ConvEmpty
	EncodeStringMap               = protojson.EncodeStringMap
	EncodeStringMap_OmitEmpty     = protojson.EncodeStringMap_OmitEmpty
	EncodeStringMap_WithEmpty     = protojson.EncodeStringMap_WithEmpty
	EncodeStringMap_ConvEmpty     = protojson.EncodeStringMap_ConvEmpty
	EncodeStringHtmlMap           = protojson.EncodeStringHtmlMap
	EncodeStringHtmlMap_OmitEmpty = protojson.EncodeStringHtmlMap_OmitEmpty
	EncodeStringHtmlMap_WithEmpty = protojson.EncodeStringHtmlMap_WithEmpty
	EncodeStringHtmlMap_ConvEmpty = protojson.EncodeStringHtmlMap_ConvEmpty
	EncodeBytesMap                = protojson.EncodeBytesMap
	EncodeBytesMap_OmitEmpty      = protojson.EncodeBytesMap_OmitEmpty
	EncodeBytesMap_WithEmpty      = protojson.EncodeBytesMap_WithEmpty
	EncodeBytesMap_ConvEmpty      = protojson.EncodeBytesMap_ConvEmpty
	EncodeEnumNameMap             = protojson.EncodeEnumNameMap
	EncodeEnumNameMap_OmitEmpty   = protojson.EncodeEnumNameMap_OmitEmpty
	EncodeEnumNameMap_WithEmpty   = protojson.EncodeEnumNameMap_WithEmpty
	EncodeEnumNameMap_ConvEmpty   = protojson.EncodeEnumNameMap_ConvEmpty
	EncodeEnumMap                 = protojson.EncodeEnumMap
	EncodeEnumMap_OmitEmpty       = protojson.EncodeEnumMap_OmitEmpty
	EncodeEnumMap_WithEmpty       = protojson.EncodeEnumMap_WithEmpty
	EncodeEnumMap_ConvEmpty       = protojson.EncodeEnumMap_ConvEmpty
	EncodeMessageMap              = protojson.EncodeMessageMap
	EncodeMessageMap_OmitEmpty    = protojson.EncodeMessageMap_OmitEmpty
	EncodeMessageMap_WithEmpty    = protojson.EncodeMessageMap_WithEmpty
	EncodeMessageMap_ConvEmpty    = protojson.EncodeMessageMap_ConvEmpty
)

func DecodeJSON(in io.Reader, val any) error {
	// 加速实现JsonCodec
	if jc, ok := val.(JsonCodec); ok {
		dec := GetDecoder(in)
		defer PutDecoder(dec)

		DecodeMessage(dec, jc)
		return dec.Close()
	}
	// 其他仍用encoding/json
	return json.NewDecoder(in).Decode(val)
}

var decoders = sync.Pool{
	New: func() interface{} {
		return NewJsonDecoder(nil, profile.DecoderBufferSize)
	},
}

var encoders = sync.Pool{
	New: func() interface{} {
		return NewJsonEncoder(nil, profile.EncoderBufferSize)
	},
}

func GetDecoder(in io.Reader) *JsonDecoder {
	return decoders.Get().(*JsonDecoder).reset(in)
}

func PutDecoder(dec *JsonDecoder) {
	decoders.Put(dec.clean())
}

func GetEncoder(out io.Writer) *JsonEncoder {
	return encoders.Get().(*JsonEncoder).reset(out)
}

func PutEncoder(enc *JsonEncoder) {
	encoders.Put(enc.clean())
}
