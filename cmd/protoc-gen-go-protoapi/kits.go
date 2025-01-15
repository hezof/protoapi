package main

func NvlS(ss ...string) string {
	for _, s := range ss {
		if s != `` {
			return s
		}
	}
	return ``
}

func NvlI[V int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](vs ...V) V {
	for _, v := range vs {
		if v != 0 {
			return v
		}
	}
	return 0
}

func Set(v string, vs ...string) []string {
	if len(vs) == 0 {
		return []string{v}
	}
	m := make(map[string]bool)
	m[v] = true
	for _, v = range vs {
		m[v] = true
	}
	rt := make([]string, 0, len(m))
	for k, _ := range m {
		rt = append(rt, k)
	}
	return rt
}
