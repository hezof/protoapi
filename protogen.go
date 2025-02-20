package protoapi

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

// Call 每个service的method生成一个相应的HttpFunc闭包, 用于适配restful/websocket/sse等
type Call func(ctx *Context, in io.Reader) (interface{}, error)

// MethodSetting 对应Service.Method的元数据
type MethodSetting struct {
	Meta          *Meta            // 方法元数据
	Call          Call             // 方法回调
	Service       *ServiceSetting  // 服务设置
	MessageExtend MessageExtend    // 消息校验插件
	FieldPlugins  []FieldPlugin    // 字段校验插件, 下标必须与Meta中的Rules对应
	FiledPatterns []*regexp.Regexp // 字段正侧表达式
}

// ServiceSetting 对应Service的元数据
type ServiceSetting struct {
	Impl     any               // service实现
	Desc     *grpc.ServiceDesc // service描述
	HttpOnly bool              // 是否仅用于HTTP
	Aspects  []ServiceAspect   // aop切面列表
	Methods  []*MethodSetting  // methods设置
}

// ServiceRegistry 由protoc-go-gen-protoapi生成的注册器. 专供Server.RegisterService()使用.
type ServiceRegistry func(impl any, aspects []ServiceAspect) *ServiceSetting

// ServiceAspect 切面接口
type ServiceAspect interface {
	// Order 切面执行顺序[主,次]. Before Advice按[major,minor]的升序执行. After Advice按[major,minor]的降序执行.
	Order() [2]int
	// Before Advice执行前置处理, 返回ctx, req作为后面节点入参. 返回err会将执行流程跳至After Advice()
	Before(set *MethodSetting, ctx context.Context, req any) (context.Context, error)
	// After 事后内容. 返回ctx, rsp, err覆盖后面的传递内容.
	After(set *MethodSetting, ctx context.Context, req, rsp any, err error) (context.Context, any, error)
}

// BeforeAspect BeforeAspect与AfterAspect实现ServiceAspect的流程逻辑.
func BeforeAspect(set *MethodSetting, ctx context.Context, req any) (int, context.Context, error) {
	var idx = -1 // 初始为-1. 可能没有aspects
	var err error
	for _, asp := range set.Service.Aspects {
		idx++
		if ctx, err = asp.Before(set, ctx, req); err != nil {
			return idx, ctx, err
		}
	}
	if mv, ok := req.(MessageValidator); ok {
		if err = mv.Validate(set, ctx); err != nil {
			return idx, ctx, err
		}
	}
	return idx, ctx, nil
}

// AfterAspect BeforeAdvice与AfterAdvice实现ServiceAspect的流程逻辑.
func AfterAspect(set *MethodSetting, idx int, ctx context.Context, req, rsp any, err error) (any, error) {
	for idx >= 0 {
		ctx, rsp, err = set.Service.Aspects[idx].After(set, ctx, req, rsp, err)
		idx--
	}
	return rsp, err
}

// MessageValidator 校验接口
type MessageValidator interface {
	Validate(set *MethodSetting, ctx context.Context) error
}

// MessageExtend message校验插件
type MessageExtend func(ctx context.Context, req any) error

// FieldPlugin field校验插件
type FieldPlugin func(ctx context.Context, key string, val any, plg *Plugin) error

// EncodeMeta 断言编码. 用于protogen传值
func EncodeMeta(meta *Meta) string {
	bs, err := proto.Marshal(meta)
	if err != nil {
		panic(fmt.Errorf("ecnode meta error: %v", err))
	}
	return base64.StdEncoding.EncodeToString(bs)
}

// DecodeMeta 断言解码. 用于protogen传值
func DecodeMeta(b64 string) *Meta {
	bs, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(fmt.Errorf("decode meta error: %v", err))
	}
	meta := new(Meta)
	err = proto.Unmarshal(bs, meta)
	if err != nil {
		panic(fmt.Errorf("decode meta error: %v", err))
	}
	return meta
}

// DecodeRequest 解码请求对象. DecodeMessage()很相似, 但不支持泛型!
func DecodeRequest(in io.Reader, req any) error {
	r := GetDecoder(in)
	defer PutDecoder(r)

	switch r.token {
	case ObjectBegin:
		if fc, ok := req.(FieldCodec); ok {
			// 已实现JsonCodec使用protojson加速解码
			r.readObject(fc)
		} else {
			// 未实现JsonCodec使用encoding/std反射解码
			err := UnmarshalJSON(r.dumpObjectOrArray(ObjectBegin), &req)
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
	if err := r.Close(); err != nil {
		return StatusError(http.StatusBadRequest, profile.DefaultBadRequestErrorCode, err.Error())
	}
	return nil
}

// EncodeResponse 编码请求对象
func EncodeResponse(out io.Writer, rsp any) error {
	w := GetEncoder(out)
	defer PutEncoder(w)

	encodeObject(w, rsp)
	return w.Close()
}

/*************************************************
* bool:
**************************************************/

func ParamBool(f func(string) (string, error), key string, ptr *bool) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	*ptr = val == "true"
	return nil
}
func ParamBoolOptional(f func(string) (string, error), key string, ptr **bool) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	if val == "true" {
		*ptr = &_true
	} else {
		*ptr = &_false
	}
	return nil
}
func ParamBoolRepeated(f func(string, bool) ([]string, error), key string, ptr *[]bool, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]bool, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		*ptr = append(*ptr, v == "true")
	}
	return nil
}
func ParamBoolMap(f func(string, bool) (map[string]string, error), key string, ptr *map[string]bool, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]bool)
	}
	for k, v := range val {
		(*ptr)[k] = v == "true"
	}
	return nil
}

/*************************************************
* int32
**************************************************/

func ParamInt32(f func(string) (string, error), key string, ptr *int32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return err
	}
	*ptr = int32(tmp)
	return nil
}
func ParamInt32Optional(f func(string) (string, error), key string, ptr **int32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return err
	}
	ret := int32(tmp)
	*ptr = &ret
	return nil
}
func ParamInt32Repeated(f func(string, bool) ([]string, error), key string, ptr *[]int32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]int32, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, int32(tmp))
	}
	return nil
}
func ParamInt32Map(f func(string, bool) (map[string]string, error), key string, ptr *map[string]int32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]int32)
	}
	for k, v := range val {
		tmp, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		(*ptr)[k] = int32(tmp)
	}
	return nil
}

/*************************************************
* int64
**************************************************/

func ParamInt64(f func(string) (string, error), key string, ptr *int64) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = tmp
	return nil
}
func ParamInt64Optional(f func(string) (string, error), key string, ptr **int64) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = &tmp
	return nil
}
func ParamInt64Repeated(f func(string, bool) ([]string, error), key string, ptr *[]int64, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]int64, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, tmp)
	}
	return nil
}
func ParamInt64Map(f func(string, bool) (map[string]string, error), key string, ptr *map[string]int64, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]int64)
	}
	for k, v := range val {
		tmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		(*ptr)[k] = tmp
	}
	return nil
}

/*************************************************
* uint32
**************************************************/

func ParamUint32(f func(string) (string, error), key string, ptr *uint32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = uint32(tmp)
	return nil
}
func ParamUint32Optional(f func(string) (string, error), key string, ptr **uint32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	ret := uint32(tmp)
	*ptr = &ret
	return nil
}
func ParamUint32Repeated(f func(string, bool) ([]string, error), key string, ptr *[]uint32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]uint32, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, uint32(tmp))
	}
	return nil
}
func ParamUint32Map(f func(string, bool) (map[string]string, error), key string, ptr *map[string]uint32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]uint32)
	}
	for k, v := range val {
		tmp, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		(*ptr)[k] = uint32(tmp)
	}
	return nil
}

/*************************************************
* uint64
**************************************************/

func ParamUint64(f func(string) (string, error), key string, ptr *uint64) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = tmp
	return nil
}
func ParamUint64Optional(f func(string) (string, error), key string, ptr **uint64) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = &tmp
	return nil
}
func ParamUint64Repeated(f func(string, bool) ([]string, error), key string, ptr *[]uint64, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]uint64, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, tmp)
	}
	return nil
}
func ParamUint64Map(f func(string, bool) (map[string]string, error), key string, ptr *map[string]uint64, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]uint64)
	}
	for k, v := range val {
		tmp, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		(*ptr)[k] = tmp
	}
	return nil
}

/*************************************************
* float
**************************************************/

func ParamFloat(f func(string) (string, error), key string, ptr *float32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = float32(tmp)
	return nil
}
func ParamFloatOptional(f func(string) (string, error), key string, ptr **float32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	ret := float32(tmp)
	*ptr = &ret
	return nil
}
func ParamFloatRepeated(f func(string, bool) ([]string, error), key string, ptr *[]float32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]float32, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, float32(tmp))
	}
	return nil
}
func ParamFloatMap(f func(string, bool) (map[string]string, error), key string, ptr *map[string]float32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]float32)
	}
	for k, v := range val {
		tmp, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		(*ptr)[k] = float32(tmp)
	}
	return nil
}

/*************************************************
* double
**************************************************/

func ParamDouble(f func(string) (string, error), key string, ptr *float64) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = tmp
	return nil
}
func ParamDoubleOptional(f func(string) (string, error), key string, ptr **float64) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = &tmp
	return nil
}
func ParamDoubleRepeated(f func(string, bool) ([]string, error), key string, ptr *[]float64, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]float64, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, tmp)
	}
	return nil
}
func ParamDoubleMap(f func(string, bool) (map[string]string, error), key string, ptr *map[string]float64, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]float64)
	}
	for k, v := range val {
		tmp, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		(*ptr)[k] = tmp
	}
	return nil
}

/*************************************************
* string
**************************************************/

func ParamString(f func(string) (string, error), key string, ptr *string) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	*ptr = val
	return nil
}
func ParamStringOptional(f func(string) (string, error), key string, ptr **string) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	*ptr = &val
	return nil
}
func ParamStringRepeated(f func(string, bool) ([]string, error), key string, ptr *[]string, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = val
	return nil
}
func ParamStringMap(f func(string, bool) (map[string]string, error), key string, ptr *map[string]string, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = val
	return nil
}

/*************************************************
* bytes
**************************************************/

func ParamBytes(f func(string) (string, error), key string, ptr *[]byte) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	*ptr, err = base64.StdEncoding.DecodeString(val)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	return nil
}
func ParamBytesOptional(f func(string) (string, error), key string, ptr *[]byte) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	*ptr, err = base64.StdEncoding.DecodeString(val)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	return nil
}
func ParamBytesRepeated(f func(string, bool) ([]string, error), key string, ptr *[][]byte, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([][]byte, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, tmp)
	}
	return nil
}
func ParamBytesMap(f func(string, bool) (map[string]string, error), key string, ptr *map[string][]byte, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string][]byte)
	}
	for k, v := range val {
		(*ptr)[k], err = base64.StdEncoding.DecodeString(v)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
	}
	return nil
}

/*************************************************
* enum
**************************************************/

func ParamEnum[Enum ~int32](f func(string) (string, error), key string, ptr *Enum) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	*ptr = Enum(tmp)
	return nil
}
func ParamEnumOptional[Enum ~int32](f func(string) (string, error), key string, ptr **Enum) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	tmp, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	ret := Enum(tmp)
	*ptr = &ret
	return nil
}
func ParamEnumRepeated[Enum ~int32](f func(string, bool) ([]string, error), key string, ptr *[]Enum, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]Enum, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		tmp, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		*ptr = append(*ptr, Enum(tmp))
	}
	return nil
}
func ParamEnumMap[Enum ~int32](f func(string, bool) (map[string]string, error), key string, ptr *map[string]Enum, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]Enum)
	}
	for k, v := range val {
		tmp, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
		}
		(*ptr)[k] = Enum(tmp)
	}
	return nil
}

/*************************************************
* enum_name
**************************************************/

func ParamEnumName[Enum ~int32](f func(string) (string, error), key string, ptr *Enum, values map[string]int32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	*ptr = Enum(values[val])
	return nil
}
func ParamEnumNameOptional[Enum ~int32](f func(string) (string, error), key string, ptr **Enum, values map[string]int32) error {
	val, err := f(key)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	} else if val == "" {
		return nil
	}
	ret := Enum(values[val])
	*ptr = &ret
	return nil
}
func ParamEnumNameRepeated[Enum ~int32](f func(string, bool) ([]string, error), key string, ptr *[]Enum, values map[string]int32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make([]Enum, 0, len(val))
	} else {
		*ptr = (*ptr)[0:0]
	}
	for _, v := range val {
		*ptr = append(*ptr, Enum(values[v]))
	}
	return nil
}
func ParamEnumNameMap[Enum ~int32](f func(string, bool) (map[string]string, error), key string, ptr *map[string]Enum, values map[string]int32, explode bool) error {
	val, err := f(key, explode)
	if err != nil {
		return StatusError(profile.DefaultBadRequestErrorStatus, profile.DefaultBadRequestErrorCode, fmt.Sprintf("[%v] %v", key, err))
	}
	if *ptr == nil {
		*ptr = make(map[string]Enum)
	}
	for k, v := range val {
		(*ptr)[k] = Enum(values[v])
	}
	return nil
}
