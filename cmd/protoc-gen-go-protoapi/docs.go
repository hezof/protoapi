package main

import (
	"encoding/json"
	"google.golang.org/protobuf/compiler/protogen"
)

// generateDocsFile 生成文档文件: *_json
func generateDocsFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.yaml`, file.GoImportPath)
	bs, err := json.MarshalIndent(meta, ``, "\t")
	if err != nil {
		gen.Error(err)
	}
	g.Write(bs)
}

type OASv2 struct {
	Swagger     string                    `yaml:"swagger,omitempty"`
	Tags        Vec[*OASv2Tag]            `json:"tags,omitempty"`
	Paths       Map[Map[*OASv2Operation]] `yaml:"paths,omitempty"`
	Definitions Vec[*OASv2Definition]     `yaml:"definitions,omitempty"`
}

type OASv2Tag struct {
	Name         string             `yaml:"name,omitempty"`
	Description  string             `yaml:"description,omitempty"`
	ExternalDocs *OASv2ExternalDocs `yaml:"externalDocs,omitempty"`
}

type OASv2ExternalDocs struct {
	Description string `yaml:"description,omitempty"`
	Url         string `yaml:"url,omitempty"`
}

type OASv2Operation struct {
	XOrder       int                  `yaml:"x-order"`
	OperationId  string               `yaml:"operationId,omitempty"`
	Summary      string               `yaml:"summary,omitempty"`
	Description  string               `yaml:"description,omitempty"`
	Tags         Vec[string]          `yaml:"tags,omitempty"`
	Parameters   Vec[*OASv2Parameter] `yaml:"parameters,omitempty"`
	Responses    Map[*OASv2Response]  `yaml:"responses,omitempty"`
	ExternalDocs *OASv2ExternalDocs   `yaml:"externalDocs,omitempty"`
}

type OASv2Parameter struct {
	XOrder      int          `yaml:"x-order"`
	Name        string       `yaml:"name,omitempty"`
	In          string       `yaml:"in,omitempty"`
	Style       string       `yaml:"style,omitempty"`
	Schema      *OASv2Schema `yaml:"schema,omitempty"`
	Description string       `yaml:"description,omitempty"`
	Required    bool         `yaml:"required,omitempty"`
}

type OASv2Response struct {
	XOrder       int                `yaml:"x-order"`
	Description  string             `yaml:"description,omitempty"`
	Schema       *OASv2Schema       `yaml:"schema,omitempty"`
	ExternalDocs *OASv2ExternalDocs `yaml:"externalDocs,omitempty"`
}

type OASv2Definition struct {
}

type OASv2Schema struct {
	Ref                  string            `yaml:"$ref,omitempty"`
	Type                 string            `yaml:"type,omitempty"`
	Format               string            `yaml:"format,omitempty"`               // 当type为scalar类型时描述格式
	Items                *OASv2Schema      `yaml:"items,omitempty"`                // 当type为array时描述元素
	AdditionalProperties *OASv2Schema      `yaml:"additionalProperties,omitempty"` // 当type为object时描述元素类型(map)
	Properties           Map[*OASv2Schema] `yaml:"properties,omitempty"`           // 当type为object时描述属性
}
