package protoapi

func QueryBool(ctx *Context, key string, ptr *bool, style Style) error {

}

func QueryBoolOptional(ctx *Context, key string, ptr **bool, style Style) error {

}

func QueryBoolRepeated(ctx *Context, key string, ptr *[]bool, style Style) error {

}

func QueryBoolMap(ctx *Context, key string, ptr *map[string]bool, style Style) error {

}

func QueryInt32(ctx *Context, key string, ptr *int32, style Style) error {

}

func QueryInt32Optional(ctx *Context, key string, ptr **int32, style Style) error {

}

func QueryInt32Repeated(ctx *Context, key string, ptr *[]int32, style Style) error {

}

func QueryInt32Map(ctx *Context, key string, ptr *map[string]int32, style Style) error {

}

func QueryInt64(ctx *Context, key string, ptr *int64, style Style) error {

}

func QueryInt64Optional(ctx *Context, key string, ptr **int64, style Style) error {

}

func QueryInt64Repeated(ctx *Context, key string, ptr *[]int64, style Style) error {

}

func QueryInt64Map(ctx *Context, key string, ptr *map[string]int64, style Style) error {

}

func QueryUint32(ctx *Context, key string, ptr *uint32, style Style) error {

}

func QueryUint32Optional(ctx *Context, key string, ptr **uint32, style Style) error {

}

func QueryUint32Repeated(ctx *Context, key string, ptr *[]uint32, style Style) error {

}

func QueryUint32Map(ctx *Context, key string, ptr *map[string]uint32, style Style) error {

}

func QueryUint64(ctx *Context, key string, ptr *uint64, style Style) error {

}

func QueryUint64Optional(ctx *Context, key string, ptr **uint64, style Style) error {

}

func QueryUint64Repeated(ctx *Context, key string, ptr *[]uint64, style Style) error {

}

func QueryUint64Map(ctx *Context, key string, ptr *map[string]uint64, style Style) error {

}

func QueryFloat32(ctx *Context, key string, ptr *float32, style Style) error {

}

func QueryFloat32Optional(ctx *Context, key string, ptr **float32, style Style) error {

}

func QueryFloat32Repeated(ctx *Context, key string, ptr *[]float32, style Style) error {

}

func QueryFloat32Map(ctx *Context, key string, ptr *map[string]float32, style Style) error {

}

func QueryDouble(ctx *Context, key string, ptr *float64, style Style) error {

}

func QueryDoubleOptional(ctx *Context, key string, ptr **float64, style Style) error {

}

func QueryDoubleRepeated(ctx *Context, key string, ptr *[]float64, style Style) error {

}

func QueryDoubleMap(ctx *Context, key string, ptr *map[string]float64, style Style) error {

}

func QueryString(ctx *Context, key string, ptr *string, style Style) error {

}

func QueryStringOptional(ctx *Context, key string, ptr **string, style Style) error {

}

func QueryStringRepeated(ctx *Context, key string, ptr *[]string, style Style) error {

}

func QueryStringMap(ctx *Context, key string, ptr *map[string]string, style Style) error {

}

func QueryBytes(ctx *Context, key string, ptr *[]byte, style Style) error {

}

func QueryBytesOptional(ctx *Context, key string, ptr *[]byte, style Style) error {

}

func QueryBytesRepeated(ctx *Context, key string, ptr *[][]byte, style Style) error {

}

func QueryBytesMap(ctx *Context, key string, ptr *map[string][]byte, style Style) error {

}

func QueryEnum[Enum ~int32](ctx *Context, key string, ptr *Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func QueryEnumOptional[Enum ~int32](ctx *Context, key string, ptr **Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func QueryEnumRepeated[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func QueryEnumMap[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func QueryEnum_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func QueryEnumOptional_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr **Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func QueryEnumRepeated_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func QueryEnumMap_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Style, names map[int32]string, values map[string]int32) error {

}
