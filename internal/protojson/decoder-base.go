package protojson

import (
	"encoding/base64"
	"encoding/json"
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
	default:
		r.expectedTokenError(String)
	}
}

func DecodeBytes(r *JsonDecoder, p *[]byte) {
	switch r.token {
	case String:
		val, err := base64.StdEncoding.DecodeString(r.readString())
		if err != nil {
			r.reportError(err)
			return
		}
		*p = val
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
	default:
		r.expectedTokenError(String)
	}
}

func DecodeEnumName[Enum ~int32](r *JsonDecoder, p *Enum, names map[int32]string, values map[string]int32) {
	switch r.token {
	case String:
		s := r.readString()
		v, ok := values[s]
		if !ok {
			r.reportError(fmt.Errorf("invalid enum: %v", s))
			return
		}
		*p = Enum(v)
	case Number:
		v := int32(r.readInt64())
		_, ok := names[v]
		if !ok {
			r.reportError(fmt.Errorf("invalid enum: %v", v))
			return
		}
		*p = Enum(v)
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
	default:
		r.expectedTokenError(String)
	}
}

func DecodeEnum[Enum ~int32](r *JsonDecoder, p *Enum, names map[int32]string, values map[string]int32) {
	switch r.token {
	case Number:
		v := int32(r.readInt64())
		_, ok := names[v]
		if !ok {
			r.reportError(fmt.Errorf("invalid enum: %v", v))
			return
		}
		*p = Enum(v)
	case String:
		s := r.readString()
		v, ok := values[s]
		if !ok {
			r.reportError(fmt.Errorf("invalid enum: %v", s))
			return
		}
		*p = Enum(v)
	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
	default:
		r.expectedTokenError(Number)
	}
}

func DecodeMessage[Message any](r *JsonDecoder, p **Message) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = new(Message)
		}
		if jc, ok := any(*p).(JsonCodec); ok {
			// 已实现JsonCodec使用protojson加速解码
			jc.DecodeJSON(r)
		} else {
			// 未实现JsonCodec使用encoding/std反射解码
			r.unreadByte() // 回退"{"
			err := json.Unmarshal(r.dumpObjectOrArray(ObjectBegin), *p)
			if err != nil {
				r.reportError(err)
			}
		}

	case Null:
		r.skipNull()
	case 0:
		r.unexpectedEndError()
	case -1:
	default:
		r.expectedTokenError(ObjectBegin)
	}
}
