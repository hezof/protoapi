package protoapi

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"io"
	"unsafe"
)

// MaximumNestingDepth limit maximum depth of nesting, as allowed by https://tools.ietf.org/html/rfc7159#section-9
const MaximumNestingDepth = 128

// MinimumBufferLength limit minimum length of buffer
const MinimumBufferLength = 1024

// MaximumErrorLength limit maximum length of error
const MaximumErrorLength = 13

// MessageDecoder Message的解码器
type MessageDecoder interface {
	proto.Message
	DecodeJSON(r *JsonDecoder)
}

// MessageEncoder Message的编码器
type MessageEncoder interface {
	proto.Message
	EncodeJSON(w *JsonEncoder)
}

// ToJson Json转换快捷方法
func ToJson(v any) string {
	if vm, ok := v.(MessageEncoder); ok {
		w := NewJsonEncoder(nil, 1024)
		vm.EncodeJSON(w)
		bs, _ := w.Close()
		return UnsafeString(bs)
	}
	bs, _ := json.Marshal(v)
	return UnsafeString(bs)
}

// UnsafeBytes string到[]byte的不安全转换
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString []byte到string的不安全转换
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

var (
	JsonDecoderBufferSize = 8 * 1024 // 默认8K
	JsonEncoderBufferSize = 8 * 1024 // 默认8K
)

func InitJsonBuffer(jsonDecoderBufferSize, jsonEncodeBufferSize int) {
	JsonDecoderBufferSize = jsonDecoderBufferSize
	JsonEncoderBufferSize = jsonEncodeBufferSize
}

func NewDecoder(in io.Reader) *JsonDecoder {
	return NewJsonDecoder(in, JsonDecoderBufferSize)
}

func NewEncoder(out io.Writer) *JsonEncoder {
	return NewJsonEncoder(out, JsonEncoderBufferSize)
}

func JsonBind(ctx *Context, val MessageDecoder, body Body) error {
	// 从pool里面取得JsonDecoder

}
