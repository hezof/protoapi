package test

import (
	"context"
	"fmt"
	"github.com/hezof/framework"
	"github.com/hezof/protoapi"
	"github.com/hezof/protoapi/demo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func conn() *grpc.ClientConn {
	con, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return con
}

var gcli = api.NewStoreClient(conn())

func TestGrpcSimple(t *testing.T) {
	rsp, err := gcli.Simple(context.Background(), book)
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestGrpcClient(t *testing.T) {
	ccli, err := gcli.Client(context.Background())
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	err = ccli.Send(book)
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	rsp, err := ccli.CloseAndRecv()
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestGrpcServer(t *testing.T) {
	scli, err := gcli.Server(context.Background(), book)
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	rsp, err := scli.Recv()
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	fmt.Println(core.ToJson(rsp))
}

func TestGrpcBidirectional(t *testing.T) {
	dcli, err := gcli.Bidirectional(context.Background())
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	err = dcli.Send(book)
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	rsp, err := dcli.Recv()
	if err != nil {
		serr := protoapi.StatusErrorFrom(err)
		t.Fatal(serr)
	}
	fmt.Println(core.ToJson(rsp))
}
