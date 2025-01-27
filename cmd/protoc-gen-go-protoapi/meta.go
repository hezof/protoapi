package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func extractFile(gen *protogen.Plugin, f *protogen.File) *FileExt {
	file := new(FileExt)
	file.Path = f.Desc.Path()
	file.Package = string(f.Desc.Package())
	file.GoPackage = string(f.GoPackageName)
	file.GoImportPath = string(f.GoImportPath)

	for _, s := range f.Enums {
		extractEnum(file, s)
	}

	for _, s := range f.Messages {
		m := extractMessage(file, s, false)
		for _, v := range m.Fields {
			if v.IsMap && v.Message.Fields[0].Kind != protoreflect.StringKind {
				gen.Error(fmt.Errorf("invalid map key type %v", v.Message.Fields[0].Kind))
				return nil
			}
		}
	}

	for _, s := range f.Services {
		extractService(file, s)
		// TODO: 校验method合法性
	}

	return file
}

func extractEnumValue(f *FileExt, e *EnumExt, s *protogen.EnumValue) *EnumValueExt {
	if s == nil {
		return nil
	}
	v := new(EnumValueExt)
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
	v.Number = int32(s.Desc.Number())
	v.Deprecated = s.Desc.Options().(*descriptorpb.EnumValueOptions).GetDeprecated()
	e.Values = append(e.Values, v)
	return v
}

func extractEnum(f *FileExt, s *protogen.Enum) *EnumExt {
	if s == nil {
		return nil
	}
	v := new(EnumExt)
	v.FilePath = s.Desc.ParentFile().Path()
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
	if rv, ok := f.Enums.Add(v.FullName, v); !ok {
		return rv
	}
	for _, s1 := range s.Values {
		extractEnumValue(f, v, s1)
	}
	v.Deprecated = s.Desc.Options().(*descriptorpb.EnumOptions).GetDeprecated()
	return v
}

func extractField(file *FileExt, message *MessageExt, s *protogen.Field) *FieldExt {
	if s == nil {
		return nil
	}
	v := new(FieldExt)
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	v.GoIdent = s.GoIdent
	v.Kind = s.Desc.Kind()
	v.IsMap = s.Desc.IsMap()
	v.IsRepeated = s.Desc.IsList()
	v.HasOptional = s.Desc.HasOptionalKeyword()
	v.Message = extractMessage(file, s.Message, s.Desc.IsMap())
	v.Prop = proto.GetExtension(s.Desc.Options(), E_Prop).(*Prop)
	v.Rule = proto.GetExtension(s.Desc.Options(), E_Rule).(*Rule)
	v.Deprecated = s.Desc.Options().(*descriptorpb.FieldOptions).GetDeprecated()
	message.Fields = append(message.Fields, v)
	return v
}

func extractMessage(file *FileExt, s *protogen.Message, isMap bool) *MessageExt {
	if s == nil {
		return nil
	}

	v := new(MessageExt)
	v.FilePath = s.Desc.ParentFile().Path()
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
	if !isMap {
		// 如果不是map字段的message, 则将其加至全局messages
		if rv, ok := file.Messages.Add(v.FullName, v); !ok {
			// 如果已经包含则跳过,否则会形成"环"
			return rv
		}
	}
	for _, f1 := range s.Fields {
		extractField(file, v, f1)
	}
	for _, s1 := range s.Enums {
		extractEnum(file, s1)
	}
	for _, s1 := range s.Messages {
		extractMessage(file, s1, false)
	}
	v.Desc = proto.GetExtension(s.Desc.Options(), E_Desc).(string)
	v.Plugin = proto.GetExtension(s.Desc.Options(), E_Plugin).(*Plugin)
	v.Deprecated = s.Desc.Options().(*descriptorpb.MessageOptions).GetDeprecated()
	return v

}

func extractMethod(file *FileExt, service *ServiceExt, s *protogen.Method) *MethodExt {
	if s == nil {
		return nil
	}
	v := new(MethodExt)
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	v.IsStreamingClient = s.Desc.IsStreamingClient()
	v.IsStreamingServer = s.Desc.IsStreamingServer()
	v.InputMessage = extractMessage(file, s.Input, false)
	v.OutputMessage = extractMessage(file, s.Output, false)
	v.Http = proto.GetExtension(s.Desc.Options(), E_Http).(*Http)
	v.Role = proto.GetExtension(s.Desc.Options(), E_Role).(*Role)
	v.Deprecated = s.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated()

	service.Methods = append(service.Methods, v)
	return v
}

func extractService(file *FileExt, s *protogen.Service) *ServiceExt {
	if s == nil {
		return nil
	}
	v := new(ServiceExt)
	v.FilePath = s.Desc.ParentFile().Path()
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	for _, m1 := range s.Methods {
		extractMethod(file, v, m1)
	}
	v.Tag = proto.GetExtension(s.Desc.Options(), E_Tag).(*Tag)
	v.HttpOnly = proto.GetExtension(s.Desc.Options(), E_HttpOnly).(bool)
	v.Deprecated = s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated()
	file.Services = append(file.Services, v)
	return v
}

type FieldExt struct {
	Name        string
	FullName    string
	GoName      string
	GoIdent     protogen.GoIdent
	Kind        protoreflect.Kind
	IsMap       bool
	IsRepeated  bool
	HasOptional bool
	Enum        *EnumExt
	Message     *MessageExt
	Prop        *Prop
	Rule        *Rule
	Deprecated  bool
}

type MessageExt struct {
	FilePath   string // 需要用来判断是否当前file
	Name       string
	FullName   string
	GoIdent    protogen.GoIdent
	Fields     []*FieldExt
	Desc       string
	Plugin     *Plugin
	Deprecated bool
}

type EnumValueExt struct {
	Name       string
	FullName   string
	GoIdent    protogen.GoIdent
	Number     int32
	Deprecated bool
}

type EnumExt struct {
	FilePath   string // 需要用来判断是否当前file
	Name       string
	FullName   string
	GoIdent    protogen.GoIdent
	Values     []*EnumValueExt
	Deprecated bool
}

type MethodExt struct {
	Name              string
	FullName          string
	GoName            string
	IsStreamingClient bool // client streaming
	IsStreamingServer bool // server streaming
	InputMessage      *MessageExt
	OutputMessage     *MessageExt
	Http              *Http
	Role              *Role
	Deprecated        bool
}

type ServiceExt struct {
	FilePath   string // 需要用来判断是否当前file
	Name       string
	FullName   string
	GoName     string
	Methods    []*MethodExt
	Tag        *Tag
	HttpOnly   bool
	Deprecated bool
}

type FileExt struct {
	Path         string
	Package      string
	GoPackage    string
	GoImportPath string
	Services     []*ServiceExt
	Messages     IdxVec[*MessageExt] /*展开成平面*/
	Enums        IdxVec[*EnumExt]    /*展开成平面*/
}

func (f *FieldExt) PropName() string {
	if f.Prop != nil && f.Prop.Name != `` {
		return f.Prop.Name
	}
	return f.Name
}

func (f *FieldExt) PropExplode() bool {
	if f.Prop != nil && f.Prop.Explode {
		return true
	}
	return false
}

func (f *FieldExt) RuleEnums() []any {
	if f.Rule != nil && f.Rule.Enum != nil {
		if f.Kind == protoreflect.StringKind {
			ret := make([]any, len(f.Rule.Enum.Str))
			for i, v := range f.Rule.Enum.Str {
				ret[i] = v
			}
			return ret
		} else {
			ret := make([]any, len(f.Rule.Enum.Int))
			for i, v := range f.Rule.Enum.Int {
				ret[i] = v
			}
			return ret
		}
	}
	return nil
}
