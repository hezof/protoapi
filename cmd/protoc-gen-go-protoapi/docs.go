package main

import (
	"encoding/json"
	"google.golang.org/protobuf/compiler/protogen"
)

// generateDocsFile 生成文档文件: *_json
func generateDocsFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.json`, file.GoImportPath)
	bs, err := json.MarshalIndent(meta, ``, "\t")
	if err != nil {
		gen.Error(err)
	}
	g.Write(bs)
}

type Openapi struct {
}
