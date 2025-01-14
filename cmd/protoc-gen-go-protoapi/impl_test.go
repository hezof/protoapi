package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestImpl(t *testing.T) {
	meta := new(FileExt)
	meta.Enums.Add("test", new(EnumExt))
	meta.Messages.Add("test", new(MessageExt))
	meta.Services.Add("test", new(ServiceExt))
	bs, err := json.MarshalIndent(meta, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
