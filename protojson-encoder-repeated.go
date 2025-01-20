package protoapi

/*************************************
	bool类型: EncodeBoolRepeated_<empty>
 *************************************/

func EncodeBoolRepeated(w *JsonEncoder, value []bool) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeBoolRepeated_OmitEmpty(w *JsonEncoder, name string, value []bool) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeBoolRepeated_WithEmpty(w *JsonEncoder, name string, value []bool) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeBoolRepeated_ConvEmpty(w *JsonEncoder, name string, value []bool) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeBool(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	int32类型: EncodeInt32Repeated_<empty>
 *************************************/

func EncodeInt32Repeated(w *JsonEncoder, value []int32) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeInt32Repeated_OmitEmpty(w *JsonEncoder, name string, value []int32) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeInt32Repeated_WithEmpty(w *JsonEncoder, name string, value []int32) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeInt32Repeated_ConvEmpty(w *JsonEncoder, name string, value []int32) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeInt32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	int64类型: EncodeInt64Repeated_<empty>
 *************************************/

func EncodeInt64Repeated(w *JsonEncoder, value []int64) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeInt64Repeated_OmitEmpty(w *JsonEncoder, name string, value []int64) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeInt64Repeated_WithEmpty(w *JsonEncoder, name string, value []int64) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeInt64Repeated_ConvEmpty(w *JsonEncoder, name string, value []int64) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeInt64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	uint32类型: EncodeUint32Repeated_<empty>
 *************************************/

func EncodeUint32Repeated(w *JsonEncoder, value []uint32) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeUint32Repeated_OmitEmpty(w *JsonEncoder, name string, value []uint32) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeUint32Repeated_WithEmpty(w *JsonEncoder, name string, value []uint32) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeUint32Repeated_ConvEmpty(w *JsonEncoder, name string, value []uint32) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeUint32(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	uin64类型: EncodeUint64Repeated_<empty>
 *************************************/

func EncodeUint64Repeated(w *JsonEncoder, value []uint64) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeUint64Repeated_OmitEmpty(w *JsonEncoder, name string, value []uint64) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeUint64Repeated_WithEmpty(w *JsonEncoder, name string, value []uint64) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeUint64Repeated_ConvEmpty(w *JsonEncoder, name string, value []uint64) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeUint64(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	float32类型: EncodeFloatRepeated_<empty>
 *************************************/

func EncodeFloatRepeated(w *JsonEncoder, value []float32) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeFloatRepeated_OmitEmpty(w *JsonEncoder, name string, value []float32) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeFloatRepeated_WithEmpty(w *JsonEncoder, name string, value []float32) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeFloatRepeated_ConvEmpty(w *JsonEncoder, name string, value []float32) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeFloat(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	float64类型: EncodeDoubleRepeated_<empty>
 *************************************/

func EncodeDoubleRepeated(w *JsonEncoder, value []float64) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeDoubleRepeated_OmitEmpty(w *JsonEncoder, name string, value []float64) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeDoubleRepeated_WithEmpty(w *JsonEncoder, name string, value []float64) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeDoubleRepeated_ConvEmpty(w *JsonEncoder, name string, value []float64) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, comma)
		for _, v := range value {
			EncodeDouble(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	string类型: EncodeStringRepeated_<escape_html>_<empty>
 *************************************/

func EncodeStringRepeated(w *JsonEncoder, value []string) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeStringRepeated_OmitEmpty(w *JsonEncoder, name string, value []string) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeStringRepeated_WithEmpty(w *JsonEncoder, name string, value []string) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeStringRepeated_ConvEmpty(w *JsonEncoder, name string, value []string) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, comma)
		for _, v := range value {
			EncodeString(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	bytes类型: EncodeBytesRepeated_<empty>
 *************************************/

func EncodeBytesRepeated(w *JsonEncoder, value [][]byte) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeBytesRepeated_OmitEmpty(w *JsonEncoder, name string, value [][]byte) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeBytesRepeated_WithEmpty(w *JsonEncoder, name string, value [][]byte) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeBytesRepeated_ConvEmpty(w *JsonEncoder, name string, value [][]byte) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeBytes(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	enum类型: OptionalEnum_<enum_as_int>_<empty>
 *************************************/

func EncodeEnumNameRepeated[E ~int32](w *JsonEncoder, value []E, names map[int32]string) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeEnumNameRepeated_OmitEmpty[E ~int32](w *JsonEncoder, name string, value []E, names map[int32]string) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeEnumNameRepeated_WithEmpty[E ~int32](w *JsonEncoder, name string, value []E, names map[int32]string) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeEnumNameRepeated_ConvEmpty[E ~int32](w *JsonEncoder, name string, value []E, names map[int32]string) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeEnumName(w, v, names)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeEnumRepeated[Enum ~int32](w *JsonEncoder, value []Enum) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeEnumRepeated_OmitEmpty[Enum ~int32](w *JsonEncoder, name string, value []Enum) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeEnumRepeated_WithEmpty[Enum ~int32](w *JsonEncoder, name string, value []Enum) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeEnumRepeated_ConvEmpty[Enum ~int32](w *JsonEncoder, name string, value []Enum) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeEnum(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

/*************************************
	message类型: EncodeMessageRepeated_<empty>
 *************************************/

func EncodeMessageRepeated[Message any](w *JsonEncoder, value []*Message) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket, rightBracket)
	default:
		w.ensure(2)
		w.buff = append(w.buff, leftBracket)
		for _, v := range value {
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
		} else {
			w.buff = append(w.buff, rightBracket)
		}
	}
}

func EncodeMessageRepeated_OmitEmpty[Message any](w *JsonEncoder, name string, value []*Message) {
	if len(value) != 0 {
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeMessageRepeated_WithEmpty[Message any](w *JsonEncoder, name string, value []*Message) {
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
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}

func EncodeMessageRepeated_ConvEmpty[Message any](w *JsonEncoder, name string, value []*Message) {
	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket, rightBracket, comma)
	default:
		w.ensure(4 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBracket)
		for _, v := range value {
			EncodeMessage(w, v)
			w.buff = append(w.buff, comma)
		}
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBracket
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBracket, comma)
		}
	}
}
