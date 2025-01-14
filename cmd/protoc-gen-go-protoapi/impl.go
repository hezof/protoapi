package main

import (
	"encoding/json"
	"fmt"
	"github.com/hezof/protoapi/cmd/protoc-gen-go-protoapi/protoapi"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func generateFile(g *protogen.Plugin, f *protogen.File) {
	file := extractFile(f)
	bs, _ := json.MarshalIndent(file, "", "	")
	fmt.Println(string(bs))
}

func extractFile(f *protogen.File) *File {
	file := new(File)
	file.Path = f.Desc.Path()
	file.Package = string(f.Desc.Package())
	file.GoPackage = string(f.GoPackageName)
	file.GoImportPath = string(f.GoImportPath)

	for _, s := range f.Enums {
		extractEnum(file, s)
	}

	for _, s := range f.Messages {
		extractMessage(file, s)
	}

	for _, s := range f.Services {
		extractService(file, s)
	}

	return file
}

func extractEnum(f *File, s *protogen.Enum) {
	v := new(Enum)
	v.File = f
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
	f.Enums.Add(v.FullName, v)
}

func extractMessage(file *File, s *protogen.Message) *Message {
	if s == nil {
		return nil
	}

	v := new(Message)
	v.File = file
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
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
		extractMessage(file, s1)
	}
	v.Schema = proto.GetExtension(s.Desc.Options(), protoapi.E_Schema).(*protoapi.Schema)
	v.Plugin = proto.GetExtension(s.Desc.Options(), protoapi.E_Plugin).(*protoapi.Plugin)

	return v

}

func extractField(file *File, message *Message, s *protogen.Field) *Field {
	if s == nil {
		return nil
	}
	v := new(Field)
	v.File = file
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	v.GoIdent = s.GoIdent
	v.Kind = s.Desc.Kind()
	v.IsMap = s.Desc.IsMap()
	v.IsRepeated = s.Desc.IsList()
	v.IsOptional = s.Desc.HasOptionalKeyword()
	v.Message = extractMessage(file, s.Message)
	v.Prop = proto.GetExtension(s.Desc.Options(), protoapi.E_Prop).(*protoapi.Prop)
	v.Rule = proto.GetExtension(s.Desc.Options(), protoapi.E_Rule).(*protoapi.Rule)
	message.Fields.Add(v.FullName, v)
	return v
}

func extractService(file *File, s *protogen.Service) *Service {
	v := new(Service)
	v.File = file
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	for _, m1 := range s.Methods {
		extractMethod(file, v, m1)
	}
	v.Tag = proto.GetExtension(s.Desc.Options(), protoapi.E_Tag).(*protoapi.Tag)
	v.HttpOnly = proto.GetExtension(s.Desc.Options(), protoapi.E_HttpOnly).(bool)
	file.Services.Add(v.FullName, v)
	return v
}

func extractMethod(file *File, service *Service, s *protogen.Method) *Method {
	v := new(Method)
	v.File = file
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoName = s.GoName
	v.IsStreamingClient = s.Desc.IsStreamingClient()
	v.IsStreamingServer = s.Desc.IsStreamingServer()
	v.InputMessage = extractMessage(file, s.Input)
	v.OutputMessage = extractMessage(file, s.Output)
	v.Http = proto.GetExtension(s.Desc.Options(), protoapi.E_Http).(*protoapi.Http)
	v.Role = proto.GetExtension(s.Desc.Options(), protoapi.E_Role).(*protoapi.Role)
	service.Methods.Add(v.FullName, v)
	return v
}
