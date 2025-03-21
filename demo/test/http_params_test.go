package test

import (
	"bytes"
	"fmt"
	"github.com/hezof/core"
	"github.com/hezof/protoapi"
	"github.com/hezof/protoapi/demo/api"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

var req = url.Values{
	"bool":       []string{"true"},
	"int32":      []string{"32"},
	"int64":      []string{"64"},
	"uint32":     []string{"32"},
	"uint64":     []string{"64"},
	"float":      []string{"32.0"},
	"double":     []string{"64.0"},
	"string":     []string{"string"},
	"bytes":      []string{BYTES},
	"genre":      []string{"2"},
	"genre_name": []string{"MAGAZINE"},

	"bool_optional":       []string{"true"},
	"int32_optional":      []string{"32"},
	"int64_optional":      []string{"64"},
	"uint32_optional":     []string{"32"},
	"uint64_optional":     []string{"64"},
	"float_optional":      []string{"32.0"},
	"double_optional":     []string{"64.0"},
	"string_optional":     []string{"string"},
	"bytes_optional":      []string{BYTES},
	"genre_optional":      []string{"2"},
	"genre_name_optional": []string{"MAGAZINE"},

	"bool_repeated":       []string{"true"},
	"int32_repeated":      []string{"32"},
	"int64_repeated":      []string{"64"},
	"uint32_repeated":     []string{"32"},
	"uint64_repeated":     []string{"64"},
	"float_repeated":      []string{"32.0"},
	"double_repeated":     []string{"64.0"},
	"string_repeated":     []string{"string"},
	"bytes_repeated":      []string{BYTES},
	"genre_repeated":      []string{"2"},
	"genre_name_repeated": []string{"MAGAZINE"},

	"bool_map":       []string{"key,true"},     // explode=false
	"int32_map":      []string{"key,32"},       // explode=false
	"int64_map":      []string{"key,64"},       // explode=false
	"uint32_map":     []string{"key,32"},       // explode=false
	"uint64_map":     []string{"key,64"},       // explode=false
	"float_map":      []string{"key,32.0"},     // explode=false
	"double_map":     []string{"key,64.0"},     // explode=false
	"string_map":     []string{"key,string"},   // explode=false
	"bytes_map":      []string{"key," + BYTES}, // explode=false
	"genre_map":      []string{"key,2"},        // explode=false
	"genre_name_map": []string{"key,MAGAZINE"}, // explode=false
}

func TestForm(t *testing.T) {

	hrsp, err := http.PostForm("http://localhost:8080/params/form", req)
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()

	data := new(bytes.Buffer)
	io.Copy(data, hrsp.Body)

	fmt.Fprintln(os.Stdout, hrsp.Status)
	fmt.Fprintln(os.Stdout, data.String())
	fmt.Fprintln(os.Stdout)

	if hrsp.StatusCode != 200 {
		t.Fatal(err)
	}
	rsp := new(api.ParamsInForm)
	err = protoapi.DecodeRequest(data, core.NormalResult(rsp))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestPath(t *testing.T) {

	hrsp, err := http.PostForm("http://localhost:8080/params/path/key,2/key,MAGAZINE", req)
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()

	data := new(bytes.Buffer)
	io.Copy(data, hrsp.Body)

	fmt.Fprintln(os.Stdout, hrsp.Status)
	fmt.Fprintln(os.Stdout, data.String())
	fmt.Fprintln(os.Stdout)

	if hrsp.StatusCode != 200 {
		t.Fatal(err)
	}
	rsp := new(api.ParamsInForm)
	err = protoapi.DecodeRequest(data, core.NormalResult(rsp))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestQuery(t *testing.T) {
	buf := new(bytes.Buffer)
	for k, vs := range req {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(strings.Join(vs, ",")))
	}
	hrsp, err := http.Post("http://localhost:8080/params/query?"+buf.String(), "multipart/form-data", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()

	data := new(bytes.Buffer)
	io.Copy(data, hrsp.Body)

	fmt.Fprintln(os.Stdout, hrsp.Status)
	fmt.Fprintln(os.Stdout, data.String())
	fmt.Fprintln(os.Stdout)

	if hrsp.StatusCode != 200 {
		t.Fatal(err)
	}
	rsp := new(api.ParamsInForm)
	err = protoapi.DecodeRequest(data, core.NormalResult(rsp))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestHeader(t *testing.T) {

	hreq, err := http.NewRequest(http.MethodPost, "http://localhost:8080/params/header", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range req {
		hreq.Header[k] = v
	}

	hrsp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()

	data := new(bytes.Buffer)
	io.Copy(data, hrsp.Body)

	fmt.Fprintln(os.Stdout, hrsp.Status)
	fmt.Fprintln(os.Stdout, data.String())
	fmt.Fprintln(os.Stdout)

	if hrsp.StatusCode != 200 {
		t.Fatal(err)
	}
	rsp := new(api.ParamsInForm)
	err = protoapi.DecodeRequest(data, core.NormalResult(rsp))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestCookie(t *testing.T) {
	hreq, err := http.NewRequest(http.MethodPost, "http://localhost:8080/params/cookie", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	for k, vs := range req {
		hreq.AddCookie(&http.Cookie{
			Name:  k,
			Value: strings.Join(vs, ","),
		})
	}

	hrsp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()

	data := new(bytes.Buffer)
	io.Copy(data, hrsp.Body)

	fmt.Fprintln(os.Stdout, hrsp.Status)
	fmt.Fprintln(os.Stdout, data.String())
	fmt.Fprintln(os.Stdout)

	if hrsp.StatusCode != 200 {
		t.Fatal(err)
	}
	rsp := new(api.ParamsInForm)
	err = protoapi.DecodeRequest(data, core.NormalResult(rsp))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(core.ToJson(rsp))
}
