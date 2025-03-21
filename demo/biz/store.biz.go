package biz

import (
	"bufio"
	"context"
	"fmt"
	"github.com/hezof/core"
	"github.com/hezof/protoapi"
	"github.com/hezof/protoapi/demo/api"
	"google.golang.org/grpc"
	"io"
	"os"
)

// StoreImplement api.Store. 这是一个商店服务示例跨行第一句.protobuf自动拼接相邻字串!跨行第二句.protobuf自动拼接相邻字串!跨行第三句.protobuf自动拼接相邻字串!
type StoreImplement struct{}

var _ api.StoreServer = (*StoreImplement)(nil)

// Simple api.Store.Simple. 通过该方法可以创建一本书籍信息
// POST /simple/book
func (ps *StoreImplement) Simple(ctx context.Context, req *api.Book) (rsp *api.Book, err error) {
	rsp = req
	return
}

// Client api.Store.Client. 通过该方法可以流式批量创建多本书籍信息
// POST /client/book
// WEBSOCKET /client/book
func (ps *StoreImplement) Client(svr grpc.ClientStreamingServer[api.Book, api.Book]) (err error) {
	req, err := svr.Recv()
	if err != nil {
		return err
	}
	rsp := req
	err = svr.SendAndClose(rsp)
	if err != nil {
		return err
	}
	return
}

// Server api.Store.muxGroupServer. 通过该方法可以流式返回指书籍信息
// POST /server/book
// WEBSOCKET /server/book
func (ps *StoreImplement) Server(req *api.Book, svr grpc.ServerStreamingServer[api.Book]) (err error) {
	rsp := req
	err = svr.Send(rsp)
	if err != nil {
		return err
	}
	return
}

// Bidirectional api.Store.Bidirectional. 通过该方法可以流式返回指书籍信息
// POST /bidirectional/book
// WEBSOCKET /bidirectional/book
func (ps *StoreImplement) Bidirectional(svr grpc.BidiStreamingServer[api.Book, api.Book]) (err error) {
	req, err := svr.Recv()
	if err != nil {
		return err
	}
	rsp := req
	err = svr.Send(rsp)
	if err != nil {
		return err
	}
	return
}

// UploadCovert api.Store.UploadCovert. 通过该方法上传封面图片
// POST /covert/upload
func (ps *StoreImplement) UploadCovert(ctx context.Context, req *api.Void) (rsp *api.Void, err error) {
	file, info, err := protoapi.ReadFile(ctx, "file")
	if err != nil {
		return
	}
	fmt.Println(core.ToJson(info))
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return
	}
	return
}

// ReviewCovert api.Store.ReviewCovert. 通过该方法预览封面图片
// GET /review/covert
func (ps *StoreImplement) ReviewCovert(ctx context.Context, req *api.Void) (rsp *api.Void, err error) {
	file, err := os.Open(`D:\Workspace\klibsrc\protoapi\demo\go.mod`)
	if err != nil {
		return
	}
	defer file.Close()
	err = protoapi.WriteFile(ctx, "go.mod", bufio.NewReader(file))
	return
}
