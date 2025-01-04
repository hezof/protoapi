package protoapi

import (
	"net/http"
	"testing"
)

type ProtoDemo struct {
	Name string `protobuf:"name" json:"name,omitempty"`
	Age  int    `protobuf:"age" json:"age,omitempty"`
}

func TestJson(t *testing.T) {
	http.ListenAndServe()
}
