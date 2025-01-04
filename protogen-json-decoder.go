package protoapi

import (
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func NewJsonDecoder(in io.Reader, size int) *JsonDecoder {
	r := &JsonDecoder{
		in:   in,
		buff: make([]byte, size),
		mark: 0,
		size: 0,
	}
	// NOTE: 务必初始首个token
	r.next()
	return r
}

func NewJsonBuffer(buf []byte) *JsonDecoder {
	r := &JsonDecoder{
		in:   nil,
		buff: buf,
		mark: 0,
		size: len(buf),
	}
	// NOTE: 务必初始首个token
	r.next()
	return r
}

type JsonDecoder struct {
	token      JsonToken // token类型
	in         io.Reader // 读入流
	buff       []byte    // 缓存区
	mark       int       // 读位置
	size       int       // 末位置
	base       int       // 基位置(base + pos)
	number     []byte    // 数值缓存区
	depth      int       // 嵌套深度
	firstError error     // 上下文错误
}

func (r *JsonDecoder) Decode(v JsonDecodc) {
	v.DecodeJSON(r)
}

// Close 关闭读流
func (r *JsonDecoder) Close() error {
	// 解析结束,不能再有其他字符
	if r.next() != 0 {
		r.invalidCharacterError()
	}
	return r.firstError
}

func (r *JsonDecoder) expectedTokenError(t JsonToken) {
	if r.firstError == nil {
		r.token = -1
		r.firstError = newParseError(r, r.mark-1, fmt.Sprintf("expected %s, but got '%c'", tokens[t], r.buff[r.mark-1]))
	}
}

func (r *JsonDecoder) expectedCharacterError(c byte) {
	if r.firstError == nil {
		r.token = -1
		r.firstError = newParseError(r, r.mark-1, fmt.Sprintf("expected '%c', but got '%c'", c, r.buff[r.mark-1]))
	}
}

func (r *JsonDecoder) invalidCharacterError2(c byte, m int) {
	if r.firstError == nil {
		r.token = -1
		if c < '\u0020' {
			r.firstError = newParseError(r, r.mark+m, fmt.Sprintf("invalid character '\\x%x'", c))
		} else {
			r.firstError = newParseError(r, r.mark+m, fmt.Sprintf("invalid character '%c'", c))
		}
	}
}

func (r *JsonDecoder) invalidCharacterError() {
	if r.firstError == nil {
		r.token = -1
		r.firstError = newParseError(r, r.mark-1, fmt.Sprintf("invalid character '%c'", r.buff[r.mark-1]))
	}
}

func (r *JsonDecoder) unexpectedEndError() {
	if r.firstError == nil {
		r.token = -1
		r.firstError = newParseError(r, r.mark-1, "unexpected end of JSON input")
	}
}

func (r *JsonDecoder) exceedMaximumNestingDepthError() {
	if r.firstError == nil {
		r.token = -1
		r.firstError = newParseError(r, r.mark-1, "exceed maximum depth of nesting")
	}
}

func (r *JsonDecoder) NextNull() bool {
	if r.token != 0 {
		return r.token == Null
	}
	return r.next() == Null
}

// more 从读入流读取缓存区
func (r *JsonDecoder) more() bool {
	if r.in == nil {
		return false
	}
	r.base += r.size
	r.mark = 0
	r.size = 0
	for r.size < cap(r.buff) {
		n, err := r.in.Read(r.buff[r.size:])
		r.size += n
		if err != nil {
			return r.size > 0
		}
	}
	return true
}

func (r *JsonDecoder) next() JsonToken {
	if r.token != 0 {
		switch r.token {
		case ObjectBegin, ObjectEnd, ArrayBegin, ArrayEnd, Comma, Colon:
			r.token = 0
		case String:
			r.skipString()
		case Number:
			r.skipNumber()
		case True:
			r.skipTrue()
		case False:
			r.skipFalse()
		case Null:
			r.skipNull()
		case -1:
			return -1
		}
	}

	for {
		for i, c := range r.buff[r.mark:r.size] {
			switch c {
			case '\u0020', '\u000A', '\u000D', '\u0009':
			case '{':
				if r.depth++; r.depth > profile.MaximumNestingDepth {
					r.exceedMaximumNestingDepthError()
					r.depth--
					return 0
				}
				r.mark += i + 1
				r.token = ObjectBegin
				return ObjectBegin
			case '}':
				r.depth--
				r.mark += i + 1
				r.token = ObjectEnd
				return ObjectEnd
			case '[':
				if r.depth++; r.depth > profile.MaximumNestingDepth {
					r.exceedMaximumNestingDepthError()
					r.depth--
					return 0
				}
				r.mark += i + 1
				r.token = ArrayBegin
				return ArrayBegin
			case ']':
				r.depth--
				r.mark += i + 1
				r.token = ArrayEnd
				return ArrayEnd
			case ',':
				r.mark += i + 1
				r.token = Comma
				return Comma
			case ':':
				r.mark += i + 1
				r.token = Colon
				return Colon
			case '"':
				r.mark += i + 1
				r.token = String
				return String
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
				r.mark += i + 1
				r.number = append(r.number[0:0], c)
				r.token = Number
				return Number
			case 't':
				r.mark += i + 1
				r.token = True
				return True
			case 'f':
				r.mark += i + 1
				r.token = False
				return False
			case 'n':
				r.mark += i + 1
				r.token = Null
				return Null
			default:
				r.mark += i + 1
				r.invalidCharacterError2(c, i)
				return 0
			}
		}
		if !r.more() {
			return 0
		}
	}
}

func (r *JsonDecoder) skipObject() {
	r.token = 0

	token := r.next()
	for token != ObjectEnd {
		if token != String {
			r.expectedCharacterError('"')
			return
		}
		token = r.next()
		if token != Colon {
			r.expectedCharacterError(':')
			return
		}
		token = r.next()
		switch token {
		case 0:
			r.unexpectedEndError()
			return
		case ObjectBegin:
			r.skipObject()
		case ObjectEnd:
			r.invalidCharacterError()
			return
		case ArrayBegin:
			r.skipArray()
		case ArrayEnd:
			r.invalidCharacterError()
			return
		case Comma:
			r.invalidCharacterError()
			return
		case Colon:
			r.invalidCharacterError()
			return
		case String:
			r.skipString()
		case Number:
			r.skipNumber()
		case True:
			r.skipTrue()
		case False:
			r.skipFalse()
		case Null:
			r.skipNull()
		}
		token = r.next()
		if token == Comma {
			token = r.next()
			if token == ObjectEnd { // ",}"是错误格式
				r.invalidCharacterError()
				return
			}
		}
	}
	return
}

func (r *JsonDecoder) skipArray() {
	r.token = 0

	token := r.next()
	for token != ArrayEnd {
		switch token {
		case 0:
			r.unexpectedEndError()
			return
		case ObjectBegin:
			r.skipObject()
		case ObjectEnd:
			r.invalidCharacterError()
			return
		case ArrayBegin:
			r.skipArray()
		case ArrayEnd:
			r.invalidCharacterError()
			return
		case Comma:
			r.invalidCharacterError()
			return
		case Colon:
			r.invalidCharacterError()
			return
		case String:
			r.skipString()
		case Number:
			r.skipNumber()
		case True:
			r.skipTrue()
		case False:
			r.skipFalse()
		case Null:
			r.skipNull()
		}
		token = r.next()
		if token == Comma {
			token = r.next()
			if token == ArrayEnd { // ",]"是错误格式
				r.invalidCharacterError()
				return
			}
		}
	}
	return
}

func (r *JsonDecoder) readInt64() int64 {
	r.readNumber()
	ret, err := strconv.ParseInt(UnsafeString(r.number), 10, 64)
	if err != nil && r.firstError == nil {
		r.firstError = err
	}
	return ret
}

func (r *JsonDecoder) readUint64() uint64 {
	r.readNumber()
	ret, err := strconv.ParseUint(UnsafeString(r.number), 10, 64)
	if err != nil && r.firstError == nil {
		r.firstError = err
	}
	return ret
}

func (r *JsonDecoder) readFloat64() float64 {
	r.readNumber()
	ret, err := strconv.ParseFloat(UnsafeString(r.number), 64)
	if err != nil && r.firstError == nil {
		r.firstError = err
	}
	return ret
}

func (r *JsonDecoder) readByte() byte {
	if r.mark < r.size || r.more() {
		c := r.buff[r.mark]
		r.mark++
		return c
	}
	return 0
}

func (r *JsonDecoder) unreadByte() {
	r.mark--
}

func (r *JsonDecoder) readU4() rune {
	var ret rune
	for i := 0; i < 4; i++ {
		c := r.readByte()
		switch {
		case c >= '0' && c <= '9':
			ret = ret*16 + rune(c-'0')
		case c >= 'a' && c <= 'f':
			ret = ret*16 + rune(c-'a'+10)
		case c >= 'A' && c <= 'F':
			ret = ret*16 + rune(c-'A'+10)
		default:
			if c != 0 {
				r.invalidCharacterError()
			} else {
				r.unexpectedEndError()
			}
			return 0
		}
	}
	return ret
}

func (r *JsonDecoder) appendEscape(str []byte, c byte) []byte {
	/*
		https://www.json.org/json-en.html
			escape
			'"'
			'\'
			'/'
			'b'
			'f'
			'n'
			'r'
			't'
			'u' hex hex hex hex
	*/
	switch c {
	case '"':
		str = append(str, '"')
	case '\\':
		str = append(str, '\\')
	case '/':
		str = append(str, '/')
	case 'b':
		str = append(str, '\b')
	case 'f':
		str = append(str, '\f')
	case 'n':
		str = append(str, '\n')
	case 'r':
		str = append(str, '\r')
	case 't':
		str = append(str, '\t')
	case 'u':
		u1 := r.readU4()
		if utf16.IsSurrogate(u1) {
			c = r.readByte()
			if c == 0 {
				return utf8.AppendRune(str, u1)
			}
			if c != '\\' {
				r.unreadByte()
				return utf8.AppendRune(str, u1)
			}
			c = r.readByte()
			if c == 0 {
				return utf8.AppendRune(str, u1)
			}
			if c != 'u' {
				str = utf8.AppendRune(str, u1)
				return r.appendEscape(str, c)
			}
			u2 := r.readU4()
			combined := utf16.DecodeRune(u1, u2)
			if combined == unicode.ReplacementChar {
				str = utf8.AppendRune(str, u1)
				str = utf8.AppendRune(str, u2)
			} else {
				str = utf8.AppendRune(str, combined)
			}
		} else {
			str = utf8.AppendRune(str, u1)
		}
	default:
		if c != 0 {
			r.invalidCharacterError()
		} else {
			r.unexpectedEndError()
		}
	}
	return str
}

func (r *JsonDecoder) skipString() {
	r.token = 0

	isEscape := false
	for {
		for i, c := range r.buff[r.mark:r.size] {
			switch {
			case c < '\u0020':
				r.invalidCharacterError2(c, i) // 记录错误后继续
			case c == '"' && !isEscape:
				r.mark += i + 1 // 跳过当前(")
				c = r.readByte()
				if !isTokenEnd(c) {
					r.invalidCharacterError2(c, i)
				}
				r.unreadByte()
				return
			case c == '\\' && !isEscape:
				isEscape = true
			default:
				isEscape = false
			}
		}
		if !r.more() {
			r.unexpectedEndError()
			return
		}
	}
}

func (r *JsonDecoder) readString() string {
	r.token = 0

	var buf []byte
__ENTRY__:
	for {
		for i, c := range r.buff[r.mark:r.size] {
			switch {
			case c < '\u0020':
				r.invalidCharacterError2(c, i) // 记录错误后继续
			case c == '"':
				mark := r.mark
				r.mark += i
				buf = append(buf, r.buff[mark:r.mark]...)
				r.mark++ // 跳过当前(")
				c = r.readByte()
				if !isTokenEnd(c) {
					r.invalidCharacterError2(c, i)
				}
				r.unreadByte()
				return UnsafeString(buf)
			case c == '\\':
				mark := r.mark
				r.mark += i
				buf = append(buf, r.buff[mark:r.mark]...)
				r.mark++ // 跳过当前(/)
				buf = r.appendEscape(buf, r.readByte())
				continue __ENTRY__
			}
		}
		buf = append(buf, r.buff[r.mark:r.size]...)
		if !r.more() {
			r.unexpectedEndError()
			return UnsafeString(buf)
		}
	}
}

func (r *JsonDecoder) skipNumber() {
	r.token = 0

	hasE := false
	afterE := false
	hasDot := false
	for {
		for i, c := range r.buff[r.mark:r.size] {
			switch {
			case c >= '0' && c <= '9':
				afterE = false
			case c == '.' && !hasDot:
				hasDot = true
			case (c == 'e' || c == 'E') && !hasE:
				hasE = true
				hasDot = true
				afterE = true
			case (c == '+' || c == '-') && afterE:
				afterE = false
			default:
				r.mark += i
				if !isTokenEnd(c) {
					r.invalidCharacterError2(c, i)
				}
				r.unreadByte()
				return
			}
		}
		if !r.more() {
			return
		}
	}
}

func (r *JsonDecoder) readNumber() {
	r.token = 0

	hasE := false
	afterE := false
	hasDot := false
	for {
		for i, c := range r.buff[r.mark:r.size] {
			switch {
			case c >= '0' && c <= '9':
				afterE = false
			case c == '.' && !hasDot:
				hasDot = true
			case (c == 'e' || c == 'E') && !hasE:
				hasE = true
				hasDot = true
				afterE = true
			case (c == '+' || c == '-') && afterE:
				afterE = false
			default:
				mark := r.mark
				r.mark += i
				r.number = append(r.number, r.buff[mark:r.mark]...)
				if !isTokenEnd(c) {
					r.invalidCharacterError2(c, i)
				}
				return
			}
		}
		r.number = append(r.number, r.buff[r.mark:r.size]...)
		if !r.more() {
			return
		}
	}
}

func (r *JsonDecoder) skipTrue() {
	r.token = 0
	// 跳过't'
	if r.readByte() != 'r' {
		r.expectedCharacterError('r')
	}
	if r.readByte() != 'u' {
		r.expectedCharacterError('u')
	}
	if r.readByte() != 'e' {
		r.expectedCharacterError('e')
	}
	if !isTokenEnd(r.readByte()) {
		r.invalidCharacterError()
	}
	r.unreadByte()
}

func (r *JsonDecoder) skipFalse() {
	r.token = 0
	//  跳过'f'
	if r.readByte() != 'a' {
		r.expectedCharacterError('a')
	}
	if r.readByte() != 'l' {
		r.expectedCharacterError('l')
	}
	if r.readByte() != 's' {
		r.expectedCharacterError('s')
	}
	if r.readByte() != 'e' {
		r.expectedCharacterError('e')
	}
	if !isTokenEnd(r.readByte()) {
		r.invalidCharacterError()
	}
	r.unreadByte()
}

func (r *JsonDecoder) skipNull() {
	r.token = 0
	// 跳过'n'
	if r.readByte() != 'u' {
		r.expectedCharacterError('u')
	}
	if r.readByte() != 'l' {
		r.expectedCharacterError('l')
	}
	if r.readByte() != 'l' {
		r.expectedCharacterError('l')
	}
	if !isTokenEnd(r.readByte()) {
		r.invalidCharacterError()
	}
	r.unreadByte()
}

// JsonToken 词汇令牌. 0表示EOF, -1表示ERROR
type JsonToken int8

const (
	ObjectBegin JsonToken = 1
	ObjectEnd   JsonToken = 2
	ArrayBegin  JsonToken = 3
	ArrayEnd    JsonToken = 4
	Comma       JsonToken = 5
	Colon       JsonToken = 6
	String      JsonToken = 7
	Number      JsonToken = 8
	True        JsonToken = 9
	False       JsonToken = 10
	Null        JsonToken = 11
)

var tokens = map[JsonToken]string{
	ObjectBegin: `<object>`,
	ObjectEnd:   `<object>`,
	ArrayBegin:  `<array>`,
	ArrayEnd:    `<array>`,
	Comma:       `','`,
	Colon:       `':'`,
	String:      `<string>`,
	Number:      `<number>`,
	True:        `<bool>`,
	False:       `<bool>`,
	Null:        `<null>`,
}

func isTokenEnd(c byte) bool {
	return c == 0 || c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '[' || c == ']' || c == '{' || c == '}' || c == ',' || c == ':'
}

// ParseError Json解码错误
type ParseError struct {
	Reason string
	Offset int
	Data   string
}

func (l *ParseError) Error() string {
	return fmt.Sprintf("%s near offset %d of '%s'", l.Reason, l.Offset, l.Data)
}

func newParseError(ctx *JsonDecoder, mark int, reason string) error {
	var data string
	if ctx.size-mark <= profile.MaximumErrorLength {
		data = string(ctx.buff[mark:ctx.size]) + "..."
	} else {
		data = string(ctx.buff[mark:mark+profile.MaximumErrorLength]) + "..."
	}
	return &ParseError{
		Reason: reason,
		Offset: ctx.base + mark,
		Data:   data,
	}
}
