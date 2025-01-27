package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func implementFieldCodec(gen *protogen.Plugin, g *protogen.GeneratedFile, m *MessageExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protoapi", GoImportPath: protogen.GoImportPath(protoapiImport)})

	g.P("func (x *", g.QualifiedGoIdent(m.GoIdent), ") DecodeField(r *protoapi.JsonDecoder, f string) {")
	g.P("switch f {")
	for _, f := range m.Fields {
		g.P("case `", f.PropName(), "`:")
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeBoolOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeBoolRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeBoolMap(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeBool(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("protoapi.DecodeEnumNameOptional(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				case f.IsRepeated:
					g.P("protoapi.DecodeEnumNameRepeated(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				case f.IsMap:
					g.P("protoapi.DecodeEnumNameMap(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				default:
					g.P("protoapi.DecodeEnumName(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("protoapi.DecodeEnumOptional(r, &x.", f.GoName, ")")
				case f.IsRepeated:
					g.P("protoapi.DecodeEnumRepeated(r, &x.", f.GoName, ")")
				case f.IsMap:
					g.P("protoapi.DecodeEnumMap(r, &x.", f.GoName, ")")
				default:
					g.P("protoapi.DecodeEnum(r, &x.", f.GoName, ")")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeInt32Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeInt32Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeInt32Map(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeInt32(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeUint32Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeUint32Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeUint32Map(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeUint32(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeInt64Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeInt64Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeInt64Map(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeInt64(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeUint64Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeUint64Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeUint64Map(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeUint64(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeFloatOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeFloatRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeFloatMap(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeFloat(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeDoubleOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeDoubleRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeDoubleMap(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeDouble(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeStringOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeStringRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeStringMap(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeString(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeBytesOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeBytesRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeBytesMap(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeBytes(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.MessageKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeMessageOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeMessageRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeMessageMap(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeMessage(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.GroupKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.DecodeMessageOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.DecodeMessageRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.DecodeMessageMap(r, &x.", f.GoName, ")")
			default:
				g.P("protoapi.DecodeMessage(r, &x.", f.GoName, ")")
			}
		}
	}
	g.P("}") // switch
	g.P("}") // func
	g.P("func (x *", g.QualifiedGoIdent(m.GoIdent), ") EncodeJSON(w *protoapi.JsonEncoder) {")
	for _, f := range m.Fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeBoolOptional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeBoolRepeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeBoolMap(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeBool(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.Prop != nil && f.Prop.EnumName {
				switch {
				case f.HasOptional:
					g.P("protoapi.EncodeEnumNameOptional(w, x.", f.GoName, ", ", enumGoName, "_name)")
				case f.IsRepeated:
					g.P("protoapi.EncodeEnumNameRepeated(w, x.", f.GoName, ", ", enumGoName, "_name)")
				case f.IsMap:
					g.P("protoapi.EncodeEnumNameMap(w, x.", f.GoName, ", ", enumGoName, "_name)")
				default:
					g.P("protoapi.EncodeEnumName(w, x.", f.GoName, ", ", enumGoName, "_name)")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("protoapi.EncodeEnumOptional(w, x.", f.GoName, ")")
				case f.IsRepeated:
					g.P("protoapi.EncodeEnumRepeated(w, x.", f.GoName, ")")
				case f.IsMap:
					g.P("protoapi.EncodeEnumMap(w, x.", f.GoName, ")")
				default:
					g.P("protoapi.EncodeEnum(w, x.", f.GoName, ")")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeInt32Optional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeInt32Repeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeInt32Map(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeInt32(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeUint32Optional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeUint32Repeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeUint32Map(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeUint32(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeInt64Optional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeInt64Repeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeInt64Map(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeInt64(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeUint64Optional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeUint64Repeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeUint64Map(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeUint64(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeFloatOptional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeFloatRepeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeFloatMap(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeFloat(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeDoubleOptional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeDoubleRepeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeDoubleMap(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeDouble(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeStringOptional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeStringRepeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeStringMap(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeString(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeBytesOptional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeBytesRepeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeBytesMap(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeBytes(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.MessageKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeMessageOptional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeMessageRepeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeMessageMap(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeMessage(w, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.GroupKind):
			switch {
			case f.HasOptional:
				g.P("protoapi.EncodeMessageOptional(w, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protoapi.EncodeMessageRepeated(w, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protoapi.EncodeMessageMap(w, x.", f.GoName, ")")
			default:
				g.P("protoapi.EncodeMessage(w, x.", f.GoName, ")")
			}
		}
	}
	g.P("}") // func
}
