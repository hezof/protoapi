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

func EncodeMessageEncoder(w *JsonEncoder, name string, encoder MessageEncoder) {
	if encoder != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		_ = encoder.EncodeJSON(w)
		if last := len(w.buff) - 1; w.buff[last] == comma {
			w.buff[last] = rightBrace
			w.buff = append(w.buff, comma)
		} else {
			w.buff = append(w.buff, rightBrace, comma)
		}
	}
}

func EncodeMessageMarshal(w *JsonEncoder, name string, data interface{}) {
	if data != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, leftBrace)
		bs, err := json.Marshal(data)
		if err != nil {
			if w.firstError != nil {
				w.firstError = err
			}
			return
		} else {
			w.ensure(len(bs))
			w.buff = append(w.buff, bs...)
		}
		w.buff = append(w.buff, rightBrace, comma)
	}
}
