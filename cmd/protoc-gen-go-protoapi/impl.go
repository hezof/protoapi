package main

import (
	"encoding/json"
	"github.com/hezof/protoapi"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func generateFile(gen *protogen.Plugin, file *protogen.File) {

	// 1. 提取数据
	meta := extractFile(file)

	// 2. 生成实现文件: *_protoapi.pb.go
	generateImplFile(gen, file, meta)

	// 3. 生成文档文件: *_protoapi.json
	generateDocsFile(gen, file, meta)

	// 4. 生成代码文件: *_protoapi.code
	generateCodeFile(gen, file, meta)
}

// generateCodeFile 生成实现文件: *_protoapi.pb.go
func generateImplFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {

}

// generateDocsFile 生成文档文件: *_protoapi.json
func generateDocsFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.json`, file.GoImportPath)
	bs, _ := json.MarshalIndent(meta, "", "\t")
	g.P(string(bs))
}

// generateCodeFile 生成代码文件: *_protoapi.code
func generateCodeFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {

}

func extractFile(f *protogen.File) *FileExt {
	file := new(FileExt)
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

func extractEnum(f *FileExt, s *protogen.Enum) *EnumExt {
	if s == nil {
		return nil
	}
	v := new(EnumExt)
	v.File = f
	v.Name = string(s.Desc.Name())
	v.FullName = string(s.Desc.FullName())
	v.GoIdent = s.GoIdent
	f.Enums.Add(v.FullName, v)
	return v
}

func extractMessage(file *FileExt, s *protogen.Message) *MessageExt {
	if s == nil {
		return nil
	}

	v := new(MessageExt)
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

func extractField(file *FileExt, message *MessageExt, s *protogen.Field) *FieldExt {
	if s == nil {
		return nil
	}
	v := new(FieldExt)
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

func extractService(file *FileExt, s *protogen.Service) *ServiceExt {
	if s == nil {
		return nil
	}
	v := new(ServiceExt)
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

func extractMethod(file *FileExt, service *ServiceExt, s *protogen.Method) *MethodExt {
	if s == nil {
		return nil
	}
	v := new(MethodExt)
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
