package protoapi

import (
	"testing"

	"google.golang.org/grpc"
)

type Demo struct {
}

func TestXxx(t *testing.T) {
	var c = new(Context)
	var v grpc.ClientStreamingServer[Demo, Demo] = &ServerStreamContext[Demo, Demo]{c}
	_ = v
}
