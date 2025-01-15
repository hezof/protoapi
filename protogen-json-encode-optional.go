package protoapi

import (
	"strconv"
)

/*
	需要支持optional的编码:
*/

/*************************************
	bool类型: EncodeBytes_<empty>
 *************************************/

func EncodeBoolOptional(w *JsonEncoder, value *bool) {
	w.ensure(5)
	switch {
	case value == nil:
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == false:
		w.buff = append(w.buff, 'f', 'a', 'l', 's', 'e')
	default:
		w.buff = append(w.buff, 't', 'r', 'u', 'e')
	}
}

func EncodeBoolOptional_OmitEmpty(w *JsonEncoder, name string, value *bool) {
	if value != nil && *value != false {
		w.ensure(9 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 't', 'r', 'u', 'e', comma)
	}
}

func EncodeBoolOptional_WithEmpty(w *JsonEncoder, name string, value *bool) {
	w.ensure(9 + len(name))
	switch {
	case value == nil:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == false:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'f', 'a', 'l', 's', 'e', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 't', 'r', 'u', 'e', comma)
	}
}

func EncodeBoolOptional_ConvEmpty(w *JsonEncoder, name string, value *bool) {
	w.ensure(9 + len(name))
	switch {
	case value == nil || *value == false:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'f', 'a', 'l', 's', 'e', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 't', 'r', 'u', 'e', comma)
	}
}

/*************************************
	int32类型: WriteInt32Optional_<empty>
 *************************************/

func EncodeInt32Optional(w *JsonEncoder, value *int32) {
	w.ensure(11)
	switch {
	case value == nil:
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == 0:
		w.buff = append(w.buff, '0')
	default:
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], int64(*value), 10)...)
	}
}

func EncodeInt32Optional_OmitEmpty(w *JsonEncoder, name string, value *int32) {
	w.ensure(15 + len(name))
	if value != nil && *value != 0 {
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], int64(*value), 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeInt32Optional_WithEmpty(w *JsonEncoder, name string, value *int32) {
	w.ensure(15 + len(name))
	switch {
	case value == nil:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], int64(*value), 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeInt32Optional_ConvEmpty(w *JsonEncoder, name string, value *int32) {
	w.ensure(15 + len(name))
	switch {
	case value == nil || *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], int64(*value), 10)...)
		w.buff = append(w.buff, comma)
	}
}

/*************************************
	int64类型: WriteInt64Optional_<empty>
 *************************************/

func EncodeInt64Optional(w *JsonEncoder, value *int64) {
	w.ensure(21)
	switch {
	case value == nil:
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == 0:
		w.buff = append(w.buff, '0')
	default:
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], *value, 10)...)
	}
}

func EncodeInt64Optional_OmitEmpty(w *JsonEncoder, name string, value *int64) {
	w.ensure(25 + len(name))
	if value != nil && *value != 0 {
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], *value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeInt64Optional_WithEmpty(w *JsonEncoder, name string, value *int64) {
	w.ensure(25 + len(name))
	switch {
	case value == nil:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], *value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeInt64Optional_ConvEmpty(w *JsonEncoder, name string, value *int64) {
	w.ensure(25 + len(name))
	switch {
	case value == nil || *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendInt(w.number[0:0], *value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

/*************************************
	uint32类型: WriteUint32Optional_<empty>
 *************************************/

func EncodeUint32Optional(w *JsonEncoder, value *uint32) {
	w.ensure(21)
	switch {
	case value == nil:
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == 0:
		w.buff = append(w.buff, '0')
	default:
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], uint64(*value), 10)...)
	}
}

func EncodeUint32Optional_OmitEmpty(w *JsonEncoder, name string, value *uint32) {
	w.ensure(25 + len(name))
	if value != nil && *value != 0 {
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], uint64(*value), 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeUint32Optional_WithEmpty(w *JsonEncoder, name string, value *uint32) {
	w.ensure(25 + len(name))
	switch {
	case value == nil:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], uint64(*value), 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeUint32Optional_ConvEmpty(w *JsonEncoder, name string, value *uint32) {
	w.ensure(25 + len(name))
	switch {
	case value == nil || *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], uint64(*value), 10)...)
		w.buff = append(w.buff, comma)
	}
}

/*************************************
	uin64类型: WriteUint64Optional_<empty>
 *************************************/

func EncodeUint64Optional(w *JsonEncoder, value *uint64) {
	w.ensure(21)
	switch {
	case value == nil:
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == 0:
		w.buff = append(w.buff, '0')
	default:
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], *value, 10)...)
	}
}

func EncodeUint64Optional_OmitEmpty(w *JsonEncoder, name string, value *uint64) {
	if value != nil && *value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], *value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeUint64Optional_WithEmpty(w *JsonEncoder, name string, value *uint64) {
	w.ensure(25 + len(name))
	switch {
	case value == nil:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], *value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeUint64Optional_ConvEmpty(w *JsonEncoder, name string, value *uint64) {
	w.ensure(25 + len(name))
	switch {
	case value == nil || *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendUint(w.number[0:0], *value, 10)...)
		w.buff = append(w.buff, comma)
	}
}

/*************************************
	float32类型: WriteFloatOptional_<empty>
 *************************************/

func EncodeFloatOptional(w *JsonEncoder, value *float32) {
	w.ensure(21)
	switch {
	case value == nil:
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == 0:
		w.buff = append(w.buff, '0')
	default:
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], float64(*value), 'g', -1, 32)...)
	}
}

func EncodeFloatOptional_OmitEmpty(w *JsonEncoder, name string, value *float32) {
	if value != nil && *value != 0 {
		w.ensure(25 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], float64(*value), 'g', -1, 32)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeFloatOptional_WithEmpty(w *JsonEncoder, name string, value *float32) {
	w.ensure(25 + len(name))
	switch {
	case value == nil:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], float64(*value), 'g', -1, 32)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeFloatOptional_ConvEmpty(w *JsonEncoder, name string, value *float32) {
	w.ensure(25 + len(name))
	switch {
	case value == nil || *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], float64(*value), 'g', -1, 32)...)
		w.buff = append(w.buff, comma)
	}
}

/*************************************
	float64类型: WriteDoubleOptional_<empty>
 *************************************/

func EncodeDoubleOptional(w *JsonEncoder, value *float64) {
	w.ensure(21)
	switch {
	case value == nil:
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == 0:
		w.buff = append(w.buff, '0')
	default:
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], *value, 'g', -1, 64)...)
	}
}

func EncodeDoubleOptional_OmitEmpty(w *JsonEncoder, name string, value *float64) {
	w.ensure(25 + len(name))
	if value != nil && *value != 0 {
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], *value, 'g', -1, 64)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeDoubleOptional_WithEmpty(w *JsonEncoder, name string, value *float64) {
	w.ensure(25 + len(name))
	switch {
	case value == nil:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], *value, 'g', -1, 64)...)
		w.buff = append(w.buff, comma)
	}
}

func EncodeDoubleOptional_ConvEmpty(w *JsonEncoder, name string, value *float64) {
	w.ensure(25 + len(name))
	switch {
	case value == nil || *value == 0:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	default:
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		w.buff = append(w.buff, strconv.AppendFloat(w.number[0:0], *value, 'g', -1, 64)...)
		w.buff = append(w.buff, comma)
	}
}

/*************************************
	string类型: WriteStringOptional_<escape_html>_<empty>
 *************************************/

func EncodeStringOptional(w *JsonEncoder, value *string) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == "":
		w.ensure(2)
		w.buff = append(w.buff, quotes, quotes)
	default:
		w.ensure(2 + len(*value))
		w.buff = append(w.buff, quotes)
		w.escape(*value, &noEscapeHtmlTable)
		w.buff = append(w.buff, quotes)
	}
}

func EncodeStringOptional_OmitEmpty(w *JsonEncoder, name string, value *string) {
	if value != nil && *value != "" {
		w.ensure(6 + len(name) + len(*value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(*value, &noEscapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeStringOptional_WithEmpty(w *JsonEncoder, name string, value *string) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == "":
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	default:
		w.ensure(6 + len(name) + len(*value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(*value, &noEscapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeStringOptional_ConvEmpty(w *JsonEncoder, name string, value *string) {
	switch {
	case value == nil || *value == "":
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	default:
		w.ensure(6 + len(name) + len(*value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(*value, &noEscapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeStringOptional_EscapeHtml(w *JsonEncoder, value *string) {
	switch {
	case value == nil:
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	case *value == "":
		w.ensure(2)
		w.buff = append(w.buff, quotes, quotes)
	default:
		w.ensure(2 + len(*value))
		w.buff = append(w.buff, quotes)
		w.escape(*value, &escapeHtmlTable)
		w.buff = append(w.buff, quotes)
	}
}

func EncodeStringOptional_EscapeHtml_OmitEmpty(w *JsonEncoder, name string, value *string) {
	if value != nil && *value != "" {
		w.ensure(6 + len(name) + len(*value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(*value, &escapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeStringOptional_EscapeHtml_WithEmpty(w *JsonEncoder, name string, value *string) {
	switch {
	case value == nil:
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	case *value == "":
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	default:
		w.ensure(6 + len(name) + len(*value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(*value, &escapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

func EncodeStringOptional_EscapeHtml_ConvEmpty(w *JsonEncoder, name string, value *string) {
	switch {
	case value == nil || *value == "":
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	default:
		w.ensure(6 + len(name) + len(*value))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes)
		w.escape(*value, &escapeHtmlTable)
		w.buff = append(w.buff, quotes, comma)
	}
}

/*************************************
	bytes类型: optional类型相同
 *************************************/

func EncodeBytesOptional(w *JsonEncoder, value []byte) {
	EncodeBytes(w, value)
}

func EncodeBytesOptional_OmitEmpty(w *JsonEncoder, name string, value []byte) {
	EncodeBytes_OmitEmpty(w, name, value)
}

func EncodeBytesOptional_WithEmpty(w *JsonEncoder, name string, value []byte) {
	EncodeBytes_WithEmpty(w, name, value)
}

func EncodeBytesOptional_ConvEmpty(w *JsonEncoder, name string, value []byte) {
	EncodeBytes_ConvEmpty(w, name, value)
}

/*************************************
	enum类型: OptionalEnum_<enum_as_int>_<empty>
 *************************************/

func EncodeEnumNameOptional[Enum ~int32](w *JsonEncoder, value *Enum, names map[int32]string) {
	if value != nil {
		EncodeEnumName(w, *value, names)
	} else {
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	}
}

func EncodeEnumNameOptional_OmitEmpty[Enum ~int32](w *JsonEncoder, name string, value *Enum, names map[int32]string) {
	if value != nil {
		EncodeEnumName_OmitEmpty(w, name, *value, names)
	}
}

func EncodeEnumNameOptional_WithEmpty[Enum ~int32](w *JsonEncoder, name string, value *Enum, names map[int32]string) {
	if value != nil {
		EncodeEnumName_WithEmpty(w, name, *value, names)
	} else {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	}
}

func EncodeEnumNameOptional_ConvEmpty[Enum ~int32](w *JsonEncoder, name string, value *Enum, names map[int32]string) {
	if value != nil {
		EncodeEnumName_ConvEmpty(w, name, *value, names)
	} else {
		w.ensure(6 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, quotes, quotes, comma)
	}
}

func EncodeEnumOptional[Enum ~int32](w *JsonEncoder, value *Enum) {
	if value != nil {
		EncodeEnum(w, *value)
	} else {
		w.ensure(4)
		w.buff = append(w.buff, 'n', 'u', 'l', 'l')
	}
}

func EncodeEnumOptional_OmitEmpty[Enum ~int32](w *JsonEncoder, name string, value *Enum) {
	if value != nil {
		EncodeEnum_OmitEmpty(w, name, *value)
	}
}

func EncodeEnumOptional_WithEmpty[Enum ~int32](w *JsonEncoder, name string, value *Enum) {
	if value != nil {
		EncodeEnum_WithEmpty(w, name, *value)
	} else {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	}
}

func EncodeEnumOptional_ConvEmpty[Enum ~int32](w *JsonEncoder, name string, value *Enum) {
	if value != nil {
		EncodeEnum_ConvEmpty(w, name, *value)
	} else {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, '0', comma)
	}
}

/*************************************
	message类型
 *************************************/

func EncodeMessageOptional[Message any](w *JsonEncoder, name string, value *Message, f func(w *JsonEncoder, m *Message)) {
	EncodeMessage(w, value, f)
}

func EncodeMessageOptional_OmitEmpty[Message any](w *JsonEncoder, name string, value *Message, f func(w *JsonEncoder, m *Message)) {
	EncodeMessage_OmitEmpty(w, name, value, f)
}

func EncodeMessageOptional_WithEmpty[Message any](w *JsonEncoder, name string, value *Message, f func(w *JsonEncoder, m *Message)) {
	EncodeMessage_WithEmpty(w, name, value, f)
}

func EncodeMessageOptional_ConvEmpty[Message any](w *JsonEncoder, name string, value *Message, f func(w *JsonEncoder, m *Message)) {
	EncodeMessage_ConvEmpty(w, name, value, f)
}
