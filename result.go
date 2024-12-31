package protoapi

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

/*
默认code错误码. 兼容http/grpc的code范围是[0,math.MaxInt32]
统一约定:
- [0,99]              表示保留错误码! (业务/扩展切勿占用)
- [100,999]           表示请求错误码! (与http status code一致)
- [1000,9999]         表示系统错误码!
- [10000,2147483647]  表示业务错误码!
*/

type StatusResult struct {
	Status  int32         `json:"-"`                 // 状态代码(http).
	Code    int32         `json:"code"`              // 错误代码. 0表示成功
	Name    string        `json:"name,omitempty"`    // 错误名称. OK表示成功
	Message string        `json:"message,omitempty"` // 错误消息.
	Details []interface{} `json:"-"`                 // 错误参数.
}

type WrapResult struct {
	CodePrefix    string // code前缀, 默认: `"code":`, 0表示成功
	NamePrefix    string // name前缀, 默认: `"name":`, OK表示成功
	DataPrefix    string // data前缀, 默认: `"data":`.
	MessagePrefix string // message前缀, 默认: `"message":`
}

var wrapResult = WrapResult{
	CodePrefix:    `"code":`,
	NamePrefix:    `"name":`,
	DataPrefix:    `"data":`,
	MessagePrefix: `"message":`,
}

func InitWrapResult(wr WrapResult) {
	wrapResult = wr
}

// Sprintf 用于规避标fmt.Sprintf()的异常问题
func Sprintf(format string, args ...interface{}) string {

	// 没有参数
	argc := len(args)
	if argc == 0 {
		return format
	}

	// 扫描参数
	wild := 0 // 通配符%数量(%%除外)
	mark := false
	for _, c := range format {
		if c == '%' {
			if mark {
				mark = false // 上一字符是%
			} else {
				mark = true // 上一字符不是%
			}
		} else {
			if mark { //上一字符是%
				wild++
				mark = false
			}
		}
	}

	// 适配处理
	if wild == 0 {
		return format
	}
	if argc == wild {
		// 参数数量相同
		return fmt.Sprintf(format, args...)
	} else if wild < argc {
		// 参数数量有多
		return fmt.Sprintf(format, args[0:wild]...)
	} else {
		// 参数数量不足
		temp := make([]interface{}, wild)
		copy(temp, args)
		return fmt.Sprintf(format, temp...)
	}
}

// 打印堆栈追踪信息,如果是"/src/runtime/"自动跳过!
func StackTrace(skip int, sep string) string {
	var sb strings.Builder
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			return sb.String()
		}
		// 过滤runtime的行项,避免错误日志过多!
		if strings.Index(file, "/src/runtime/") == -1 {
			if skip > 0 {
				skip--
			} else {
				if sb.Len() > 0 {
					sb.WriteString(sep)
				}
				sb.WriteString(file)
				sb.WriteByte(':')
				sb.WriteString(strconv.Itoa(line))
			}
		}
	}
}
