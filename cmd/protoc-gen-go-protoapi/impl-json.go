package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func implementFieldCodec(g *protogen.GeneratedFile, m *MessageExt) {
	implementDecodeField(g, m)
	implementEncodeField(g, m)
}

func implementDecodeField(g *protogen.GeneratedFile, m *MessageExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protojson", GoImportPath: protogen.GoImportPath(protojsonImportPath)})

	g.P("func (x *", g.QualifiedGoIdent(m.GoIdent), ") DecodeField(r *protojson.JsonDecoder, f string) {")
	g.P("switch f {")
	for _, f := range m.Fields {
		g.P("case `", f.Name, "`:")
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeBoolOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeBoolRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeBoolMap(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeBool(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.EnumName {
				switch {
				case f.HasOptional:
					g.P("protojson.DecodeEnumNameOptional(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				case f.IsRepeated:
					g.P("protojson.DecodeEnumNameRepeated(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				case f.IsMap:
					g.P("protojson.DecodeEnumNameMap(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				default:
					g.P("protojson.DecodeEnumName(r, &x.", f.GoName, ", ", enumGoName, "_value)")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("protojson.DecodeEnumOptional(r, &x.", f.GoName, ")")
				case f.IsRepeated:
					g.P("protojson.DecodeEnumRepeated(r, &x.", f.GoName, ")")
				case f.IsMap:
					g.P("protojson.DecodeEnumMap(r, &x.", f.GoName, ")")
				default:
					g.P("protojson.DecodeEnum(r, &x.", f.GoName, ")")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeInt32Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeInt32Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeInt32Map(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeInt32(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeUint32Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeUint32Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeUint32Map(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeUint32(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeInt64Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeInt64Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeInt64Map(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeInt64(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeUint64Optional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeUint64Repeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeUint64Map(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeUint64(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeFloatOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeFloatRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeFloatMap(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeFloat(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeDoubleOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeDoubleRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeDoubleMap(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeDouble(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeStringOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeStringRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeStringMap(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeString(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeBytesOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeBytesRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeBytesMap(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeBytes(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.MessageKind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeMessageOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeMessageRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeMessageMap(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeMessage(r, &x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.GroupKind):
			switch {
			case f.HasOptional:
				g.P("protojson.DecodeMessageOptional(r, &x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.DecodeMessageRepeated(r, &x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.DecodeMessageMap(r, &x.", f.GoName, ")")
			default:
				g.P("protojson.DecodeMessage(r, &x.", f.GoName, ")")
			}
		}
	}
	g.P("}") // switch
	g.P("}") // func
}

func implementEncodeField(g *protogen.GeneratedFile, m *MessageExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protojson", GoImportPath: protogen.GoImportPath(protojsonImportPath)})

	g.P("func (x *", g.QualifiedGoIdent(m.GoIdent), ") EncodeField(w *protojson.JsonEncoder) {")
	for _, f := range m.Fields {
		switch {
		case IsKind(f, protoreflect.BoolKind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeBoolOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeBoolRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeBoolMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeBool", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.EnumKind):
			enumGoName := g.QualifiedGoIdent(f.Enum.GoIdent)
			if f.EnumName {
				switch {
				case f.HasOptional:
					g.P("protojson.EncodeEnumNameOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ", ", enumGoName, "_name)")
				case f.IsRepeated:
					g.P("protojson.EncodeEnumNameRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ", ", enumGoName, "_name)")
				case f.IsMap:
					g.P("protojson.EncodeEnumNameMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ", ", enumGoName, "_name)")
				default:
					g.P("protojson.EncodeEnumName", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ", ", enumGoName, "_name)")
				}
			} else {
				switch {
				case f.HasOptional:
					g.P("protojson.EncodeEnumOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
				case f.IsRepeated:
					g.P("protojson.EncodeEnumRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
				case f.IsMap:
					g.P("protojson.EncodeEnumMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
				default:
					g.P("protojson.EncodeEnum", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
				}
			}
		case IsKind(f, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeInt32Optional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeInt32Repeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeInt32Map", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeInt32", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint32Kind, protoreflect.Fixed32Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeUint32Optional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeUint32Repeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeUint32Map", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeUint32", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeInt64Optional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeInt64Repeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeInt64Map", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeInt64", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.Uint64Kind, protoreflect.Fixed64Kind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeUint64Optional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeUint64Repeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeUint64Map", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeUint64", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.FloatKind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeFloatOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeFloatRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeFloatMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeFloat", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.DoubleKind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeDoubleOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeDoubleRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeDoubleMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeDouble", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.StringKind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeStringOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeStringRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeStringMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeString", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.BytesKind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeBytesOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeBytesRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeBytesMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeBytes", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.MessageKind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeMessageOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeMessageRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeMessageMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeMessage", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		case IsKind(f, protoreflect.GroupKind):
			switch {
			case f.HasOptional:
				g.P("protojson.EncodeMessageOptional", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsRepeated:
				g.P("protojson.EncodeMessageRepeated", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			case f.IsMap:
				g.P("protojson.EncodeMessageMap", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			default:
				g.P("protojson.EncodeMessage", propZeroSuffix(f), "(w,`", f.Name, "`, x.", f.GoName, ")")
			}
		}
	}
	g.P("}") // func
}

func propZeroSuffix(f *FieldExt) string {
	switch f.Zero {
	case Zero_omit_empty:
		return "_OmitEmpty"
	case Zero_with_empty:
		return "_WithEmpty"
	case Zero_conv_empty:
		return "_ConvEmpty"
	default:
		panic("invalid zero: " + f.Zero.String())
	}
}
