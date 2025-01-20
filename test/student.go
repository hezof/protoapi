package test

import (
	"github.com/hezof/protoapi"
)

type Sex int32

var Sex_Names = map[int32]string{
	0: "MALE",
	1: "FEMALE",
	2: "UNKNOWN",
}

var Sex_Values = map[string]int32{
	"MALE":    0,
	"FEMALE":  1,
	"UNKNOWN": 2,
}

type Student struct {
	Name  string              `json:"name,omitempty"`
	Sex   Sex                 `json:"sex"`
	Age   uint32              `json:"age,omitempty"`
	Score float32             `json:"score,omitempty"`
	Class []*Student          `json:"class,omitempty"`
	Index map[string]*Student `json:"index,omitempty"`
	Data  []byte              `json:"data,omitempty"`
}

type StdStudent Student

// Student_MessageEncoder  NOTE: 用于jsonpb的MessageEncoder
func (s *Student) EncodeField(w *protoapi.JsonEncoder) {
	protoapi.EncodeString_OmitEmpty(w, `name`, s.Name)
	protoapi.EncodeEnum_OmitEmpty(w, `sex`, s.Sex)
	protoapi.EncodeUint32_OmitEmpty(w, `age`, s.Age)
	protoapi.EncodeFloat_OmitEmpty(w, `score`, s.Score)
	protoapi.EncodeMessageRepeated_OmitEmpty(w, `class`, s.Class)
	protoapi.EncodeMessageMap_OmitEmpty(w, `index`, s.Index)
	protoapi.EncodeBytes_OmitEmpty(w, `data`, s.Data)
}

// Student_MessageDecoder NOTE: 用于jsonpb的FieldDecoder
func (s *Student) DecodeField(r *protoapi.JsonDecoder, f string) {
	switch f {
	case `name`:
		protoapi.DecodeString(r, &s.Name)
	case `sex`:
		protoapi.DecodeEnum(r, &s.Sex, Sex_Names, Sex_Values)
	case `age`:
		protoapi.DecodeUint32(r, &s.Age)
	case `score`:
		protoapi.DecodeFloat(r, &s.Score)
	case `class`:
		protoapi.DecodeMessageRepeated(r, &s.Class)
	case `index`:
		protoapi.DecodeMessageMap(r, &s.Index)
	case `data`:
		protoapi.DecodeBytes(r, &s.Data)
	}
}

var _ protoapi.FieldCodec = (*Student)(nil)
