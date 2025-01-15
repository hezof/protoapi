package main

import (
	"encoding/json"
	"google.golang.org/protobuf/compiler/protogen"
)

// generateCodeFile 生成代码文件: *_protoapi.code
func generateCodeFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.code`, file.GoImportPath)

	bs, err := json.MarshalIndent(meta, ``, `	`)
	if err != nil {
		gen.Error(err)
	}
	_, err = g.Write(bs)
	if err != nil {
		gen.Error(err)
	}
}
