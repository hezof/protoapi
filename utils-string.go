package protoapi

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// UnsafeBytes string到[]byte的不安全转换
// For more Details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString []byte到string的不安全转换
// For more Details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
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

func NvlD(vs ...time.Duration) time.Duration {
	for _, v := range vs {
		if v != 0 {
			return v
		}
	}
	return 0
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

// FullMethod is the full RPC method string, i.e., /package.component/method.
func FullMethod(meta *Meta) string {
	return "/" + meta.Package + "." + meta.Service + "/" + meta.Method
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

func as(vs []string) []any {
	ret := make([]any, len(vs))
	for i, v := range vs {
		ret[i] = v
	}
	return ret
}
