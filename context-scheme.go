package protoapi

import (
	"encoding/json"
	"fmt"
	"github.com/hezof/protoapi/kits"

	"io"
	"net/url"
	"reflect"
	"strings"
)

type values struct {
	json  map[string]json.RawMessage
	form  map[string][]string
	path  *Params
	query map[string][]string
}

func parseValues(c *Context) (*values, error) {
	var ret = new(values)
	var contentType string
	if vs := c.Request.Header["Content-Type"]; len(vs) > 0 {
		contentType = vs[0]
	}
	if c.Request.ContentLength > 0 {
		if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") || strings.HasPrefix(contentType, "multipart/form-data") {
			// form
			if c.Request.PostForm == nil {
				err := c.Request.ParseMultipartForm(c.Handler.FormMaxMemory)
				if err != nil {
					return nil, err
				}
			}
			ret.form = c.Request.PostForm
		} else {
			// json
			err := json.NewDecoder(c.Request.Body).Decode(&ret.json)
			if err != nil && err != io.EOF {
				return nil, err
			}
		}
	}
	ret.path = &c.params
	if c.query == nil {
		if c.Request.URL.RawQuery != "" {
			c.query, _ = url.ParseQuery(c.Request.URL.RawQuery)
		}
	}
	ret.query = c.query

	return ret, nil
}

/*
Scheme

	支持将Request的json/form/path/query等参数绑定到Struct结构.

处理规则:
  - 基础类型(bool/intN/unitN/floatN/complexN/string/[]byte/duration/time)从string适配转换
  - 分片类型(slice)适配多参数
  - 特殊类型(array, struct, map, slice, ptr)当作json内容处理
  - 其他类型转换错误

覆盖顺序(从高到低):
  - json (application/json)
  - form (application/x-www-form-urlencoded或multipart/form-data)
  - path (params)
  - query

需要注意: json字符串赋值不做类型适配处理(例如从string转为int)
*/
func (ctx *Context) Scheme(dst interface{}, tag string) error {

	/*************************************
	 * 检查与分析参数合法性
	 *************************************/
	dstTyp := reflect.TypeOf(dst)
	if dstTyp == nil || dstTyp.Kind() != reflect.Ptr {
		return kits.ErrInvalidStructPointer
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.IsNil() {
		return kits.ErrInvalidMemoryOrNilPointer
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
		return kits.ErrInvalidStructPointer
	}
	/*************************************
	 * 解析上下文数值
	 *************************************/
	vals, err := parseValues(ctx)
	if err != nil {
		return err
	}
	_, err = adaptStruct(vals, dstTyp, &dstVal, tag)
	return err
}

func adaptStruct(vs *values, dstTyp reflect.Type, dstVal *reflect.Value, tag string) (set bool, err error) {
	/*************************************
	 * 反射字段列表
	 *************************************/
__NEXT__:
	for _, fld := range kits.GetStructField(dstTyp) {
		fldKey := fld.Tags[tag]
		if fldKey == "-" {
			continue __NEXT__
		}
		fldTyp := fld.Type
		fldVal := dstVal.Field(fld.Index)
		if fld.Anonymous {
			if fldTyp.Kind() == reflect.Ptr {
				fldTyp = fldTyp.Elem()
				if fldVal.IsNil() {
					tmpVal := reflect.New(fldTyp)
					subVal := tmpVal.Elem()
					subSet := false
					if subSet, err = adaptStruct(vs, fldTyp, &subVal, tag); err != nil {
						return
					} else if subSet {
						fldVal.Set(tmpVal)
					}
				} else {
					fldVal = fldVal.Elem()
					if _, err = adaptStruct(vs, fldTyp, &fldVal, tag); err != nil {
						return
					}
				}
				set = true
			} else {
				if _, err = adaptStruct(vs, fldTyp, &fldVal, tag); err != nil {
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
		// 取值顺序: json/form/path/query
		if js, ok := vs.json[fldKey]; ok {
			if err = json.Unmarshal(js, fldVal.Addr().Interface()); err != nil {
				return
			} else {
				set = true
			}
			continue __NEXT__
		}

		if ss, ok := vs.form[fldKey]; ok {
			if fldTyp.Kind() == reflect.Slice && fldTyp != kits.RTypeBytes {
				if err = adaptSlice(ss, fldTyp, &fldVal); err != nil {
					return
				} else {
					set = true
				}
			} else {
				if err = adaptValue(ss[0], fldTyp, &fldVal); err != nil {
					return
				} else {
					set = true
				}
			}
			continue __NEXT__
		}

		if ps, ok := vs.path.Get(fldKey); ok {
			if fldTyp.Kind() == reflect.Slice && fldTyp != kits.RTypeBytes {
				if err = adaptSlice([]string{ps}, fldTyp, &fldVal); err != nil {
					return
				} else {
					set = true
				}
			} else {
				if err = adaptValue(ps, fldTyp, &fldVal); err != nil {
					return
				} else {
					set = true
				}
			}
			continue __NEXT__
		}

		if ss, ok := vs.query[fldKey]; ok {
			if fldTyp.Kind() == reflect.Slice && fldTyp != kits.RTypeBytes {
				if err = adaptSlice(ss, fldTyp, &fldVal); err != nil {
					return
				} else {
					set = true
				}
			} else {
				if err = adaptValue(ss[0], fldTyp, &fldVal); err != nil {
					return
				} else {
					set = true
				}
			}
			continue __NEXT__
		}
	}
	return
}

func adaptSlice(ss []string, dstTyp reflect.Type, dstVal *reflect.Value) error {

	// 优化
	subTyp := dstTyp.Elem()
	if subTyp == kits.RTypeString {
		dstVal.Set(reflect.ValueOf(ss))
		return nil
	}

	orgLen := len(ss)
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
	for i, v := range ss {
		typ := dstTyp.Elem()
		val := dstVal.Index(i)
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			if val.IsNil() {
				val.Set(reflect.New(typ))
			}
			val = val.Elem()
		}
		if err := adaptValue(v, typ, &val); err != nil {
			return err
		}
	}
	return nil
}

func adaptValue(org string, dstTyp reflect.Type, dstVal *reflect.Value) error {

	// 优化几种常见特殊的类型
	if dstTyp == kits.RTypeDuration {
		if x, ok := kits.Duration(org); ok {
			dstVal.SetInt(int64(x))
			return nil
		}
	} else if dstTyp == kits.RTypeTime {
		if x, ok := kits.Time(org); ok {
			dstVal.Set(reflect.ValueOf(x))
			return nil
		}
	} else if dstTyp == kits.RTypeBytes {
		if x, ok := kits.Bytes(org); ok {
			dstVal.Set(reflect.ValueOf(x))
			return nil
		}
	}

	switch dstTyp.Kind() {
	case reflect.String:
		dstVal.SetString(org)
		return nil
	case reflect.Bool:
		if x, ok := kits.Bool(org); ok {
			dstVal.SetBool(x)
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if x, ok := kits.Int64(org); ok {
			dstVal.SetInt(x)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if x, ok := kits.Uint64(org); ok {
			dstVal.SetUint(x)
			return nil
		}
	case reflect.Uintptr:
		if x, ok := kits.Uintptr(org); ok {
			dstVal.Set(reflect.ValueOf(x))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if x, ok := kits.Float64(org); ok {
			dstVal.SetFloat(x)
			return nil
		}
	case reflect.Complex64, reflect.Complex128:
		if x, ok := kits.Complex128(org); ok {
			dstVal.SetComplex(x)
			return nil
		}
	case reflect.Array, reflect.Struct, reflect.Map, reflect.Slice, reflect.Ptr:
		// 其他当成json处理
		if err := json.Unmarshal([]byte(org), dstVal.Addr().Interface()); err != nil {
			return err
		} else {
			return nil
		}
	}

	return fmt.Errorf("convert type error: %T -> %v", org, dstTyp)
}
