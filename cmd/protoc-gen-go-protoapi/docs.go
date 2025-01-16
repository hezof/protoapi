package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v2"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// generateDocsFile 生成文档文件: *_protoapi.yaml
func generateDocsFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.yaml`, file.GoImportPath)

	doc := &OASv2Doc{
		Swagger:     "2.0",
		Paths:       make(OASv2PathMap),
		Definitions: make(DefinitionMap),
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

				sse := m.Http.Result == Http_events

				p := &OASv2Operation{Responses: make(OASv2ResponseMap)}
				p.OperationId = m.FullName
				p.Summary = NvlS(m.Http.Name, m.FullName)
				p.Description = NvlS(m.Http.Desc, m.FullName)
				p.Tags = Set(tag.Name, m.Http.Tags...)
				p.Parameters = parameters(m.Http, m.InputMessage)
				p.Responses = responses(m.Http, m.OutputMessage)
				p.Deprecated = m.Deprecated

				if err := doc.Paths.Add(http.MethodGet, m.Http.Get, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodPut, m.Http.Put, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodPost, m.Http.Post, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodDelete, m.Http.Delete, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodOptions, m.Http.Options, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodHead, m.Http.Head, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodPatch, m.Http.Patch, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodTrace, m.Http.Trace, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(http.MethodConnect, m.Http.Connect, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
					return
				}
				if err := doc.Paths.Add(`websocket`, m.Http.Websocket, *p, sse); err != nil {
					gen.Error(fmt.Errorf(`%v %v`, m.FullName, err))
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

func fieldDescription(f *FieldExt) string {
	if f.Prop != nil && f.Prop.Desc != `` {
		return f.Prop.Desc
	}
	return f.FullName
}

func definition(m *MessageExt) *Definition {
	s := &Definition{OASv2Schema: OASv2Schema{Properties: make(OASv2SchemaMap)}}
	s.Description = NvlS(m.Desc, m.FullName)
	s.Type = `object`
	for _, f := range m.Fields {
		// path参数都是必需的
		if (f.Prop != nil && f.Prop.In == Prop_path) || (f.Rule != nil && f.Rule.Required != nil) {
			s.Required = append(s.Required, f.Name)
		}
		var p *OASv2Schema
		switch {
		case f.IsRepeated:
			p = &OASv2Schema{
				Type:        `array`,
				Description: fieldDescription(f),
				Items:       field(f, true),
			}
		case f.IsMap:
			p = &OASv2Schema{
				Type:                 `object`,
				Description:          fieldDescription(f),
				AdditionalProperties: field(f.Message.Fields[1], true),
			}
		default:
			p = field(f, false)
		}
		// 允许覆盖字段名称
		s.Properties.Add(fname(f), p)
	}
	return s
}

func responses(h *Http, m *MessageExt) OASv2ResponseMap {
	s := make(OASv2ResponseMap)
	{
		r := new(OASv2Response)
		r.Description = `请求成功`
		if h.Result == Http_simple {
			r.Schema = &OASv2Schema{
				Type:       `object`,
				Properties: make(OASv2SchemaMap),
			}
			r.Schema.Properties.Add(`code`, &OASv2Schema{
				Description: `0`,
				Type:        `integer`,
				Format:      `int32`,
			})
			r.Schema.Properties.Add(`data`, &OASv2Schema{
				Description: `结果数据.`,
				Ref:         ref(m.FullName),
			})
		} else {
			r.Schema = &OASv2Schema{Ref: ref(m.FullName)}
		}
		s.Add(strconv.FormatUint(uint64(NvlI(h.Status, http.StatusOK)), 10), r)
	}
	{
		for _, err := range h.Errors {
			r := new(OASv2Response)
			r.Description = err.Message
			r.Schema = &OASv2Schema{
				Type:       `object`,
				Properties: make(OASv2SchemaMap),
			}
			if err.Code != 0 {
				r.Schema.Properties.Add(`code`, &OASv2Schema{
					Description: strconv.FormatUint(uint64(err.Code), 10),
					Type:        `integer`,
					Format:      `int32`,
				})
			}
			if err.Name != `` {
				r.Schema.Properties.Add(`name`, &OASv2Schema{
					Description: err.Name,
					Type:        `string`,
				})
			}
			if err.Message != `` {
				r.Schema.Properties.Add(`message`, &OASv2Schema{
					Description: err.Message,
					Type:        `string`,
				})
			}
			s.Add(strconv.FormatUint(uint64(err.Status), 10), r)
		}
	}
	{
		r := new(OASv2Response)
		r.Description = `请求错误`
		r.Schema = &OASv2Schema{
			Type:       `object`,
			Properties: make(OASv2SchemaMap),
		}
		r.Schema.Properties.Add(`code`, &OASv2Schema{
			Description: `错误代码. 约定0表示成功, 其他表示错误.`,
			Type:        `integer`,
			Format:      `int32`,
		})
		r.Schema.Properties.Add(`name`, &OASv2Schema{
			Description: `错误名称.可选`,
			Type:        `string`,
		})
		r.Schema.Properties.Add(`message`, &OASv2Schema{
			Description: `错误消息.可选`,
			Type:        `string`,
		})
		s.Add(`default`, r)
	}
	return s
}

func parameters(h *Http, m *MessageExt) []*OASv2Parameter {

	total := len(m.Fields)
	other := 0 // 除了body外的其他位置
	for _, f := range m.Fields {
		// 默认都是在body
		if f.Prop != nil && f.Prop.In != Prop_body {
			other++
		}
	}

	ss := make([]*OASv2Parameter, 0, len(m.Fields))
	if other > 0 {
		call := func(f *FieldExt) {
			// 默认都是在body
			if f.Prop != nil && f.Prop.In != Prop_body {
				p := new(OASv2Parameter)
				p.Name = fname(f)
				switch f.Prop.In {
				case Prop_path:
					p.In = `path`
				case Prop_query:
					p.In = `query`
				case Prop_header:
					p.In = `header`
				case Prop_cookie:
					p.In = `cookie`
				}
				switch {
				case f.IsRepeated:
					p.OASv2Schema = OASv2Schema{
						Type:        `array`,
						Description: fieldDescription(f),
						Items:       field(f, true),
					}
				case f.IsMap:
					p.OASv2Schema = OASv2Schema{
						Type:                 `object`,
						Description:          fieldDescription(f),
						AdditionalProperties: field(f.Message.Fields[1], true),
					}
				default:
					p.OASv2Schema = *field(f, false)
				}
				switch f.Prop.Style {
				case Prop_simple:
					p.Style = `simple`
				case Prop_form:
					p.Style = `form`
				case Prop_json:
					p.Style = `json`
				}
				p.Explode = true // 默认是true
				p.Required = f.Rule != nil && f.Rule.Required != nil
				if f.Prop.In == Prop_path {
					p.Required = true // path参数必需的
				}
				p.Deprecated = f.Deprecated
				ss = append(ss, p)
			}
		}
		// path
		for _, f := range m.Fields {
			if f.Prop != nil && f.Prop.In == Prop_path {
				call(f)
			}
		}
		// query
		for _, f := range m.Fields {
			if f.Prop != nil && f.Prop.In == Prop_query {
				call(f)
			}
		}
		// header
		for _, f := range m.Fields {
			if f.Prop != nil && f.Prop.In == Prop_header {
				call(f)
			}
		}
		// cookie
		for _, f := range m.Fields {
			if f.Prop != nil && f.Prop.In == Prop_cookie {
				call(f)
			}
		}
	}
	if other < total {
		// body处理
		p := new(OASv2Parameter)
		p.Name = `body`
		p.In = `body`
		p.Required = true // body参数必须的
		if other == 0 {
			p.Schema = &OASv2Schema{Ref: ref(m.FullName)}
		} else {
			p.Schema = &OASv2Schema{Properties: make(OASv2SchemaMap)}
			for _, f := range m.Fields {
				// 默认都是在body
				if f.Prop == nil || f.Prop.In == Prop_body {
					p.Schema.Properties.Add(fname(f), field(f, false))
				}
			}
		}
		ss = append(ss, p)
	}
	return ss
}

func fname(f *FieldExt) string {
	if f.Prop != nil && f.Prop.Name != `` {
		return f.Prop.Name
	}
	return f.Name
}

func field(f *FieldExt, sub bool) *OASv2Schema {
	s := new(OASv2Schema)
	if !sub {
		s.Description = fieldDescription(f)
	}
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
		// 需要区分map
		if f.IsMap {
			s.Type = `object`
			s.AdditionalProperties = field(f.Message.Fields[1], true)
		} else {
			s.Ref = ref(f.Message.FullName)
		}
	case protoreflect.GroupKind:
		if f.IsMap {
			s.Type = `object`
			s.AdditionalProperties = field(f.Message.Fields[1], true)
		} else {
			s.Ref = ref(f.Message.FullName)
		}
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
			// 避免覆盖enum的枚举
			if s.Enum == nil {
				s.Enum = enums(f.Kind, f.Rule.Enum.Val)
			}
		}
		if f.Rule.Pattern != nil {
			s.Pattern = f.Rule.Pattern.Val
		}
	}
	s.Deprecated = f.Deprecated
	return s
}

func enums(k protoreflect.Kind, s string) []interface{} {
	var rt []interface{}
	for _, v := range strings.Split(s, `,`) {
		v = strings.TrimSpace(v)
		if v != `` {
			switch k {
			case protoreflect.BoolKind:
				rt = append(rt, v == `true`)
			case protoreflect.Int32Kind:
				vi, _ := strconv.ParseInt(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Sint32Kind:
				vi, _ := strconv.ParseInt(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Uint32Kind:
				vi, _ := strconv.ParseUint(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Int64Kind:
				vi, _ := strconv.ParseInt(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Sint64Kind:
				vi, _ := strconv.ParseInt(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Uint64Kind:
				vi, _ := strconv.ParseUint(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Sfixed32Kind:
				vi, _ := strconv.ParseInt(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Fixed32Kind:
				vi, _ := strconv.ParseUint(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.FloatKind:
				vi, _ := strconv.ParseFloat(v, 32)
				rt = append(rt, vi)
			case protoreflect.Sfixed64Kind:
				vi, _ := strconv.ParseInt(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.Fixed64Kind:
				vi, _ := strconv.ParseUint(v, 10, 64)
				rt = append(rt, vi)
			case protoreflect.DoubleKind:
				vi, _ := strconv.ParseFloat(v, 64)
				rt = append(rt, vi)
			case protoreflect.StringKind:
				rt = append(rt, v)
			}
		}
	}
	return rt
}

func ref(fullName string) string {
	return fmt.Sprintf(`#/definitions/%s`, fullName)
}

/****************************************************************
 * swagger 2.0 数据结构及辅助结构
 ****************************************************************/

type OASv2PathMap map[string]*OASv2Path

var OASv2PathOrder int

func (m OASv2PathMap) Add(method, path string, op OASv2Operation, sse bool) error {

	// 忽略path为空的方法
	if path == `` {
		return nil
	}

	method = strings.ToUpper(method)
	// WEBSOCKET会自动忽略sse
	if sse && method != `WEBSOCKET` {
		op.Description = fmt.Sprintf(`%v %v [SSE]<br/><br/>%v`, method, path, op.Description)
	} else {
		op.Description = fmt.Sprintf(`%v %v<br/><br/>%v`, method, path, op.Description)
	}

	val := m[path]
	if val == nil {
		OASv2PathOrder++
		val = &OASv2Path{
			XOrder: OASv2PathOrder,
		}
	}
	switch method {
	case http.MethodGet:
		if val.Get != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Get = &op
	case http.MethodPut:
		if val.Put != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Put = &op
	case http.MethodPost:
		if val.Post != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Post = &op
	case http.MethodDelete:
		if val.Delete != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Delete = &op
	case http.MethodOptions:
		if val.Options != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Options = &op
	case http.MethodHead:
		if val.Head != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Head = &op
	case http.MethodPatch:
		if val.Patch != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Patch = &op
	case http.MethodConnect:
		if val.Connect != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Connect = &op
	case http.MethodTrace:
		if val.Trace != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Trace = &op
	case `WEBSOCKET`:
		if val.Get != nil {
			return fmt.Errorf(`%v %v duplicated`, method, path)
		}
		val.Get = &op
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

type DefinitionMap map[string]*Definition

var DefinitionOrder int

func (m DefinitionMap) Add(k string, v *Definition) {
	DefinitionOrder++
	v.XOrder = DefinitionOrder

	m[k] = v
}

func (m DefinitionMap) MarshalYAML() (interface{}, error) {
	var items []yaml.MapItem
	for k, v := range m {
		items = append(items, yaml.MapItem{Key: k, Value: v})
	}
	sort.SliceStable(items, func(i, j int) bool {
		vi := items[i].Value.(*Definition)
		vj := items[j].Value.(*Definition)
		if vi.XOrder < vj.XOrder {
			return true
		}
		return false
	})
	return items, nil
}

var _ yaml.Marshaler = DefinitionMap{}

type OASv2Doc struct {
	Swagger     string        `yaml:"swagger,omitempty"`
	Tags        []*OASv2Tag   `yaml:"tags,omitempty"`
	Paths       OASv2PathMap  `yaml:"paths,omitempty"`
	Definitions DefinitionMap `yaml:"definitions,omitempty"`
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
	Schema      *OASv2Schema `yaml:"schema,omitempty"` // If in is "body":
	OASv2Schema `yaml:",inline"`                       // if in is not "body"
	Style       string       `yaml:"style,omitempty"`
	Explode     bool         `yaml:"explode,omitempty"`
	Required    bool         `yaml:"required,omitempty"`
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
	Description          string         `yaml:"description,omitempty"`
	Deprecated           bool           `yaml:"deprecated,omitempty"`
}

type Definition struct {
	OASv2Schema `yaml:",inline"`
	Required    []string `yaml:"required,omitempty"`
}
