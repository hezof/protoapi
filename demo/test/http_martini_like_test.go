package test

import (
	"fmt"
	"github.com/hezof/protoapi"
	"github.com/hezof/protoapi/demo/mdl"
	"net/http"
	"testing"
)

func TestGET(t *testing.T) {
	var rsp *mdl.Student
	err := cli.GET("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, protoapi.NormalResult(&rsp), http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Println(protoapi.ToJson(rsp))
}

func TestPUT(t *testing.T) {
	var rsp *mdl.Student
	err := cli.PUT("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, protoapi.NormalResult(&rsp), http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Println(protoapi.ToJson(rsp))
}

func TestPOST(t *testing.T) {
	var rsp *mdl.Student
	err := cli.POST("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, protoapi.NormalResult(&rsp), http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Println(protoapi.ToJson(rsp))
}

func TestDELETE(t *testing.T) {
	var rsp *mdl.Student
	err := cli.DELETE("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, protoapi.NormalResult(&rsp), http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Println(protoapi.ToJson(rsp))
}

func TestPATCH(t *testing.T) {
	var rsp *mdl.Student
	err := cli.PATCH("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, protoapi.NormalResult(&rsp), http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Println(protoapi.ToJson(rsp))
}

func TestHEAD(t *testing.T) {
	//var rsp *mdl.Student
	err := cli.HEAD("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, nil, http.StatusOK)
	if err != nil {
		panic(err)
	}
	//fmt.Println(protoapi.ToJson(rsp))
}

func TestOPTIONS(t *testing.T) {
	var rsp *mdl.Student
	err := cli.OPTIONS("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, protoapi.NormalResult(&rsp), http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Println(protoapi.ToJson(rsp))
}

func TestCONNECT(t *testing.T) {
	var rsp *mdl.Student
	err := cli.CONNECT("/student", &mdl.Student{
		Id:    1,
		Name:  "张三",
		Age:   8,
		Class: "甲班",
		Scores: map[string]float32{
			"体育": 5,
			"数学": 4,
			"语文": 4,
			"英语": 1,
		},
	}, protoapi.NormalResult(&rsp), http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Println(protoapi.ToJson(rsp))
}
