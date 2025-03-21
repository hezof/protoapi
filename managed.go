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

type ServerComponent struct {
	Server
}

func (m *ServerComponent) GetTarget() core.ManagedTarget {
	return m.Server
}

func (m *ServerComponent) SetTarget(target core.ManagedTarget) {
	m.Server = target.(Server)
}

var _ Server = (*ServerComponent)(nil)
var _ core.ManagedComponent = (*ServerComponent)(nil)

/********************************************
 * 托管工厂
 ********************************************/

type ServerFactory struct {
}

func (s *ServerFactory) Manage(n string) core.ManagedComponent {
	return new(ServerComponent)
}

func (s *ServerFactory) Create(c *core.ManagedConfig) (core.ManagedTarget, error) {
	cfg := new(Config)
	err := core.SimpleStructBinder.MapStruct(c.Value, &cfg, "")
	if err != nil {
		return nil, err
	}
	return NewServer(cfg), nil
}

func (s *ServerFactory) Destroy(t core.ManagedTarget) error {
	// server会一直阻塞
	return nil
}

var _ core.ManagedFactory = (*ServerFactory)(nil)

func init() {
	core.Register("protoapi", new(ServerFactory))
}
