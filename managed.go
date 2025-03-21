package protoapi

import (
	"github.com/hezof/core"
	"google.golang.org/grpc"
	"net/http"
)

/********************************************
 * 对外接口
 ********************************************/

type GroupRouter interface {
	BasePath() string
	Group(path string, hs ...HandleFunc) GroupRouter
	Use(filters ...HandleFunc) GroupRouter
	HandleFunc(method string, path string, hs ...HandleFunc) GroupRouter
	Handle(hd *Handler) GroupRouter
	Any(path string, f ...HandleFunc) GroupRouter
	GET(path string, f ...HandleFunc) GroupRouter
	POST(path string, f ...HandleFunc) GroupRouter
	PUT(path string, f ...HandleFunc) GroupRouter
	DELETE(path string, f ...HandleFunc) GroupRouter
	HEAD(path string, f ...HandleFunc) GroupRouter
	PATCH(path string, f ...HandleFunc) GroupRouter
	OPTIONS(path string, f ...HandleFunc) GroupRouter
	CONNECT(path string, f ...HandleFunc) GroupRouter
	TRACE(path string, f ...HandleFunc) GroupRouter
	StaticFile(path string, file string) GroupRouter
	Static(prefix string, root string) GroupRouter
	StaticFS(prefix string, fs http.FileSystem) GroupRouter
}

type Server interface {
	GroupRouter
	Config() *Config
	RegisterService(registry ServiceRegistry, implement interface{}, aspects ...ServiceAspect) error
	ListenAndServe() (err error)
	OnInit(fs ...func())
	OnReady(fs ...func())
	OnExit(fs ...func())
	OnRegisterGrpcService(f func(name, addr string, undo bool))
	OnRegisterHttpService(f func(name, addr string, undo bool))
	ServiceAspect(vs ...ServiceAspect)
	ServicePlugin(vs ...ServicePlugin)
	RequestPlugin(vs ...RequestPlugin)
	GrpcServerOption(vs ...grpc.ServerOption)
	HttpServerOption(vs ...HandleFunc)
	GrpcServerInvoke(vs ...func(server *grpc.Server))
	GrpcPanicFunc(f GrpcPanicFunc)
	HttpPanicFunc(f HandleFunc)
	HttpPanic(h *Handler)
	HttpNotFoundFunc(f HandleFunc)
	HttpNotFound(h *Handler)
}

/********************************************
 * 托管组件
 ********************************************/

type ServerManagedComponent struct {
	Server
}

func (m *ServerManagedComponent) GetTarget() core.ManagedTarget {
	return m.Server
}

func (m *ServerManagedComponent) SetTarget(target core.ManagedTarget) {
	m.Server = target.(Server)
}

var _ Server = (*ServerManagedComponent)(nil)
var _ core.ManagedComponent = (*ServerManagedComponent)(nil)

/********************************************
 * 托管工厂
 ********************************************/

type ServerManagedFactory struct {
}

func (s ServerManagedFactory) Manage(name string) core.ManagedComponent {
	return new(ServerManagedComponent)
}

func (s ServerManagedFactory) Create(c *core.ManagedConfig) (core.ManagedTarget, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServerManagedFactory) Destroy(v core.ManagedTarget) error {
	return nil
}

var _ core.ManagedFactory = (*ServerManagedFactory)(nil)
