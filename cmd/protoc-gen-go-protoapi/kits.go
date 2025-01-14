package main

type Vec[V any] []V

func (vs *Vec[V]) Add(v V) {
	*vs = append(*vs, v)
}

type Item[V any] struct {
	Key string
	Val V
}
type Map[V any] []*Item[[]V]

func (ms *Map[V]) Add(k string, v V) {
	for _, it := range *ms {
		if it.Key == k {
			it.Val = append(it.Val, v)
			return
		}
	}
	it := new(Item[[]V])
	it.Val = append(it.Val, v)
	*ms = append(*ms, it)
}
