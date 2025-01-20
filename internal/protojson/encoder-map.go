package protojson

/*************************************
	bool类型: MapBoolMap_<empty>
 *************************************/

func EncodeBoolMap(w *JsonEncoder, value map[string]bool) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeBoolMap_OmitEmpty(w *JsonEncoder, name string, value map[string]bool) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeBoolMap_WithEmpty(w *JsonEncoder, name string, value map[string]bool) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeBoolMap_ConvEmpty(w *JsonEncoder, name string, value map[string]bool) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	int32类型: MapInt32Map_<empty>
 *************************************/

func EncodeInt32Map(w *JsonEncoder, value map[string]int32) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeInt32Map_OmitEmpty(w *JsonEncoder, name string, value map[string]int32) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeInt32Map_WithEmpty(w *JsonEncoder, name string, value map[string]int32) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeInt32Map_ConvEmpty(w *JsonEncoder, name string, value map[string]int32) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	int64类型: MapInt64Map_<empty>
 *************************************/

func EncodeInt64Map(w *JsonEncoder, value map[string]int64) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeInt64Map_OmitEmpty(w *JsonEncoder, name string, value map[string]int64) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeInt64Map_WithEmpty(w *JsonEncoder, name string, value map[string]int64) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeInt64Map_ConvEmpty(w *JsonEncoder, name string, value map[string]int64) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	uint32类型: MapUint32Map_<empty>
 *************************************/

func EncodeUint32Map(w *JsonEncoder, value map[string]uint32) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeUint32Map_OmitEmpty(w *JsonEncoder, name string, value map[string]uint32) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeUint32Map_WithEmpty(w *JsonEncoder, name string, value map[string]uint32) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeUint32Map_ConvEmpty(w *JsonEncoder, name string, value map[string]uint32) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	uin64类型: MapUint64Map_<empty>
 *************************************/

func EncodeUint64Map(w *JsonEncoder, value map[string]uint64) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeUint64Map_OmitEmpty(w *JsonEncoder, name string, value map[string]uint64) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeUint64Map_WithEmpty(w *JsonEncoder, name string, value map[string]uint64) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeUint64Map_ConvEmpty(w *JsonEncoder, name string, value map[string]uint64) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	float32类型: MapFloatMap_<empty>
 *************************************/

func EncodeFloatMap(w *JsonEncoder, value map[string]float32) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeFloatMap_OmitEmpty(w *JsonEncoder, name string, value map[string]float32) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeFloatMap_WithEmpty(w *JsonEncoder, name string, value map[string]float32) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeFloatMap_ConvEmpty(w *JsonEncoder, name string, value map[string]float32) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	float64类型: MapDoubleMap_<empty>
 *************************************/

func EncodeDoubleMap(w *JsonEncoder, value map[string]float64) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeDoubleMap_OmitEmpty(w *JsonEncoder, name string, value map[string]float64) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeDoubleMap_WithEmpty(w *JsonEncoder, name string, value map[string]float64) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeDoubleMap_ConvEmpty(w *JsonEncoder, name string, value map[string]float64) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	string类型: MapStringMap_<escap_html>_<empty>
 *************************************/

func EncodeStringMap(w *JsonEncoder, value map[string]string) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeStringMap_OmitEmpty(w *JsonEncoder, name string, value map[string]string) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeStringMap_WithEmpty(w *JsonEncoder, name string, value map[string]string) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeStringMap_ConvEmpty(w *JsonEncoder, name string, value map[string]string) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeStringHtmlMap(w *JsonEncoder, value map[string]string) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeStringHtml(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeStringHtmlMap_OmitEmpty(w *JsonEncoder, name string, value map[string]string) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeStringHtml(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeStringHtmlMap_WithEmpty(w *JsonEncoder, name string, value map[string]string) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeStringHtml(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeStringHtmlMap_ConvEmpty(w *JsonEncoder, name string, value map[string]string) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeStringHtml(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	bytes类型: MapBytesMap_<empty>
 *************************************/

func EncodeBytesMap(w *JsonEncoder, value map[string][]byte) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeBytesMap_OmitEmpty(w *JsonEncoder, name string, value map[string][]byte) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeBytesMap_WithEmpty(w *JsonEncoder, name string, value map[string][]byte) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeBytesMap_ConvEmpty(w *JsonEncoder, name string, value map[string][]byte) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	enum类型: MapEnumMap_<enum_as_int>_<empty>
 *************************************/

func EncodeEnumNameMap[Enum ~int32](w *JsonEncoder, value map[string]Enum, names map[int32]string) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeEnumNameMap_OmitEmpty[Enum ~int32](w *JsonEncoder, name string, value map[string]Enum, names map[int32]string) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeEnumNameMap_WithEmpty[Enum ~int32](w *JsonEncoder, name string, value map[string]Enum, names map[int32]string) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeEnumNameMap_ConvEmpty[Enum ~int32](w *JsonEncoder, name string, value map[string]Enum, names map[int32]string) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeEnumMap[Enum ~int32](w *JsonEncoder, value map[string]Enum) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeEnumMap_OmitEmpty[Enum ~int32](w *JsonEncoder, name string, value map[string]Enum) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeEnumMap_WithEmpty[Enum ~int32](w *JsonEncoder, name string, value map[string]Enum) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeEnumMap_ConvEmpty[Enum ~int32](w *JsonEncoder, name string, value map[string]Enum) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

/*************************************
	message类型: MapMessageMap_<empty>
 *************************************/

func EncodeMessageMap[Message any](w *JsonEncoder, value map[string]*Message) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace, rightBrace)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	}
}

func EncodeMessageMap_OmitEmpty[Message any](w *JsonEncoder, name string, value map[string]*Message) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeMessageMap_WithEmpty[Message any](w *JsonEncoder, name string, value map[string]*Message) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeMessageMap_ConvEmpty[Message any](w *JsonEncoder, name string, value map[string]*Message) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		for k, v := range value {
			EncodeString(w, k)
			w.buff = append(w.buff, colon)
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}
