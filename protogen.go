package protoapi

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
)

/*
protogen.api.go定义下述plugins的API承诺.
- protoc-gen-go-protoapi,
- protoc-gen-go-protojson,
- protoc-gen-go-validator,
- protoc-gen-go-openapi
*/

// Func 每个service的method生成一个相应的HttpFunc闭包, 用于适配restful/websocket/sse等
type Func func(ctx *Context, in io.Reader) (proto.Message, error)

// Registry protoc-go-gen-protoapi生成的注册器. 专供Server.RegisterService()使用.
type Registry func(impl interface{}, aspects []Aspect) *ServiceMeta

type Body int32

const (
	JsonBody  Body = 0  // 解析body使用application/json
	ProtoJson Body = 1  // 解析body使用google.golang.org/protobuf/encoding/protojson
	FormBody  Body = 2  // 解析body使用multipart/form-data或application/x-www-form-urlencoded
	OmitBody  Body = 15 // 忽略解析body
)

type Result int32

const (
	SimpleResult Result = 0 // 结果使用Result包裹
	UnwrapResult Result = 1 // 结果不用Result包裹
	EventsResult Result = 2 // 结果使用Server-Send-Events包裹
)

type Style int32

const (
	SimpleStyle Style = 0
	FormStyle   Style = 1
	JsonStyle   Style = 2
)

// HttpMeta http配置元数据
type HttpMeta struct {
	Get       string // GET请求
	Put       string // PUT请求
	Post      string // POST请求
	Delete    string // DELETE请求
	Options   string // OPTIONS请求
	Head      string // HEAD请求
	Patch     string // PATCH请求
	Trace     string // TRACE请求
	Connect   string // CONNECT请求
	Websocket string // WS请求(可能与GET冲突)
	Body      Body   // body解析方式. 默认json!
	Status    int32  // 成功响应状态码
	Result    Result // 是否"不"包裹Result! 可用WrapperFactory()指定Wrapper实例!
}

// RoleMeta role配置元数据
type RoleMeta struct {
	Code uint64 // 角色标识
	Name string // 角色名称
	Desc string // 角色描述
}

// MethodMeta 对应Service.Method的元数据
type MethodMeta struct {
	Parent       *ServiceMeta // 父服务
	Method       string       // proto方法名称
	FullMethod   string       // proto方法全称:  /package.service/method
	ClientStream bool         // client stream
	ServerStream bool         // server stream
	Http         *HttpMeta    // protoapi.http设置
	Role         *RoleMeta    // protoapi.role设置
	Func         Func         // 方法函数
}

// ServiceMeta 对应Service的元数据
type ServiceMeta struct {
	Impl     interface{}       // service实现
	Desc     *grpc.ServiceDesc // service描述
	Package  string            // proto包名称
	Service  string            // proto服务名称
	HttpOnly bool              // 是否仅用于HTTP
	Aspects  []*Aspect         // aop切面列表
	Methods  []*MethodMeta     // methods设置
}

// Aspect 切面接口
type Aspect interface {
	// Order 切面执行顺序[主,次]. Before Advice按[major,minor]的升序执行. After Advice按[major,minor]的降序执行.
	Order() [2]int
	// Before Advice执行前置处理, 返回ctx, req作为后面节点入参. 返回err会将执行流程跳至After Advice()
	Before(meta *MethodMeta, ctx context.Context, req proto.Message) (context.Context, error)
	// After 事后内容. 返回ctx, rsp, err覆盖后面的传递内容.
	After(meta *MethodMeta, ctx context.Context, req, rsp proto.Message, err error) (context.Context, proto.Message, error)
}

// Validator 校验接口
type Validator interface {
	Validate(ctx context.Context) error
}
