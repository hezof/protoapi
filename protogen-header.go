package protoapi

func HeaderBool(ctx *Context, key string, ptr *bool, style Style) error {

}

func HeaderBoolOptional(ctx *Context, key string, ptr **bool, style Style) error {

}

func HeaderBoolRepeated(ctx *Context, key string, ptr *[]bool, style Style) error {

}

func HeaderBoolMap(ctx *Context, key string, ptr *map[string]bool, style Style) error {

}

func HeaderInt32(ctx *Context, key string, ptr *int32, style Style) error {

}

func HeaderInt32Optional(ctx *Context, key string, ptr **int32, style Style) error {

}

func HeaderInt32Repeated(ctx *Context, key string, ptr *[]int32, style Style) error {

}

func HeaderInt32Map(ctx *Context, key string, ptr *map[string]int32, style Style) error {

}

func HeaderInt64(ctx *Context, key string, ptr *int64, style Style) error {

}

func HeaderInt64Optional(ctx *Context, key string, ptr **int64, style Style) error {

}

func HeaderInt64Repeated(ctx *Context, key string, ptr *[]int64, style Style) error {

}

func HeaderInt64Map(ctx *Context, key string, ptr *map[string]int64, style Style) error {

}

func HeaderUint32(ctx *Context, key string, ptr *uint32, style Style) error {

}

func HeaderUint32Optional(ctx *Context, key string, ptr **uint32, style Style) error {

}

func HeaderUint32Repeated(ctx *Context, key string, ptr *[]uint32, style Style) error {

}

func HeaderUint32Map(ctx *Context, key string, ptr *map[string]uint32, style Style) error {

}

func HeaderUint64(ctx *Context, key string, ptr *uint64, style Style) error {

}

func HeaderUint64Optional(ctx *Context, key string, ptr **uint64, style Style) error {

}

func HeaderUint64Repeated(ctx *Context, key string, ptr *[]uint64, style Style) error {

}

func HeaderUint64Map(ctx *Context, key string, ptr *map[string]uint64, style Style) error {

}

func HeaderFloat32(ctx *Context, key string, ptr *float32, style Style) error {

}

func HeaderFloat32Optional(ctx *Context, key string, ptr **float32, style Style) error {

}

func HeaderFloat32Repeated(ctx *Context, key string, ptr *[]float32, style Style) error {

}

func HeaderFloat32Map(ctx *Context, key string, ptr *map[string]float32, style Style) error {

}

func HeaderDouble(ctx *Context, key string, ptr *float64, style Style) error {

}

func HeaderDoubleOptional(ctx *Context, key string, ptr **float64, style Style) error {

}

func HeaderDoubleRepeated(ctx *Context, key string, ptr *[]float64, style Style) error {

}

func HeaderDoubleMap(ctx *Context, key string, ptr *map[string]float64, style Style) error {

}

func HeaderString(ctx *Context, key string, ptr *string, style Style) error {

}

func HeaderStringOptional(ctx *Context, key string, ptr **string, style Style) error {

}

func HeaderStringRepeated(ctx *Context, key string, ptr *[]string, style Style) error {

}

func HeaderStringMap(ctx *Context, key string, ptr *map[string]string, style Style) error {

}

func HeaderBytes(ctx *Context, key string, ptr *[]byte, style Style) error {

}

func HeaderBytesOptional(ctx *Context, key string, ptr *[]byte, style Style) error {

}

func HeaderBytesRepeated(ctx *Context, key string, ptr *[][]byte, style Style) error {

}

func HeaderBytesMap(ctx *Context, key string, ptr *map[string][]byte, style Style) error {

}

func HeaderEnum[Enum ~int32](ctx *Context, key string, ptr *Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func HeaderEnumOptional[Enum ~int32](ctx *Context, key string, ptr **Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func HeaderEnumRepeated[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func HeaderEnumMap[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func HeaderEnum_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func HeaderEnumOptional_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr **Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func HeaderEnumRepeated_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func HeaderEnumMap_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Style, names map[int32]string, values map[string]int32) error {

}
