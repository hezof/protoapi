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

type OASv2Doc struct {
	Swagger     string                  `yaml:"swagger,omitempty"`
	Tags        []*OASv2Tag             `yaml:"tags,omitempty"`
	Paths       map[string]*OASv2Path   `yaml:"paths,omitempty"`
	Definitions map[string]*OASv2Schema `yaml:"definitions,omitempty"`
}

type OASv2Tag struct {
	Deprecated  bool   `yaml:"deprecated,omitempty"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type OASv2Path struct {
	Get     *OASv2Operation `json:"get,omitempty"`
	Put     *OASv2Operation `json:"put,omitempty"`
	Post    *OASv2Operation `json:"post,omitempty"`
	Delete  *OASv2Operation `json:"delete,omitempty"`
	Options *OASv2Operation `json:"options,omitempty"`
	Head    *OASv2Operation `json:"head,omitempty"`
	Patch   *OASv2Operation `json:"patch,omitempty"`
}

type OASv2Operation struct {
	XOrder      int                       `yaml:"x-order,omitempty"`
	OperationId string                    `yaml:"operationId,omitempty"`
	Summary     string                    `yaml:"summary,omitempty"`
	Description string                    `yaml:"description,omitempty"`
	Tags        []string                  `yaml:"tags,omitempty"`
	Parameters  []*OASv2Parameter         `yaml:"parameters,omitempty"`
	Responses   map[string]*OASv2Response `yaml:"responses,omitempty"`
	Deprecated  bool                      `yaml:"deprecated,omitempty"`
}

type OASv2Parameter struct {
	Name        string       `yaml:"name,omitempty"`
	In          string       `yaml:"in,omitempty"`
	Style       string       `yaml:"style,omitempty"`
	Schema      *OASv2Schema `yaml:"schema,omitempty"`
	Description string       `yaml:"description,omitempty"`
	Required    bool         `yaml:"required,omitempty"`
	Deprecated  bool         `yaml:"deprecated,omitempty"`
}

type OASv2Response struct {
	XOrder      int          `yaml:"x-order,omitempty"`
	Description string       `yaml:"description,omitempty"`
	Schema      *OASv2Schema `yaml:"schema,omitempty"`
}

type OASv2Schema struct {
	XOrder               int                     `yaml:"x-order,omitempty"`
	Ref                  string                  `yaml:"$ref,omitempty"`
	Type                 string                  `yaml:"type,omitempty"`
	Format               string                  `yaml:"format,omitempty"`               // 当type为scalar类型时描述格式
	Items                *OASv2Schema            `yaml:"items,omitempty"`                // 当type为array时描述元素
	AdditionalProperties *OASv2Schema            `yaml:"additionalProperties,omitempty"` // 当type为object时描述元素类型(map)
	Properties           map[string]*OASv2Schema `yaml:"properties,omitempty"`           // 当type为object时描述属性
	Required             []string                `yaml:"required,omitempty"`
	Maximum              *float64                `json:"maximum,omitempty"`
	ExclusiveMaximum     bool                    `json:"exclusiveMaximum,omitempty"`
	Minimum              *float64                `json:"minimum,omitempty"`
	ExclusiveMinimum     bool                    `json:"exclusiveMinimum,omitempty"`
	MaxLength            *int64                  `json:"maxLength,omitempty"`
	MinLength            *int64                  `json:"minLength,omitempty"`
	Pattern              string                  `json:"pattern,omitempty"`
	MaxItems             *int64                  `json:"maxItems,omitempty"`
	MinItems             *int64                  `json:"minItems,omitempty"`
	Enum                 []interface{}           `json:"enum,omitempty"`
	Deprecated           bool                    `yaml:"deprecated,omitempty"`
}
