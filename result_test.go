package protoapi

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	sr.Data = "abc"
	in := strings.NewReader(bs)
	err := DecodeJSON(in, sr)
	if err != nil {
		panic(err)
	}
	fmt.Println(sr.Data)
}

func TestType(t *testing.T) {
	var v ***int
	rv := reflect.ValueOf(&v)
	rt := rv.Type()
	for rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		if rv.IsNil() {
			if rv.CanSet() {
				rv.Set(reflect.New(rt))
			}
		}
		rv = rv.Elem()
	}
	if rv.CanSet() {
		rv.SetInt(123)
	}
	fmt.Println(***v)
}

func TestUnmarshal(t *testing.T) {
	var v ***int
	err := json.Unmarshal([]byte(`123`), v)
	if err != nil {
		panic(err)
	}
}
