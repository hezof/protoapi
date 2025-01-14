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
	protogen.go定义protogen生成代码使用的API!
	所有API承诺向后兼容! 否则影响protogen历史版本...
*/

// Call 每个service的method生成一个相应的HttpFunc闭包, 用于适配restful/websocket/sse等
type Call func(ctx *Context, in io.Reader) (any, error)

// ServiceRegistry protoc-go-gen-protoapi生成的注册器. 专供Server.RegisterService()使用.
type ServiceRegistry func() *ServiceSetting

// MethodSetting 对应Service.Method的元数据
type MethodSetting struct {
	parent          *ServiceSetting // 父节点设置
	fullMethod      string          // 方法全名
	Package         string          // 包名
	Service         string          // 服务名
	Method          string          // 方法名
	StreamingClient bool            // streaming client
	StreamingServer bool            // streaming server
	Http                            // protoapi.http元数据
	Role                            // protoapi.role元数据
	Call                            // 方法函数
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

// MessageValidator message校验接口
type MessageValidator interface {
	Validate(ctx context.Context) error
}

// MessageValidatePlugin message校验插件
type MessageValidatePlugin func(ctx context.Context, req any, err *Error) error

// FieldValidatePlugin field校验插件
type FieldValidatePlugin func(ctx context.Context, key string, val any, err *Error) error

func AssertMessageValidatePlugin(el string) MessageValidatePlugin {
	name, args := CompilePluginExpression(el)
	if p := globalMessageValidatePluginProvider[name]; p != nil {
		return p(args)
	}
	panic(fmt.Sprintf("assert message validate plugin failed: %s", el))
}

func AssertFieldValidatePlugin(el string) FieldValidatePlugin {
	name, args := CompilePluginExpression(el)
	if p := globalFieldValidatePluginProvider[name]; p != nil {
		return p(args)
	}
	panic(fmt.Sprintf("assert field validate plugin failed: %s", el))
}

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
