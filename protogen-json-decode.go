package protoapi

import (
	"encoding/base64"
	"fmt"
)

func DecodeBool(r *JsonDecoder, p *bool) {
	switch r.token {
	case True:
		r.skipTrue()
		*p = true
	case False:
		r.skipFalse()
		*p = false
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(True)
	}
}

func DecodeInt32(r *JsonDecoder, p *int32) {
	switch r.token {
	case Number:
		*p = int32(r.readInt64())
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(Number)
	}
}

func DecodeInt64(r *JsonDecoder, p *int64) {
	switch r.token {
	case Number:
		*p = r.readInt64()
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(Number)
	}
}

func DecodeUint32(r *JsonDecoder, p *uint32) {
	switch r.token {
	case Number:
		*p = uint32(r.readUint64())
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(Number)
	}
}

func DecodeUint64(r *JsonDecoder, p *uint64) {
	switch r.token {
	case Number:
		*p = r.readUint64()
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(Number)
	}
}

func DecodeFloat(r *JsonDecoder, p *float32) {
	switch r.token {
	case Number:
		*p = float32(r.readFloat64())
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(Number)
	}
}

func DecodeDouble(r *JsonDecoder, p *float64) {
	switch r.token {
	case Number:
		*p = r.readFloat64()
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(Number)
	}
}

func DecodeString(r *JsonDecoder, p *string) {
	switch r.token {
	case String:
		*p = r.readString()
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(String)
	}
}

func DecodeBytes(r *JsonDecoder, p *[]byte) {
	switch r.token {
	case String:
		val, err := base64.StdEncoding.DecodeString(r.readString())
		if err != nil && r.firstError == nil {
			r.firstError = err
		}
		*p = val
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(String)
	}
}

func DecodeEnum[Enum ~int32](r *JsonDecoder, p *Enum, names map[int32]string, values map[string]int32) {
	switch r.token {
	case String:
		s := r.readString()
		v, ok := values[s]
		if !ok && r.firstError == nil {
			r.firstError = fmt.Errorf("invalid enum: %s", s)
			return
		}
		*p = Enum(v)
	case Number:
		v := int32(r.readInt64())
		_, ok := names[v]
		if !ok && r.firstError == nil {
			r.firstError = fmt.Errorf("invalid enum: %d", v)
			return
		}
		*p = Enum(v)
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(String)
	}
}

func DecodeEnum_EnumAsInt[Enum ~int32](r *JsonDecoder, p *Enum, names map[int32]string, values map[string]int32) {
	switch r.token {
	case Number:
		v := int32(r.readInt64())
		_, ok := names[v]
		if !ok && r.firstError == nil {
			r.firstError = fmt.Errorf("invalid enum: %d", v)
			return
		}
		*p = Enum(v)
	case String:
		s := r.readString()
		v, ok := values[s]
		if !ok && r.firstError == nil {
			r.firstError = fmt.Errorf("invalid enum: %s", s)
			return
		}
		*p = Enum(v)
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(String)
	}
}

func DecodeMessage[Message any](r *JsonDecoder, p **Message, h func(r *JsonDecoder, m *Message, k string)) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = new(Message)
		}
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
				h(r, *p, k)
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
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
		return
	default:
		r.expectedTokenError(ObjectBegin)
	}
}
