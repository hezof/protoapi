package protojson

import (
	"fmt"
	"github.com/hezof/protoapi/internal/json" // 处理escapeHtml的问题
	"reflect"
)

/*
JsonCodec 核心接口, 实现struct的解码与编码.
该接口用于加速proto.Message的JSON解码/编码速度.
*/
type JsonCodec interface {
	DecodeJSON(r *JsonDecoder)
	EncodeJSON(w *JsonEncoder)
}

// VisitObject 用于DecodeJSON内部遍历的辅助函数
func VisitObject(r *JsonDecoder, f func(r *JsonDecoder, k string)) {
	r.token = 0 // 指示next()执行"step info"而不是"step over"
	t := r.next()
	if t == ObjectEnd {
		return
	}
	for {
		if t != String {
			r.expectedTokenError(String)
			return
		}
		k := r.readString()
		if r.next() != Colon {
			r.expectedTokenError(Colon)
			return
		}

		switch r.next() {
		case 0:
			r.unexpectedEndError()
			return
		case Null:
			r.skipNull()
		default:
			f(r, k)
		}

		t = r.next()
		switch t {
		case Comma:
			t = r.next()
			if t == ObjectEnd {
				r.invalidCharacterError()
				return
			}
		case ObjectEnd:
			return
		default:
			r.invalidCharacterError()
			return
		}
	}
}

func DecodeAny(r *JsonDecoder, val any) {
	// 加速实现JsonCodec
	if jc, ok := val.(JsonCodec); ok {
		jc.DecodeJSON(r)
		return
	}
	// 其他仍用encoding/json
	switch r.token {
	case ObjectBegin:
		err := json.Unmarshal(r.dumpObjectOrArray(ObjectBegin), val)
		if err != nil {
			r.reportError(err)
		}
	case ObjectEnd:
		r.invalidCharacterError()
	case ArrayBegin:
		err := json.Unmarshal(r.dumpObjectOrArray(ArrayBegin), val)
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

func EncodeAny(w *JsonEncoder, val any) {
	if jc, ok := val.(JsonCodec); ok {
		jc.EncodeJSON(w)
	} else {
		bs, err := json.Marshal(val)
		if err != nil {
			w.reportError(err)
		}
		_, _ = w.Write(bs)
	}
}

func EncodeAny_OmitEmpty(w *JsonEncoder, name string, val any) {
	if val != nil {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		EncodeAny(w, val)
		w.buff = append(w.buff, comma)
	}
}

func EncodeAny_WithEmpty(w *JsonEncoder, name string, val any) {
	if val != nil {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		EncodeAny(w, val)
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	}
}
