package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
)

// generateCodeFile 生成代码文件: *_protoapi.code
func generateCodeFile(gen *protogen.Plugin, file *protogen.File, meta *FileExt) {
	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+`_protoapi.code`, file.GoImportPath)

	qualifiedGoIdent := genQualifiedGoIdentFunc(file)

	streaming := false
	for _, ps := range meta.Services {
		for _, pm := range ps.Methods {
			if pm.IsStreamingClient || pm.IsStreamingServer {
				streaming = true
				break
			}
		}
	}

	g.P(`import (`)
	g.P(`	"context"`)
	if streaming {
		g.P(`	"google.golang.org/grpc"`)
	}
	g.P(`)`)
	for _, ps := range meta.Services {
		g.P()
		g.P(serviceTitle(ps))
		if *requireUnimplemented {
			g.P(`type `, ps.GoName, "Implement struct {")
			g.P(`    *`, meta.GoPackage, `.`, `Unimplemented`, ps.GoName, "Server")
			g.P(`}`)
		} else {
			g.P(`type `, ps.GoName, "Implement struct {}")
		}
		g.P()
		g.P(`var _ `, meta.GoPackage, `.`, ps.GoName, `Server = (*`, ps.GoName, "Implement)(nil)")
		for _, pm := range ps.Methods {
			g.P()
			g.P(methodTitle(pm))
			if pm.Http != nil {
				if pm.Http.Get != "" {
					g.P(`// GET `, pm.Http.Get, sse(pm))
				}
				if pm.Http.Put != "" {
					g.P(`// PUT `, pm.Http.Put, sse(pm))
				}
				if pm.Http.Post != "" {
					g.P(`// POST `, pm.Http.Post, sse(pm))
				}
				if pm.Http.Delete != "" {
					g.P(`// DELETE `, pm.Http.Delete, sse(pm))
				}
				if pm.Http.Options != "" {
					g.P(`// OPTIONS `, pm.Http.Options, sse(pm))
				}
				if pm.Http.Head != "" {
					g.P(`// HEAD `, pm.Http.Head, sse(pm))
				}
				if pm.Http.Patch != "" {
					g.P(`// PATCH `, pm.Http.Patch, sse(pm))
				}
				if pm.Http.Trace != "" {
					g.P(`// TRACE `, pm.Http.Trace, sse(pm))
				}
				if pm.Http.Connect != "" {
					g.P(`// CONNECT `, pm.Http.Connect, sse(pm))
				}
				if pm.Http.Websocket != "" {
					g.P(`// WEBSOCKET `, pm.Http.Websocket) // WEBSOCKET会覆盖SSE
				}
			}
			// streaming泛型接口要求grpc v1.64.0+ 及protoc-gen-go-grpc v1.4.0+
			switch {
			case !pm.IsStreamingClient && !pm.IsStreamingServer:
				g.P(`func (ps *`, ps.GoName, "Implement) ", pm.GoName, `(ctx context.Context, req *`, qualifiedGoIdent(pm.InputMessage.GoIdent), `) (rsp *`, qualifiedGoIdent(pm.OutputMessage.GoIdent), `, err error) {`)
				g.P(`    return`)
				g.P(`}`)
			case pm.IsStreamingClient && !pm.IsStreamingServer:
				g.P(`func (ps *`, ps.GoName, "Implement) ", pm.GoName, `(ctx context.Context, svr grpc.ClientStreamingServer[`, qualifiedGoIdent(pm.InputMessage.GoIdent), `,`, qualifiedGoIdent(pm.OutputMessage.GoIdent), `]) (err error) {`)
				g.P(`    return`)
				g.P(`}`)
			case !pm.IsStreamingClient && pm.IsStreamingServer:
				g.P(`func (ps *`, ps.GoName, "Implement) ", pm.GoName, `(ctx context.Context, req *`, qualifiedGoIdent(pm.InputMessage.GoIdent), `, svr grpc.ServerStreamingServer[`, qualifiedGoIdent(pm.OutputMessage.GoIdent), `]) (err error) {`)
				g.P(`    return`)
				g.P(`}`)
			case pm.IsStreamingClient && pm.IsStreamingServer:
				g.P(`func (ps *`, ps.GoName, "Implement) ", pm.GoName, `(ctx context.Context, svr grpc.BidiStreamingServer[`, qualifiedGoIdent(pm.InputMessage.GoIdent), `,`, qualifiedGoIdent(pm.OutputMessage.GoIdent), `]) (err error) {`)
				g.P(`    return`)
				g.P(`}`)
			}
		}
	}
	g.P()

}

// 直接调用原生GeneratedFile.QualifiedGoIdent()会导致"unused imports"问题, 此处使用plugin复制手段确保不会加到
func genQualifiedGoIdentFunc(file *protogen.File) func(ident protogen.GoIdent) string {

	g := new(protogen.Plugin).NewGeneratedFile("", file.GoImportPath)
	return func(ident protogen.GoIdent) string {
		if ident.GoImportPath == file.GoImportPath {
			return string(file.GoPackageName) + "." + ident.GoName
		} else {
			return g.QualifiedGoIdent(ident)
		}
	}
}

func serviceTitle(ps *ServiceExt) string {
	if ps.Tag == nil {
		return fmt.Sprintf(`// %v %v.`, ps.GoName, ps.FullName)
	} else {
		return fmt.Sprintf(`// %v %v. %v`, ps.GoName, ps.FullName, ps.Tag.Desc)
	}
}

func methodTitle(pm *MethodExt) string {
	if pm.Http == nil {
		return fmt.Sprintf(`// %v %v.`, pm.GoName, pm.FullName)
	} else {
		return fmt.Sprintf(`// %v %v. %v`, pm.GoName, pm.FullName, pm.Http.Desc)
	}
}

func sse(pm *MethodExt) string {
	if pm.Http != nil && pm.Http.Result == Http_events {
		return ` [SSE]`
	}
	return ``
}
