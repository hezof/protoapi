package main

import (
	"github.com/hezof/protoapi"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FieldExt struct {
	File       *FileExt
	Name       string
	FullName   string
	GoName     string
	GoIdent    protogen.GoIdent
	Kind       protoreflect.Kind
	IsMap      bool
	IsRepeated bool
	IsOptional bool
	Message    *MessageExt
	Prop       *protoapi.Prop
	Rule       *protoapi.Rule
}

type MessageExt struct {
	File     *FileExt
	Name     string
	FullName string
	GoIdent  protogen.GoIdent

	Fields IdxVec[*FieldExt]
	Schema *protoapi.Schema
	Plugin *protoapi.Plugin
}

type EnumExt struct {
	File     *FileExt
	Name     string
	FullName string
	GoIdent  protogen.GoIdent
}

type MethodExt struct {
	File              *FileExt
	Name              string
	FullName          string
	GoName            string
	IsStreamingClient bool // client streaming
	IsStreamingServer bool // server streaming
	InputMessage      *MessageExt
	OutputMessage     *MessageExt
	Http              *protoapi.Http
	Role              *protoapi.Role
}

type ServiceExt struct {
	File     *FileExt
	Name     string
	FullName string
	GoName   string
	Methods  IdxVec[*MethodExt]
	Tag      *protoapi.Tag
	HttpOnly bool
}

type FileExt struct {
	Path         string
	Package      string
	GoPackage    string
	GoImportPath string
	/*展开成平面*/
	Enums    IdxVec[*EnumExt]
	Messages IdxVec[*MessageExt]
	Services IdxVec[*ServiceExt]
}

type IdxVec[V any] struct {
	Idx map[string]V
	Vec []V
}

func (m IdxVec[V]) Add(k string, v V) (V, bool) {
	if m.Idx == nil {
		m.Idx = make(map[string]V)
	}
	if vl, ok := m.Idx[k]; ok {
		return vl, false
	}
	m.Idx[k] = v
	m.Vec = append(m.Vec, v)
	return v, true
}
