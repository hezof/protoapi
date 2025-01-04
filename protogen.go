package protoapi

import (
	"context"
	"encoding/base64"
	"fmt"
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

// Call 每个service的method生成一个相应的HttpFunc闭包, 用于适配restful/websocket/sse等
type Call func(ctx *Context, in io.Reader) (any, error)

// ServiceRegistry protoc-go-gen-protoapi生成的注册器. 专供Server.RegisterService()使用.
type ServiceRegistry func(impl interface{}, aspects []ServiceAspect) *ServiceSetting

// MethodSetting 对应Service.Method的元数据
type MethodSetting struct {
	parent     *ServiceSetting // 父节点设置
	Package    string          // 包名
	Service    string          // 服务名
	Method     string          // 方法名
	FullMethod string          // 方法全名
	Http                       // protoapi.http元数据
	Role                       // protoapi.role元数据
	Call                       // 方法函数
}

// ServiceSetting 对应Service的元数据
type ServiceSetting struct {
	Desc     *grpc.ServiceDesc // service描述
	Impl     any               // service实现
	HttpOnly bool              // 是否仅用于HTTP
	Aspects  []ServiceAspect   // aop切面列表
	Methods  []*MethodSetting  // methods设置
}

// ServiceAspect 切面接口
type ServiceAspect interface {
	// Order 切面执行顺序[主,次]. Before Advice按[major,minor]的升序执行. After Advice按[major,minor]的降序执行.
	Order() [2]int
	// Before Advice执行前置处理, 返回ctx, req作为后面节点入参. 返回err会将执行流程跳至After Advice()
	Before(meta *MethodSetting, ctx context.Context, req any) (context.Context, error)
	// After 事后内容. 返回ctx, rsp, err覆盖后面的传递内容.
	After(meta *MethodSetting, ctx context.Context, req, rsp any, err error) (context.Context, any, error)
}

// MessageValidator 校验接口
type MessageValidator interface {
	Validate(ctx context.Context) error
}

type MessageValidatePlugin func(ctx context.Context, req any, err *Error) error

// AssertEncode 断言编码. 用于protogen传值
func AssertEncode(msg proto.Message) string {
	bs, err := proto.Marshal(msg)
	if err != nil {
		panic(fmt.Errorf("assert ecnode error: %v", err))
	}
	return base64.StdEncoding.EncodeToString(bs)
}

// AssertDecode 断言解码. 用于protogen传值
func AssertDecode(b64 string, msg proto.Message) {
	bs, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(fmt.Errorf("assert decode error: %v", err))
	}
	err = proto.Unmarshal(bs, msg)
	if err != nil {
		panic(fmt.Errorf("assert decode error: %v", err))
	}
}
