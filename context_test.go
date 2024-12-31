package protoapi

import (
	"google.golang.org/grpc/status"
	"testing"
)

type Demo struct {
}

func PrintDemo(d *Demo) {

}

func Print[T any](t T, p func(d T)) {
	p(t)
}

func handle() any {
	return new(Demo)
}

func TestDecode(t *testing.T) {
	status.FromError()
}
