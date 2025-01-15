package main

import (
	"encoding/json"
	"google.golang.org/protobuf/compiler/protogen"
)

// generateCodeFile 生成实现文件: *_protoapi.pb.go
func generateImplFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.pb.go`, file.GoImportPath)

	bs, err := json.MarshalIndent(meta, ``, `	`)
	if err != nil {
		gen.Error(err)
	}
	_, err = g.Write(bs)
	if err != nil {
		gen.Error(err)
	}
}
