package protoapi

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
)

/*
	需要支持default的编码:
*/

/*************************************
	bool类型: WriteBool_<empty>
 *************************************/

func EncodeBool(w *JsonEncoder, value bool) {
	if value {
		w.ensure(4)
		w.buff = append(w.buff, 't', 'r', 'u', 'e')
	} else {
		w.ensure(5)
		w.buff = append(w.buff, 'f', 'a', 'l', 's', 'e')
	}
}

func EncodeBool_OmitEmpty(w *JsonEncoder, name string, value bool) {
	if value {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 't', 'r', 'u', 'e', comma)
	}
}

func EncodeBool_WithEmpty(w *JsonEncoder, name string, value bool) {
	if value {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 't', 'r', 'u', 'e', comma)
	} else {
		w.ensure(9 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'f', 'a', 'l', 's', 'e', comma)
	}
}

func EncodeBool_ConvEmpty(w *JsonEncoder, name string, value bool) {
	EncodeBool_WithEmpty(w, name, value)
}

/*************************************
	int32类型: EncodeInt32_<empty>
 *************************************/

func EncodeInt32(w *JsonEncoder, value int32) {
	EncodeInt64(w, int64(value))
}

func EncodeInt32_OmitEmpty(w *JsonEncoder, name string, value int32) {
	EncodeInt64_OmitEmpty(w, name, int64(value))
}

func EncodeInt32_WithEmpty(w *JsonEncoder, name string, value int32) {
	EncodeInt64_WithEmpty(w, name, int64(value))
}

func EncodeInt32_ConvEmpty(w *JsonEncoder, name string, value int32) {
	EncodeInt64_ConvEmpty(w, name, int64(value))
}

/*************************************
	int64类型: EncodeInt64_<empty>
 *************************************/

func EncodeInt64(w *JsonEncoder, value int64) {
	if value != 0 {
		w.ensure(21)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], value, 10)...)
	} else {
		w.ensure(1)
		w.buff = append(w.buff, '0')
	}
}

func EncodeInt64_OmitEmpty(w *JsonEncoder, name string, value int64) {
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeInt64_WithEmpty(w *JsonEncoder, name string, value int64) {
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], value, 10)...)
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	}
}

func EncodeInt64_ConvEmpty(w *JsonEncoder, name string, value int64) {
	EncodeInt64_WithEmpty(w, name, value)
}

/*************************************
	uint32类型: EncodeUint32_<empty>
 *************************************/

func EncodeUint32(w *JsonEncoder, value uint32) {
	EncodeUint64(w, uint64(value))
}

func EncodeUint32_OmitEmpty(w *JsonEncoder, name string, value uint32) {
	EncodeUint64_OmitEmpty(w, name, uint64(value))
}

func EncodeUint32_WithEmpty(w *JsonEncoder, name string, value uint32) {
	EncodeUint64_WithEmpty(w, name, uint64(value))
}

func EncodeUint32_ConvEmpty(w *JsonEncoder, name string, value uint32) {
	EncodeUint64_WithEmpty(w, name, uint64(value))
}

/*************************************
	uin64类型: EncodeUint64_<empty>
 *************************************/

func EncodeUint64(w *JsonEncoder, value uint64) {
	if value != 0 {
		w.ensure(21)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], value, 10)...)
	} else {
		w.ensure(1)
		w.buff = append(w.buff, '0')
	}
}

func EncodeUint64_OmitEmpty(w *JsonEncoder, name string, value uint64) {
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeUint64_WithEmpty(w *JsonEncoder, name string, value uint64) {
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], value, 10)...)
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	}
}

func EncodeUint64_ConvEmpty(w *JsonEncoder, name string, value uint64) {
	EncodeUint64_WithEmpty(w, name, value)
}

/*************************************
	float32类型: EncodeFloat_<empty>
 *************************************/

func EncodeFloat(w *JsonEncoder, value float32) {
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		if w.firstError == nil {
			w.firstError = fmt.Errorf("unsupported value: %f", value)
		}
		return
	}
	if value != 0 {
		w.ensure(21)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], float64(value), 'g', -1, 32)...)
	} else {
		w.ensure(1)
		w.buff = append(w.buff, '0')
	}
}

func EncodeFloat_OmitEmpty(w *JsonEncoder, name string, value float32) {
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		if w.firstError == nil {
			w.firstError = fmt.Errorf("unsupported value: %f", value)
		}
		return
	}
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], float64(value), 'g', -1, 32)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeFloat_WithEmpty(w *JsonEncoder, name string, value float32) {
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		if w.firstError == nil {
			w.firstError = fmt.Errorf("unsupported value: %f", value)
		}
		return
	}
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], float64(value), 'g', -1, 32)...)
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	}
}

func EncodeFloat_ConvEmpty(w *JsonEncoder, name string, value float32) {
	EncodeFloat_WithEmpty(w, name, value)
}

/*************************************
	float64类型: EncodeDouble_<empty>
 *************************************/

func EncodeDouble(w *JsonEncoder, value float64) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		if w.firstError == nil {
			w.firstError = fmt.Errorf("unsupported value: %f", value)
		}
		return
	}
	w.ensure(21)
	if value != 0 {
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], value, 'g', -1, 64)...)
	} else {
		w.buff = append(w.buff, '0')
	}
}

func EncodeDouble_OmitEmpty(w *JsonEncoder, name string, value float64) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		if w.firstError == nil {
			w.firstError = fmt.Errorf("unsupported value: %f", value)
		}
		return
	}
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], value, 'g', -1, 64)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeDouble_WithEmpty(w *JsonEncoder, name string, value float64) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		if w.firstError == nil {
			w.firstError = fmt.Errorf("unsupported value: %f", value)
		}
		return
	}
	if value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], value, 'g', -1, 64)...)
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	}
}

func EncodeDouble_ConvEmpty(w *JsonEncoder, name string, value float64) {
	EncodeDouble_WithEmpty(w, name, value)
}

/*************************************
	string类型: EncodeString_<escape_html>_<empty>
 *************************************/

func EncodeString(w *JsonEncoder, value string) {
	if value != "" {
		w.ensure(2 + len(value))
		w.buff = append(w.buff, quotes)
		w.escape(value, &noEscapeHtmlTable)
		w.buff = append(w.buff, quotes)
	} else {
		w.ensure(2)
		w.buff = append(w.buff, quotes, quotes)
	}
}

func EncodeString_OmitEmpty(w *JsonEncoder, name string, value string) {
	if value != "" {
		w.ensure(4 + len(name) + len(value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(value, &noEscapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeString_WithEmpty(w *JsonEncoder, name string, value string) {
	if value != "" {
		w.ensure(6 + len(name) + len(value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(value, &noEscapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	} else {
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	}
}

func EncodeString_ConvEmpty(w *JsonEncoder, name string, value string) {
	EncodeString_WithEmpty(w, name, value)
}

func EncodeString_EscapeHtml(w *JsonEncoder, value string) {
	if value != "" {
		w.ensure(2 + len(value))
		w.buff = append(w.buff, quotes)
		w.escape(value, &escapeHtmlTable)
		w.buff = append(w.buff, quotes)
	} else {
		w.ensure(2)
		w.buff = append(w.buff, quotes, quotes)
	}
}

func EncodeString_EscapeHtml_OmitEmpty(w *JsonEncoder, name string, value string) {
	if value != "" {
		w.ensure(6 + len(name) + len(value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(value, &escapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeString_EscapeHtml_WithEmpty(w *JsonEncoder, name string, value string) {
	if value != "" {
		w.ensure(6 + len(name) + len(value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(value, &escapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	} else {
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	}
}

func EncodeString_EscapeHtml_ConvEmpty(w *JsonEncoder, name string, value string) {
	EncodeString_EscapeHtml_WithEmpty(w, name, value)
}

/*************************************
	bytes类型: WriteBytes_<empty>
 *************************************/

func EncodeBytes(w *JsonEncoder, value []byte) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case len(value) == 0:
		w.ensure(2)
		w.buff = append(w.buff, quotes, quotes)
	default:
		w.ensure(2 + len(value))
		w.buff = append(w.buff, quotes)
		w.base64(value)
		w.buff = append(w.buff, quotes)
	}
}

func EncodeBytes_OmitEmpty(w *JsonEncoder, name string, value []byte) {
	if len(value) != 0 {
		w.ensure(6 + len(name) + len(value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.base64(value)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeBytes_WithEmpty(w *JsonEncoder, name string, value []byte) {

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
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	default:
		w.ensure(6 + len(name) + base64.StdEncoding.EncodedLen(len(value)))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.base64(value)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeBytes_ConvEmpty(w *JsonEncoder, name string, value []byte) {

	switch {
	case value == nil || len(value) == 0:
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	default:
		w.ensure(6 + len(name) + base64.StdEncoding.EncodedLen(len(value)))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.base64(value)
		w.buff = append(w.buff, quotes, comma)
	}
}

/*************************************
	enum类型: EncodeEnum_<enum_as_int>_<empty>
 *************************************/

func EncodeEnumName[Enum ~int32](w *JsonEncoder, value Enum, names map[int32]string) {
	w.ensure(10)
	w.buff = append(w.buff, quotes)
	w.buff = append(w.buff, names[int32(value)]...)
	w.buff = append(w.buff, quotes)
}

func EncodeEnumName_OmitEmpty[Enum ~int32](w *JsonEncoder, name string, value Enum, names map[int32]string) {
	// Enum比较特殊, 0对应的name非空
	EncodeEnumName_WithEmpty(w, name, value, names)
}

func EncodeEnumName_WithEmpty[Enum ~int32](w *JsonEncoder, name string, value Enum, names map[int32]string) {
	w.ensure(10 + len(name))
	w.buff = append(w.buff, quotes)
	w.buff = append(w.buff, name...)
	w.buff = append(w.buff, quotes, colon, quotes)
	w.buff = append(w.buff, names[int32(value)]...)
	w.buff = append(w.buff, quotes, comma)
}

func EncodeEnumName_ConvEmpty[Enum ~int32](w *JsonEncoder, name string, value Enum, names map[int32]string) {
	EncodeEnumName_WithEmpty(w, name, value, names)
}

func EncodeEnum[Enum ~int32](w *JsonEncoder, value Enum) {
	EncodeInt32(w, int32(value))
}

func EncodeEnum_OmitEmpty[Enum ~int32](w *JsonEncoder, name string, value Enum) {
	EncodeInt32_OmitEmpty(w, name, int32(value))
}

func EncodeEnum_WithEmpty[Enum ~int32](w *JsonEncoder, name string, value Enum) {
	EncodeInt32_WithEmpty(w, name, int32(value))
}

func EncodeEnum_ConvEmpty[Enum ~int32](w *JsonEncoder, name string, value Enum) {
	EncodeInt32_ConvEmpty(w, name, int32(value))
}

/*************************************
	message类型: WriteMessage_<empty>
 *************************************/

func EncodeMessage[Message any](w *JsonEncoder, value *Message, f func(w *JsonEncoder, m *Message)) {
	if value != nil {
		w.ensure(2)
		w.buff = append(w.buff, leftBrace)
		f(w, value)
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
		} else {
			w.buff = append(w.buff, rightBrace)
		}
	} else {
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	}
}

func EncodeMessage_OmitEmpty[Message any](w *JsonEncoder, name string, value *Message, f func(w *JsonEncoder, m *Message)) {
	if value != nil {
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		f(w, value)
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeMessage_WithEmpty[Message any](w *JsonEncoder, name string, value *Message, f func(w *JsonEncoder, m *Message)) {
	if value != nil {
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		f(w, value)
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	} else {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	}
}

func EncodeMessage_ConvEmpty[Message any](w *JsonEncoder, name string, value *Message, f func(w *JsonEncoder, m *Message)) {
	if value != nil {
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		f(w, value)
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	} else {
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace, rightBrace, comma)
	}
}
