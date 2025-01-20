package protojson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/*
JsonCodec 核心接口, 实现struct的解码与编码.
该接口用于加速proto.Message的JSON解码/编码速度.
*/
type JsonCodec interface {
	DecodeJSON(r *JsonDecoder, k string)
	EncodeJSON(w *JsonEncoder)
}
