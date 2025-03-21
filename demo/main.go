package main

import (
	"context"
	"fmt"
	"github.com/hezof/core"
	"github.com/hezof/log"
	"github.com/hezof/protoapi"
	"github.com/hezof/protoapi/demo/api"
	"github.com/hezof/protoapi/demo/biz"
	"github.com/hezof/protoapi/demo/mdl"
)

func main() {

	defer log.Flush()

	// 0. 安装插件
	protoapi.InstallMessageExtendProvider("message", func(args []string) protoapi.MessageExtend {
		return func(ctx context.Context, req any) error {
			if book, ok := req.(*api.Book); ok {
				fmt.Println(core.ToJson(book))
			}
			return nil
		}
	})

	protoapi.InstallFieldPluginProvider("check_sale", func(args []string) protoapi.FieldPlugin {
		return func(ctx context.Context, key string, val any, plg *protoapi.Plugin) error {
			fmt.Println(key, core.ToJson(val))
			return nil
		}
	})

	svr := protoapi.NewServer(&protoapi.Config{
		HttpAddr: ":8080",
		GrpcAddr: ":9090",
	})

	// 1. 测试martini-like api
	svr.GET("/student", student)
	svr.PUT("/student", student)
	svr.POST("/student", student)
	svr.DELETE("/student", student)
	svr.PATCH("/student", student)
	svr.HEAD("/student", student)
	svr.OPTIONS("/student", student)
	svr.CONNECT("/student", student)
	svr.Static("/demo", `D:\Workspace\klibsrc\protoapi\demo`)
	svr.StaticFile("/go.mod", `D:\Workspace\klibsrc\protoapi\demo\go.mod`)

	svr.RegisterService(api.ParamsRegistry, new(biz.ParamsImplement))
	if err := svr.RegisterService(api.StoreRegistry, new(biz.StoreImplement)); err != nil {
		log.Error("register error: %v", err)
	}

	if err := svr.ListenAndServe(); err != nil {
		panic(err)
	}
}

func student(ctx *protoapi.Context) {
	var student *mdl.Student
	if err := ctx.Scheme(&student, "json"); err != nil {
		log.Error("Scheme error: %v", err)
		return
	}
	fmt.Println("/demo/get: " + core.ToJson(student))
	if err := ctx.WriteApplyResult(student); err != nil {
		log.Error("WriteApplyResult error: %v", err)
		return
	}
}
