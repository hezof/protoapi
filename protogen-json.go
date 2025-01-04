package protoapi

import (
	"io"
)

// JsonCodec JSON解码器
type JsonCodec interface {
	DecodeJSON(r *JsonDecoder)
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

	EncodeAny(enc, val)
	_, err := enc.Close()
	return err
}
