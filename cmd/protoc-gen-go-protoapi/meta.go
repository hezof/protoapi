package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FieldExt struct {
	Name       string
	FullName   string
	GoName     string
	GoIdent    protogen.GoIdent
	Kind       protoreflect.Kind
	IsMap      bool
	IsRepeated bool
	IsOptional bool
	Message    *MessageExt
	Prop       *Prop
	Rule       *Rule
}

type MessageExt struct {
	FilePath string // 需要用来判断是否当前file
	Name     string
	FullName string
	GoIdent  protogen.GoIdent
	Fields   IdxVec[*FieldExt]
	Schema   *Schema
	Plugin   *Plugin
}

type EnumExt struct {
	FilePath string // 需要用来判断是否当前file
	Name     string
	FullName string
	GoIdent  protogen.GoIdent
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
}

type ServiceExt struct {
	FilePath string // 需要用来判断是否当前file
	Name     string
	FullName string
	GoName   string
	Methods  IdxVec[*MethodExt]
	Tag      *Tag
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
	Idx map[string]V `json:"-,omitempty"`
	Vec []V
}

func (m *IdxVec[V]) Add(k string, v V) (V, bool) {
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
