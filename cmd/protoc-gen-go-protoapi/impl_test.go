package main

import (
	"fmt"
	"testing"
)

type Demo struct {
	Ps *string
}

func TestPlugin(t *testing.T) {
	d := new(Demo)
	d.Ps = new(string)
	*d.Ps = `b`
	fmt.Println(d.Ps == nil || *d.Ps < `a`)
}
