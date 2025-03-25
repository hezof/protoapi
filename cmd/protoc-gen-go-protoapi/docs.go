package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// generateDocsFile 生成文档文件: *_protoapi.yaml
func generateDocsFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.yaml`, file.GoImportPath)

	doc := NewOASv2Doc()

	messages := IdxVec[*MessageExt]{}
	for _, s := range meta.Services {
		var tag *OASv2Tag
		for _, m := range s.Methods {
			if m.Http != nil {

				if tag == nil {
					tag = NewOASv2Tag()
					if s.Info != nil {
						tag.Name = NvlS(s.Info.Name, s.FullName)
						tag.Description = NvlS(s.Info.Desc, s.FullName)
					} else {
						tag.Name = s.FullName
						tag.Description = s.FullName
					}
					tag.Deprecated = s.Deprecated
				}

				sse := m.Http.Result == Result_events

				p := NewOASv2Operation()
				p.OperationId = m.FullName
				p.Summary = NvlS(m.Http.Name, m.FullName)
				p.Description = NvlS(m.Http.Desc, m.FullName)
				p.Tags = Set(tag.Name)
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
				defineMessage(&messages, m.InputMessage)
				defineMessage(&messages, m.OutputMessage)
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

func defineMessage(ms *IdxVec[*MessageExt], m *MessageExt) {
	ms.Add(m.FullName, m)
	for _, f := range m.Fields {
		if f.Message != nil {
			defineMessage(ms, f.Message)
		}
	}
}

func fieldDescription(f *FieldExt) string {
	if f.Desc != `` {
		return f.Desc
	}
	return f.FullName
}

func definition(m *MessageExt) *Definition {
	s := NewDefinition()
	s.Description = NvlS(m.Schema, m.FullName)
	s.Type = `object`
	for _, f := range m.Fields {
		// path参数都是必需的
		if (f.In == In_path) || (f.Rule != nil && f.Rule.Required != nil) {
			s.Required = append(s.Required, f.Name)
		}
		var p *OASv2Schema
		switch {
		case f.IsRepeated:
			p = NewOASv2Schema()
			p.Type = `array`
			p.Description = fieldDescription(f)
			p.Items = field(f, true)
		case f.IsMap:
			p = NewOASv2Schema()
			p.Type = `object`
			p.Description = fieldDescription(f)
			p.Items = field(f, true)
		default:
			p = field(f, false)
		}
		// 允许覆盖字段名称
		s.Properties.Add(f.Name, p)
	}
	return s
}

func responses(h *Http, m *MessageExt) OASv2ResponseMap {
	s := make(OASv2ResponseMap)
	{
		r := NewOASv2Response()
		r.Description = `请求成功`
		if h.Result == Result_normal {
			r.Schema = NewOASv2Schema()
			r.Schema.Type = `object`
			r.Schema.Properties = make(OASv2SchemaMap)
			r.Schema.Properties.Add(`code`, &OASv2Schema{
				XOrder:      1,
				Description: `0`,
				Type:        `integer`,
				Format:      `int32`,
			})
			r.Schema.Properties.Add(`data`, &OASv2Schema{
				XOrder:      2,
				Description: NvlS(m.Schema, m.FullName),
				Ref:         ref(m.FullName),
			})
		} else {
			r.Schema = NewOASv2Schema()
			r.Schema.Ref = ref(m.FullName)
		}
		s.Add(strconv.FormatUint(uint64(NvlI(h.Status, http.StatusOK)), 10), r)
	}
	{
		for _, err := range h.Errors {
			r := NewOASv2Response()
			r.Description = err.Message
			r.Schema = NewOASv2Schema()
			r.Schema.Type = `object`
			r.Schema.Properties = make(OASv2SchemaMap)
			if err.Code != 0 {
				r.Schema.Properties.Add(`code`, &OASv2Schema{
					XOrder:      1,
					Description: strconv.FormatUint(uint64(err.Code), 10),
					Type:        `integer`,
					Format:      `int32`,
				})
			}
			if err.Name != `` {
				r.Schema.Properties.Add(`name`, &OASv2Schema{
					XOrder:      2,
					Description: err.Name,
					Type:        `string`,
				})
			}
			if err.Message != `` {
				r.Schema.Properties.Add(`message`, &OASv2Schema{
					XOrder:      3,
					Description: err.Message,
					Type:        `string`,
				})
			}
			s.Add(strconv.FormatUint(uint64(err.Status), 10), r)
		}
	}
	{
		r := NewOASv2Response()
		r.Description = `请求错误`
		r.Schema = NewOASv2Schema()
		r.Schema.Type = `object`
		r.Schema.Properties = make(OASv2SchemaMap)
		r.Schema.Properties.Add(`code`, &OASv2Schema{
			XOrder:      1,
			Description: `错误代码. 约定0表示成功, 其他表示错误.`,
			Type:        `integer`,
			Format:      `int32`,
		})
		r.Schema.Properties.Add(`name`, &OASv2Schema{
			XOrder:      2,
			Description: `错误名称.可选`,
			Type:        `string`,
		})
		r.Schema.Properties.Add(`message`, &OASv2Schema{
			XOrder:      3,
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
		if f.In != In_body {
			other++
		}
	}

	ss := make([]*OASv2Parameter, 0, len(m.Fields))
	if other > 0 {
		call := func(f *FieldExt) {
			// 默认都是在body
			if f.In != In_body {
				p := NewOASv2Parameter()
				p.Name = f.Name
				switch f.In {
				case In_path:
					p.In = `path`
				case In_query:
					p.In = `query`
				case In_header:
					p.In = `header`
				case In_cookie:
					p.In = `cookie`
				}
				switch {
				case f.IsRepeated:
					p.OASv2Schema = OASv2Schema{
						Type:             `array`,
						Description:      fieldDescription(f),
						Items:            field(f, true),
						CollectionFormat: If(f.Explode, `multi`, `csv`),
					}
				case f.IsMap:
					p.OASv2Schema = OASv2Schema{
						Type:                 `object`,
						Description:          fieldDescription(f),
						AdditionalProperties: field(f, true),
						CollectionFormat:     If(f.Explode, `multi`, `csv`),
					}
				default:
					p.OASv2Schema = *field(f, false)
				}
				p.Required = f.Rule != nil && f.Rule.Required != nil
				if f.In == In_path {
					p.Required = true // path参数必需的
				}
				p.Deprecated = f.Deprecated
				ss = append(ss, p)
			}
		}
		// path
		for _, f := range m.Fields {
			if f.In == In_path {
				call(f)
			}
		}
		// query
		for _, f := range m.Fields {
			if f.In == In_query {
				call(f)
			}
		}
		// header
		for _, f := range m.Fields {
			if f.In == In_header {
				call(f)
			}
		}
		// cookie
		for _, f := range m.Fields {
			if f.In == In_cookie {
				call(f)
			}
		}
	}
	if other < total {
		// body处理
		p := NewOASv2Parameter()
		p.Name = `body`
		p.In = `body`
		p.Required = true // body参数必须的
		if other == 0 {
			p.Schema = NewOASv2Schema()
			p.Schema.Ref = ref(m.FullName)
		} else {
			p.Schema = NewOASv2Schema()
			p.Schema.Properties = make(OASv2SchemaMap)
			for _, f := range m.Fields {
				// 默认都是在body
				if f.In == In_body {
					p.Schema.Properties.Add(f.Name, field(f, false))
				}
			}
		}
		ss = append(ss, p)
	}
	return ss
}

func field(f *FieldExt, sub bool) *OASv2Schema {
	s := NewOASv2Schema()
	if !sub {
		s.Description = fieldDescription(f)
	}
	switch f.Kind {
	case protoreflect.BoolKind:
		s.Type = `boolean`
	case protoreflect.EnumKind:
		if f.EnumName {
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
	case protoreflect.BytesKind:
		s.Type = `string`
		s.Format = `binary`
	case protoreflect.MessageKind:
		// 需要区分map
		if f.IsMap {
			s.Type = `object`
			s.AdditionalProperties = field(f, true)
			s.CollectionFormat = If(f.Explode, `multi`, `csv`)
		} else {
			s.Ref = ref(f.Message.FullName)
		}
	case protoreflect.GroupKind:
		panic(fmt.Sprintf("invalid kind of %v: group", f.Name))
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
				s.Enum = f.RuleEnums()
			}
		}
		if f.Rule.Pattern != nil {
			s.Pattern = f.Rule.Pattern.Val
		}
	}
	s.Deprecated = f.Deprecated
	return s
}

func ref(fullName string) string {
	return fmt.Sprintf(`#/definitions/%s`, fullName)
}

/****************************************************************
 * swagger 2.0 数据结构及辅助结构
 ****************************************************************/

type OASv2PathMap map[string]*OASv2Path

func (m OASv2PathMap) Add(method, path string, op OASv2Operation, sse bool) error {

	// 忽略path为空的方法
	if path == `` {
		return nil
	}

	method = strings.ToUpper(method)
	op.OperationId = op.OperationId + "." + method
	// WEBSOCKET会自动忽略sse
	if sse && method != `WEBSOCKET` {
		op.Description = fmt.Sprintf(`%v %v [SSE]<br/><br/>%v`, method, path, op.Description)
	} else {
		op.Description = fmt.Sprintf(`%v %v<br/><br/>%v`, method, path, op.Description)
	}

	val := m[path]
	if val == nil {
		val = NewOASv2Path()
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

func (m OASv2SchemaMap) Add(k string, v *OASv2Schema) {
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

func (m OASv2ResponseMap) Add(k string, v *OASv2Response) {
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

func (m DefinitionMap) Add(k string, v *Definition) {
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

func NewOASv2Doc() *OASv2Doc {
	return &OASv2Doc{
		Swagger:     "2.0",
		Paths:       make(OASv2PathMap),
		Definitions: make(DefinitionMap),
	}
}

type OASv2Tag struct {
	XOrder      int    `yaml:"x-order,omitempty"`
	Deprecated  bool   `yaml:"deprecated,omitempty"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

var OASv2TagOrder int

func NewOASv2Tag() *OASv2Tag {
	OASv2TagOrder++
	return &OASv2Tag{
		XOrder: OASv2TagOrder,
	}
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

var OASv2PathOrder int

func NewOASv2Path() *OASv2Path {
	OASv2PathOrder++
	return &OASv2Path{
		XOrder: OASv2PathOrder,
	}
}

type OASv2Operation struct {
	XOrder      int               `yaml:"x-order,omitempty"`
	OperationId string            `yaml:"operationId,omitempty"`
	Summary     string            `yaml:"summary,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Tags        []string          `yaml:"tags,omitempty"`
	Parameters  []*OASv2Parameter `yaml:"parameters,omitempty"`
	Responses   OASv2ResponseMap  `yaml:"responses,omitempty"`
	Produces    []string          `yaml:"produces,omitempty"`
	Deprecated  bool              `yaml:"deprecated,omitempty"`
}

var OASv2OperationOrder int

func NewOASv2Operation() *OASv2Operation {
	OASv2OperationOrder++
	return &OASv2Operation{
		XOrder:    OASv2OperationOrder,
		Responses: make(OASv2ResponseMap),
		Produces:  []string{"application/json"},
	}
}

type OASv2Parameter struct {
	OASv2Schema `yaml:",inline"` // if in is not "body"
	Name        string           `yaml:"name,omitempty"`
	In          string           `yaml:"in,omitempty"`
	Schema      *OASv2Schema     `yaml:"schema,omitempty"` // If in is "body":
	Required    bool             `yaml:"required,omitempty"`
}

var OASv2ParameterOrder int

func NewOASv2Parameter() *OASv2Parameter {
	OASv2ParameterOrder++
	return &OASv2Parameter{
		OASv2Schema: OASv2Schema{
			XOrder: OASv2ParameterOrder,
		},
	}
}

type OASv2Response struct {
	XOrder      int          `yaml:"x-order,omitempty"`
	Description string       `yaml:"description,omitempty"`
	Schema      *OASv2Schema `yaml:"schema,omitempty"`
}

var OASv2ResponseOrder int

func NewOASv2Response() *OASv2Response {
	OASv2ResponseOrder++
	return &OASv2Response{
		XOrder: OASv2ResponseOrder,
	}
}

type OASv2Schema struct {
	XOrder               int            `yaml:"x-order,omitempty"`
	Ref                  string         `yaml:"$ref,omitempty"`
	Type                 string         `yaml:"type,omitempty"`
	Format               string         `yaml:"format,omitempty"`               // 当type为scalar时的解析方式
	Items                *OASv2Schema   `yaml:"items,omitempty"`                // 当type为array时描述元素
	CollectionFormat     string         `yaml:"collectionFormat,omitempty"`     // 当type为array时的解析方式
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

var OASv2SchemaOrder int

func NewOASv2Schema() *OASv2Schema {
	OASv2SchemaOrder++
	return &OASv2Schema{
		XOrder: OASv2SchemaOrder,
	}
}

type Definition struct {
	OASv2Schema `yaml:",inline"`
	Required    []string `yaml:"required,omitempty"`
}

var DefinitionOrder int

func NewDefinition() *Definition {
	DefinitionOrder++
	return &Definition{
		OASv2Schema: OASv2Schema{
			XOrder:     DefinitionOrder,
			Properties: make(OASv2SchemaMap),
		},
	}
}
