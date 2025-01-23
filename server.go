package protoapi

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"ksogit.kingsoft.net/kgo/log"
	"net"
	"net/http"
)

/**
- http服务器
- grpc服务器
- http-grpc网关
*/

type Server struct {
	mux                                                      // 多路
	config                *Config                            // 配置
	settings              map[string]*MethodSetting          // 方法设置(过程数据)
	onInit                []func()                           // 初始前回调
	onReady               []func()                           // 启动前回调
	onExit                []func()                           // 退出前回调
	onRegisterGrpcService func(name, addr string, undo bool) // grpc服务注册回调
	onRegisterHttpService func(name, addr string, undo bool) // http服务注册回调
	grpcPanicFunc         GrpcPanicFunc                      // Grpc panic函数
	/*
		下述过程数据在服务最终启动前会被清除释放内存!
	*/
	*_group                                       // 分组路由(过程数据)
	_serviceSetting   []*ServiceSetting           // 服务设置(过程数据)
	_serviceAspect    []ServiceAspect             // 服务校验器(过程参数)
	_servicePlugin    []ServicePlugin             // 服务插件(过程数据)
	_requestPlugin    []RequestPlugin             // 请求插件(过程数据)
	_grpcServerOption []grpc.ServerOption         // grpc拦截器(过程数据)
	_httpServerOption []HandleFunc                // http拦截器(过程数据)
	_grpcServerInvoke []func(server *grpc.Server) // grpc服务器回调(过程数据)
}

// 清理运行过程不再需要的中间变量
func (s *Server) clean() {
	s._group = nil
	s._serviceSetting = nil
	s._serviceAspect = nil
	s._servicePlugin = nil
	s._requestPlugin = nil
	s._grpcServerOption = nil
	s._httpServerOption = nil
	s._grpcServerInvoke = nil
	_globalRegister = nil
}

func NewServer(c *Config) *Server {
	return &Server{
		mux: mux{
			httpPanic:    defaultHttpPanicHandler,
			httpNotFound: defaultHttpNotFoundHandler,
		},
		config:        mergeConfig(c),
		settings:      make(map[string]*MethodSetting),
		_group:        newGroup(""),
		grpcPanicFunc: defaultGrpcPanicFunc,
	}
}

/**********************************************
	元数据加载
 **********************************************/

func (s *Server) Config() *Config {
	return s.config
}

/**********************************************
	生命周期回调
 **********************************************/

// OnInit 服务器启动时回调
func (s *Server) OnInit(fs ...func()) {
	s.onInit = append(s.onInit, fs...)
}

// OnReady 服务器完成时回调
func (s *Server) OnReady(fs ...func()) {
	s.onReady = append(s.onReady, fs...)
}

// OnExit 服务器退出时回调
func (s *Server) OnExit(fs ...func()) {
	s.onExit = append(s.onExit, fs...)
}

// OnRegisterGrpcService 注册GRPC服务回调
func (s *Server) OnRegisterGrpcService(f func(name, addr string, undo bool)) {
	s.onRegisterGrpcService = f
}

// OnRegisterHttpService 注册HTTP服务回调
func (s *Server) OnRegisterHttpService(f func(name, addr string, undo bool)) {
	s.onRegisterHttpService = f
}

/**********************************************
	配置添加
 **********************************************/

func (s *Server) ServiceAspect(vs ...ServiceAspect) {
	s._serviceAspect = append(s._serviceAspect, vs...)
}

func (s *Server) ServicePlugin(vs ...ServicePlugin) {
	s._servicePlugin = append(s._servicePlugin, vs...)
}

func (s *Server) RequestPlugin(vs ...RequestPlugin) {
	s._requestPlugin = append(s._requestPlugin, vs...)
}

func (s *Server) GrpcServerOption(vs ...grpc.ServerOption) {
	s._grpcServerOption = append(s._grpcServerOption, vs...)
}

func (s *Server) HttpServerOption(vs ...HandleFunc) {
	s._httpServerOption = append(s._httpServerOption, vs...)
}

func (s *Server) GrpcServerInvoke(vs ...func(server *grpc.Server)) {
	s._grpcServerInvoke = append(s._grpcServerInvoke, vs...)
}

func (s *Server) GrpcPanicFunc(f GrpcPanicFunc) {
	// 不能为空
	if f != nil {
		s.grpcPanicFunc = f
	}
}

/**********************************************
	服务注册
 **********************************************/

func (s *Server) RegisterService(registry ServiceRegistry, implement interface{}, aspects ...ServiceAspect) *Server {
	aspects = orderServiceAspects(s._serviceAspect, aspects)
	serviceSetting := registry(implement, aspects)
	for _, methodSetting := range serviceSetting.Methods {
		methodSetting.Service = serviceSetting
	}
	s._serviceSetting = append(s._serviceSetting, serviceSetting)
	return s
}

/**********************************************
	启动监听
 **********************************************/

func (s *Server) ListenAndServe() (err error) {

	/********************************************************
	* 如果没有配置grpc或http地址则自动结束服务初始化流程!
	 ********************************************************/
	if s.config.GrpcAddr == "" && s.config.HttpAddr == "" {
		return
	}

	/*******************************************************
	* 全局组件注册
	* 该步骤主要用于分离服务层与访问层,从而支持更灵活的组件部署模式
	 *******************************************************/
	for _, c := range _globalRegister {
		s.RegisterService(c.Registry, c.Implement, c.Aspects...)
	}

	/********************************************************
	* 处理ServiceSetting的插件(统计/修改/添加/删除)等
	 ********************************************************/
	for _, p := range s._servicePlugin {
		p(&s._serviceSetting)
	}

	/********************************************************
	* 添加MethodSetting相应的AccessSetting以及MethodGroup
	 ********************************************************/
	for _, ss := range s._serviceSetting {
		for _, ms := range ss.Methods {
			s.settings[FullMethod(ms.Meta)] = ms
			// 添加method相应的RequestSetting
			if ms.Meta.Http.Get != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodGet,
					Path:        ms.Meta.Http.Get,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Put != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodPut,
					Path:        ms.Meta.Http.Put,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Post != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodPost,
					Path:        ms.Meta.Http.Post,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Delete != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodDelete,
					Path:        ms.Meta.Http.Delete,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Options != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodOptions,
					Path:        ms.Meta.Http.Options,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Head != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodHead,
					Path:        ms.Meta.Http.Head,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Patch != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodPatch,
					Path:        ms.Meta.Http.Patch,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Trace != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodTrace,
					Path:        ms.Meta.Http.Trace,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Connect != "" {
				s._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodConnect,
					Path:        ms.Meta.Http.Connect,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Websocket != "" {
				if s.mux.upgrader == nil {
					s.mux.upgrader = newWebsocketUpgrader(s.config)
				}
				s._group.Handle(&Handler{
					Setting:      ms,
					Method:       http.MethodGet,
					Path:         ms.Meta.Http.Websocket,
					Status:       ms.Meta.Http.Status,
					Result:       ms.Meta.Http.Result,
					HandleChain:  []HandleFunc{WebsocketHandleFunc},
					BodyMaxBytes: -1, // 如果是Websocket长链接则自动忽略BodyMaxBytes参数
				})
			}
		}
	}

	/********************************************************
	* 扁平化router的AccessSetting: map{httpMethod: map{httpPath:requestSetting}}
	* 后面应用AccessOption时实现proxy, access, cache等的拦截逻辑
	 ********************************************************/
	_requestSetting := s._group.Flatten()

	/********************************************************
	* 处理RequestSetting的插件(缓存,日志,代理)等
	 ********************************************************/
	for _, p := range s._requestPlugin {
		p(_requestSetting)
	}

	/********************************************************
	* 将router的内容正式添加到mux. router是预处理, mux才是真路由
	* 拦截器顺序:
	*	1. options
	*	2. plugins
	*	3. filters
	*	4. handles
	* 同时设置所有handler的默认值
	 ********************************************************/
	for method, pathSetting := range _requestSetting {
		for path, setting := range pathSetting {
			handler := setting.Handler
			handler.BodyMaxBytes = NvlI(setting.Handler.BodyMaxBytes, s.config.HttpBodyMaxBytes)
			handler.FormMaxMemory = NvlI(setting.Handler.FormMaxMemory, s.config.HttpFormMaxMemory)
			handler.HandleChain = Join(s._httpServerOption, setting.Plugins, setting.Filters, setting.Handler.HandleChain)
			s.mux.route(method, path, handler) // 正式注册到mux
		}
	}
	s.mux.initServeHTTP()

	/********************************************************
	* 初始与退出钩子机制.
	 ********************************************************/
	if len(s.onInit) > 0 {
		for _, f := range s.onInit {
			protect(f)
		}
	}
	if len(s.onExit) > 0 {
		// 此处defer没用闭包
		defer func(exit []func()) {
			for _, f := range exit {
				protect(f)
			}
		}(s.onExit)
	}

	/********************************************************
	* 声明grpc,http,websocket所需的控件并确保退出关闭
	 ********************************************************/
	var (
		grpcServer   *grpc.Server
		grpcListener net.Listener
		httpServer   *http.Server
		httpListener net.Listener
	)
	defer func(grpcListener, httpListener net.Listener) {
		if grpcListener != nil {
			_ = grpcListener.Close()
		}
		if httpListener != nil {
			_ = httpListener.Close()
		}
	}(grpcListener, httpListener)

	// 启动grpc服务器
	if s.config.GrpcAddr != "" {
		var opts []grpc.ServerOption
		if len(s._grpcServerOption) > 0 {
			opts = append(opts, s._grpcServerOption...)
		}
		if s.config.GrpcKeepAlive > 0 {
			opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{Time: s.config.GrpcKeepAlive}))
		}
		if s.config.GrpcKeepAlivePolicy > 0 {
			opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{MinTime: s.config.GrpcKeepAlivePolicy, PermitWithoutStream: true}))
		}
		// 很关键: 必须确保"boostrap interceptor"位于最后位置, 即必须是grpc server最后添加的ServerOption
		grpcServer = grpc.NewServer(append(opts,
			grpc.ChainUnaryInterceptor(s.generateBootstrapUnaryInterceptor),
			grpc.ChainStreamInterceptor(s.generateBootstrapStreamInterceptor),
		)...)
		for _, ps := range s._serviceSetting {
			// v0.9.9+支持service仅用于http
			if !ps.HttpOnly {
				grpcServer.RegisterService(ps.Desc, ps.Impl)
			}
		}
		for _, invoke := range s._grpcServerInvoke {
			invoke(grpcServer)
		}
		// 注册grpc服务.暂时没有保护机制,由回调函数确保panic安全.
		if s.onRegisterGrpcService != nil && s.config.Name != "" {
			s.onRegisterGrpcService(s.config.Name, s.config.GrpcAddr, false)
			defer s.onRegisterGrpcService(s.config.Name, s.config.GrpcAddr, true)
		}
	}

	if s.config.HttpAddr != "" {
		// 如果还要其他http server配置,请在Config添加
		httpServer = &http.Server{
			Addr:         s.config.HttpAddr,
			Handler:      &s.mux,
			ReadTimeout:  s.config.HttpReadTimeout,
			WriteTimeout: s.config.HttpWriteTimeout,
			IdleTimeout:  s.config.HttpIdleTimeout,
		}
		if s.config.HttpKeepAlive > 0 {
			httpServer.SetKeepAlivesEnabled(true)
		} else if s.config.HttpKeepAlive < 0 {
			httpServer.SetKeepAlivesEnabled(false)
		}

		// 注册http服务.暂时没有保护机制,由回调函数确保panic安全.
		if s.onRegisterHttpService != nil && s.config.Name != "" {
			s.onRegisterHttpService(s.config.Name, s.config.HttpAddr, false)
			defer s.onRegisterHttpService(s.config.Name, s.config.HttpAddr, true)
		}
	}

	// 启动GRPC服务器
	if grpcServer != nil {
		grpcListener, err = graceGrpcListener(s.config.GrpcAddr)
		if err != nil {
			log.Error("grace grpc listener error: %v", err)
			return
		}
		// 异步避免阻塞
		go func(grpcServer *grpc.Server, grpcListener net.Listener) {
			if xrr := grpcServer.Serve(grpcListener); err != nil {
				log.Error("grpc serve error: %v", xrr)
			}
		}(grpcServer, grpcListener)
	}

	// 启动HTTP服务器
	if httpServer != nil {
		httpListener, err = graceHttpListener(s.config.HttpAddr, s.config.HttpKeepAlive)
		if err != nil {
			log.Error("grace http listener error: %v", err)
			return
		}
		// 异步避免阻塞
		go func(httpServer *http.Server, httpListener net.Listener) {
			if s.config.HttpCertFile != "" {
				if xrr := httpServer.ServeTLS(httpListener, s.config.HttpCertFile, s.config.HttpKeyFile); xrr != nil {
					log.Error("http serve tls error: %v", xrr)
				}
			} else {
				if xrr := httpServer.Serve(httpListener); xrr != nil {
					log.Error("http serve error: %v", xrr)
				}
			}
		}(httpServer, httpListener)
	}

	// 整理不再使用的内存变量
	s.clean()

	// 正常启动回调机制
	if len(s.onReady) > 0 {
		for _, f := range s.onReady {
			protect(f)
		}
	}
	// 等待信号,优雅关闭或重启服务. 明确关闭HTTP服务器,返回response添加Connection:closed
	graceShutdownOrRestart(grpcServer, grpcListener, httpServer, httpListener, &s.mux.closed)

	return
}

/**********************************************
* 流式拦截链尾部控制整体执行流程, 包括ErrorResult及Localize处理.
* 流式拦截链尾部必须位于拦截链尾位置(即通过grpc.ChainStreamInterceptor设置).
 **********************************************/
func (s *Server) generateBootstrapStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	var ctx = ss.Context()
	var setting = s.settings[info.FullMethod]
	if setting == nil {
		return StatusError(http.StatusNotFound, uint32(codes.NotFound), "Meta not found: %v", info.FullMethod)
	}
	defer func(meta *MethodSetting, ctx context.Context, grpcPanicFunc GrpcPanicFunc) {
		if p := recover(); p != nil {
			err = grpcPanicFunc(meta, ctx, p)
		}
	}(setting, ctx, s.grpcPanicFunc)
	var aspects = setting.Service.Aspects

	var idx = -1
	for _, s := range aspects {
		idx++
		if ctx, err = s.Before(setting, ctx, nil); err != nil {
			goto __AFTER__
		}
	}

	// 业务调用
	err = handler(srv, ss)

__AFTER__:
	for idx >= 0 {
		ctx, _, err = aspects[idx].After(setting, ctx, nil, nil, err)
		idx--
	}

	// 错误转换(grpc默认关闭i18n追求更快性能)
	if err != nil && s.config.GrpcI18nError {
		err = i18nGrpcError(ctx, err)
	}

	// 结果返回
	return
}

/**********************************************
* 一元拦截链尾部控制整体执行流程, 包括ErrorResult及Localize处理.
* 一元拦截链尾部必须位于拦截链尾位置(即通过grpc.ChainUnaryInterceptor设置).
 **********************************************/
func (s *Server) generateBootstrapUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {

	var setting = s.settings[info.FullMethod]
	if setting == nil {
		return nil, StatusError(http.StatusNotFound, uint32(codes.NotFound), "Meta not found: %v", info.FullMethod)
	}
	defer func(meta *MethodSetting, ctx context.Context, grpcPanicFunc GrpcPanicFunc) {
		if p := recover(); p != nil {
			err = grpcPanicFunc(meta, ctx, p)
		}
	}(setting, ctx, s.grpcPanicFunc)
	var aspects = setting.Service.Aspects

	var idx = -1
	for _, s := range aspects {
		idx++
		if ctx, err = s.Before(setting, ctx, req); err != nil {
			goto __AFTER__
		}
	}

	// 语法校验
	if vd, ok := req.(MessageValidator); ok {
		if err = vd.Validate(setting, ctx); err != nil {
			goto __AFTER__
		}
	}
	// 业务调用
	rsp, err = handler(ctx, req)

__AFTER__:
	for idx >= 0 {
		ctx, rsp, err = aspects[idx].After(setting, ctx, req, rsp, err)
		idx--
	}

	// 错误转换(grpc默认关闭i18n追求更快性能)
	if err != nil && s.config.GrpcI18nError {
		err = i18nGrpcError(ctx, err)
	}

	// 结果返回
	return
}

func i18nGrpcError(c context.Context, err error) error {

	if result, ok := err.(*StatusResult); ok {
		if md, ok := metadata.FromIncomingContext(c); ok {
			if vs, ok := md["accept-language"]; ok {
				if resMap := fastGetResMapByAcceptLanguage(vs[0]); resMap != nil {
					if rs, ok := resMap[result.Code]; ok {
						if len(result.Details) == 0 {
							result.Message = rs.Message
						} else {
							result.Message = Sprintf(rs.Message, result.Details...)
						}
					}
				}
			}
		}
	} else if sta, ok := status.FromError(err); ok {
		if md, ok := metadata.FromIncomingContext(c); ok {
			if vs, ok := md["accept-language"]; ok {
				if resMap := fastGetResMapByAcceptLanguage(vs[0]); resMap != nil {
					if rs, ok := resMap[result.Code]; ok {
						if len(result.Details) == 0 {
							err = status.Error(sta.Code(), rs.Message)
						} else {
							err = status.Error(sta.Code(), Sprintf(rs.Message, result.Details...))
						}
					}
				}
			}
		}
	}

	return err
}

/**********************************************
	默认实现
 **********************************************/

var defaultHttpPanicHandler = &Handler{
	HandleChain: []HandleFunc{
		func(ctx *Context) {
			log.Error("panic: %+v\n%v", ctx.panic, StackTrace(2, "\n"))
			_ = ctx.WriteErrorResult(StatusError(http.StatusInternalServerError, http.StatusInternalServerError, "internal server error: %+v", ctx.panic))
		},
	},
}

var defaultHttpNotFoundHandler = &Handler{
	HandleChain: []HandleFunc{
		func(ctx *Context) {
			_ = ctx.WriteErrorResult(StatusError(http.StatusNotFound, http.StatusNotFound, "not found"))
		},
	},
}

// GrpcPanicFunc grpc panic处理函数
type GrpcPanicFunc func(meta *MethodSetting, ctx context.Context, p interface{}) error

func defaultGrpcPanicFunc(meta *MethodSetting, ctx context.Context, p interface{}) error {
	log.Error("panic: %+v\n%v", p, StackTrace(2, "\n"))
	return status.Error(codes.Internal, fmt.Sprintf("panic: %v", p))
}

func protect(f func()) {
	defer func() {
		if per := recover(); per != nil {
			log.Error("panic: %+v\n%v", per, StackTrace(1, "\n"))
		}
	}()
	f()
}
