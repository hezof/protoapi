package protoapi

import (
	"github.com/hezof/protojson"
	"io"
)

var (
	JsonDecoderBufferSize = 8 * 1024 // 默认8K
	JsonEncoderBufferSize = 8 * 1024 // 默认8K
)

func InitJsonBuffer(jsonDecoderBufferSize, jsonEncodeBufferSize int) {
	JsonDecoderBufferSize = jsonDecoderBufferSize
	JsonEncoderBufferSize = jsonEncodeBufferSize
}

func NewDecoder(in io.Reader) *protojson.JsonDecoder {
	return protojson.NewJsonDecoder(in, JsonDecoderBufferSize)
}

func NewEncoder(out io.Writer) *protojson.JsonEncoder {
	return protojson.NewJsonEncoder(out, JsonEncoderBufferSize)
}

func JsonBind(ctx *Context, val protojson.MessageDecoder, body Body) error {
	// 从pool里面取得JsonDecoder

}
