package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "v0.9.9" // 与发版相同

var requireUnimplemented *bool
var useGenericStreams *bool

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-protoapi %v\n", version)
		return
	}

	var flags flag.FlagSet
	requireUnimplemented = flags.Bool("require_unimplemented_servers", true, "set to false to match legacy behavior")
	useGenericStreams = flags.Bool("use_generic_streams_experimental", true, "set to true to use generic types for streaming client and server objects; this flag is EXPERIMENTAL and may be changed or removed in a future release")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) | uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		gen.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO2
		gen.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {

	// 1. 提取数据
	meta := extractFile(file)

	// 2. 生成实现文件: *_pb.go
	generateImplFile(gen, file, meta)

	// 3. 生成文档文件: *_json
	generateDocsFile(gen, file, meta)

	// 4. 生成代码文件: *_code
	generateCodeFile(gen, file, meta)
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
	v.FilePath = s.Desc.ParentFile().Path()
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
	v.FilePath = s.Desc.ParentFile().Path()
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
	v.Schema = proto.GetExtension(s.Desc.Options(), E_Schema).(*Schema)
	v.Plugin = proto.GetExtension(s.Desc.Options(), E_Plugin).(*Plugin)

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
	v.IsOptional = s.Desc.HasOptionalKeyword()
	v.Message = extractMessage(file, s.Message)
	v.Prop = proto.GetExtension(s.Desc.Options(), E_Prop).(*Prop)
	v.Rule = proto.GetExtension(s.Desc.Options(), E_Rule).(*Rule)
	message.Fields.Add(v.FullName, v)
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
	file.Services.Add(v.FullName, v)
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
	v.InputMessage = extractMessage(file, s.Input)
	v.OutputMessage = extractMessage(file, s.Output)
	v.Http = proto.GetExtension(s.Desc.Options(), E_Http).(*Http)
	v.Role = proto.GetExtension(s.Desc.Options(), E_Role).(*Role)
	service.Methods.Add(v.FullName, v)
	return v
}
