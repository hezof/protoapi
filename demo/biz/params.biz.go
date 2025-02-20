package biz

import (
	"context"
	"github.com/hezof/protoapi/demo/api"
)

// ParamsImplement api.Params. 参数服务用于测试通过form/path/query/header/cookie传递的参数.该服务仅用于http!
type ParamsImplement struct{}

var _ api.ParamsServer = (*ParamsImplement)(nil)

// Form api.Params.form. 通过form/body传递参数
// POST /params/form
func (ps *ParamsImplement) Form(ctx context.Context, req *api.ParamsInForm) (rsp *api.ParamsInForm, err error) {
	rsp = req
	return
}

// Path api.Params.path. 通过path传递参数
// POST /params/path
func (ps *ParamsImplement) Path(ctx context.Context, req *api.ParamsInPath) (rsp *api.ParamsInPath, err error) {
	rsp = req
	return
}

// Query api.Params.query. 通过path传递参数
// POST /params/query
func (ps *ParamsImplement) Query(ctx context.Context, req *api.ParamsInQuery) (rsp *api.ParamsInQuery, err error) {
	rsp = req
	return
}

// Header api.Params.header. 通过header传递参数
// POST /params/header
func (ps *ParamsImplement) Header(ctx context.Context, req *api.ParamsInHeader) (rsp *api.ParamsInHeader, err error) {
	rsp = req
	return
}

// Cookie api.Params.cookie. 通过cookie传递参数
// POST /params/cookie
func (ps *ParamsImplement) Cookie(ctx context.Context, req *api.ParamsInCookie) (rsp *api.ParamsInCookie, err error) {
	rsp = req
	return
}
