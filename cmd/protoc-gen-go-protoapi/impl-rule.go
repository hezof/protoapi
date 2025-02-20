package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func shouldMessageValidator(m *MessageExt) bool {
	if m.Extend != "" {
		return true
	}
	for _, f := range m.Fields {
		if r := f.Rule; r != nil {
			if r.Required != nil {
				return true
			}
			if r.Minimum != nil {
				return true
			}
			if r.Maximum != nil {
				return true
			}
			if r.MinLength != nil {
				return true
			}
			if r.MaxLength != nil {
				return true
			}
			if r.MinItems != nil {
				return true
			}
			if r.MaxItems != nil {
				return true
			}
			if r.Enum != nil {
				return true
			}
			if r.Pattern != nil {
				return true
			}
			if r.Plugin != nil {
				return true
			}
		}
	}
	return false
}

func implementMessageValidator(g *protogen.GeneratedFile, m *MessageExt) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "protoapi", GoImportPath: protogen.GoImportPath(protoapiImport)})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "context", GoImportPath: "context"})
	g.P("func (x *", g.QualifiedGoIdent(m.GoIdent), ") Validate(set *protoapi.MethodSetting, ctx context.Context) error {")
	if m.Extend != "" {
		g.P("if err:=set.MessageExtend(ctx, x); err != nil {")
		g.P("return err")
		g.P("}")
	}
	for i, f := range m.Fields {
		if r := f.Rule; r != nil {
			if r.Required != nil {
				validateRequired(g, i, f)
			}
			if r.Minimum != nil {
				validateMinimum(g, i, f)
			}
			if r.Maximum != nil {
				validateMaximum(g, i, f)
			}
			if r.MinLength != nil {
				validateMinLength(g, i, f)
			}
			if r.MaxLength != nil {
				validateMaxLength(g, i, f)
			}
			if r.MinItems != nil {
				validateMinItems(g, i, f)
			}
			if r.MaxItems != nil {
				validateMaxItems(g, i, f)
			}
			if r.Enum != nil {
				validateEnum(g, i, f)
			}
			if r.Pattern != nil {
				validatePattern(g, i, f)
			}
			if r.Plugin != nil {
				g.P("if err:=set.FieldPlugins[", i, "](ctx, `", f.Name, "`, x.", f.GoName, ", set.Meta.FieldRules[", i, "].Plugin); err != nil {")
				g.P("return err")
				g.P("}")
			}
		}
	}
	g.P("return nil")
	g.P("}")
}

func validateRequired(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		required的定义是各值非0值.
	*/
	switch {
	case f.IsRepeated:
		fallthrough
	case f.IsMap:
		g.P("if x.", f.GoName, " == nil {")
		g.P("return set.Meta.FieldRules[", i, "].Required")
		g.P("}")
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
			if f.HasOptional {
				g.P("if x.", f.GoName, " == nil {")
				g.P("return set.Meta.FieldRules[", i, "].Required")
				g.P("}")
			}
		case protoreflect.EnumKind:
			if f.HasOptional {
				g.P("if x.", f.GoName, " == nil {")
				g.P("return set.Meta.FieldRules[", i, "].Required")
				g.P("}")
			}
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.DoubleKind:
			if f.HasOptional {
				g.P("if x.", f.GoName, " == nil {")
				g.P("return set.Meta.FieldRules[", i, "].Required")
				g.P("}")
			}
		case protoreflect.StringKind:
			if f.HasOptional {
				g.P("if x.", f.GoName, " == nil {")
				g.P("return set.Meta.FieldRules[", i, "].Required")
				g.P("}")
			} else {
				// 严格地讲, 空串是有值的. 但习惯性将空串当成无值看待!
				g.P("if x.", f.GoName, " == `` {")
				g.P("return set.Meta.FieldRules[", i, "].Required")
				g.P("}")
			}
		case protoreflect.BytesKind:
			g.P("if x.", f.GoName, " == nil {")
			g.P("return set.Meta.FieldRules[", i, "].Required")
			g.P("}")
		case protoreflect.MessageKind:
			fallthrough
		case protoreflect.GroupKind:
			g.P("if x.", f.GoName, " == nil {")
			g.P("return set.Meta.FieldRules[", i, "].Required")
			g.P("}")
		}
	}
}

func validateMinimum(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		minimum针对enum/number类型. 其他类型自动忽略
	*/
	switch {
	case f.IsRepeated:
	case f.IsMap:
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
			fallthrough
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.DoubleKind:
			if f.HasOptional {
				// 只比较非nil值, 否则nil算是上界, 还是下界? 扯开了难有定论
				g.P("if x.", f.GoName, " != nil {")
				if f.Rule.Minimum.Exclusive {
					g.P("if *x.", f.GoName, " <= ", f.Rule.Minimum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Minimum.Err")
					g.P("}")
				} else {
					g.P("if *x.", f.GoName, " < ", f.Rule.Minimum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Minimum.Err")
					g.P("}")
				}
				g.P("}")
			} else {
				if f.Rule.Minimum.Exclusive {
					g.P("if x.", f.GoName, " <= ", f.Rule.Minimum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Minimum.Err")
					g.P("}")
				} else {
					g.P("if x.", f.GoName, " < ", f.Rule.Minimum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Minimum.Err")
					g.P("}")
				}
			}
		case protoreflect.StringKind:
		case protoreflect.BytesKind:
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}

func validateMaximum(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		maximum针对enum/number类型. 其他类型自动忽略
	*/
	switch {
	case f.IsRepeated:
	case f.IsMap:
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
			fallthrough
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.DoubleKind:
			if f.HasOptional {
				// 只比较非nil值, 否则nil算是上界, 还是下界? 扯开了难有定论
				g.P("if x.", f.GoName, " != nil {")
				if f.Rule.Maximum.Exclusive {
					g.P("if *x.", f.GoName, " >= ", f.Rule.Maximum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Maximum.Err")
					g.P("}")
				} else {
					g.P("if *x.", f.GoName, " > ", f.Rule.Maximum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Maximum.Err")
					g.P("}")
				}
				g.P("}")
			} else {
				if f.Rule.Maximum.Exclusive {
					g.P("if x.", f.GoName, " >= ", f.Rule.Maximum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Maximum.Err")
					g.P("}")
				} else {
					g.P("if x.", f.GoName, " > ", f.Rule.Maximum.Val, " {")
					g.P("return set.Meta.FieldRules[", i, "].Maximum.Err")
					g.P("}")
				}
			}
		case protoreflect.StringKind:
		case protoreflect.BytesKind:
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}

func validateMinLength(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		min_length针对string/bytes类型. 其他类型自动忽略
	*/
	switch {

	case f.IsRepeated:
	case f.IsMap:
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.DoubleKind:
		case protoreflect.StringKind:
			if f.HasOptional {
				// 只比较非nil值, 否则nil算是上界, 还是下界? 扯开了难有定论
				g.P("if x.", f.GoName, " != nil {")
				g.P("if len(*x.", f.GoName, ") < ", f.Rule.MinLength.Val, " {")
				g.P("return set.Meta.FieldRules[", i, "].MinLength.Err")
				g.P("}")
				g.P("}")
			} else {
				g.P("if len(x.", f.GoName, ") < ", f.Rule.MinLength.Val, " {")
				g.P("return set.Meta.FieldRules[", i, "].MinLength.Err")
				g.P("}")
			}
		case protoreflect.BytesKind:
			g.P("if x.", f.GoName, " != nil {")
			g.P("if len(x.", f.GoName, ") < ", f.Rule.MinLength.Val, " {")
			g.P("return set.Meta.FieldRules[", i, "].MinLength.Err")
			g.P("}")
			g.P("}")
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}

func validateMaxLength(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		max_length针对string/bytes类型. 其他类型自动忽略
	*/
	switch {
	case f.IsRepeated:
	case f.IsMap:
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.DoubleKind:
		case protoreflect.StringKind:
			if f.HasOptional {
				// 只比较非nil值, 否则nil算是上界, 还是下界? 扯开了难有定论
				g.P("if x.", f.GoName, " != nil {")
				g.P("if len(*x.", f.GoName, ") > ", f.Rule.MaxLength.Val, " {")
				g.P("return set.Meta.FieldRules[", i, "].MaxLength.Err")
				g.P("}")
				g.P("}")
			} else {
				g.P("if len(x.", f.GoName, ") > ", f.Rule.MaxLength.Val, " {")
				g.P("return set.Meta.FieldRules[", i, "].MaxLength.Err")
				g.P("}")
			}
		case protoreflect.BytesKind:
			g.P("if x.", f.GoName, " != nil {")
			g.P("if len(x.", f.GoName, ") > ", f.Rule.MaxLength.Val, " {")
			g.P("return set.Meta.FieldRules[", i, "].MaxLength.Err")
			g.P("}")
			g.P("}")
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}

func validateMinItems(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		min_items针对repeated/map类型. 其他类型自动忽略
	*/
	switch {
	case f.IsRepeated:
		fallthrough
	case f.IsMap:
		g.P("if len(x.", f.GoName, ") < ", f.Rule.MinItems.Val, " {")
		g.P("return set.Meta.FieldRules[", i, "].MinItems.Err")
		g.P("}")
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.DoubleKind:
		case protoreflect.StringKind:
		case protoreflect.BytesKind:
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}

func validateMaxItems(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		max_items针对repeated/map类型. 其他类型自动忽略
	*/
	switch {
	case f.IsRepeated:
		fallthrough
	case f.IsMap:
		g.P("if len(x.", f.GoName, ") > ", f.Rule.MaxItems.Val, " {")
		g.P("return set.Meta.FieldRules[", i, "].MaxItems.Err")
		g.P("}")
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
			protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
			protoreflect.FloatKind, protoreflect.DoubleKind:
		case protoreflect.StringKind:
		case protoreflect.BytesKind:
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}

func validateEnum(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		enum针对enum/number/string类型. 其他类型自动忽略
	*/
	switch {
	case f.IsRepeated:
	case f.IsMap:
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
			fallthrough
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			if len(f.Rule.Enum.Int) > 0 {
				if f.HasOptional {
					g.P("if x.", f.GoName, " != nil {")
					g.P("switch *x.", f.GoName, " {")
					for _, v := range f.Rule.Enum.Int {
						g.P(fmt.Sprintf(`case %v:`, v))
					}
					g.P("default: return set.Meta.FieldRules[", i, "].Enum.Err ")
					g.P("}")
					g.P("}")
				} else {
					g.P("switch x.", f.GoName, " {")
					for _, v := range f.Rule.Enum.Int {
						g.P(fmt.Sprintf(`case %v:`, v))
					}
					g.P("default: return set.Meta.FieldRules[", i, "].Enum.Err ")
					g.P("}")
				}
			}
		case protoreflect.StringKind:
			if len(f.Rule.Enum.Int) > 0 {
				if f.HasOptional {
					g.P("if x.", f.GoName, " != nil {")
					g.P("switch *x.", f.GoName, " {")
					for _, v := range f.Rule.Enum.Str {
						g.P(fmt.Sprintf(`case %q:`, v)) // 不要简单地拼接"..."
					}
					g.P("default: return set.Meta.FieldRules[", i, "].Enum.Err ")
					g.P("}")
					g.P("}")
				} else {
					g.P("if x.", f.GoName, " != `` {")
					g.P("switch x.", f.GoName, " {")
					for _, v := range f.Rule.Enum.Str {
						g.P(fmt.Sprintf(`case %q:`, v)) // 不要简单地拼接"..."
					}
					g.P("default: return set.Meta.FieldRules[", i, "].Enum.Err ")
					g.P("}")
					g.P("}")
				}
			}
		case protoreflect.BytesKind:
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}

func validatePattern(g *protogen.GeneratedFile, i int, f *FieldExt) {
	/*
		pattern针对string/bytes类型. 其他类型自动忽略
	*/
	switch {
	case f.IsRepeated:
	case f.IsMap:
	default:
		switch f.Kind {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		case protoreflect.FloatKind, protoreflect.DoubleKind:
		case protoreflect.StringKind:
			if f.HasOptional {
				g.P("if x.", f.GoName, " != nil {")
				g.P("if !set.FieldPatterns[", i, "].MatchString(*x.", f.GoName, ") {")
				g.P("return set.Meta.FieldRules[", i, "].Pattern.Err")
				g.P("}")
				g.P("}")
			} else {
				g.P("if x.", f.GoName, " != `` {")
				g.P("if !set.FieldPatterns[", i, "].MatchString(x.", f.GoName, ") {")
				g.P("return set.Meta.FieldRules[", i, "].Pattern.Err")
				g.P("}")
				g.P("}")
			}
		case protoreflect.BytesKind:
			g.P("if x.", f.GoName, " != nil {")
			g.P("if !set.FieldPatterns[", i, "].Match(x.", f.GoName, ") {")
			g.P("return set.Meta.FieldRules[", i, "].Pattern.Err")
			g.P("}")
			g.P("}")
		case protoreflect.MessageKind:
		case protoreflect.GroupKind:
		}
	}
}
