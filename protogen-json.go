package protoapi

import (
	"encoding/json"
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
	if jc, ok := val.(JsonCodec); ok {
		dec := GetDecoder(in)
		defer PutDecoder(dec)

		jc.DecodeJSON(dec)
		return dec.Close()
	} else {
		return json.NewDecoder(in).Decode(val)
	}
}

func EncodeJSON(out io.Writer, val any) error {
	if jc, ok := val.(JsonCodec); ok {
		enc := GetEncoder(out)
		defer PutEncoder(enc)

		jc.EncodeJSON(enc)
		return enc.Close()
	} else {
		enc := json.NewEncoder(out)
		enc.SetEscapeHTML(false)
		return enc.Encode(val)
	}
}

func EncodeJSON_OmitEmpty(w *JsonEncoder, name string, val any) {
	if val != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		if jc, ok := val.(JsonCodec); ok {
			jc.EncodeJSON(w)
		} else {
			enc := json.NewEncoder(w)
			enc.SetEscapeHTML(false)
			err := enc.Encode(val)
			if err != nil {
				if w.firstError == nil {
					w.firstError = err
				}
				return
			}
		}
		w.buff = append(w.buff, comma)
	}
}

func EncodeJSON_WithEmpty(w *JsonEncoder, name string, val any) {
	if val != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		if jc, ok := val.(JsonCodec); ok {
			jc.EncodeJSON(w)
		} else {
			enc := json.NewEncoder(w)
			enc.SetEscapeHTML(false)
			err := enc.Encode(val)
			if err != nil {
				if w.firstError == nil {
					w.firstError = err
				}
				return
			}
		}
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	}
}

// ToJson Json转换快捷方法
func ToJson(v any) string {
	if enc, ok := v.(JsonCodec); ok {
		out := NewJsonEncoder(nil, 1024)
		enc.EncodeJSON(out)
		_ = out.Close()
		return UnsafeString(out.buff)
	} else {
		bs, _ := json.Marshal(v)
		return UnsafeString(bs)
	}
}
