package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func implementServiceRegistry(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protoapi", GoImportPath: protogen.GoImportPath(protoapiImport)})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "io", GoImportPath: "io"})

	g.P("func ", service.GoName, "Registry(impl interface{}, aspects []protoapi.ServiceAspect) *protoapi.ServiceSetting {")
	g.P("ret := new(protoapi.ServiceSetting)")
	g.P("ret.Impl = impl")
	g.P("ret.Desc = &", service.GoName, "_ServiceDesc") // FIXBUG: protoc-gen-go-grpc@v1.1.0后变成public类型
	g.P("ret.HttpOnly = ", service.HttpOnly)
	g.P("ret.Aspects = aspects")
	for _, method := range service.Methods {
		g.P("ret.Methods = append(ret.Method, &protoapi.MethodSetting{")
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

			g.P("Call: func(ctx *protoapi.Context, in io.Reader) (rsp interface{}, err error) {")
			g.P("var set = ctx.Handler.Setting")
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
				parseFormParams(g, path)
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
			g.P("for _, asp := range set.Service.Aspects {")
			g.P("idx++")
			g.P("if ctx, err = asp.Before(set, ctx, req); err != nil {")
			g.P("goto __AFTER__")
			g.P("}")
			g.P("}")
			// 3. 执行message validator(如果有的话)
			g.P("if mv, ok := req.(protoapi.MessageValidator); ok {")
			g.P("if err = mv.Validate(set, ctx); err != nil {")
			g.P("goto __AFTER__")
			g.P("}")
			g.P("}")
			// 4. 执行service调用
			g.P("rsp, err = set.Service.Impl.(", service.GoName, ").", method.GoName, "(ctx, req)")
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
				g.P("if err := protoapi.ParamBoolOptional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(ctx.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(ctx.FormValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(ctx.FormValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(ctx.FormValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
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
				g.P("if err := protoapi.ParamBoolOptional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(ctx.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(ctx.PathValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(ctx.PathValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(ctx.PathValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
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
				g.P("if err := protoapi.ParamBoolOptional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(ctx.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(ctx.QueryValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(ctx.QueryValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(ctx.QueryValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
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
				g.P("if err := protoapi.ParamBoolOptional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(ctx.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(ctx.HeaderValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(ctx.HeaderValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(ctx.HeaderValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
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
				g.P("if err := protoapi.ParamBoolOptional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBoolRepeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBoolMap(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBool(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumNameOptional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumNameRepeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumNameMap(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value, ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnumName(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, ", ", enumGoName, "_value); err != nil {")
					g.P("return err")
					g.P("}")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("if err := protoapi.ParamEnumOptional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsRepeated:
					g.P("if err := protoapi.ParamEnumRepeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				case f.IsMap:
					g.P("if err := protoapi.ParamEnumMap(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ",", f.PropExplode(), "); err != nil {")
					g.P("return err")
					g.P("}")
				default:
					g.P("if err := protoapi.ParamEnum(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
					g.P("return err")
					g.P("}")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt32Optional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt32Repeated(ctx.FormRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt32Map(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt32(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint32Optional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint32Repeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint32Map(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint32(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamInt64Optional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamInt64Repeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamInt64Map(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamInt64(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamUint64Optional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamUint64Repeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamUint64Map(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamUint64(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamFloatOptional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamFloatRepeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamFloatMap(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamFloat(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamDoubleOptional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamDoubleRepeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamDoubleMap(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamDouble(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamStringOptional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamStringRepeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamStringMap(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamString(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("if err := protoapi.ParamBytesOptional(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsRepeated:
				g.P("if err := protoapi.ParamBytesRepeated(ctx.CookieValueRepeated, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			case f.IsMap:
				g.P("if err := protoapi.ParamBytesMap(ctx.CookieValueMap, `", f.PropName(), "`, &req.", f.GoName, ", ", f.PropExplode(), "); err != nil {")
				g.P("return err")
				g.P("}")
			default:
				g.P("if err := protoapi.ParamBytes(ctx.CookieValue, `", f.PropName(), "`, &req.", f.GoName, "); err != nil {")
				g.P("return err")
				g.P("}")
			}
		case IsKind(f, protoreflect.MessageKind):
		case IsKind(f, protoreflect.GroupKind):
		}
	}
}
