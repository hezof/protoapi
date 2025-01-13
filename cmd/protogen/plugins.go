package main

import "path/filepath"

var Plugins = []*Plugin{
	{
		Mode:    HttpGetProtoc,
		Name:    "protoc",
		Module:  "com/google/protobuf/protoc",
		Version: "v3.21.12", // debian 12仓库版本
	},
	{
		Mode:    GoGetSrc,
		Name:    "include",
		Module:  "github.com/hezof/protoapi/cmd/include",
		Version: "v0.5.0",
	},
	{
		Mode:    GoGetBin,
		Name:    "protoc-gen-go",
		Module:  "google.golang.org/protobuf/cmd/protoc-gen-go",
		Version: "v1.36.2",
	},
	{
		Mode:    GoGetBin,
		Name:    "protoc-gen-go-grpc",
		Module:  "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		Version: "v1.5.1",
	},
	{
		Mode:    GoGetBin,
		Name:    "protoc-gen-go-protoapi",
		Module:  "github.com/hezof/protoapi/cmd/protoc-gen-go-protoapi",
		Version: "v0.5.0",
	},
}

type Mode uint8

const (
	GoGetBin      Mode = 0
	GoGetSrc      Mode = 1
	HttpGetProtoc Mode = 2
)

type Plugin struct {
	Mode    Mode
	Name    string
	Module  string
	Version string
}

func (p *Plugin) FullName() string {
	fname := filepath.Base(p.Module) + `_` + p.Version[1:]
	if p.Mode != GoGetSrc {
		fname += goexe()
	}
	return fname
}

func FullName(module, version string, mode Mode) string {
	fname := filepath.Base(module) + `_` + version[1:]
	if mode != GoGetSrc {
		fname += goexe()
	}
	return fname
}
