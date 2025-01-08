package protoapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
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
	// 至少传个指针呀
	if val == nil {
		return errors.New("decode to nil")
	}
	// 加速实现JsonCodec
	if jc, ok := val.(JsonCodec); ok {
		dec := GetDecoder(in)
		defer PutDecoder(dec)

		jc.DecodeJSON(dec)
		return dec.Close()
	}
	// 其他仍用encoding/json
	return json.NewDecoder(in).Decode(val)
}

func DecodeAny(r *JsonDecoder, val any) {
	// 至少传个指针呀
	if val == nil {
		r.reportError(errors.New("decode to nil"))
		return
	}
	// 加速实现JsonCodec
	if jc, ok := val.(JsonCodec); ok {
		jc.DecodeJSON(r)
		return
	}
	// 其他仍用encoding/json
	switch r.token {
	case ObjectBegin:
		err := json.Unmarshal(r.dumpObjectOrArray(), val)
		if err != nil {
			r.reportError(err)
		}
	case ObjectEnd:
		r.invalidCharacterError()
	case ArrayBegin:
		err := json.Unmarshal(r.dumpObjectOrArray(), val)
		if err != nil {
			r.reportError(err)
		}
	case ArrayEnd:
		r.invalidCharacterError()
	case Comma:
		r.invalidCharacterError()
	case Colon:
		r.invalidCharacterError()
	case String:
		rv := reflect.ValueOf(val)
		rt := rv.Type()
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
			if rv.IsNil() {
				if rv.CanSet() {
					rv.Set(reflect.New(rt))
				}
			}
			rv = rv.Elem()
		}
		switch rt.Kind() {
		case reflect.String:
			if rv.CanSet() {
				rv.SetString(r.readString())
			}
		default:
			r.reportError(fmt.Errorf("cannot unmarshal string into type %v", rv.Kind()))
		}
	case Number:
		rv := reflect.ValueOf(val)
		rt := rv.Type()
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
			if rv.IsNil() {
				if rv.CanSet() {
					rv.Set(reflect.New(rt))
				}
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if rv.CanSet() {
				rv.SetInt(r.readInt64())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if rv.CanSet() {
				rv.SetUint(r.readUint64())
			}
		case reflect.Float32, reflect.Float64:
			if rv.CanSet() {
				rv.SetFloat(r.readFloat64())
			}
		default:
			r.reportError(fmt.Errorf("cannot unmarshal number into type %v", rv.Kind()))
		}
	case True:
		r.skipTrue()
		rv := reflect.ValueOf(val)
		rt := rv.Type()
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
			if rv.IsNil() {
				if rv.CanSet() {
					rv.Set(reflect.New(rt))
				}
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Bool:
			if rv.CanSet() {
				rv.SetBool(true)
			}
		default:
			r.reportError(fmt.Errorf("cannot unmarshal true into type %v", rv.Kind()))
		}
	case False:
		r.skipFalse()
		rv := reflect.ValueOf(val)
		rt := rv.Type()
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
			if rv.IsNil() {
				if rv.CanSet() {
					rv.Set(reflect.New(rt))
				}
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Bool:
			if rv.CanSet() {
				rv.SetBool(false)
			}
		default:
			r.reportError(fmt.Errorf("cannot unmarshal false into type %v", rv.Kind()))
		}
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
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

func EncodeAny_OmitEmpty(w *JsonEncoder, name string, val any) {
	if val != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		if jc, ok := val.(JsonCodec); ok {
			jc.EncodeJSON(w)
		} else {
			bs, err := json.Marshal(val)
			if err != nil {
				if w.firstError == nil {
					w.firstError = err
				}
				return
			}
			_, _ = w.Write(bs)
		}
		w.buff = append(w.buff, comma)
	}
}

func EncodeAny_WithEmpty(w *JsonEncoder, name string, val any) {
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
