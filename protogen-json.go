package protoapi

import (
	"google.golang.org/protobuf/proto"
	"io"
)

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

func DecodeJSON(in io.Reader, val any) error {
	// 从pool里面取得JsonDecoder

}

func EncodeJSON(out io.Writer, val any) error {
	enc := GetEncoder(out)
	defer PutEncoder(enc)

	if me, ok := val.(MessageEncoder); ok {
		me.EncodeJSON(enc)
	}
	_, err := enc.Close()
	return err
}
