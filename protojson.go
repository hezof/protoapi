package protoapi

import "github.com/hezof/protoapi/internal/encoding/json"

var (
	MarshalJSON   = json.Marshal
	UnmarshalJSON = json.Unmarshal
)
