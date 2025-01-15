package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

type Map map[string]any

func (m Map) MarshalYAML() (interface{}, error) {
	var items []yaml.MapItem
	for k, v := range m {
		items = append(items, yaml.MapItem{
			Key:   k,
			Value: v,
		})
	}
	return items, nil

}

var _ yaml.Marshaler = Map{}

func TestDocs(t *testing.T) {
	m := make(Map)
	m["a"] = 1
	m["b"] = 2
	bs, err := yaml.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
