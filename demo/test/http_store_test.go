package test

import (
	"bytes"
	"fmt"
	"github.com/hezof/framework"
	"github.com/hezof/protoapi"
	"github.com/hezof/protoapi/demo/api"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"
)

var book = &api.Book{
	Sale:        true,
	Stars:       5,
	Publication: time.Now().Unix(),
	Pages:       32,
	Stock:       64,
	Price:       32,
	Discount:    0.64,
	Name:        "苏格拉底的救赎",
	Cover:       []byte("🤣🤣🤣🤣🤣"),
	Genre:       GENRE,
	GenreName:   GENRE,
	Author: &api.Author{
		No:    1122,
		Name:  "吴天",
		Email: "wutian@xx.com",
	},
	Isbn: "1-3-1-4-159",
}

func TestHttpSimple(t *testing.T) {
	rsp := new(api.Book)
	err := cli.POST("/simple/book", book, protoapi.NormalResult(rsp), 200)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestHttpClient(t *testing.T) {
	rsp := new(api.Book)
	err := cli.POST("/client/book", book, protoapi.UnwrapResult(rsp), 200)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestHttpServer(t *testing.T) {
	rsp := new(api.Book)
	err := cli.POST("/server/book", book, protoapi.UnwrapResult(rsp), 200)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestHttpDirectional(t *testing.T) {
	rsp := new(api.Book)
	err := cli.POST("/bidirectional/book", book, protoapi.UnwrapResult(rsp), 200)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}
func TestHttpUploadCovert(t *testing.T) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	part, err := writer.CreateFormFile("file", "上传文件")
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(`D:\Workspace\klibsrc\protoapi\demo\api\store.proto`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	hrsp, err := http.Post("http://localhost:8080/covert/upload", "multipart/form-data;boundary="+writer.Boundary(), buffer)
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()

	io.Copy(os.Stdout, hrsp.Body)
}
func TestHttpReviewCovert(t *testing.T) {
	hrsp, err := http.Get("http://localhost:8080/review/covert")
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()

	io.Copy(os.Stdout, hrsp.Body)

}
