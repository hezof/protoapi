package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hezof/protoapi"
	"github.com/mailru/easyjson"
	"io"
	"testing"
	"time"
)

func testData() *Student {
	count++
	s := new(Student)
	s.Name = `"\\这是一个小小的测试\\😹"`
	s.Sex = Sex(1)
	s.Age = uint32(count)
	s.Score = 3.14
	s.Class = []*Student{{Sex: 1}}
	s.Index = make(map[string]*Student)
	s.Index["s"] = &Student{Sex: 2}
	s.Data = []byte("this is a test")
	return s
}

var count int

func testStdData() *StdStudent {
	count++
	s := new(StdStudent)
	s.Name = `"test"`
	s.Sex = Sex(1)
	s.Age = uint32(count)
	s.Score = 3.14159
	s.Class = []*Student{new(Student)}
	s.Index = make(map[string]*Student)
	s.Index["s"] = new(Student)
	s.Data = []byte("this is a test")
	return s
}

//var _ Codec = (*Student)(nil)

func _TestEncodeWithBuffer() []byte {
	s := testData()
	w := protoapi.NewJsonEncoder(nil, 1024)
	protoapi.EncodeMessage(w, s)
	err := w.Close()
	if err != nil {
		panic(err)
	}
	return w.Buffer()
}

func _TestEncodeWithWriter() []byte {
	s := testData()
	out := bytes.NewBuffer(make([]byte, 0, 1024))
	w := protoapi.NewJsonEncoder(out, 1024)
	protoapi.EncodeMessage(w, s)
	err := w.Close()
	if err != nil {
		panic(err)
	}
	return out.Bytes()
}

func _TestDecodeWithBuffer(bs []byte) *Student {
	r := protoapi.NewJsonBuffer(bs)
	var s *Student
	protoapi.DecodeMessage(r, &s)
	err := r.Close()
	if err != nil {
		panic(err)
	}
	return s
}

func _TestDecodeWithReader(in io.Reader) *Student {
	r := protoapi.NewJsonDecoder(in, 1024)
	var s *Student
	protoapi.DecodeMessage(r, &s)
	err := r.Close()
	if err != nil {
		panic(err)
	}
	return s
}

func TestIt(t *testing.T) {
	bs := _TestEncodeWithWriter()
	fmt.Println(string(bs))
	s := _TestDecodeWithBuffer(bs)
	fmt.Println(protoapi.ToJson(s))
}

func _TestStdEncode() []byte {
	s := testStdData()
	bs, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return bs
}

func _TestStdDecode(bs []byte) {
	var s Student
	err := json.Unmarshal(bs, &s)
	if err != nil {
		panic(err)
	}
	_ = &s
}

func _TestEasyEncode() {
	s := testData()
	bs, err := easyjson.Marshal(s)
	if err != nil {
		panic(err)
	}
	_ = bs
}

func _TestEasyDecode(bs []byte) {
	var s Student
	err := easyjson.Unmarshal(bs, &s)
	if err != nil {
		panic(err)
	}
	_ = &s
}

var bs = _TestStdEncode()

func BenchmarkEncodeWithBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_TestEncodeWithBuffer()
		_TestDecodeWithBuffer(bs)
	}
}

func BenchmarkEncodeWithWriter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_TestEncodeWithWriter()
		_TestDecodeWithReader(bytes.NewReader(bs))
	}
}

func BenchmarkStdEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_TestStdEncode()
		_TestStdDecode(bs)
	}
}

func BenchmarkEasyEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_TestEasyEncode()
		_TestEasyDecode(bs)
	}
}

const bN = 10000000

func TestEncodeWithBuffer(t *testing.T) {
	start := time.Now()
	for i := 0; i < bN; i++ {
		_TestEncodeWithBuffer()
	}
	used := time.Now().Sub(start).Nanoseconds()
	fmt.Println("TestEncodeWithBuffer: ", used)
}

func TestEncodeWithWriter(t *testing.T) {
	start := time.Now()
	for i := 0; i < bN; i++ {
		_TestEncodeWithWriter()
	}
	used := time.Now().Sub(start).Nanoseconds()
	fmt.Println("TestEncodeWithWriter: ", used)
}

func TestStdEncode(t *testing.T) {
	start := time.Now()
	for i := 0; i < bN; i++ {
		_TestStdEncode()
	}
	used := time.Now().Sub(start).Nanoseconds()
	fmt.Println("TestStdEncode: ", used)
}

func TestEasyEncode(t *testing.T) {
	start := time.Now()
	for i := 0; i < bN; i++ {
		_TestEasyEncode()
	}
	used := time.Now().Sub(start).Nanoseconds()
	fmt.Println("TestEasyEncode: ", used)
}

const errorJson = `{"name":"\"\\\\这是一个小小的测试\\\\😹\"","sex":"FEMALE","age":2,"score":3.14,"class":[{"sex":"FEMALE"}},"index":{"s":{"sex":"UNKNOWN"}},"data":"dGhpcyBpcyBhIHRlc3Q="}`

func TestErrorJson(t *testing.T) {
	s := _TestDecodeWithBuffer([]byte(errorJson))
	fmt.Println(protoapi.ToJson(s))
}

type Demo struct {
}
