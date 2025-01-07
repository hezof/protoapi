package kits

import (
	"fmt"
	"reflect"
	"time"
)

func Slice(v interface{}) ([]interface{}, bool) {
	switch v := v.(type) {
	case nil:
		return nil, true
	case []bool:
		return SliceBool(v), true
	case []int:
		return SliceInt(v), true
	case []int8:
		return SliceInt8(v), true
	case []int16:
		return SliceInt16(v), true
	case []int32:
		return SliceInt32(v), true
	case []int64:
		return SliceInt64(v), true
	case []uint:
		return SliceUint(v), true
	case []uint8:
		return SliceUint8(v), true
	case []uint16:
		return SliceUint16(v), true
	case []uint32:
		return SliceUint32(v), true
	case []uint64:
		return SliceUint64(v), true
	case []float32:
		return SliceFloat32(v), true
	case []float64:
		return SliceFloat64(v), true
	case []string:
		return SliceString(v), true
	case [][]byte:
		return SliceBytes(v), true
	case []time.Time:
		return SliceTime(v), true
	case []time.Duration:
		return SliceDuration(v), true
	}
	for {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr:
			rv = rv.Elem()
		case reflect.Slice, reflect.Array:
			n := rv.Len()
			rt := make([]interface{}, n)
			for i := 0; i < n; i++ {
				rt[i] = rv.Index(i).Interface()
			}
			return rt, true
		default:
			return nil, false
		}
	}
}

func ToSlice(v interface{}) []interface{} {
	if rt, ok := Slice(v); ok {
		return rt
	}
	panic(fmt.Sprintf("convert slice error: %#v", v))
}

func SliceBool(v []bool) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceInt(v []int) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceInt8(v []int8) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceInt16(v []int16) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceInt32(v []int32) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceInt64(v []int64) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceUint(v []uint) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceUint8(v []uint8) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceUint16(v []uint16) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceUint32(v []uint32) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceUint64(v []uint64) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceFloat32(v []float32) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceFloat64(v []float64) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceString(v []string) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceBytes(v [][]byte) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceTime(v []time.Time) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}

func SliceDuration(v []time.Duration) []interface{} {
	n := len(v)
	rt := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rt[i] = v[i]
	}
	return rt
}
