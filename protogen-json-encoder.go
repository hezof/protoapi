package protoapi

import (
	"encoding/json"
	"io"
	"unicode/utf8"
)

func NewJsonEncoder(out io.Writer, size int) *JsonEncoder {
	if size < profile.MinimumBufferLength {
		size = profile.MinimumBufferLength
	}
	return &JsonEncoder{
		out:  out,
		buff: make([]byte, 0, size),
	}
}

type JsonEncoder struct {
	out        io.Writer
	buff       []byte
	number     [32]byte // 数值缓存区
	firstError error    // 上下文错误
}

func (w *JsonEncoder) Encode(v MessageEncoder) {
	v.EncodeJSON(w)
}

func (w *JsonEncoder) ensure(n int) {
	if w.out != nil && cap(w.buff)-len(w.buff) < n {
		_, err := w.out.Write(w.buff)
		if err != nil {
			if w.firstError == nil {
				w.firstError = err
			}
			return
		}
		w.buff = w.buff[0:0]
	}
}

// Close 关闭写流, 并返回剩余buff
func (w *JsonEncoder) Close() ([]byte, error) {
	if w.firstError != nil {
		return nil, w.firstError
	}
	if w.out != nil && len(w.buff) > 0 {
		_, err := w.out.Write(w.buff)
		return nil, err
	}
	return w.buff, nil
}

func (w *JsonEncoder) escape(s string, escapeTable *[128]bool) {

	// Portions of the string that contain no escapes are appended as
	// byte slices.

	p := 0 // last non-escape symbol

	for i := 0; i < len(s); {
		c := s[i]

		if c < utf8.RuneSelf {
			if escapeTable[c] {
				// single-width character, no escaping is required
				i++
				continue
			}

			w.buff = append(w.buff, s[p:i]...)
			switch c {
			case '\t':
				w.buff = append(w.buff, `\t`...)
			case '\r':
				w.buff = append(w.buff, `\r`...)
			case '\n':
				w.buff = append(w.buff, `\n`...)
			case '\\':
				w.buff = append(w.buff, `\\`...)
			case '"':
				w.buff = append(w.buff, `\"`...)
			default:
				w.buff = append(w.buff, '\\', 'u', '0', '0', chars[c>>4], chars[c&0xf])
			}

			i++
			p = i
			continue
		}

		// broken utf
		runeValue, runeWidth := utf8.DecodeRuneInString(s[i:])
		if runeValue == utf8.RuneError && runeWidth == 1 {
			w.buff = append(w.buff, s[p:i]...)
			w.buff = append(w.buff, '\\', 'u', 'f', 'f', 'f', 'd')
			i++
			p = i
			continue
		}

		// jsonp stuff - tab separator and line separator
		if runeValue == '\u2028' || runeValue == '\u2029' {
			w.buff = append(w.buff, s[p:i]...)
			w.buff = append(w.buff, '\\', 'u', '2', '0', '2', chars[runeValue&0xf])
			i += runeWidth
			p = i
			continue
		}
		i += runeWidth
	}
	w.buff = append(w.buff, s[p:]...)
}

func (w *JsonEncoder) base64(in []byte) {
	si := 0
	n := (len(in) / 3) * 3

	for si < n {
		// Convert 3x 8bit source bytes into 4 bytes
		val := uint(in[si+0])<<16 | uint(in[si+1])<<8 | uint(in[si+2])
		w.buff = append(w.buff, encode[val>>18&0x3F], encode[val>>12&0x3F], encode[val>>6&0x3F], encode[val&0x3F])
		si += 3
	}

	remain := len(in) - si
	if remain == 0 {
		return
	}

	// Add the remaining small block
	val := uint(in[si+0]) << 16
	if remain == 2 {
		val |= uint(in[si+1]) << 8
	}

	w.buff = append(w.buff, encode[val>>18&0x3F], encode[val>>12&0x3F])

	switch remain {
	case 2:
		w.buff = append(w.buff, encode[val>>6&0x3F], byte(padChar))
	case 1:
		w.buff = append(w.buff, byte(padChar), byte(padChar))
	}
}

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
	escapeHtmlTable   = table(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '&', '<', '>', '\\')
	noEscapeHtmlTable = table(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '\\')
)

func EncodeAny(w *JsonEncoder, value any) {
	switch value := value.(type) {
	case bool:
		EncodeBool(w, value)
	case int:
		EncodeInt64(w, int64(value))
	case int8:
		EncodeInt32(w, int32(value))
	case int16:
		EncodeInt32(w, int32(value))
	case int32:
		EncodeInt32(w, value)
	case int64:
		EncodeInt64(w, value)
	case uint:
		EncodeUint64(w, uint64(value))
	case uint8:
		EncodeUint32(w, uint32(value))
	case uint16:
		EncodeUint32(w, uint32(value))
	case uint32:
		EncodeUint32(w, value)
	case uint64:
		EncodeUint64(w, value)
	case float32:
		EncodeFloat(w, value)
	case float64:
		EncodeDouble(w, value)
	case string:
		EncodeString(w, value)
	case []byte:
		EncodeBytes(w, value)
	case MessageEncoder:
		value.EncodeJSON(w)
	default:
		bs, err := json.Marshal(value)
		if err != nil {
			if w.firstError == nil {
				w.firstError = err
			}
			return
		} else {
			w.ensure(len(bs))
			w.buff = append(w.buff, bs...)
		}
	}
}

func EncodeAny_OmitEmpty(w *JsonEncoder, name string, value any) {
	if value != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		EncodeAny(w, value)
		w.buff = append(w.buff, comma)
	}
}

func EncodeAny_WithEmpty(w *JsonEncoder, name string, value any) {
	if value != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		EncodeAny(w, value)
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(9 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	}
}

// ToJson Json转换快捷方法
func ToJson(v any) string {
	if enc, ok := v.(MessageEncoder); ok {
		out := NewJsonEncoder(nil, 1024)
		enc.EncodeJSON(out)
		bs, _ := out.Close()
		return UnsafeString(bs)
	}
	bs, _ := json.Marshal(v)
	return UnsafeString(bs)
}
