package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func implementServiceRegistry(g *protogen.GeneratedFile, file *FileExt, service *ServiceExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protoapi", GoImportPath: ProtoapiModule})
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

			g.P("Call: func(ctx *protoapi.Context, in io.Reader) (interface{}, error) {")
			g.P("req := new(", g.QualifiedGoIdent(method.InputMessage.GoIdent), ")")
			// 1. 根据body解码request
			if len(body) > 0 {
				switch method.Http.Body {
				case Http_json:
					g.P("if err := protoapi.JsonBody(in, req); err != nil {")
					g.P("return nil, err")
					g.P("}")
				case Http_form:
					for _, f := range body {
						switch {
						case IsKind(f, protoreflect.BoolKind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.ParamBoolOptional(ctx.FormBool, `", fname(f), "`, &req.", f.GoName, "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.ParamBoolRepeated(ctx.FormBoolSlice, `", fname(f), "`, &req.", f.GoName, ", ); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathBoolMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.ParamBool(ctx.FormBool, `", fname(f), "`, &req.", f.GoName, "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.EnumKind):
							enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
							if f.Prop != nil && f.Prop.EnumName {
								switch {
								case f.HasOptional:
									g.P("if err := protoapi.PathEnumNameOptional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_value); err != nil {")
									g.P("return err")
									g.P("}")
								case f.IsRepeated:
									g.P("if err := protoapi.PathEnumNameRepeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_value); err != nil {")
									g.P("return err")
									g.P("}")
								case f.IsMap:
									g.P("if err := protoapi.PathEnumNameMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_value); err != nil {")
									g.P("return err")
									g.P("}")
								default:
									g.P("if err := protoapi.PathEnumName(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_value); err != nil {")
									g.P("return err")
									g.P("}")
								}
							} else {
								switch {
								case f.HasOptional:
									g.P("if err := protoapi.PathEnumOptional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_name); err != nil {")
									g.P("return err")
									g.P("}")
								case f.IsRepeated:
									g.P("if err := protoapi.PathEnumRepeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_name); err != nil {")
									g.P("return err")
									g.P("}")
								case f.IsMap:
									g.P("if err := protoapi.PathEnumMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_name); err != nil {")
									g.P("return err")
									g.P("}")
								default:
									g.P("if err := protoapi.PathEnum(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), enumGoName, "_name); err != nil {")
									g.P("return err")
									g.P("}")
								}
							}
						case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathInt32Optional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathInt32Repeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathInt32Map(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathInt32(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathUint32Optional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathUint32Repeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathUint32Map(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathUint32(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathInt64Optional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathInt64Repeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathInt64Map(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathInt64(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathUint64Optional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathUint64Repeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathUint64Map(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathUint64(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.FloatKind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathFloatOptional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathFloatRepeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathFloatMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathFloat(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.DoubleKind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathDoubleOptional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathDoubleRepeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathDoubleMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathDouble(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.StringKind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathStringOptional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathStringRepeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathStringMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathString(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.BytesKind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathBytesOptional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathBytesRepeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathBytesMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathBytes(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.MessageKind):
							switch {
							case f.HasOptional:
								g.P("if err := protoapi.PathMessageOptional(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsRepeated:
								g.P("if err := protoapi.PathBytesRepeated(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							case f.IsMap:
								g.P("if err := protoapi.PathBytesMap(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							default:
								g.P("if err := protoapi.PathBytes(ctx, `", fname(f), "`, &req.", f.GoName, ", ", style(f), "); err != nil {")
								g.P("return err")
								g.P("}")
							}
						case IsKind(f, protoreflect.GroupKind):
						}
					}
				case Http_omit:
				}
			}
			if len(path) > 0 {

			}
			if len(query) > 0 {

			}
			if len(header) > 0 {

			}
			if len(cookie) > 0 {

			}
			// 2. 执行aspects的before advice

			// 3. 执行message validator(如果有的话)

			// 4. 执行service调用

			// 5. 执行aspects的after advice

			// 6. 返回response

			g.P("},") // call
		}
		g.P("})") // append
	}
	g.P("}")
}
