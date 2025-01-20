package protojson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func DecodeAny(r *JsonDecoder, val any) {
	// 其他仍用encoding/json
	switch r.token {
	case ObjectBegin:
		if d, ok := val.(JsonCodec); ok {
			// 已实现JsonCodec使用protojson加速解码
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
					d.DecodeJSON(r, k)
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
		} else {
			rt := reflect.TypeOf(val)
			if rt.Kind() == reflect.Map {

			} else {
				r.unreadByte()
				err := json.Unmarshal(r.dumpObjectOrArray(ObjectBegin), val)
				if err != nil {
					r.reportError(err)
				}
			}
		}
	case ObjectEnd:
		r.invalidCharacterError()
	case ArrayBegin:
		r.unreadByte()
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
