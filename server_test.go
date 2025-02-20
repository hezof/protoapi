package protoapi

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestServer_ListenAndServe(t *testing.T) {
	data := []byte(`{"a": {"b": "c"}}`)

	var obj map[string]any
	err := json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(obj)
}
