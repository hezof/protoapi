package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func main() {

	var (
		pBase   string
		pPath   string
		pAddr   string
		pLogo   string
		pTitle  string
		primary string
	)
	flag.StringVar(&pBase, "base", "/", "基础路径")
	flag.StringVar(&pPath, "path", "docs", "文档路径")
	flag.StringVar(&pAddr, "addr", ":8080", "监听地址")
	flag.StringVar(&pLogo, "logo", "", "favicon.ico")
	flag.StringVar(&pTitle, "title", "文档标题", "文档标题")
	flag.StringVar(&primary, "primary", "", "主文档")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println(`Usage: protodoc [options] <yaml_dir|yaml_file> [...]`)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 去重
	files := make(map[string]bool)
	for _, arg := range args {
		info, err := os.Stat(arg)
		if info == nil || os.IsNotExist(err) {
			continue
		}
		if info.IsDir() {
			filepath.Walk(arg, func(arg string, info fs.FileInfo, err error) error {
				if err != nil || info.IsDir() || strings.HasPrefix(info.Name(), ".") || !strings.HasSuffix(info.Name(), ".yaml") {
					return nil
				}
				abs, _ := filepath.Abs(arg)
				files[abs] = true
				return nil
			})
		} else {
			if strings.HasPrefix(info.Name(), ".") || !strings.HasSuffix(info.Name(), ".yaml") {
				continue
			}
			abs, _ := filepath.Abs(arg)
			files[abs] = true
		}
	}

	args = make([]string, 0, len(files))
	for file := range files {
		args = append(args, file)
	}

	obj, err := Mixin(primary, args)
	if err != nil {
		SysLog(Error, "swagger mixin file error: %v", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		SysLog(Error, "swagger marshal file error: %v", err)
		os.Exit(2)
	}

	docServer := new(http.Server)
	docServer.SetKeepAlivesEnabled(true)
	docServer.Addr = pAddr
	docServer.Handler = middleware.Spec(pBase, data, Redoc(middleware.RedocOpts{
		BasePath: pBase,
		Path:     pPath,
		SpecURL:  path.Join(pBase, "swagger.json"),
		RedocURL: path.Join(pBase, "swagger.js"),
		Title:    pTitle,
	}, Other(pLogo)))

	if err = docServer.ListenAndServe(); err != nil {
		SysLog(Error, "swagger serve error: %v", err)
		os.Exit(3)
	}
}

func Other(logo string) http.Handler {
	assets, _ := fs.Sub(embedfs, "assets")
	fileServer := http.FileServer(http.FS(assets))
	if logo == "" {
		return fileServer
	} else {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.URL.Path == "/favicon.ico" {
				http.ServeFile(writer, request, logo)
			} else {
				fileServer.ServeHTTP(writer, request)
			}
		})
	}
}

// Mixin 不落地实现
func Mixin(primaryFile string, mixinFiles []string) (*XOrderTable, error) {

	var primary *XOrderTable

	if primaryFile != "" {
		primary, _ = spec(primaryFile)
	} else {
		primary = NewXOrderTable(false)
	}

	if swagger := primary.Get("swagger"); swagger == nil {
		primary.items = append(primary.items, &XOrderItem{
			XOrder: 0,
			Key:    "swagger",
			Val:    "2.0",
		})
	} else {
		swagger.Val = "2.0"
		swagger.XOrder = 0
	}

	if info := primary.Get("info"); info == nil {
		primary.items = append(primary.items, &XOrderItem{
			XOrder: 1,
			Key:    "info",
			Val: &XOrderTable{
				array: false,
				items: []*XOrderItem{
					{XOrder: 1, Key: "title", Val: ""},
					{XOrder: 2, Key: "version", Val: "v1.0.0"},
				},
			},
		})
	}

	for _, mixinFile := range mixinFiles {
		mixin, err := spec(mixinFile)
		if err != nil {
			return nil, err
		}
		merge(primary, mixin)
	}

	tags := primary.Get("tags")
	pval := tags.Val.(*XOrderTable)
	// 按照名称升序
	sort.SliceStable(pval.items, func(i, j int) bool {
		itemi := pval.items[i].Val.(*XOrderTable).Get("name")
		itemj := pval.items[j].Val.(*XOrderTable).Get("name")

		switch {
		case itemi == nil:
			return true
		case itemj == nil:
			return false
		default:
			return itemi.Val.(string) < itemj.Val.(string)
		}
	})
	for i, v := range pval.items {
		v.XOrder = i + 1
		// 覆盖掉x-order
		if tbl, ok := v.Val.(*XOrderTable); ok {
			if ord := tbl.Get("x-order"); ord != nil {
				ord.Val = v.XOrder
			}
		}
	}
	return primary, nil
}

func merge(primary *XOrderTable, mixin *XOrderTable) {
	/** 合并内容只有3部分
	 * - tags
	 * - paths
	 * - definitions
	 */
	{
		key := "tags"
		mitem := mixin.Get(key)
		if mitem == nil {
			return
		}
		pitem := primary.Get(key)
		if pitem == nil {
			pitem = new(XOrderItem)
			primary.Add(pitem)
		}
		pitem.Mix(mitem)
	}
	{
		key := "paths"
		mitem := mixin.Get(key)
		if mitem == nil {
			return
		}
		pitem := primary.Get(key)
		if pitem == nil {
			pitem = new(XOrderItem)
			primary.Add(pitem)
		}
		pitem.Mix(mitem)
	}
	{
		key := "definitions"
		mitem := mixin.Get(key)
		if mitem == nil {
			return
		}
		pitem := primary.Get(key)
		if pitem == nil {
			pitem = new(XOrderItem)
			primary.Add(pitem)
		}
		pitem.Mix(mitem)
	}
}

func spec(file string) (*XOrderTable, error) {

	object := make(Object)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &object)
	if err != nil {
		return nil, err
	}

	return convObject(object), nil
}

func convObject(object Object) *XOrderTable {
	table := NewXOrderTable(false)
	for k, v := range object {
		item := new(XOrderItem)
		item.Key = k
		switch v := v.(type) {
		case Object:
			item.Val = convObject(v)
			item.XOrder = ToInt(v["x-order"])
		case Array:
			item.Val = convArray(v)
		default:
			item.Val = v
		}
		table.items = append(table.items, item)
	}
	return table
}

func convArray(array Array) *XOrderTable {
	table := NewXOrderTable(true)
	for i, v := range array {
		item := new(XOrderItem)
		item.XOrder = i
		switch v := v.(type) {
		case Object:
			item.Val = convObject(v)
			item.XOrder = ToInt(v["x-order"])
		case Array:
			item.Val = convArray(v)
		default:
			item.Val = v
		}
		table.items = append(table.items, item)
	}
	return table
}

func Redoc(opts middleware.RedocOpts, next http.Handler) http.Handler {
	opts.EnsureDefaults()

	pth := path.Join(opts.BasePath, opts.Path)
	tmpl := template.Must(template.New("redoc").Parse(redocTemplate))

	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, opts)
	b := buf.Bytes()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pth {
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			rw.WriteHeader(http.StatusOK)

			_, _ = rw.Write(b)
			return
		}

		if next == nil {
			rw.Header().Set("Content-Type", "text/plain")
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write([]byte(fmt.Sprintf("%q not found", pth)))
			return
		}
		next.ServeHTTP(rw, r)
	})
}

const redocTemplate = `<!DOCTYPE html>
<html>
 <head>
   <link rel="shortcut icon" type="image/x-icon" href="/favicon.ico?" />
   <title>{{ .Title }}</title>
   <!-- needed for adaptive design -->
   <meta charset="utf-8"/>
   <meta name="viewport" content="width=device-width, initial-scale=1"/>
   <!--
   ReDoc doesn't change outer page styles
   -->
   <style>
     body {
       margin: 0;
       padding: 0;
     }
   </style>
 </head>
 <body>
   <redoc spec-url="{{ .SpecURL }}"></redoc>
   <script src="{{ .RedocURL }}"> </script>
 </body>
</html>
`
