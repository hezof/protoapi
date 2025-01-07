package kits

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	LayoutDatetime       = "2006-01-02 15:04:05"
	LayoutDatetimeT      = "2006-01-02T15:04:05"
	LayoutDateTimeLength = len(LayoutDatetime)
)

func Bool(v interface{}) (bool, bool) {

	if v == nil {
		return false, true
	}
	if vl, ok := v.(bool); ok {
		return vl, true
	}

	switch v := v.(type) {
	case string:
		return v == "true", true
	case int:
		return v != 0, true
	case int8:
		return v != 0, true
	case int16:
		return v != 0, true
	case int32:
		return v != 0, true
	case int64:
		return v != 0, true
	case uint:
		return v != 0, true
	case uint8:
		return v != 0, true
	case uint16:
		return v != 0, true
	case uint32:
		return v != 0, true
	case uint64:
		return v != 0, true
	case float32:
		return v != 0, true
	case float64:
		return v != 0, true
	}
	return false, false
}

func String(v interface{}) (string, bool) {

	if v == nil {
		return "", true
	}
	if vl, ok := v.(string); ok {
		return vl, true
	}
	switch v := v.(type) {
	case []byte:
		return string(v), true
	case int:
		return strconv.FormatInt(int64(v), 10), true
	case uint:
		return strconv.FormatUint(uint64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case uint64:
		return strconv.FormatUint(uint64(v), 10), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case int16:
		return strconv.FormatInt(int64(v), 10), true
	case uint16:
		return strconv.FormatUint(uint64(v), 10), true
	case int8:
		return strconv.FormatInt(int64(v), 10), true
	case uint8:
		return strconv.FormatUint(uint64(v), 10), true
	case bool:
		if v {
			return "true", true
		} else {
			return "false", true
		}
	case time.Time:
		return v.Format(LayoutDatetime), true
	case time.Duration:
		return strconv.FormatInt(int64(v), 10), true
	case fmt.Stringer:
		return v.String(), true
	case fmt.Formatter:
		return fmt.Sprint(v), true
	}

	return "", false
}

func ToString(v interface{}) string {
	if rt, ok := String(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert string error: %#v", v))
}

func Int(v interface{}) (int, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(int); ok {
		return vl, true
	}

	switch v := v.(type) {
	case uint:
		return int(v), true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case uint64:
		return int(v), true
	case int32:
		return int(v), true
	case float32:
		return int(v), true
	case uint32:
		return int(v), true
	case int16:
		return int(v), true
	case uint16:
		return int(v), true
	case int8:
		return int(v), true
	case uint8:
		return int(v), true
	case string:
		if rt, err := strconv.ParseInt(v, 10, 0); err == nil {
			return int(rt), true
		} else {
			return 0, true
		}
	case time.Duration:
		return int(v), true
	}
	return 0, false
}

func ToInt(v interface{}) int {
	if rt, ok := Int(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert int error: %#v", v))
}

func Int32(v interface{}) (int32, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(int32); ok {
		return vl, true
	}

	switch v := v.(type) {
	case int:
		return int32(v), true
	case uint:
		return int32(v), true
	case int64:
		return int32(v), true
	case float64:
		return int32(v), true
	case uint64:
		return int32(v), true
	case float32:
		return int32(v), true
	case uint32:
		return int32(v), true
	case int16:
		return int32(v), true
	case uint16:
		return int32(v), true
	case int8:
		return int32(v), true
	case uint8:
		return int32(v), true
	case string:
		if rt, err := strconv.ParseInt(v, 10, 32); err == nil {
			return int32(rt), true
		}
	case time.Duration:
		return int32(v), true
	}
	return 0, false
}

func ToInt32(v interface{}) int32 {
	if rt, ok := Int32(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert int32 error: %#v", v))
}

func Int64(v interface{}) (int64, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(int64); ok {
		return vl, true
	}

	switch v := v.(type) {
	case int:
		return int64(v), true
	case uint:
		return int64(v), true
	case float64:
		return int64(v), true
	case uint64:
		return int64(v), true
	case int32:
		return int64(v), true
	case float32:
		return int64(v), true
	case uint32:
		return int64(v), true
	case int16:
		return int64(v), true
	case uint16:
		return int64(v), true
	case int8:
		return int64(v), true
	case uint8:
		return int64(v), true
	case string:
		if rt, err := strconv.ParseInt(v, 10, 64); err == nil {
			return int64(rt), true
		}
	case time.Duration:
		return int64(v), true
	}
	return 0, false
}

func ToInt64(v interface{}) int64 {
	if rt, ok := Int64(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert int64 error: %#v", v))
}

func Uint(v interface{}) (uint, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(uint); ok {
		return vl, true
	}

	switch v := v.(type) {
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	case float64:
		return uint(v), true
	case uint64:
		return uint(v), true
	case int32:
		return uint(v), true
	case float32:
		return uint(v), true
	case uint32:
		return uint(v), true
	case int16:
		return uint(v), true
	case uint16:
		return uint(v), true
	case int8:
		return uint(v), true
	case uint8:
		return uint(v), true
	case string:
		if rt, err := strconv.ParseUint(v, 10, 0); err == nil {
			return uint(rt), true
		}
	case time.Duration:
		return uint(v), true
	}
	return 0, false
}

func ToUint(v interface{}) uint {
	if rt, ok := Uint(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert uint error: %#v", v))
}

func Uint32(v interface{}) (uint32, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(uint32); ok {
		return vl, true
	}

	switch v := v.(type) {
	case int:
		return uint32(v), true
	case uint:
		return uint32(v), true
	case int64:
		return uint32(v), true
	case float64:
		return uint32(v), true
	case uint64:
		return uint32(v), true
	case int32:
		return uint32(v), true
	case float32:
		return uint32(v), true
	case int16:
		return uint32(v), true
	case uint16:
		return uint32(v), true
	case int8:
		return uint32(v), true
	case uint8:
		return uint32(v), true
	case string:
		if rt, err := strconv.ParseUint(v, 10, 32); err == nil {
			return uint32(rt), true
		}
	case time.Duration:
		return uint32(v), true
	}
	return 0, false
}

func ToUint32(v interface{}) uint32 {
	if rt, ok := Uint32(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert uint32 error: %#v", v))
}

func Uint64(v interface{}) (uint64, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(uint64); ok {
		return vl, true
	}

	switch v := v.(type) {
	case int:
		return uint64(v), true
	case uint:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case float64:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case float32:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case int16:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case int8:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case string:
		if rt, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint64(rt), true
		}
	case time.Duration:
		return uint64(v), true
	}
	return 0, false
}

func ToUint64(v interface{}) uint64 {
	if rt, ok := Uint64(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert uint64 error: %#v", v))
}

func Float32(v interface{}) (float32, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(float32); ok {
		return vl, true
	}

	switch v := v.(type) {
	case int:
		return float32(v), true
	case uint:
		return float32(v), true
	case int64:
		return float32(v), true
	case float64:
		return float32(v), true
	case uint64:
		return float32(v), true
	case int32:
		return float32(v), true
	case uint32:
		return float32(v), true
	case int16:
		return float32(v), true
	case uint16:
		return float32(v), true
	case int8:
		return float32(v), true
	case uint8:
		return float32(v), true
	case string:
		if rt, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(rt), true
		}
	case time.Duration:
		return float32(v), true
	}
	return 0, false
}

func ToFloat32(v interface{}) float32 {
	if rt, ok := Float32(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert float32 error: %#v", v))
}

func Float64(v interface{}) (float64, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(float64); ok {
		return vl, true
	}

	switch v := v.(type) {
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case uint:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint32:
		return float64(v), true
	case int16:
		return float64(v), true
	case uint16:
		return float64(v), true
	case int8:
		return float64(v), true
	case uint8:
		return float64(v), true
	case string:
		if rt, err := strconv.ParseFloat(v, 64); err == nil {
			return rt, true
		}
	case time.Duration:
		return float64(v), true
	}
	return 0, false
}

func ToFloat64(v interface{}) float64 {
	if rt, ok := Float64(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert float64 error: %#v", v))
}

var ZeroTime = time.Time{}

func Time(v interface{}) (time.Time, bool) {

	if v == nil {
		return ZeroTime, true
	}
	if vl, ok := v.(time.Time); ok {
		return vl, true
	}
	if vl, ok := v.(*time.Time); ok {
		return *vl, true
	}

	switch v := v.(type) {
	case int64:
		return time.Unix(v, 0), true
	case uint64:
		return time.Unix(int64(v), 0), true
	case float64:
		return time.Unix(int64(v), 0), true
	case string:
		// 适配date/time串的多种格式
		var layout string
		if strings.IndexByte(v, 'T') != -1 {
			layout = LayoutDatetimeT
		} else {
			layout = LayoutDatetime
		}
		if length := len(v); length == LayoutDateTimeLength {
			if rt, err := time.ParseInLocation(layout, v, time.Local); err == nil {
				return rt, true
			}
		} else if length > LayoutDateTimeLength {
			if rt, err := time.ParseInLocation(layout, v[0:LayoutDateTimeLength], time.Local); err == nil {
				return rt, true
			}
		} else {
			if rt, err := time.ParseInLocation(layout[0:length], v, time.Local); err == nil {
				return rt, true
			}
		}
	}
	return ZeroTime, false
}

func ToTime(v interface{}) time.Time {
	if rt, ok := Time(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert time error: %#v", v))
}

func Duration(v interface{}) (time.Duration, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(time.Duration); ok {
		return vl, true
	}
	if vl, ok := v.(int64); ok {
		return time.Duration(vl), true
	}

	switch v := v.(type) {
	case float64:
		return time.Duration(v), true
	case uint64:
		return time.Duration(v), true
	case int32:
		return time.Duration(v), true
	case float32:
		return time.Duration(v), true
	case uint32:
		return time.Duration(v), true
	case int:
		return time.Duration(v), true
	case uint:
		return time.Duration(v), true
	case int16:
		return time.Duration(v), true
	case uint16:
		return time.Duration(v), true
	case int8:
		return time.Duration(v), true
	case uint8:
		return time.Duration(v), true
	case string:
		if rt, err := time.ParseDuration(v); err == nil {
			return rt, true
		}

	}
	return 0, false
}

func ToDuration(v interface{}) time.Duration {
	if rt, ok := Duration(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert duration error: %#v", v))
}

func Complex128(v interface{}) (complex128, bool) {

	if v == nil {
		return 0, true
	}
	if vl, ok := v.(complex128); ok {
		return vl, true
	}
	if vl, ok := v.(complex64); ok {
		return complex128(vl), true
	}
	return 0, false
}

func Bytes(v interface{}) ([]byte, bool) {

	if v == nil {
		return nil, true
	}
	if vl, ok := v.([]byte); ok {
		return vl, true
	}
	if vl, ok := v.(string); ok {
		return []byte(vl), true
	}
	return nil, false
}

func Uintptr(v interface{}) (uintptr, bool) {
	if v == nil {
		return 0, true
	}
	if vl, ok := v.(uintptr); ok {
		return vl, true
	}
	if vl, ok := v.(unsafe.Pointer); ok {
		return uintptr(vl), true
	}
	return 0, false
}

func UnsafePointer(v interface{}) (unsafe.Pointer, bool) {
	if v == nil {
		return nil, true
	}
	if vl, ok := v.(unsafe.Pointer); ok {
		return vl, true
	}
	if vl, ok := v.(uintptr); ok {
		return unsafe.Pointer(vl), true
	}
	return nil, false
}
