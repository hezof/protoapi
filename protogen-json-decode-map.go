package protoapi

func DecodeBoolMap(r *JsonDecoder, p *map[string]bool) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]bool)
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
			r.next()

			var v bool
			DecodeBool(r, &v)
			(*p)[k] = v

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

func DecodeInt32Map(r *JsonDecoder, p *map[string]int32) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]int32)
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
			r.next()

			var v int32
			DecodeInt32(r, &v)
			(*p)[k] = v

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

func DecodeInt64Map(r *JsonDecoder, p *map[string]int64) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]int64)
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
			r.next()

			var v int64
			DecodeInt64(r, &v)
			(*p)[k] = v

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

func DecodeUint32Map(r *JsonDecoder, p *map[string]uint32) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]uint32)
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
			r.next()

			var v uint32
			DecodeUint32(r, &v)
			(*p)[k] = v

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

func DecodeUint64Map(r *JsonDecoder, p *map[string]uint64) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]uint64)
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
			r.next()

			var v uint64
			DecodeUint64(r, &v)
			(*p)[k] = v

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

func DecodeFloat32Map(r *JsonDecoder, p *map[string]float32) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]float32)
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
			r.next()

			var v float32
			DecodeFloat(r, &v)
			(*p)[k] = v

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

func DecodeDoubleMap(r *JsonDecoder, p *map[string]float64) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]float64)
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
			r.next()

			var v float64
			DecodeDouble(r, &v)
			(*p)[k] = v

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

func DecodeStringMap(r *JsonDecoder, p *map[string]string) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]string)
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
			r.next()

			var v string
			DecodeString(r, &v)
			(*p)[k] = v

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

func DecodeBytesMap(r *JsonDecoder, p *map[string][]byte) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string][]byte)
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
			r.next()

			var v []byte
			DecodeBytes(r, &v)
			(*p)[k] = v

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

func DecodeEnumMap[Enum ~int32](r *JsonDecoder, p *map[string]Enum, names map[int32]string, values map[string]int32) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]Enum)
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
			r.next()

			var v Enum
			DecodeEnum(r, &v, names, values)
			(*p)[k] = v

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

func DecodeEnumMap_EnumAsInt[Enum ~int32](r *JsonDecoder, p *map[string]Enum, names map[int32]string, values map[string]int32) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]Enum)
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
			r.next()

			var v Enum
			DecodeEnum_EnumAsInt(r, &v, names, values)
			(*p)[k] = v

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

func DecodeMessageMap[Message any](r *JsonDecoder, p *map[string]*Message, h func(r *JsonDecoder, m *Message, f string)) {
	switch r.token {
	case ObjectBegin:
		if *p == nil {
			*p = make(map[string]*Message)
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
			r.next()

			var v *Message
			DecodeMessage(r, &v, h)
			(*p)[k] = v

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
