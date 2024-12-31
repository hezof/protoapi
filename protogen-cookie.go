package protoapi

func CookieBool(ctx *Context, key string, ptr *bool, style Style) error {

}

func CookieBoolOptional(ctx *Context, key string, ptr **bool, style Style) error {

}

func CookieBoolRepeated(ctx *Context, key string, ptr *[]bool, style Style) error {

}

func CookieBoolMap(ctx *Context, key string, ptr *map[string]bool, style Style) error {

}

func CookieInt32(ctx *Context, key string, ptr *int32, style Style) error {

}

func CookieInt32Optional(ctx *Context, key string, ptr **int32, style Style) error {

}

func CookieInt32Repeated(ctx *Context, key string, ptr *[]int32, style Style) error {

}

func CookieInt32Map(ctx *Context, key string, ptr *map[string]int32, style Style) error {

}

func CookieInt64(ctx *Context, key string, ptr *int64, style Style) error {

}

func CookieInt64Optional(ctx *Context, key string, ptr **int64, style Style) error {

}

func CookieInt64Repeated(ctx *Context, key string, ptr *[]int64, style Style) error {

}

func CookieInt64Map(ctx *Context, key string, ptr *map[string]int64, style Style) error {

}

func CookieUint32(ctx *Context, key string, ptr *uint32, style Style) error {

}

func CookieUint32Optional(ctx *Context, key string, ptr **uint32, style Style) error {

}

func CookieUint32Repeated(ctx *Context, key string, ptr *[]uint32, style Style) error {

}

func CookieUint32Map(ctx *Context, key string, ptr *map[string]uint32, style Style) error {

}

func CookieUint64(ctx *Context, key string, ptr *uint64, style Style) error {

}

func CookieUint64Optional(ctx *Context, key string, ptr **uint64, style Style) error {

}

func CookieUint64Repeated(ctx *Context, key string, ptr *[]uint64, style Style) error {

}

func CookieUint64Map(ctx *Context, key string, ptr *map[string]uint64, style Style) error {

}

func CookieFloat32(ctx *Context, key string, ptr *float32, style Style) error {

}

func CookieFloat32Optional(ctx *Context, key string, ptr **float32, style Style) error {

}

func CookieFloat32Repeated(ctx *Context, key string, ptr *[]float32, style Style) error {

}

func CookieFloat32Map(ctx *Context, key string, ptr *map[string]float32, style Style) error {

}

func CookieDouble(ctx *Context, key string, ptr *float64, style Style) error {

}

func CookieDoubleOptional(ctx *Context, key string, ptr **float64, style Style) error {

}

func CookieDoubleRepeated(ctx *Context, key string, ptr *[]float64, style Style) error {

}

func CookieDoubleMap(ctx *Context, key string, ptr *map[string]float64, style Style) error {

}

func CookieString(ctx *Context, key string, ptr *string, style Style) error {

}

func CookieStringOptional(ctx *Context, key string, ptr **string, style Style) error {

}

func CookieStringRepeated(ctx *Context, key string, ptr *[]string, style Style) error {

}

func CookieStringMap(ctx *Context, key string, ptr *map[string]string, style Style) error {

}

func CookieBytes(ctx *Context, key string, ptr *[]byte, style Style) error {

}

func CookieBytesOptional(ctx *Context, key string, ptr *[]byte, style Style) error {

}

func CookieBytesRepeated(ctx *Context, key string, ptr *[][]byte, style Style) error {

}

func CookieBytesMap(ctx *Context, key string, ptr *map[string][]byte, style Style) error {

}

func CookieEnum[Enum ~int32](ctx *Context, key string, ptr *Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func CookieEnumOptional[Enum ~int32](ctx *Context, key string, ptr **Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func CookieEnumRepeated[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func CookieEnumMap[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func CookieEnum_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func CookieEnumOptional_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr **Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func CookieEnumRepeated_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *[]Enum, style Style, names map[int32]string, values map[string]int32) error {

}

func CookieEnumMap_EnumAsInt[Enum ~int32](ctx *Context, key string, ptr *map[string]Enum, style Style, names map[int32]string, values map[string]int32) error {

}
