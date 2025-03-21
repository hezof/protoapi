package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func implementServiceRegistry(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protoapi", GoImportPath: protogen.GoImportPath(protojsonImportPath)})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "io", GoImportPath: "io"})

	g.P("func ", service.GoName, "Registry(impl interface{}, aspects []protoapi.ServiceAspect) *protoapi.ServiceSetting {")
	g.P("ret := new(protoapi.ServiceSetting)")
	g.P("ret.Impl = impl")
	g.P("ret.Desc = &", service.GoName, "_ServiceDesc") // FIXBUG: protoc-gen-go-grpc@v1.1.0后变成public类型
	g.P("ret.HttpOnly = ", service.HttpOnly)
	g.P("ret.Aspects = aspects")
	for _, method := range service.Methods {
		g.P("ret.Methods = append(ret.Methods, &protoapi.MethodSetting{")
		g.P("Meta: protoapi.DecodeMeta(`", EncodeMeta(CreateMeta(file, service, method)), "`),")
		if method.Http != nil {
			switch {
			case !method.IsStreamingClient && !method.IsStreamingServer:
				simpleRpcCall(g, file, service, method)
			case method.IsStreamingClient && !method.IsStreamingServer:
				clientStreamingRpcCall(g, file, service, method)
			case !method.IsStreamingClient && method.IsStreamingServer:
				serverStreamingRpcCall(g, file, service, method)
			case method.IsStreamingClient && method.IsStreamingServer:
				bidirectionalStreamingRpcCall(g, file, service, method)
			}
		}
		g.P("})") // append
	}
	g.P("return ret")
	g.P("}") // func
}

func simpleRpcCall(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt, method *MethodExt) {

	var (
		body   []*FieldExt
		path   []*FieldExt
		query  []*FieldExt
		header []*FieldExt
		cookie []*FieldExt
	)

	for _, f := range method.InputMessage.Fields {
		switch f.In {
		case In_body:
			body = append(body, f)
		case In_path:
			path = append(path, f)
		case In_query:
			query = append(query, f)
		case In_header:
			header = append(header, f)
		case In_cookie:
			cookie = append(cookie, f)
		}
	}

	g.P("Call: func(pc *protoapi.Context, in io.Reader) (rsp interface{}, err error) {")
	g.P("set := pc.Handler.Setting")
	g.P("req := new(", g.QualifiedGoIdent(method.InputMessage.GoIdent), ")")
	// 1. 根据body解码request
	if len(body) > 0 {
		switch method.Http.Body {
		case Body_json:
			g.P("if err = protoapi.DecodeRequest(in, req); err != nil {")
			g.P("return nil, err")
			g.P("}")
		case Body_form:
			parseFormParams(g, body)
		case Body_omit:
			// nothing
		}
	}
	if len(path) > 0 {
		parsePathParams(g, path)
	}
	if len(query) > 0 {
		parseQueryParams(g, query)
	}
	if len(header) > 0 {
		parseHeaderParams(g, header)
	}
	if len(cookie) > 0 {
		parseCookieParams(g, cookie)
	}
	// 2. 执行aspects的before advice
	g.P("idx, ctx, err := protoapi.BeforeAspect(set, pc, req)")
	// 3. 执行service逻辑
	g.P("if err == nil {")
	g.P("rsp, err = set.Service.Impl.(", service.GoName, "Server).", method.GoName, "(ctx, req)")
	g.P("}") // if
	// 6. 返回response.
	g.P("return protoapi.AfterAspect(set, idx, ctx, req, rsp, err)")
	g.P("},") // call
}

func clientStreamingRpcCall(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt, method *MethodExt) {
	g.P("Call: func(pc *protoapi.Context, in io.Reader) (rsp interface{}, err error) {")
	// 1. 引用setting
	g.P("set := pc.Handler.Setting")
	// 2. 执行aspects的before advice. 请求与响应都是nil.
	g.P("idx, ctx, err := protoapi.BeforeAspect(set, pc, nil)")
	// 3. 执行service逻辑
	g.P("if err == nil {")
	g.P("err = set.Service.Impl.(", service.GoName, "Server).", method.GoName, "(protoapi.StreamServer[", g.QualifiedGoIdent(method.InputMessage.GoIdent), ",", g.QualifiedGoIdent(method.OutputMessage.GoIdent), "](pc))")
	g.P("}") // if
	// 4. 返回response.
	g.P("return protoapi.AfterAspect(set, idx, ctx, nil, nil, err)")
	g.P("},") // call
}

func serverStreamingRpcCall(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt, method *MethodExt) {
	var (
		body   []*FieldExt
		path   []*FieldExt
		query  []*FieldExt
		header []*FieldExt
		cookie []*FieldExt
	)

	for _, f := range method.InputMessage.Fields {
		switch f.In {
		case In_body:
			body = append(body, f)
		case In_path:
			path = append(path, f)
		case In_query:
			query = append(query, f)
		case In_header:
			header = append(header, f)
		case In_cookie:
			cookie = append(cookie, f)
		}
	}

	g.P("Call: func(pc *protoapi.Context, in io.Reader) (rsp interface{}, err error) {")
	g.P("set := pc.Handler.Setting")
	g.P("req := new(", g.QualifiedGoIdent(method.InputMessage.GoIdent), ")")
	// 1. 根据body解码request
	if len(body) > 0 {
		switch method.Http.Body {
		case Body_json:
			g.P("if err = protoapi.DecodeRequest(in, req); err != nil {")
			g.P("return nil, err")
			g.P("}")
		case Body_form:
			parseFormParams(g, body)
		case Body_omit:
			// nothing
		}
	}
	if len(path) > 0 {
		parsePathParams(g, path)
	}
	if len(query) > 0 {
		parseQueryParams(g, path)
	}
	if len(header) > 0 {
		parseHeaderParams(g, header)
	}
	if len(cookie) > 0 {
		parseCookieParams(g, cookie)
	}
	// 2. 执行aspects的before advice
	g.P("idx, ctx, err := protoapi.BeforeAspect(set, pc, req)")
	// 3. 执行service逻辑
	g.P("if err == nil {")
	g.P("err = set.Service.Impl.(", service.GoName, "Server).", method.GoName, "(req, protoapi.StreamServer[", g.QualifiedGoIdent(method.InputMessage.GoIdent), ",", g.QualifiedGoIdent(method.OutputMessage.GoIdent), "](pc))")
	g.P("}") // if
	// 6. 返回response.
	g.P("return protoapi.AfterAspect(set, idx, ctx, req, nil, err)")
	g.P("},") // call
}

func bidirectionalStreamingRpcCall(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt, method *MethodExt) {
	g.P("Call: func(pc *protoapi.Context, in io.Reader) (rsp interface{}, err error) {")
	// 1. 引用setting
	g.P("set := pc.Handler.Setting")
	// 2. 执行aspects的before advice. 请求与响应都是nil.
	g.P("idx, ctx, err := protoapi.BeforeAspect(set, pc, nil)")
	// 3. 执行service逻辑
	g.P("if err == nil {")
	g.P("err = set.Service.Impl.(", service.GoName, "Server).", method.GoName, "(protoapi.StreamServer[", g.QualifiedGoIdent(method.InputMessage.GoIdent), ",", g.QualifiedGoIdent(method.OutputMessage.GoIdent), "](pc))")
	g.P("}") // if
	// 4. 返回response.
	g.P("return protoapi.AfterAspect(set, idx, ctx, nil, nil, err)")
	g.P("},") // call
}

func parseFormParams(g *protogen.GeneratedFile, fields []*FieldExt) {
	for _, f := range fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBoolOptional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBoolRepeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBoolMap(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBool(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumNameOptional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumNameRepeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumNameMap(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnumName(pc.FormValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumOptional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumRepeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumMap(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ",", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnum(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt32Optional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt32Repeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt32Map(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt32(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint32Optional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint32Repeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint32Map(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint32(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt64Optional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt64Repeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt64Map(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt64(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint64Optional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint64Repeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint64Map(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint64(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamFloatOptional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamFloatRepeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamFloatMap(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamFloat(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamDoubleOptional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamDoubleRepeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamDoubleMap(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamDouble(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamStringOptional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamStringRepeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamStringMap(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamString(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBytesOptional(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBytesRepeated(pc.FormValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBytesMap(pc.FormValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBytes(pc.FormValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.MessageKind):
		case IsKind(f, protoreflect.GroupKind):
		}
	}
}

func parsePathParams(g *protogen.GeneratedFile, fields []*FieldExt) {
	for _, f := range fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBoolOptional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBoolRepeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBoolMap(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBool(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumNameOptional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumNameRepeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumNameMap(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnumName(pc.PathValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumOptional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumRepeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumMap(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ",", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnum(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt32Optional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt32Repeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt32Map(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt32(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint32Optional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint32Repeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint32Map(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint32(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt64Optional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt64Repeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt64Map(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt64(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint64Optional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint64Repeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint64Map(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint64(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamFloatOptional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamFloatRepeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamFloatMap(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamFloat(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamDoubleOptional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamDoubleRepeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamDoubleMap(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamDouble(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamStringOptional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamStringRepeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamStringMap(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamString(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBytesOptional(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBytesRepeated(pc.PathValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBytesMap(pc.PathValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBytes(pc.PathValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.MessageKind):
		case IsKind(f, protoreflect.GroupKind):
		}
	}
}

func parseQueryParams(g *protogen.GeneratedFile, fields []*FieldExt) {
	for _, f := range fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBoolOptional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBoolRepeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBoolMap(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBool(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumNameOptional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumNameRepeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumNameMap(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnumName(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumOptional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumRepeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumMap(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ",", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnum(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt32Optional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt32Repeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt32Map(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt32(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint32Optional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint32Repeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint32Map(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint32(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt64Optional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt64Repeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt64Map(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt64(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint64Optional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint64Repeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint64Map(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint64(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamFloatOptional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamFloatRepeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamFloatMap(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamFloat(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamDoubleOptional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamDoubleRepeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamDoubleMap(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamDouble(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamStringOptional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamStringRepeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamStringMap(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamString(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBytesOptional(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBytesRepeated(pc.QueryValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBytesMap(pc.QueryValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBytes(pc.QueryValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.MessageKind):
		case IsKind(f, protoreflect.GroupKind):
		}
	}
}

func parseHeaderParams(g *protogen.GeneratedFile, fields []*FieldExt) {
	for _, f := range fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBoolOptional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBoolRepeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBoolMap(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBool(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumNameOptional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumNameRepeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumNameMap(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnumName(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumOptional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumRepeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumMap(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ",", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnum(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt32Optional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt32Repeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt32Map(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt32(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint32Optional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint32Repeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint32Map(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint32(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt64Optional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt64Repeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt64Map(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt64(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint64Optional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint64Repeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint64Map(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint64(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamFloatOptional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamFloatRepeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamFloatMap(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamFloat(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamDoubleOptional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamDoubleRepeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamDoubleMap(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamDouble(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamStringOptional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamStringRepeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamStringMap(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamString(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBytesOptional(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBytesRepeated(pc.HeaderValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBytesMap(pc.HeaderValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBytes(pc.HeaderValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.MessageKind):
		case IsKind(f, protoreflect.GroupKind):
		}
	}
}

func parseCookieParams(g *protogen.GeneratedFile, fields []*FieldExt) {
	for _, f := range fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBoolOptional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBoolRepeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBoolMap(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBool(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumNameOptional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumNameRepeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumNameMap(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnumName(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err = protoapi.ParamEnumOptional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err = protoapi.ParamEnumRepeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err = protoapi.ParamEnumMap(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ",", f.Explode, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err = protoapi.ParamEnum(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt32Optional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt32Repeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt32Map(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt32(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint32Optional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint32Repeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint32Map(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint32(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamInt64Optional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamInt64Repeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamInt64Map(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamInt64(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamUint64Optional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamUint64Repeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamUint64Map(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamUint64(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamFloatOptional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamFloatRepeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamFloatMap(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamFloat(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamDoubleOptional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamDoubleRepeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamDoubleMap(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamDouble(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamStringOptional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamStringRepeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamStringMap(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamString(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err = protoapi.ParamBytesOptional(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err = protoapi.ParamBytesRepeated(pc.CookieValueRepeated, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err = protoapi.ParamBytesMap(pc.CookieValueMap, `", f.Name, "`, &req.", f.GoName, ", ", f.Explode, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err = protoapi.ParamBytes(pc.CookieValue, `", f.Name, "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.MessageKind):
		case IsKind(f, protoreflect.GroupKind):
		}
	}
}
