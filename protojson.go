package protoapi

import "github.com/hezof/protoapi/internal/encoding/json"

var (
	MarshalJSON   = json.Marshal
	UnmarshalJSON = json.Unmarshal
)

const (
	leftBrace    byte = '{'
	rightBrace   byte = '}'
	leftBracket  byte = '['
	rightBracket byte = ']'
	comma        byte = ','
	colon        byte = ':'
	quotes       byte = '"'
)

const encode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const padChar = '='

const chars = "0123456789abcdef"

func table(falseValues ...int) [128]bool {
	ret := [128]bool{}

	for i := 0; i < 128; i++ {
		ret[i] = true
	}

	for _, v := range falseValues {
		ret[v] = false
	}

	return ret
}

var (
	//escapeHtmlTable   = table(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '&', '<', '>', '\\')
	noEscapeHtmlTable = table(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '\\')
)
