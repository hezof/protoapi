package protoapi

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

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

// UnsafeBytes string到[]byte的不安全转换
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString []byte到string的不安全转换
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
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

// StackTrace 打印堆栈追踪信息,如果是"/src/runtime/"自动跳过!
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

func NvlI[I int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](vs ...I) I {
	for _, v := range vs {
		if v != 0 {
			return v
		}
	}
	return 0
}

func NvlS(vs ...string) string {
	for _, v := range vs {
		if v != `` {
			return v
		}
	}
	return ``
}

func NvlB(vs ...bool) bool {
	for _, v := range vs {
		if v {
			return v
		}
	}
	return false
}

func Join(cs ...[]HandleFunc) []HandleFunc {
	sum := 0
	for _, c := range cs {
		sum += len(c)
	}
	if sum == 0 {
		return nil
	}
	ct := make([]HandleFunc, 0, sum)
	for _, c := range cs {
		ct = append(ct, c...)
	}
	return ct
}

// LookPath 查找路径资源, 规则如下:
// 1. 若"路径"是绝对路径,返回"路径"及"是否存在".
// 2. 若"路径"在启动目录存在,返回"启动目录+路径"及true.
// 3. 若"路径"在当前目录存在,返回"当前目录+路径"及true.
// 4. 上述步骤失败, 返回"路径"及false.
func LookPath(path string) (string, bool) {

	// 如果是绝对路径,直接返回
	if filepath.IsAbs(path) {
		fi, err := os.Stat(path)
		return path, fi != nil || os.IsExist(err)
	}

	loc, _ := exec.LookPath(os.Args[0])
	res := filepath.Join(filepath.Dir(loc), path)
	if fi, err := os.Stat(res); fi == nil || os.IsNotExist(err) {
		dir, _ := os.Getwd()
		res = filepath.Join(dir, path)
		if fi, err = os.Stat(res); fi == nil || os.IsNotExist(err) {
			return path, false
		}
	}
	return res, true
}

// fullMethod is the full RPC method string, i.e., /package.service/method.
func fullMethod(packageName, serviceName, methodName string) string {
	return "/" + packageName + "." + serviceName + "/" + methodName
}

func orderServiceAspects(v1 []ServiceAspect, v2 []ServiceAspect) []ServiceAspect {
	n1, n2 := len(v1), len(v2)
	if n1 == 0 && n2 == 0 {
		return nil
	}

	vs := make([]ServiceAspect, n1+n2)
	copy(vs, v1)
	copy(vs[n1:], v2)
	// 根据Order[0]与Order[1]排序
	sort.SliceStable(vs, func(i, j int) bool {
		ai := vs[i]
		aj := vs[j]
		if ai.Order()[0] > aj.Order()[0] {
			return false
		} else if ai.Order()[0] < aj.Order()[0] {
			return true
		} else {
			if ai.Order()[1] > aj.Order()[1] {
				return false
			} else {
				return true
			}
		}
	})
	return vs
}
