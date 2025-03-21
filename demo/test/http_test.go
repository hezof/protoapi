package test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hezof/framework"
	"github.com/hezof/protoapi/demo/api"
	"math"
	"os"
	"strconv"
	"testing"
)

var cli = base.NewJsonRpcClient("http://localhost:8080", &base.HttpConfig{Debug: os.Stdout}, nil, nil, nil)

var (
	TRUE             = true
	FALSE            = false
	INT32  int32     = 32
	INT64  int64     = 64
	UINT32 uint32    = 32
	UINT64 uint64    = 64
	FLOAT  float32   = 123.4
	DOUBLE float64   = 123.4
	STRING string    = "🤣🤣🤣🤣🤣"
	BYTES  string    = base64.StdEncoding.EncodeToString([]byte("🤣🤣🤣🤣🤣"))
	GENRE  api.Genre = api.Genre_MAGAZINE
)

func TestBytes(t *testing.T) {
	//fmt.Println(BYTES)
	//bs, err := base64.StdEncoding.DecodeString(BYTES)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(bs))
	m := map[string][]byte{
		"bytes": []byte("🤣🤣🤣🤣🤣"),
	}
	d, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(d))
}

func TestStatus(t *testing.T) {
	fmt.Println(math.MaxInt32)
	var v = "3380709036"
	n, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(n)
}
