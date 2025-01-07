package kits

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

var (
	ErrInvalidStructPointer      = errors.New("invalid struct pointer")        // 非法结构体指针
	ErrInvalidStructTablePointer = errors.New("invalid struct table pointer")  // 非法结构体分片指针
	ErrInvalidStructSlicePointer = errors.New("invalid struct slice pointer")  // 非法结构体映射指针
	ErrInvalidMemoryOrNilPointer = errors.New("invalid memory or nil pointer") // 不可用内存或空指针
	SimpleBinder                 = &StructBinder{Fields: GenStructField}       // struct解析无缓存
	CachedBinder                 = &StructBinder{Fields: GetStructField}       // struct解析带缓存
)

var (
	RTypeBool          = reflect.TypeOf(false)
	RTypeInt           = reflect.TypeOf(int(0))
	RTypeInt8          = reflect.TypeOf(int8(0))
	RTypeInt16         = reflect.TypeOf(int16(0))
	RTypeInt32         = reflect.TypeOf(int32(0))
	RTypeInt64         = reflect.TypeOf(int64(0))
	RTypeUint          = reflect.TypeOf(uint(0))
	RTypeUint8         = reflect.TypeOf(uint8(0))
	RTypeUint16        = reflect.TypeOf(uint16(0))
	RTypeUint32        = reflect.TypeOf(uint32(0))
	RTypeUint64        = reflect.TypeOf(uint64(0))
	RTypeFloat32       = reflect.TypeOf(float32(0.0))
	RTypeFloat64       = reflect.TypeOf(0.0)
	RTypeTime          = reflect.TypeOf(time.Time{})
	RTypeDuration      = reflect.TypeOf(time.Duration(0))
	RTypeBytes         = reflect.TypeOf(([]byte)(nil))
	RTypeRunes         = reflect.TypeOf(([]rune)(nil))
	RTypeString        = reflect.TypeOf("")
	RTypeGenericStruct = reflect.TypeOf((map[string]interface{})(nil))
	RTypeGenericSlice  = reflect.TypeOf(([]interface{})(nil))
)

type StructField struct {
	Name      string            // 名称
	Type      reflect.Type      // 类型
	Tags      map[string]string // tag
	Index     int               // 下标值. unexported字段会忽略掉
	Anonymous bool              // is an embedded field
}

type StructBinder struct {
	Fields func(typ reflect.Type) []*StructField
}

/*
MapSlice2StructSlice 将map分片转换成struct分片
- org: map[string]interface{}分片, []map[string]interface{}
- dst: struct分片, []T
- tag: 用于关联的tag名字
*/
func (sb StructBinder) MapSlice2StructSlice(org []interface{}, dst interface{}, tag string) error {
	// 忽略空参数
	var orgLen = len(org)
	if orgLen == 0 {
		return nil
	}

	// 适配结构体
	var dstTyp = reflect.TypeOf(dst)
	if dstTyp == nil || dstTyp.Kind() != reflect.Ptr {
		return ErrInvalidStructSlicePointer
	}

	var dstVal = reflect.ValueOf(dst)
	if dstVal.IsNil() {
		return ErrInvalidMemoryOrNilPointer
	}

	dstTyp = dstTyp.Elem()
	dstVal = dstVal.Elem()
	for dstTyp.Kind() == reflect.Ptr {
		dstTyp = dstTyp.Elem()
		if dstVal.IsNil() {
			dstVal.Set(reflect.New(dstTyp))
		}
		dstVal = dstVal.Elem()
	}

	// 元素可能是struct或*struct, 不宜作出判断
	if dstTyp.Kind() != reflect.Slice {
		return ErrInvalidStructSlicePointer
	}

	if n := dstVal.Cap(); n >= orgLen {
		// 容量足够, 但要重置元素
		dstVal.SetLen(orgLen)
		zero := reflect.Zero(dstTyp.Elem())
		for i := 0; i < n; i++ {
			dstVal.Field(i).Set(zero)
		}
	} else {
		// 容量不足
		dstVal.Set(reflect.MakeSlice(dstTyp, orgLen, orgLen))
	}

	for i, ov := range org {
		// 忽略空值
		if ov == nil {
			continue
		}
		// 要求是map
		typ := dstTyp.Elem()
		mv, ok := ov.(map[string]interface{})
		if !ok {
			return fmt.Errorf("convert type error: %T -> %v", ov, typ)
		}
		val := dstVal.Index(i)
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			if val.IsNil() {
				val.Set(reflect.New(typ))
			}
			val = val.Elem()
		}
		if typ.Kind() != reflect.Struct {
			return ErrInvalidStructSlicePointer
		}
		if _, err := sb.AdaptStruct(mv, typ, &val, tag); err != nil {
			return err
		}
	}
	return nil
}

/*
MapTable2StructTable 将map键值表转换成struct键值表
- org: map[string]interface{}键值表, map[string]map[string]interface{}
- dst: struct键值表, map[string]T
- tag: 用于关联的tag名字
*/
func (sb StructBinder) MapTable2StructTable(org map[string]interface{}, dst interface{}, tag string) error {
	// 忽略空参数
	var orgLen = len(org)
	if orgLen == 0 {
		return nil
	}

	// 适配结构体
	var dstTyp = reflect.TypeOf(dst)
	if dstTyp == nil || dstTyp.Kind() != reflect.Ptr {
		return ErrInvalidStructSlicePointer
	}

	var dstVal = reflect.ValueOf(dst)
	if dstVal.IsNil() {
		return ErrInvalidMemoryOrNilPointer
	}

	dstTyp = dstTyp.Elem()
	dstVal = dstVal.Elem()
	for dstTyp.Kind() == reflect.Ptr {
		dstTyp = dstTyp.Elem()
		if dstVal.IsNil() {
			dstVal.Set(reflect.New(dstTyp))
		}
		dstVal = dstVal.Elem()
	}

	// 元素可能是struct或*struct, 不宜作出判断.
	if dstTyp.Kind() != reflect.Map || dstTyp.Key().Kind() != reflect.String {
		return ErrInvalidStructTablePointer
	}

	dstVal.Set(reflect.MakeMapWithSize(dstTyp, orgLen))
	for k, ov := range org {
		// 要求非nil
		if ov == nil {
			continue
		}
		// 要求是map
		typ := dstTyp.Elem()
		mv, ok := ov.(map[string]interface{})
		if !ok {
			return fmt.Errorf("convert type error: %T -> %v", ov, typ)
		}
		val := reflect.New(typ).Elem()
		ref := &val
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			if ref.IsNil() {
				ref.Set(reflect.New(typ))
			}
			tmp := ref.Elem()
			ref = &tmp
		}
		if typ.Kind() != reflect.Struct {
			return ErrInvalidStructTablePointer
		}
		if _, err := sb.AdaptStruct(mv, typ, ref, tag); err != nil {
			return err
		}
		dstVal.SetMapIndex(reflect.ValueOf(k), val)
	}
	return nil
}

/*
Map2Struct 将map转换成struct
- org: map[string]interface{}
- dst: struct
- tag: 用于关联的tag名字
*/
func (sb StructBinder) Map2Struct(org map[string]interface{}, dst interface{}, tag string) error {

	if org == nil {
		return nil
	}

	var dstTyp = reflect.TypeOf(dst)
	if dstTyp == nil || dstTyp.Kind() != reflect.Ptr {
		return ErrInvalidStructPointer
	}

	var dstVal = reflect.ValueOf(dst)
	if dstVal.IsNil() {
		return ErrInvalidMemoryOrNilPointer
	}

	dstTyp = dstTyp.Elem()
	dstVal = dstVal.Elem()
	for dstTyp.Kind() == reflect.Ptr {
		dstTyp = dstTyp.Elem()
		if dstVal.IsNil() {
			dstVal.Set(reflect.New(dstTyp))
		}
		dstVal = dstVal.Elem()
	}
	if dstTyp.Kind() != reflect.Struct {
		return ErrInvalidStructPointer
	}

	_, err := sb.AdaptStruct(org, dstTyp, &dstVal, tag)
	return err
}

/*
AdaptStruct 调用者必须确保org不为nil
*/
func (sb StructBinder) AdaptStruct(org map[string]interface{}, dstTyp reflect.Type, dstVal *reflect.Value, tag string) (set bool, err error) {

__NEXT__:
	for _, fld := range sb.Fields(dstTyp) {
		fldKey := fld.Tags[tag]
		if fldKey == "-" {
			continue __NEXT__
		}
		fldTyp := fld.Type
		fldVal := dstVal.Field(fld.Index)
		if fld.Anonymous {
			if fldTyp.Kind() == reflect.Ptr {
				fldTyp = fldTyp.Elem()
				tmpVal := reflect.New(fldTyp)
				subVal := tmpVal.Elem()
				subSet := false
				if subSet, err = sb.AdaptStruct(org, fldTyp, &subVal, tag); err != nil {
					return
				} else if subSet {
					set = true
					fldVal.Set(tmpVal)
				}
			} else {
				if _, err = sb.AdaptStruct(org, fldTyp, &fldVal, tag); err != nil {
					return
				} else {
					set = true
				}
			}
			continue __NEXT__
		}
		if fldKey == "" {
			fldKey = fld.Name
		}
		if val := org[fldKey]; val != nil {
			_set := false
			if _set, err = sb.AdaptValue(val, fldTyp, &fldVal, tag); err != nil {
				return
			} else if _set {
				set = true
			}
		}
	}
	return
}

/*
AdaptValue 调用者必须确保org不为nil
*/
func (sb StructBinder) AdaptValue(org interface{}, dstTyp reflect.Type, dstVal *reflect.Value, tag string) (bool, error) {

	switch dstTyp.Kind() {
	case reflect.Bool:
		if x, ok := Bool(org); ok {
			dstVal.SetBool(x)
			return true, nil
		}
	case reflect.Int:
		if x, ok := Int64(org); ok {
			dstVal.SetInt(x)
			return true, nil
		}
	case reflect.Int8:
		if x, ok := Int64(org); ok {
			dstVal.SetInt(x)
			return true, nil
		}
	case reflect.Int16:
		if x, ok := Int64(org); ok {
			dstVal.SetInt(x)
			return true, nil
		}
	case reflect.Int32:
		if x, ok := Int64(org); ok {
			dstVal.SetInt(x)
			return true, nil
		}
	case reflect.Int64:
		if dstTyp == RTypeDuration {
			if x, ok := Duration(org); ok {
				dstVal.SetInt(int64(x))
				return true, nil
			}
		} else {
			if x, ok := Int64(org); ok {
				dstVal.SetInt(x)
				return true, nil
			}
		}
	case reflect.Uint:
		if x, ok := Uint64(org); ok {
			dstVal.SetUint(x)
			return true, nil
		}
	case reflect.Uint8:
		if x, ok := Uint64(org); ok {
			dstVal.SetUint(x)
			return true, nil
		}
	case reflect.Uint16:
		if x, ok := Uint64(org); ok {
			dstVal.SetUint(x)
			return true, nil
		}
	case reflect.Uint32:
		if x, ok := Uint64(org); ok {
			dstVal.SetUint(x)
			return true, nil
		}
	case reflect.Uint64:
		if x, ok := Uint64(org); ok {
			dstVal.SetUint(x)
			return true, nil
		}
	case reflect.Uintptr:
		if x, ok := Uintptr(org); ok {
			dstVal.Set(reflect.ValueOf(x))
			return true, nil
		}
	case reflect.Float32:
		if x, ok := Float64(org); ok {
			dstVal.SetFloat(x)
			return true, nil
		}
	case reflect.Float64:
		if x, ok := Float64(org); ok {
			dstVal.SetFloat(x)
			return true, nil
		}
	case reflect.Complex64:
		if x, ok := Complex128(org); ok {
			dstVal.SetComplex(x)
			return true, nil
		}
	case reflect.Complex128:
		if x, ok := Complex128(org); ok {
			dstVal.SetComplex(x)
			return true, nil
		}
	case reflect.Array:
		orgTyp := reflect.TypeOf(org)
		if orgTyp == dstTyp || orgTyp.AssignableTo(dstTyp) {
			dstVal.Set(reflect.ValueOf(org))
			return true, nil
		}
	case reflect.Chan:
		orgTyp := reflect.TypeOf(org)
		if orgTyp == dstTyp || orgTyp.AssignableTo(dstTyp) {
			dstVal.Set(reflect.ValueOf(org))
			return true, nil
		}
	case reflect.Func:
		orgTyp := reflect.TypeOf(org)
		if orgTyp == dstTyp || orgTyp.AssignableTo(dstTyp) {
			dstVal.Set(reflect.ValueOf(org))
			return true, nil
		}
	case reflect.Interface:
		orgTyp := reflect.TypeOf(org)
		if orgTyp == dstTyp || orgTyp.AssignableTo(dstTyp) {
			dstVal.Set(reflect.ValueOf(org))
			return true, nil
		}
	case reflect.Map:
		orgTyp := reflect.TypeOf(org)
		if orgTyp == dstTyp || orgTyp.AssignableTo(dstTyp) {
			dstVal.Set(reflect.ValueOf(org))
			return true, nil
		}
		// map[string]interface到map[string]struct
		if orgTyp == RTypeGenericStruct {
			if dstTyp.Key().Kind() == reflect.String {

				// 需要判空
				if dstVal.IsNil() {
					dstVal.Set(reflect.MakeMap(dstTyp))
				}

				var set bool
				typ := dstTyp.Elem()
				// 已经判断orgTyp是GenericStruct
				for k, v := range org.(map[string]interface{}) {
					if v == nil {
						continue
					}
					val := reflect.New(typ).Elem()
					if _set, err := sb.AdaptValue(v, typ, &val, tag); err != nil {
						return false, err
					} else if _set {
						dstVal.SetMapIndex(reflect.ValueOf(k), val)
						set = true
					}
				}
				return set, nil
			}
		}
	case reflect.Ptr:

		typ := dstTyp.Elem()
		if dstVal.IsNil() {
			// 如果是空指针,则分配新值. 并判断是否有set操作, 有的话才设置
			tmp := reflect.New(typ)
			val := tmp.Elem()
			if _set, err := sb.AdaptValue(org, typ, &val, tag); err != nil {
				return false, err
			} else if _set {
				dstVal.Set(tmp) // 中间值是个指针
				return true, nil
			} else {
				return false, nil
			}
		} else {
			val := dstVal.Elem()
			return sb.AdaptValue(org, typ, &val, tag)
		}

	case reflect.Slice:
		// []byte
		if dstTyp == RTypeBytes {
			if x, ok := Bytes(org); ok {
				dstVal.Set(reflect.ValueOf(x))
				return true, nil
			}
		} else {

			orgTyp := reflect.TypeOf(org)
			if orgTyp == dstTyp || orgTyp.AssignableTo(dstTyp) {
				dstVal.Set(reflect.ValueOf(org))
				return true, nil
			}
			// []interface{}
			if orgTyp == RTypeGenericSlice {
				// 已经判断orgTyp是GenericSlice
				orgVal := org.([]interface{})
				orgLen := len(orgVal)
				if n := dstVal.Cap(); n > orgLen {
					// 重置
					dstVal.SetLen(orgLen)
					zero := reflect.Zero(dstTyp.Elem())
					for i := 0; i < orgLen; i++ {
						dstVal.Index(i).Set(zero)
					}
				} else {
					dstVal.Set(reflect.MakeSlice(dstTyp, orgLen, orgLen))
				}
				var set bool
				for i, v := range orgVal {
					if v == nil {
						continue
					}
					typ := dstTyp.Elem()
					val := dstVal.Index(i)
					for typ.Kind() == reflect.Ptr {
						typ = typ.Elem()
						if val.IsNil() {
							val.Set(reflect.New(typ))
						}
						val = val.Elem()
					}
					if _set, err := sb.AdaptValue(v, typ, &val, tag); err != nil {
						return false, err
					} else if _set {
						set = true
					}
				}
				return set, nil
			}
		}

	case reflect.String:
		if x, ok := String(org); ok {
			dstVal.SetString(x)
			return true, nil
		}
	case reflect.Struct:
		// time.Time
		if dstTyp == RTypeTime {
			if x, ok := Time(org); ok {
				dstVal.Set(reflect.ValueOf(x))
				return true, nil
			}
		} else {
			orgTyp := reflect.TypeOf(org)
			if orgTyp == dstTyp || orgTyp.AssignableTo(dstTyp) {
				dstVal.Set(reflect.ValueOf(org))
				return true, nil
			}
			// map[string]interface{}
			if orgTyp == RTypeGenericStruct {
				// 已经判断orgTyp是GenericStruct
				if _set, err := sb.AdaptStruct(org.(map[string]interface{}), dstTyp, dstVal, tag); err != nil {
					return false, err
				} else {
					return _set, nil
				}
			}
		}
	case reflect.UnsafePointer:
		if x, ok := UnsafePointer(org); ok {
			dstVal.Set(reflect.ValueOf(x))
			return true, nil
		}
	}
	return false, fmt.Errorf("convert type error: %T -> %v", org, dstTyp)
}

/************************************************
 * 全局默认采用带缓存Binder
 ************************************************/

func MapSlice2StructSlice(org []interface{}, dst interface{}, tag string) error {
	return CachedBinder.MapSlice2StructSlice(org, dst, tag)
}

func MapTable2StructTable(org map[string]interface{}, dst interface{}, tag string) error {
	return CachedBinder.MapTable2StructTable(org, dst, tag)
}

func Map2Struct(org map[string]interface{}, dst interface{}, tag string) error {
	return CachedBinder.Map2Struct(org, dst, tag)
}

/************************************************
 * struct辅助结构
 ************************************************/
var (
	fieldsCache sync.Map
)

func GetStructField(typ reflect.Type) []*StructField {
	ret, ok := fieldsCache.Load(typ)
	if !ok {
		ret = GenStructField(typ)
		fieldsCache.Store(typ, ret)
	}
	return ret.([]*StructField)
}

func GenStructField(typ reflect.Type) (ret []*StructField) {
	n := typ.NumField()
	ret = make([]*StructField, 0, n)
	for i := 0; i < n; i++ {
		fld := typ.Field(i)
		//  It is empty for upper case (exported) field names.
		if fld.PkgPath == "" || fld.Anonymous {
			// break unexported fields
			ret = append(ret, &StructField{
				Name:      fld.Name,
				Type:      fld.Type,
				Tags:      AllTag(string(fld.Tag)),
				Index:     i,
				Anonymous: fld.Anonymous,
			})
		}
	}
	n = len(ret)
	ret = ret[0:n:n] // trim
	return
}

// 实现tag解析: key:"v1,v2,...", 提取key:v1即可
func AllTag(str string) map[string]string {
	ret := make(map[string]string)
_ITER:
	for i, n := 0, len(str); i < n; i++ {
		// 查找非空白
		for IsSpace(str[i]) {
			if i++; i >= n {
				break _ITER
			}
		}
		mark := i

		// 查找冒号
		for str[i] != ':' {
			if i++; i >= n {
				break _ITER
			}
		}
		end := i - 1
		for IsSpace(str[end]) {
			if end--; end == mark {
				break
			}
		}
		key := str[mark : end+1]

		// 查找左引号
		for str[i] != '"' && str[i] != '`' {
			if i++; i >= n {
				break _ITER
			}
		}
		i++ // 引号下位
		mark = i

		// 查找右引号
		for str[i] != '"' && str[i] != '`' {
			if i++; i >= n {
				break _ITER
			}
		}
		end = i

		// 解析a,b,c里面的a
		pos := mark
		for pos < end && str[pos] != ',' {
			pos++
		}
		val := str[mark:pos]
		ret[key] = val
	}
	return ret
}
func IsSpace(b byte) bool {
	switch b {
	case '\t':
		return true
	case '\n':
		return true
	case '\v':
		return true
	case '\f':
		return true
	case '\r':
		return true
	case ' ':
		return true
	case 0x85:
		return true
	case 0xA0:
		return true
	}
	return false
}
