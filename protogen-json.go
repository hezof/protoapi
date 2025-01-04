package protoapi

import (
	"io"
	"sync"
)

// JsonCodec JSON解码器
type JsonCodec interface {
	DecodeJSON(r *JsonDecoder)
	EncodeJSON(w *JsonEncoder)
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

func DecodeJSON(in io.Reader, val any) error {
	dec := GetDecoder(in)
	defer PutDecoder(dec)

	DecodeAny(dec, val)
	return dec.Close()
}

func EncodeJSON(out io.Writer, val any) error {
	enc := GetEncoder(out)
	defer PutEncoder(enc)

	EncodeAny(enc, val)
	return enc.Close()
}
