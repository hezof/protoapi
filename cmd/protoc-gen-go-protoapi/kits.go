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

type IdxVec[V any] struct {
	Idx map[string]V
	Vec []V
}

func (m *IdxVec[V]) Add(k string, v V) (V, bool) {
	if m.Idx == nil {
		m.Idx = make(map[string]V)
	}
	if vl, ok := m.Idx[k]; ok {
		return vl, false
	}
	m.Idx[k] = v
	m.Vec = append(m.Vec, v)
	return v, true
}
