package protoapi

func PathBool(ctx *Context, key string, ptr *bool, style Prop_Style) error {

}

func PathBoolOptional(ctx *Context, key string, ptr **bool, style Prop_Style) error {

}

func PathBoolRepeated(ctx *Context, key string, ptr *[]bool, style Prop_Style) error {

}

func PathBoolMap(ctx *Context, key string, ptr *map[string]bool, style Prop_Style) error {

}

func PathInt32(ctx *Context, key string, ptr *int32, style Prop_Style) error {

}

func PathInt32Optional(ctx *Context, key string, ptr **int32, style Prop_Style) error {

}

func PathInt32Repeated(ctx *Context, key string, ptr *[]int32, style Prop_Style) error {

}

func PathInt32Map(ctx *Context, key string, ptr *map[string]int32, style Prop_Style) error {

}

func PathInt64(ctx *Context, key string, ptr *int64, style Prop_Style) error {

}

func PathInt64Optional(ctx *Context, key string, ptr **int64, style Prop_Style) error {

}

func PathInt64Repeated(ctx *Context, key string, ptr *[]int64, style Prop_Style) error {

}

func PathInt64Map(ctx *Context, key string, ptr *map[string]int64, style Prop_Style) error {

}

func PathUint32(ctx *Context, key string, ptr *uint32, style Prop_Style) error {

}

func PathUint32Optional(ctx *Context, key string, ptr **uint32, style Prop_Style) error {

}

func PathUint32Repeated(ctx *Context, key string, ptr *[]uint32, style Prop_Style) error {

}

func PathUint32Map(ctx *Context, key string, ptr *map[string]uint32, style Prop_Style) error {

}

func PathUint64(ctx *Context, key string, ptr *uint64, style Prop_Style) error {

}

func PathUint64Optional(ctx *Context, key string, ptr **uint64, style Prop_Style) error {

}

func PathUint64Repeated(ctx *Context, key string, ptr *[]uint64, style Prop_Style) error {

}

func PathUint64Map(ctx *Context, key string, ptr *map[string]uint64, style Prop_Style) error {

}

func PathFloat32(ctx *Context, key string, ptr *float32, style Prop_Style) error {

}

func PathFloat32Optional(ctx *Context, key string, ptr **float32, style Prop_Style) error {

}

func PathFloat32Repeated(ctx *Context, key string, ptr *[]float32, style Prop_Style) error {

}

func PathFloat32Map(ctx *Context, key string, ptr *map[string]float32, style Prop_Style) error {

}

func PathDouble(ctx *Context, key string, ptr *float64, style Prop_Style) error {

}

func PathDoubleOptional(ctx *Context, key string, ptr **float64, style Prop_Style) error {

}

func PathDoubleRepeated(ctx *Context, key string, ptr *[]float64, style Prop_Style) error {

}

func PathDoubleMap(ctx *Context, key string, ptr *map[string]float64, style Prop_Style) error {

}

func PathString(ctx *Context, key string, ptr *string, style Prop_Style) error {

}

func PathStringOptional(ctx *Context, key string, ptr **string, style Prop_Style) error {

}

func PathStringRepeated(ctx *Context, key string, ptr *[]string, style Prop_Style) error {

}

func PathStringMap(ctx *Context, key string, ptr *map[string]string, style Prop_Style) error {

}

func PathBytes(ctx *Context, key string, ptr *[]byte, style Prop_Style) error {

}

func PathBytesOptional(ctx *Context, key string, ptr *[]byte, style Prop_Style) error {

}

func PathBytesRepeated(ctx *Context, key string, ptr *[][]byte, style Prop_Style) error {

}

func PathBytesMap(ctx *Context, key string, ptr *map[string][]byte, style Prop_Style) error {

}

func PathEnum[Enum ~int32](ctx *Context, key string, ptr *Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}

func PathEnumOptional[Enum ~int32](ctx *Context, key string, ptr **Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}

func PathEnumRepeated[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}

func PathEnumMap[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}

func PathEnum_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}

func PathEnumOptional_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr **Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}

func PathEnumRepeated_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}

func PathEnumMap_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Prop_Style, names map[int32]string, values map[string]int32) error {

}
