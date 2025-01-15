package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"gopkg.in/yaml.v2"
	"sort"
)

// generateDocsFile 生成文档文件: *_json
func generateDocsFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.yaml`, file.GoImportPath)

	doc := &OASv2Doc{
		Swagger: "2.0",
	}

	for _, serviceExt := range meta.Services.Vec {
		doc.Tags = append(doc.Tags, &OASv2Tag{
			Deprecated: serviceExt.Deprecated,
			Name:       serviceExt.Name,
		})
	}
	bs, err := yaml.Marshal(doc)
	if err != nil {
		gen.Error(err)
	}
	_, err = g.Write(bs)
	if err != nil {
		gen.Error(err)
	}
}

type OASv2Doc struct {
	Swagger     string         `yaml:"swagger,omitempty"`
	Tags        []*OASv2Tag    `yaml:"tags,omitempty"`
	Paths       OASv2PathMap   `yaml:"paths,omitempty"`
	Definitions OASv2SchemaMap `yaml:"definitions,omitempty"`
}

type OASv2Tag struct {
	Deprecated  bool   `yaml:"deprecated,omitempty"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type OASv2Path struct {
	XOrder  int             `yaml:"x-order,omitempty"`
	Get     *OASv2Operation `json:"get,omitempty"`
	Put     *OASv2Operation `json:"put,omitempty"`
	Post    *OASv2Operation `json:"post,omitempty"`
	Delete  *OASv2Operation `json:"delete,omitempty"`
	Options *OASv2Operation `json:"options,omitempty"`
	Head    *OASv2Operation `json:"head,omitempty"`
	Patch   *OASv2Operation `json:"patch,omitempty"`
}

type OASv2Operation struct {
	XOrder      int               `yaml:"x-order,omitempty"`
	OperationId string            `yaml:"operationId,omitempty"`
	Summary     string            `yaml:"summary,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Tags        []string          `yaml:"tags,omitempty"`
	Parameters  []*OASv2Parameter `yaml:"parameters,omitempty"`
	Responses   OASv2ResponseMap  `yaml:"responses,omitempty"`
	Deprecated  bool              `yaml:"deprecated,omitempty"`
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
	XOrder               int            `yaml:"x-order,omitempty"`
	Ref                  string         `yaml:"$ref,omitempty"`
	Type                 string         `yaml:"type,omitempty"`
	Format               string         `yaml:"format,omitempty"`               // 当type为scalar类型时描述格式
	Items                *OASv2Schema   `yaml:"items,omitempty"`                // 当type为array时描述元素
	AdditionalProperties *OASv2Schema   `yaml:"additionalProperties,omitempty"` // 当type为object时描述元素类型(map)
	Properties           OASv2SchemaMap `yaml:"properties,omitempty"`           // 当type为object时描述属性
	Required             []string       `yaml:"required,omitempty"`
	Maximum              *float64       `json:"maximum,omitempty"`
	ExclusiveMaximum     bool           `json:"exclusiveMaximum,omitempty"`
	Minimum              *float64       `json:"minimum,omitempty"`
	ExclusiveMinimum     bool           `json:"exclusiveMinimum,omitempty"`
	MaxLength            *int64         `json:"maxLength,omitempty"`
	MinLength            *int64         `json:"minLength,omitempty"`
	Pattern              string         `json:"pattern,omitempty"`
	MaxItems             *int64         `json:"maxItems,omitempty"`
	MinItems             *int64         `json:"minItems,omitempty"`
	Enum                 []interface{}  `json:"enum,omitempty"`
	Deprecated           bool           `yaml:"deprecated,omitempty"`
}

type OASv2PathMap map[string]*OASv2Path

func (m OASv2PathMap) MarshalYAML() (interface{}, error) {
	var items []yaml.MapItem
	for k, v := range m {
		items = append(items, yaml.MapItem{Key: k, Value: v})
	}
	sort.SliceStable(items, func(i, j int) bool {
		vi := items[i].Value.(*OASv2Path)
		vj := items[j].Value.(*OASv2Path)
		if vi.XOrder < vj.XOrder {
			return true
		}
		return false
	})
	return items, nil
}

var _ yaml.Marshaler = OASv2PathMap{}

type OASv2SchemaMap map[string]*OASv2Schema

func (m OASv2SchemaMap) MarshalYAML() (interface{}, error) {
	var items []yaml.MapItem
	for k, v := range m {
		items = append(items, yaml.MapItem{Key: k, Value: v})
	}
	sort.SliceStable(items, func(i, j int) bool {
		vi := items[i].Value.(*OASv2Schema)
		vj := items[j].Value.(*OASv2Schema)
		if vi.XOrder < vj.XOrder {
			return true
		}
		return false
	})
	return items, nil
}

var _ yaml.Marshaler = OASv2SchemaMap{}

type OASv2ResponseMap map[string]*OASv2Response

func (m OASv2ResponseMap) MarshalYAML() (interface{}, error) {
	var items []yaml.MapItem
	for k, v := range m {
		items = append(items, yaml.MapItem{Key: k, Value: v})
	}
	sort.SliceStable(items, func(i, j int) bool {
		vi := items[i].Value.(*OASv2Response)
		vj := items[j].Value.(*OASv2Response)
		if vi.XOrder < vj.XOrder {
			return true
		}
		return false
	})
	return items, nil
}

var _ yaml.Marshaler = OASv2ResponseMap{}
