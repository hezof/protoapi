package protoapi

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"regexp"
)

// Call 每个service的method生成一个相应的HttpFunc闭包, 用于适配restful/websocket/sse等
type Call func(ctx *Context, in io.Reader) (interface{}, error)

// MethodSetting 对应Service.Method的元数据
type MethodSetting struct {
	Meta          *Meta            // 方法元数据
	Call          Call             // 方法回调
	Service       *ServiceSetting  // 服务设置
	MessagePlugin MessagePlugin    // 消息校验插件
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
	Before(setting *MethodSetting, ctx context.Context, req any) (context.Context, error)
	// After 事后内容. 返回ctx, rsp, err覆盖后面的传递内容.
	After(setting *MethodSetting, ctx context.Context, req, rsp any, err error) (context.Context, any, error)
}

// MessageValidator 校验接口
type MessageValidator interface {
	Validate(set *MethodSetting, ctx context.Context) *Error
}

// MessagePlugin message校验插件
type MessagePlugin func(ctx context.Context, req any, plg *Plugin) *Error

// FieldPlugin field校验插件
type FieldPlugin func(ctx context.Context, key string, val any, plg *Plugin) *Error

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
			r.unreadByte() // 回退"{"
			err := UnmarshalJSON(r.dumpObjectOrArray(ObjectBegin), req)
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
	return r.Close()
}

// EncodeResponse 编码请求对象
func EncodeResponse(out io.Writer, rsp any) error {
	w := GetEncoder(out)
	defer PutEncoder(w)

	encodeObject(w, rsp)
	return w.Close()
}

/*************************************************
* bool
**************************************************/

func ParamBool(f func(key string) string, key string, ptr *bool) error {
	return nil
}
func ParamBoolOptional(f func(key string) string, key string, ptr **bool) error {
	return nil
}
func ParamBoolRepeated(f func(key string) []string, key string, ptr *[]bool) error {
	return nil
}
func ParamBoolMap(f func(key string) map[string]string, key string, ptr *map[string]bool) error {
	return nil
}

/*************************************************
* int32
**************************************************/

func ParamInt32(f func(key string) string, key string, ptr *int32) error {
	return nil
}
func ParamInt32Optional(f func(key string) string, key string, ptr **int32) error {
	return nil
}
func ParamInt32Repeated(f func(key string) []string, key string, ptr *[]int32) error {
	return nil
}
func ParamInt32Map(f func(key string) map[string]string, key string, ptr *map[string]int32) error {
	return nil
}

/*************************************************
* int64
**************************************************/

func ParamInt64(f func(key string) string, key string, ptr *int64) error {
	return nil
}
func ParamInt64Optional(f func(key string) string, key string, ptr **int64) error {
	return nil
}
func ParamInt64Repeated(f func(key string) []string, key string, ptr *[]int64) error {
	return nil
}
func ParamInt64Map(f func(key string) map[string]string, key string, ptr *map[string]int64) error {
	return nil
}

/*************************************************
* uint32
**************************************************/

func ParamUint32(f func(key string) string, key string, ptr *uint32) error {
	return nil
}
func ParamUint32Optional(f func(key string) string, key string, ptr **uint32) error {
	return nil
}
func ParamUint32Repeated(f func(key string) []string, key string, ptr *[]uint32) error {
	return nil
}
func ParamUint32Map(f func(key string) map[string]string, key string, ptr *map[string]uint32) error {
	return nil
}

/*************************************************
* uint64
**************************************************/

func ParamUint64(f func(key string) string, key string, ptr *uint64) error {
	return nil
}
func ParamUint64Optional(f func(key string) string, key string, ptr **uint64) error {
	return nil
}
func ParamUint64Repeated(f func(key string) []string, key string, ptr *[]uint64) error {
	return nil
}
func ParamUint64Map(f func(key string) map[string]string, key string, ptr *map[string]uint64) error {
	return nil
}

/*************************************************
* float
**************************************************/

func ParamFloat(f func(key string) string, key string, ptr *float32) error {
	return nil
}
func ParamFloatOptional(f func(key string) string, key string, ptr **float32) error {
	return nil
}
func ParamFloatRepeated(f func(key string) []string, key string, ptr *[]float32) error {
	return nil
}
func ParamFloatMap(f func(key string) map[string]string, key string, ptr *map[string]float32) error {
	return nil
}

/*************************************************
* double
**************************************************/

func ParamDouble(f func(key string) string, key string, ptr *float64) error {
	return nil
}
func ParamDoubleOptional(f func(key string) string, key string, ptr **float64) error {
	return nil
}
func ParamDoubleRepeated(f func(key string) []string, key string, ptr *[]float64) error {
	return nil
}
func ParamDoubleMap(f func(key string) map[string]string, key string, ptr *map[string]float64) error {
	return nil
}

/*************************************************
* string
**************************************************/

func ParamString(f func(key string) string, key string, ptr *string) error {
	return nil
}
func ParamStringOptional(f func(key string) string, key string, ptr **string) error {
	return nil
}
func ParamStringRepeated(f func(key string) []string, key string, ptr *[]string) error {
	return nil
}
func ParamStringMap(f func(key string) map[string]string, key string, ptr *map[string]string) error {
	return nil
}

/*************************************************
* bytes
**************************************************/

func ParamBytes(f func(key string) string, key string, ptr *[]byte) error {
	return nil
}
func ParamBytesOptional(f func(key string) string, key string, ptr *[]byte) error {
	return nil
}
func ParamBytesRepeated(f func(key string) []string, key string, ptr *[][]byte) error {
	return nil
}
func ParamBytesMap(f func(key string) map[string]string, key string, ptr *map[string][]byte) error {
	return nil
}

/*************************************************
* enum
**************************************************/

func ParamEnum[Enum ~int32](f func(key string) string, key string, ptr *Enum, names map[int32]string) error {
	return nil
}
func ParamEnumOptional[Enum ~int32](f func(key string) string, key string, ptr **Enum, names map[int32]string) error {
	return nil
}
func ParamEnumRepeated[Enum ~int32](f func(key string) []string, key string, ptr *[]Enum, names map[int32]string) error {
	return nil
}
func ParamEnumMap[Enum ~int32](f func(key string) map[string]string, key string, ptr *map[string]Enum, names map[int32]string) error {
	return nil
}

/*************************************************
* enum_name
**************************************************/

func ParamEnumName[Enum ~int32](f func(key string) string, key string, ptr *Enum, values map[string]int32) error {
	return nil
}
func ParamEnumNameOptional[Enum ~int32](f func(key string) string, key string, ptr **Enum, values map[string]int32) error {
	return nil
}
func ParamEnumNameRepeated[Enum ~int32](f func(key string) []string, key string, ptr *[]Enum, values map[string]int32) error {
	return nil
}
func ParamEnumNameMap[Enum ~int32](f func(key string) map[string]string, key string, ptr *map[string]Enum, values map[string]int32) error {
	return nil
}
