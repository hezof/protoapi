package protojson

import (
	"fmt"
	"testing"
)

func TestDecoder(t *testing.T) {
	var s *JsonDecoder
	if js, ok := any(s).(JsonCodec); ok {
		js.DecodeJSON(nil, "a")
	} else {
		fmt.Println("not")
	}
}
