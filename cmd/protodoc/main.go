package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"gopkg.in/yaml.v2"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	flag.StringVar(&pPath, "arg", "docs", "文档路径")
	flag.StringVar(&pAddr, "addr", ":8080", "监听地址")
	flag.StringVar(&pLogo, "logo", "", "favicon.ico")
	flag.StringVar(&pTitle, "title", "文档标题", "文档标题")
	flag.StringVar(&primary, "primary", "", "文档模板")

	flag.Parse()
	args := flag.Args()

	// 去重
	files := make(map[string]bool)
	for _, arg := range args {
		info, err := os.Stat(arg)
		if info == nil || os.IsNotExist(err) {
			continue
		}
		if info.IsDir() {
			filepath.Walk(arg, func(path string, info fs.FileInfo, err error) error {
				if err != nil || info.IsDir() || strings.HasPrefix(info.Name(), ".") || !strings.HasSuffix(info.Name(), ".yaml") {
					return nil
				}
				abs, _ := filepath.Abs(path)
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
	for file, _ := range files {
		args = append(args, file)
	}
	// 默认第一个作为primary
	if primary == "" {
		primary = args[0]
		args = args[1:]
	}

	data, err := OrderJson(Mixin(primary, args))
	if err != nil {
		SysLog(Error, "swagger marshal json error: %v", err)
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
func Mixin(primaryFile string, mixinFiles []string) (io.Reader, error) {

	primaryDoc, err := loads.Spec(primaryFile)
	if err != nil {
		return nil, err
	}
	primary := primaryDoc.Spec()

	var mixins []*spec.Swagger
	for _, mixinFile := range mixinFiles {
		mixinFile = WithAutoXOrder(mixinFile)
		mixin, lerr := loads.Spec(mixinFile)
		if lerr != nil {
			return nil, lerr
		}
		mixins = append(mixins, mixin.Spec())
	}

	_ = analysis.Mixin(primary, mixins...)
	analysis.FixEmptyResponseDescriptions(primary)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(primary)

	return buf, err
}

// Extensions supported by go-swagger
const (
	xOrder = "x-order"
	object = "object"
)

// WithAutoXOrder amends the spec to specify property order as they appear
// in the spec (supports yaml documents only).
func WithAutoXOrder(specPath string) string {
	lookFor := func(ele interface{}, key string) (yaml.MapSlice, bool) {
		if slice, ok := ele.(yaml.MapSlice); ok {
			for _, v := range slice {
				if v.Key == key {
					if slice, ok := v.Value.(yaml.MapSlice); ok {
						return slice, ok
					}
				}
			}
		}
		return nil, false
	}

	var addXOrder func(interface{})
	addXOrder = func(element interface{}) {
		if props, ok := lookFor(element, "properties"); ok {
			for i, prop := range props {
				if pSlice, ok := prop.Value.(yaml.MapSlice); ok {
					isObject := false
					xOrderIndex := -1 // find if x-order already exists

					for i, v := range pSlice {
						if v.Key == "type" && v.Value == object {
							isObject = true
						}
						if v.Key == xOrder {
							xOrderIndex = i
							break
						}
					}

					if xOrderIndex > -1 { // override existing x-order
						pSlice[xOrderIndex] = yaml.MapItem{Key: xOrder, Value: i}
					} else { // append new x-order
						pSlice = append(pSlice, yaml.MapItem{Key: xOrder, Value: i})
					}
					prop.Value = pSlice
					props[i] = prop

					if isObject {
						addXOrder(pSlice)
					}
				}
			}
		}
	}

	yamlDoc, err := swag.YAMLData(specPath)
	if err != nil {
		panic(err)
	}

	if defs, ok := lookFor(yamlDoc, "definitions"); ok {
		for _, def := range defs {
			addXOrder(def.Value)
		}
	}

	addXOrder(yamlDoc)

	out, err := yaml.Marshal(yamlDoc)
	if err != nil {
		panic(err)
	}

	tmpDir, err := os.MkdirTemp("", "go-swagger-")
	if err != nil {
		panic(err)
	}

	tmpFile := filepath.Join(tmpDir, filepath.Base(specPath))
	if err := os.WriteFile(tmpFile, out, 0600); err != nil {
		panic(err)
	}
	return tmpFile
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
