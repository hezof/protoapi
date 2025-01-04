package protoapi

type ProtoDemo struct {
	Name string `protobuf:"name" json:"name,omitempty"`
	Age  int    `protobuf:"age" json:"age,omitempty"`
}
