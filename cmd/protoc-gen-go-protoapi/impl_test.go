package main

import (
	"fmt"
	"hash/crc64"
	"strconv"
	"testing"
)

func TestPlugin(t *testing.T) {
	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	h.Write([]byte("Demo.Test_"))
	fmt.Println(strconv.FormatUint(h.Sum64(), 36))
}
