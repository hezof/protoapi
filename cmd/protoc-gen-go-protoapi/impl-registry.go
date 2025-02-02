package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func implementServiceRegistry(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "context", GoImportPath: "context"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protoapi", GoImportPath: protogen.GoImportPath(protoapiImport)})
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

			var (
				body   []*FieldExt
				path   []*FieldExt
				query  []*FieldExt
				header []*FieldExt
				cookie []*FieldExt
			)

			for _, f := range method.InputMessage.Fields {
				if f.Prop == nil {
					body = append(body, f)
				} else {
					switch f.Prop.In {
					case Prop_body:
						body = append(body, f)
					case Prop_path:
						path = append(path, f)
					case Prop_query:
						query = append(query, f)
					case Prop_header:
						header = append(header, f)
					case Prop_cookie:
						cookie = append(cookie, f)
					}
				}
			}

			g.P("Call: func(pc *protoapi.Context, in io.Reader) (rsp interface{}, err error) {")
			g.P("var set = pc.Handler.Setting")
			g.P("var req = new(", g.QualifiedGoIdent(method.InputMessage.GoIdent), ")")
			// 1. 根据body解码request
			if len(body) > 0 {
				switch method.Http.Body {
				case Http_json:
					g.P("if err := protoapi.DecodeRequest(in, req); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case Http_form:
					parseFormParams(g, body)
				case Http_omit:
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
			g.P("var idx = -1")
			g.P("var ctx context.Context = pc")
			g.P("for _, asp := range set.Service.Aspects {")
			g.P("idx++")
			g.P("if ctx, err = asp.Before(set, ctx, req); err != nil {")
			g.P("goto __AFTER__")
			g.P("}")
			g.P("}")
			// 3. 执行message validator(如果有的话)
			g.P("if mv, ok := any(req).(protoapi.MessageValidator); ok {")
			g.P("if err = mv.Validate(set, ctx); err != nil {")
			g.P("goto __AFTER__")
			g.P("}")
			g.P("}")
			// 4. 执行service调用
			g.P("rsp, err = set.Service.Impl.(", service.GoName, "Server).", method.GoName, "(ctx, req)")
			// 5. 执行aspects的after advice
			g.P("__AFTER__:")
			g.P("for idx>=0 {")
			g.P("ctx, rsp, err = set.Service.Aspects[idx].After(set, ctx, req, rsp, err)")
			g.P("idx--")
			g.P("}")
			// 6. 返回response. rsp与err都是结果变量
			g.P("return")
			g.P("},") // call
		}
		g.P("})") // append
	}
	g.P("}") // func
}

func parseFormParams(g *protogen.GeneratedFile, fields []*FieldExt) {
	for _, f := range fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBoolOptional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(pc.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(pc.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(pc.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(pc.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
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
				g.P("if err := protoapi.ParamBoolOptional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(pc.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(pc.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(pc.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(pc.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
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
				g.P("if err := protoapi.ParamBoolOptional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(pc.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(pc.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(pc.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(pc.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
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
				g.P("if err := protoapi.ParamBoolOptional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(pc.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(pc.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(pc.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(pc.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
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
				g.P("if err := protoapi.ParamBoolOptional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return nil, err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(pc.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(pc.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(pc.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(pc.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return nil, err")
				g.P("}")
			}
		case IsKind(f, protoreflect.MessageKind):
		case IsKind(f, protoreflect.GroupKind):
		}
	}
}
