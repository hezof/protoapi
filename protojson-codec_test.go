package protoapi

import (
	"fmt"
	"testing"
)

func TestDecoder(t *testing.T) {
	var s *JsonDecoder
	if js, ok := any(s).(JsonCodec); ok {
		js.DecodeJSON(nil)
	} else {
		fmt.Println("not")
	}
}
