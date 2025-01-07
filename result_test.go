package protoapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type Data struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

var sr = &StatusResult{
	Code:    1111,
	Name:    "1111",
	Message: "消息",
	Data: &Data{
		Name: "myname",
		Age:  40,
	},
}

func TestStatusResult_EncodeJSON(t *testing.T) {
	//sb := new(strings.Builder)
	bs, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}

var bs = `{"code":1111,"name":"1111","message":"消息","data":{"name":"myname","age":40}}`

func TestStatusResult_DecodeJSON(t *testing.T) {
	sr := new(StatusResult)
	in := strings.NewReader(bs)
	err := DecodeJSON(in, sr)
	if err != nil {
		panic(err)
	}
	fmt.Println(sr)
}
