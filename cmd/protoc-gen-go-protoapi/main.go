package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "v0.9.9" // 与发版相同

var requireUnimplemented bool
var useGenericStreams bool
var protoapiModule string

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-protoapi %v\n", version)
		return
	}

	var flags flag.FlagSet
	flags.BoolVar(&requireUnimplemented, "require_unimplemented_servers", true, "set to false to match legacy behavior")
	flags.BoolVar(&useGenericStreams, "use_generic_streams_experimental", true, "set to true to use generic types for streaming client and server objects; this flag is EXPERIMENTAL and may be changed or removed in a future release")
	flags.StringVar(&protoapiModule, "protoapi_module", "github.com/hezof/protoapi", "protoapi module")

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

func generateFile(gen *protogen.Plugin, src *protogen.File) {

	// 1. 提取数据
	file := extractFile(gen, src)
	if file == nil {
		return
	}

	// 2. 生成实现文件: *_pb.go
	generateImplFile(gen, src, file)

	// 3. 生成文档文件: *_json
	generateDocsFile(gen, src, file)

	// 4. 生成代码文件: *_code
	generateCodeFile(gen, src, file)
}
