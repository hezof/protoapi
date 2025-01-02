package protoapi

import (
	"fmt"
)
import "testing"

type Demo struct {
	Name string `json:"name"`
}

func TestContext(t *testing.T) {
	s := NewServer(&Config{
		HttpAddr: ":80",
	})

	s.Any("/test/:param", func(ctx *Context) {
		fmt.Println("Accept-Language: ", ctx.HeaderValue("Accept-Language"))
		fmt.Println("param: ", ctx.ParamValue("param"))
		fmt.Println("query: ", ctx.QueryValue("query"))
		var d *Demo
		err := ctx.BindJson(&d)
		if err != nil {
			ctx.WriteErrorResult(err)
			return
		}
		ctx.WriteApplyResult(d)
		panic("test")
	})

	if err := s.ListenAndServe(); err != nil {
		fmt.Errorf("error: %+v\n", err)
	}
}

func TestSprintf(t *testing.T) {
	fmt.Println(DecodeBase64RawURL(`dGVzdDI`))
}

func TestFullPath(t *testing.T) {
	var v interface{}
	if vt, ok := v.(*Demo); ok {
		fmt.Println(v)
	} else {
		fmt.Println(vt)
	}
}

type ProtoDemo struct {
	Name string `protobuf:"name" json:"name,omitempty"`
	Age  int    `protobuf:"age" json:"age,omitempty"`
}

func TestJson(t *testing.T) {
	var d = &ProtoDemo{
		Name: "测试员",
	}

	var bs []byte

	bs, _ = iteratorIgnoreOmitempty.Marshal(d)
	fmt.Printf("%s\n", bs)

	bs, _ = iterator.Marshal(d)
	fmt.Printf("%s\n", bs)
}
