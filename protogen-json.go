package protoapi

import (
	"google.golang.org/protobuf/proto"
	"io"
)

// MessageDecoder Message的解码器
type MessageDecoder interface {
	proto.Message
	DecodeJSON(r *JsonDecoder) error
}

// MessageEncoder Message的编码器
type MessageEncoder interface {
	proto.Message
	EncodeJSON(w *JsonEncoder) error
}

func GetDecoder(in io.Reader) *JsonDecoder {
	// TODO: 使用pool缓存
	return NewJsonDecoder(in, profile.DecoderBufferSize)
}

func GetEncoder(out io.Writer) *JsonEncoder {
	// TODO: 使用pool缓存
	return NewJsonEncoder(out, profile.EncoderBufferSize)
}

func PutDecoder(dec *JsonDecoder) {

}

func PutEncoder(enc *JsonEncoder) {

}

func DecodeJSON(ctx *Context, val MessageDecoder, body Body) error {
	// 从pool里面取得JsonDecoder

}
