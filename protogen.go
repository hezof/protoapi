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
type Call func(ctx *Context, in io.Reader) (any, error)

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
	Methods  []*MethodSetting  // methods设置
	Aspects  []ServiceAspect   // aop切面列表
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
	Validate(set *MethodSetting, ctx context.Context) error
}

// MessagePlugin message校验插件
type MessagePlugin func(ctx context.Context, req any, plg *Plugin) error

// FieldPlugin field校验插件
type FieldPlugin func(ctx context.Context, key string, val any, plg *Plugin) error

func AssertMessagePlugin(el string) MessagePlugin {
	name, args := CompilePluginExpression(el)
	if p := globalMessageValidatePluginProvider[name]; p != nil {
		return p(args)
	}
	panic(fmt.Sprintf("assert message validate plugin failed: %s", el))
}

func AssertFieldPlugin(el string) FieldPlugin {
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
