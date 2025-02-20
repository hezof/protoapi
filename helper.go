package protoapi

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

var ErrInvalidContextToReadFile = errors.New("invalid context to read file")
var ErrInvalidContextToWriteFile = errors.New("invalid context to write file")

// ReadFile 配合body=file使用! 专门用于HTTP上传数据! 无法用于GRPC场景!
func ReadFile(ctx context.Context, key string) (multipart.File, *multipart.FileHeader, error) {
	if wc, ok := ctx.(*Context); ok {
		return wc.Request.FormFile(key)
	}
	return nil, nil, ErrInvalidContextToReadFile
}

// WriteFile 专门用于HTTP上传文件! 无法用于GRPC场景!
// 使用WriteFile()/WriterStream()的Service Method实现必须确保返回rsp为nil(即无法用于grpc调用)!
func WriteFile(ctx context.Context, name string, data io.Reader) error {
	if wc, ok := ctx.(*Context); ok {
		if wc.mux.closed != 0 {
			wc.ResponseWriter.Header()["Connection"] = closeConnection
		}
		wc.ResponseWriter.Header()["Content-Type"] = streamContentType
		wc.ResponseWriter.Header()["Content-Disposition"] = []string{"attachment;filename=" + url.QueryEscape(name)} // 防止中文乱码
		wc.ResponseWriter.WriteStatus(wc.Handler.Status)
		if _, err := io.Copy(wc.ResponseWriter.ResponseWriter, data); err != nil {
			return err
		}
		return nil
	}
	return ErrInvalidContextToWriteFile
}

var ErrInvalidContextToReadStream = errors.New("invalid context to read stream")
var ErrInvalidContextToWriteStream = errors.New("invalid context to write stream")

// ReadStream 配合body=stream使用! 专门用于HTTP上传数据! 无法用于GRPC场景!
func ReadStream(ctx context.Context) (io.ReadCloser, error) {
	if wc, ok := ctx.(*Context); ok {
		return wc.Request.Body, nil
	}
	return nil, ErrInvalidContextToReadStream
}

// WriteStream 专门用于HTTP上传文件! 无法用于GRPC场景!
// 使用DownFile()/WriterStream()的Service Method实现必须确保返回rsp为nil(即无法用于grpc调用)!
func WriteStream(ctx context.Context, headers ...http.Header) (io.Writer, error) {
	if wc, ok := ctx.(*Context); ok {
		if wc.mux.closed != 0 {
			wc.ResponseWriter.Header()["Connection"] = closeConnection
		}
		for _, h := range headers {
			for k, vs := range h {
				wc.ResponseWriter.Header()[k] = vs
			}
		}
		wc.ResponseWriter.WriteStatus(wc.Handler.Status)
		return wc.ResponseWriter.ResponseWriter, nil
	}
	return nil, ErrInvalidContextToWriteStream
}
