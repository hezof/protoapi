package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

const (
	Error = "E"
	Warn  = "W"
	Info  = "I"
)

// swagger 2.0 schema: https://swagger.io/specification/v2/
var schema = []string{
	"swagger",
	"info",
	"host",
	"basePath",
	"schemes",
	"consumes",
	"produces",
	"paths",
	"definitions",
	"parameters",
	"responses",
	"securityDefinitions",
	"security",
	"tags",
	"externalDocs",
}

func SysLog(level string, format string, args ...interface{}) {
	fmt.Fprintln(os.Stdout, time.Now().Format("2006/01/02 15:04:05"), "["+level+"]", "-", fmt.Sprintf(format, args...))
}

func SysLoc(path string) string {
	loc, _ := exec.LookPath(os.Args[0])

	if ret := filepath.Join(filepath.Dir(loc), path); Exists(ret) {
		return ret
	}
	cwd, _ := os.Getwd()
	if ret := filepath.Join(cwd, path); Exists(ret) {
		return ret
	}
	return path
}

func Exists(path string) bool {
	stat, err := os.Stat(path)
	return stat != nil || os.IsExist(err)
}

/*
后话:
实现swagger的yaml化时理解不深,完全是靠拼凑字符串的形式实现的!
包括protoapi-swagger针对x-order排序从而达到ui有序的效果!
更优雅应该基于swagger语法抽象AST来做. 这二块功能写得实现很狗屎^_^!
文档的东东凑合就行吧~~~ 有兴趣的同学可以优化提交mr!
*/

// OrderJson 针对"paths.operations", "paths", "definitions", "definitions.*.properties"不作泛解析!
func OrderJson(in io.Reader, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	var swagger map[string]interface{}
	if err = json.NewDecoder(in).Decode(&swagger); err != nil {
		return nil, err
	}

	if _paths, ok := swagger["paths"].(map[string]interface{}); ok {
		pathsEntries := make([]*entry, 0, len(_paths))
		for path, operations := range _paths {
			if _operations, ok := operations.(map[string]interface{}); ok {
				var operationsTag string //  来自首个operation的tags的首个元素
				var operationsOrd int    // 来自首个operation的x-order
				operationsEntries := make([]*entry, 0, len(_operations))
				for method, operation := range _operations {
					_tag := tag(operation)
					_ord := xorder(operation)
					operationsEntries = append(operationsEntries, &entry{
						key: method,
						val: operation,
						tag: _tag,
						ord: _ord,
					})
					if operationsTag == "" { // 取第一个元素
						operationsTag = _tag
						operationsOrd = _ord
					}
				}
				pathsEntries = append(pathsEntries, &entry{
					key: path,
					val: toOrder(operationsEntries),
					tag: operationsTag,
					ord: operationsOrd,
				})
			}

		}
		swagger["paths"] = toOrder(pathsEntries)
	}

	if _definitions, ok := swagger["definitions"].(map[string]interface{}); ok {
		for _, definition := range _definitions {
			if _definition, ok := definition.(map[string]interface{}); ok {
				if _properties, ok := _definition["properties"].(map[string]interface{}); ok {
					propertiesEntries := make([]*entry, 0, len(_properties))
					for name, value := range _properties {
						propertiesEntries = append(propertiesEntries, &entry{
							key: name,
							val: value,
							ord: xorder(value),
						})
					}
					_definition["properties"] = toOrder(propertiesEntries)
				}
			}
		}
	}

	// 按照swagger schema的预定义顺序输出
	var src = new(bytes.Buffer)
	src.WriteByte('{')
	for _, name := range schema {
		if val, ok := swagger[name]; ok {
			if src.Len() > 1 { // 第一个是"{"
				src.WriteByte(',')
			}
			src.Write(toJson(name))
			src.WriteByte(':')
			src.Write(toJson(val))
		}
	}
	src.WriteByte('}')
	var dst = new(bytes.Buffer)
	if err = json.Indent(dst, src.Bytes(), "", "  "); err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}

func toOrder(entries []*entry) json.RawMessage {

	sort.Slice(entries, func(i, j int) bool {
		ti, tj := entries[i].tag, entries[j].tag
		if ti == tj {
			oi, oj := entries[i].ord, entries[j].ord
			return oi < oj
		} else {
			return ti < tj
		}
	})

	var data = new(bytes.Buffer)
	data.WriteByte('{')
	for i, e := range entries {
		if i > 0 {
			data.WriteByte(',')
		}
		data.Write(toJson(e.key))
		data.WriteByte(':')
		data.Write(toJson(e.val))
	}
	data.WriteByte('}')
	return data.Bytes()
}

func xorder(value interface{}) int {
	if m, ok := value.(map[string]interface{}); ok {
		if v, ok := m["x-order"]; ok {
			switch v := v.(type) {
			case int:
				return v
			case int64:
				return int(v)
			case float64:
				return int(math.Floor(v))
			}
		}
	}
	return 0
}

func tag(value interface{}) string {
	if m, ok := value.(map[string]interface{}); ok {
		if vs, ok := m["tags"].([]interface{}); ok && len(vs) > 0 {
			s, _ := vs[0].(string)
			return s
		}
	}
	return ""
}

func toJson(v interface{}) []byte {
	bs, _ := json.Marshal(v)
	return bs
}

type entry struct {
	key string
	val interface{}
	tag string
	ord int
}
