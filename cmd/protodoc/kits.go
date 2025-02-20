package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	Error = "E"
	Warn  = "W"
	Info  = "I"
)

func SysLog(level string, format string, args ...interface{}) {
	fmt.Fprintln(os.Stdout, time.Now().Format("2006/01/02 15:04:05"), "["+level+"]", "-", fmt.Sprintf(format, args...))
}

type XOrderItem struct {
	XOrder int
	Key    any
	Val    any
}

func (X *XOrderItem) Mix(Y *XOrderItem) {
	if X.XOrder == 0 {
		X.XOrder = Y.XOrder
	}
	if X.Key == nil {
		X.Key = Y.Key
	}
	if Y.Val == nil {
		return
	}
	yval, ok := Y.Val.(*XOrderTable)
	if !ok {
		return
	}
	if X.Val == nil {
		val := NewXOrderTable(yval.array)
		for _, item := range yval.items {
			val.Add(item)
		}
		X.Val = val
		return
	}
	xval, ok := X.Val.(*XOrderTable)
	if !ok {
		return
	}
	xval.items = append(xval.items, yval.items...)
}

type XOrderTable struct {
	array bool
	items []*XOrderItem
}

func NewXOrderTable(array bool) *XOrderTable {
	return &XOrderTable{
		array: array,
	}
}

func (X *XOrderTable) Add(item *XOrderItem) {
	X.items = append(X.items, item)
}

func (X *XOrderTable) Get(key any) *XOrderItem {
	for _, item := range X.items {
		if item.Key == key {
			return item
		}
	}
	return nil
}

func (X *XOrderTable) MarshalJSON() ([]byte, error) {

	sort.Sort(X)

	buf := new(bytes.Buffer)
	if X.array {
		buf.WriteByte('[')
		for _, item := range X.items {
			if item == nil {
				continue
			}
			if buf.Len() > 1 {
				buf.WriteByte(',')
			}
			bs, _ := json.Marshal(item.Val)
			buf.Write(bs)
		}
		buf.WriteByte(']')
	} else {
		buf.WriteByte('{')
		for _, item := range X.items {
			if item == nil || item.Key == "x-order" {
				continue
			}
			if buf.Len() > 1 {
				buf.WriteByte(',')
			}
			bs, _ := json.Marshal(item.Key)
			buf.Write(bs)
			buf.WriteByte(':')
			bs, _ = json.Marshal(item.Val)
			buf.Write(bs)
		}
		buf.WriteByte('}')
	}
	return buf.Bytes(), nil
}
func (X *XOrderTable) Len() int {
	return len(X.items)
}

func (X *XOrderTable) Less(i, j int) bool {
	xi := X.items[i]
	xj := X.items[j]

	switch {
	case xi == nil:
		return true
	case xj == nil:
		return false
	default:
		return xi.XOrder < xj.XOrder
	}
}

func (X *XOrderTable) Swap(i, j int) {
	X.items[i], X.items[j] = X.items[j], X.items[i]
}

var _ json.Marshaler = (*XOrderTable)(nil)
var _ sort.Interface = (*XOrderTable)(nil)

type Object = map[string]any
type Array = []any

func ToInt(v any) int {
	switch v := v.(type) {
	case int:
		return v
	case int64:
		return int(v)
	default:
		ret, _ := strconv.Atoi(fmt.Sprint(v))
		return ret
	}
}
