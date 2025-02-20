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
			if v.IsMap {
				// key必须是string
				if v.Message.Fields[0].Kind != protoreflect.StringKind {
					gen.Error(fmt.Errorf("invalid map key type %v", v.Message.Fields[0].Kind))
					return nil
				}
				// val需换成
				v.Kind = v.Message.Fields[1].Kind
				v.Enum = v.Message.Fields[1].Enum
				v.Message = v.Message.Fields[1].Message
			}
		}
	}

	for _, s := range f.Services {
		extractService(file, s)
	}

	return file
}

func extractEnumValue(file *FileExt, e *EnumExt, s *protogen.EnumValue) *EnumValueExt {
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

func extractEnum(file *FileExt, s *protogen.Enum) *EnumExt {
	if s == nil {
		return nil
	}
	v := new(EnumExt)
	v.Local = file.Path == s.Desc.ParentFile().Path()
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
	if rv, ok := file.Enums.Add(v.FullName, v); !ok {
		return rv
	}
	for _, s1 := range s.Values {
		extractEnumValue(file, v, s1)
	}
	v.Deprecated = s.Desc.Options().(*descriptorpb.EnumOptions).GetDeprecated()
	return v
}

func extractField(file *FileExt, message *MessageExt, s *protogen.Field) *FieldExt {
	if s == nil {
		return nil
	}
	v := new(FieldExt)
	v.Name = NvlS(proto.GetExtension(s.Desc.Options(), E_Name).(string), string(s.Desc.Name()))
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	v.GoIdent = s.GoIdent
	v.IsMap = s.Desc.IsMap()
	v.IsRepeated = s.Desc.IsList()
	v.HasOptional = s.Desc.HasOptionalKeyword()
	v.Kind = s.Desc.Kind()
	v.Enum = extractEnum(file, s.Enum)
	v.Message = extractMessage(file, s.Message, v.IsMap)
	v.Desc = proto.GetExtension(s.Desc.Options(), E_Desc).(string)
	v.Zero = proto.GetExtension(s.Desc.Options(), E_Zero).(Zero)
	v.EnumName = proto.GetExtension(s.Desc.Options(), E_EnumName).(bool)
	v.In = proto.GetExtension(s.Desc.Options(), E_In).(In)
	v.Explode = proto.GetExtension(s.Desc.Options(), E_Explode).(bool)
	v.Deprecated = s.Desc.Options().(*descriptorpb.FieldOptions).GetDeprecated()
	v.Rule = &Meta_Rule{
		Name:      v.Name,
		Required:  proto.GetExtension(s.Desc.Options(), E_Required).(*Error),
		Maximum:   proto.GetExtension(s.Desc.Options(), E_Maximum).(*Maximum),
		Minimum:   proto.GetExtension(s.Desc.Options(), E_Minimum).(*Minimum),
		MinLength: proto.GetExtension(s.Desc.Options(), E_MinLength).(*MinLength),
		MaxLength: proto.GetExtension(s.Desc.Options(), E_MaxLength).(*MaxLength),
		MinItems:  proto.GetExtension(s.Desc.Options(), E_MinItems).(*MinItems),
		MaxItems:  proto.GetExtension(s.Desc.Options(), E_MaxItems).(*MaxItems),
		Enum:      proto.GetExtension(s.Desc.Options(), E_Enum).(*Enum),
		Pattern:   proto.GetExtension(s.Desc.Options(), E_Pattern).(*Pattern),
		Plugin:    proto.GetExtension(s.Desc.Options(), E_Plugin).(*Plugin),
	}
	message.Fields = append(message.Fields, v)
	return v
}

func extractMessage(file *FileExt, s *protogen.Message, isMap bool) *MessageExt {
	if s == nil {
		return nil
	}

	v := new(MessageExt)
	v.Local = !isMap && file.Path == s.Desc.ParentFile().Path()
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
	// 如果不是map字段的message, 则将其加至全局messages
	if rv, ok := file.Messages.Add(v.FullName, v); !ok {
		// 如果已经包含则跳过,否则会形成"环"
		return rv
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
	v.Schema = proto.GetExtension(s.Desc.Options(), E_Schema).(string)
	v.Extend = proto.GetExtension(s.Desc.Options(), E_Extend).(string)
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
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	for _, m1 := range s.Methods {
		extractMethod(file, v, m1)
	}
	v.Info = proto.GetExtension(s.Desc.Options(), E_Info).(*Info)
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
	Desc        string
	Zero        Zero
	EnumName    bool
	In          In
	Explode     bool
	Deprecated  bool
	Rule        *Meta_Rule
}

type MessageExt struct {
	Local      bool
	Name       string
	FullName   string
	GoIdent    protogen.GoIdent
	Fields     []*FieldExt
	Schema     string
	Extend     string
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
	Local      bool
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
	Name       string
	FullName   string
	GoName     string
	Methods    []*MethodExt
	Info       *Info
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
