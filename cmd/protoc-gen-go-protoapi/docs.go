package main

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v2"
	"net/http"
	"sort"
	"strings"
)

// generateDocsFile 生成文档文件: *_protoapi.yaml
func generateDocsFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.yaml`, file.GoImportPath)

	doc := &OASv2Doc{
		Swagger:     "2.0",
		Paths:       make(OASv2PathMap),
		Definitions: make(OASv2SchemaMap),
	}

	messages := IdxVec[*MessageExt]{}
	for _, s := range meta.Services {
		var tag *OASv2Tag
		for _, m := range s.Methods {
			if m.Http != nil {

				if tag == nil {
					tag = new(OASv2Tag)
					if s.Tag != nil {
						tag.Name = NvlS(s.Tag.Name, s.FullName)
						tag.Description = NvlS(s.Tag.Desc, s.FullName)
					} else {
						tag.Name = s.FullName
						tag.Description = s.FullName
					}
					tag.Deprecated = s.Deprecated
				}

				p := &OASv2Operation{Responses: make(OASv2ResponseMap)}
				p.OperationId = m.FullName
				p.Summary = NvlS(m.Http.Name, m.FullName)
				p.Description = NvlS(m.Http.Desc, m.FullName)
				p.Tags = Set(tag.Name, m.Http.Tags...)
				p.Parameters = parameters(m.InputMessage)
				p.Responses = responses(m.Http, m.OutputMessage)
				p.Deprecated = m.Deprecated
				if err := doc.Paths.Add(http.MethodGet, m.Http.Get, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodPut, m.Http.Put, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodPost, m.Http.Post, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodDelete, m.Http.Delete, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodOptions, m.Http.Options, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodHead, m.Http.Head, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodPatch, m.Http.Patch, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodTrace, m.Http.Trace, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodConnect, m.Http.Connect, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(`websocket`, m.Http.Websocket, p); err != nil {
					gen.Error(fmt.Errorf(`%v http error %v`, m.FullName, err))
					return
				}

				// definitions定义使用
				messages.Add(m.InputMessage.FullName, m.InputMessage)
				messages.Add(m.OutputMessage.FullName, m.OutputMessage)
			}
		}
		if tag != nil {
			doc.Tags = append(doc.Tags, tag)
		}
	}

	for _, m := range messages.Vec {
		doc.Definitions.Add(m.FullName, definition(m))
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

func definition(m *MessageExt) *OASv2Schema {
	s := &OASv2Schema{Properties: make(OASv2SchemaMap)}
	s.Type = `object`
	var required []string
	for _, f := range m.Fields {
		// path参数都是必需的
		if (f.Prop != nil && f.Prop.In == Prop_path) || (f.Rule != nil && f.Rule.Required != nil) {
			required = append(required, f.Name)
		}
		p := new(OASv2Schema)
		switch {
		case f.IsRepeated:
			p.Type = `array`
			p.Items = gpf(f)
		case f.IsMap:
		default:

		}
		s.Properties.Add(f.Name, p)
	}
	return s
}

func responses(h *Http, m *MessageExt) OASv2ResponseMap {

}

func parameters(m *MessageExt) []*OASv2Parameter {

}

func gpf(f *FieldExt) *OASv2Schema {
	s := new(OASv2Schema)
	switch f.Kind {
	case protoreflect.BoolKind:
		s.Type = `boolean`
	case protoreflect.EnumKind:
		if f.Prop != nil && f.Prop.EnumName {
			s.Type = `string`
			for _, e := range f.Enum.Values {
				s.Enum = append(s.Enum, e.Name)
			}
		} else {
			s.Type = `integer`
			s.Format = `int32`
			for _, e := range f.Enum.Values {
				s.Enum = append(s.Enum, e.Number)
			}
		}
	case protoreflect.Int32Kind:
		s.Type = `integer`
		s.Format = `int32`
	case protoreflect.Sint32Kind:
		s.Type = `integer`
		s.Format = `int32`
	case protoreflect.Uint32Kind:
		s.Type = `integer`
		s.Format = `int32`
	case protoreflect.Int64Kind:
		s.Type = `integer`
		s.Format = `int64`
	case protoreflect.Sint64Kind:
		s.Type = `integer`
		s.Format = `int64`
	case protoreflect.Uint64Kind:
		s.Type = `integer`
		s.Format = `int64`
	case protoreflect.Sfixed32Kind:
		s.Type = `integer`
		s.Format = `int32`
	case protoreflect.Fixed32Kind:
		s.Type = `integer`
		s.Format = `int32`
	case protoreflect.FloatKind:
		s.Type = `number`
		s.Format = `float`
	case protoreflect.Sfixed64Kind:
		s.Type = `integer`
		s.Format = `int64`
	case protoreflect.Fixed64Kind:
		s.Type = `integer`
		s.Format = `int64`
	case protoreflect.DoubleKind:
		s.Type = `number`
		s.Format = `double`
	case protoreflect.StringKind:
		s.Type = `string`
		if f.Prop != nil {
			s.Format = f.Prop.Format
		}
	case protoreflect.BytesKind:
		s.Type = `string`
		s.Format = `binary`
	case protoreflect.MessageKind:
		s.Ref = fmt.Sprintf(`"#/definitions/%s"`, f.Message.FullName)
	case protoreflect.GroupKind:
		s.Ref = fmt.Sprintf(`"#/definitions/%s"`, f.Message.FullName)
	}
	if f.Rule != nil {
		if f.Rule.Minimum != nil {
			s.Minimum = float64(f.Rule.Minimum.Val)
			s.ExclusiveMinimum = f.Rule.Minimum.Exclusive
		}
		if f.Rule.Maximum != nil {
			s.Maximum = float64(f.Rule.Maximum.Val)
			s.ExclusiveMaximum = f.Rule.Maximum.Exclusive
		}
		if f.Rule.MinLength != nil {
			s.MinLength = f.Rule.MinLength.Val
		}
		if f.Rule.MaxLength != nil {
			s.MaxLength = f.Rule.MaxLength.Val
		}
		if f.Rule.MinItems != nil {
			s.MinItems = f.Rule.MinItems.Val
		}
		if f.Rule.MaxItems != nil {
			s.MaxItems = f.Rule.MaxItems.Val
		}
		if f.Rule.Enum != nil {
			if s.Enum == nil {
				for _, v := range strings.Split(f.Rule.Enum.Val, `,`)
			}
		}
		if f.Rule.Pattern != nil {

		}
	}
	s.Deprecated = f.Deprecated
}

/****************************************************************
 * swagger 2.0 数据结构及辅助结构
 ****************************************************************/

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
	Get     *OASv2Operation `yaml:"get,omitempty"`
	Put     *OASv2Operation `yaml:"put,omitempty"`
	Post    *OASv2Operation `yaml:"post,omitempty"`
	Delete  *OASv2Operation `yaml:"delete,omitempty"`
	Options *OASv2Operation `yaml:"options,omitempty"`
	Head    *OASv2Operation `yaml:"head,omitempty"`
	Patch   *OASv2Operation `yaml:"patch,omitempty"`
	Connect *OASv2Operation `yaml:"connect,omitempty"`
	Trace   *OASv2Operation `yaml:"trace,omitempty"`
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
	Maximum              float64        `yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool           `yaml:"exclusiveMaximum,omitempty"`
	Minimum              float64        `yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool           `yaml:"exclusiveMinimum,omitempty"`
	MaxLength            int64          `yaml:"maxLength,omitempty"`
	MinLength            int64          `yaml:"minLength,omitempty"`
	Pattern              string         `yaml:"pattern,omitempty"`
	MaxItems             int64          `yaml:"maxItems,omitempty"`
	MinItems             int64          `yaml:"minItems,omitempty"`
	Enum                 []interface{}  `yaml:"enum,omitempty"`
	Deprecated           bool           `yaml:"deprecated,omitempty"`
}

type OASv2PathMap map[string]*OASv2Path

var OASv2PathOrder int

func (m OASv2PathMap) Add(method, path string, op *OASv2Operation) error {

	if path == `` {
		return errors.New("path is empty")
	}

	val := m[path]
	if val == nil {
		OASv2PathOrder++
		val = &OASv2Path{
			XOrder: OASv2PathOrder,
		}
	}
	switch method := strings.ToLower(method); method {
	case `get`:
		if val.Get != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Get = op
	case `put`:
		if val.Put != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Put = op
	case `post`:
		if val.Post != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Post = op
	case `delete`:
		if val.Delete != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Delete = op
	case `options`:
		if val.Options != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Options = op
	case `head`:
		if val.Head != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Head = op
	case `patch`:
		if val.Patch != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Patch = op
	case `connect`:
		if val.Connect != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Connect = op
	case `trace`:
		if val.Trace != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Trace = op
	case `websocket`:
		if val.Get != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Get = op
	}
	m[path] = val
	return nil
}

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

var OASv2SchemaOrder int

func (m OASv2SchemaMap) Add(k string, v *OASv2Schema) {
	OASv2SchemaOrder++
	v.XOrder = OASv2SchemaOrder

	m[k] = v
}

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

var OASv2ResponseOrder int

func (m OASv2ResponseMap) Add(k string, v *OASv2Response) {
	OASv2ResponseOrder++
	v.XOrder = OASv2ResponseOrder
	m[k] = v
}

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
