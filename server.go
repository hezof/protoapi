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
func (svr *Server) clean() {
	svr._group = nil
	svr._serviceSetting = nil
	svr._serviceAspect = nil
	svr._servicePlugin = nil
	svr._requestPlugin = nil
	svr._grpcServerOption = nil
	svr._httpServerOption = nil
	svr._grpcServerInvoke = nil
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

func (svr *Server) Config() *Config {
	return svr.config
}

/**********************************************
	生命周期回调
 **********************************************/

// OnInit 服务器启动时回调
func (svr *Server) OnInit(fs ...func()) {
	svr.onInit = append(svr.onInit, fs...)
}

// OnReady 服务器完成时回调
func (svr *Server) OnReady(fs ...func()) {
	svr.onReady = append(svr.onReady, fs...)
}

// OnExit 服务器退出时回调
func (svr *Server) OnExit(fs ...func()) {
	svr.onExit = append(svr.onExit, fs...)
}

// OnRegisterGrpcService 注册GRPC服务回调
func (svr *Server) OnRegisterGrpcService(f func(name, addr string, undo bool)) {
	svr.onRegisterGrpcService = f
}

// OnRegisterHttpService 注册HTTP服务回调
func (svr *Server) OnRegisterHttpService(f func(name, addr string, undo bool)) {
	svr.onRegisterHttpService = f
}

/**********************************************
	配置添加
 **********************************************/

func (svr *Server) ServiceAspect(vs ...ServiceAspect) {
	svr._serviceAspect = append(svr._serviceAspect, vs...)
}

func (svr *Server) ServicePlugin(vs ...ServicePlugin) {
	svr._servicePlugin = append(svr._servicePlugin, vs...)
}

func (svr *Server) RequestPlugin(vs ...RequestPlugin) {
	svr._requestPlugin = append(svr._requestPlugin, vs...)
}

func (svr *Server) GrpcServerOption(vs ...grpc.ServerOption) {
	svr._grpcServerOption = append(svr._grpcServerOption, vs...)
}

func (svr *Server) HttpServerOption(vs ...HandleFunc) {
	svr._httpServerOption = append(svr._httpServerOption, vs...)
}

func (svr *Server) GrpcServerInvoke(vs ...func(server *grpc.Server)) {
	svr._grpcServerInvoke = append(svr._grpcServerInvoke, vs...)
}

func (svr *Server) GrpcPanicFunc(f GrpcPanicFunc) {
	// 不能为空
	if f != nil {
		svr.grpcPanicFunc = f
	}
}

/**********************************************
	服务注册
 **********************************************/

func (svr *Server) RegisterService(registry ServiceRegistry, implement interface{}, aspects ...ServiceAspect) *Server {
	aspects = orderServiceAspects(svr._serviceAspect, aspects)
	serviceSetting := registry(implement, aspects)
	for _, methodSetting := range serviceSetting.Methods {
		methodSetting.Service = serviceSetting
	}
	svr._serviceSetting = append(svr._serviceSetting, serviceSetting)
	return svr
}

/**********************************************
	启动监听
 **********************************************/

func (svr *Server) ListenAndServe() (err error) {

	/********************************************************
	* 如果没有配置grpc或http地址则自动结束服务初始化流程!
	 ********************************************************/
	if svr.config.GrpcAddr == "" && svr.config.HttpAddr == "" {
		return
	}

	/*******************************************************
	* 全局组件注册
	* 该步骤主要用于分离服务层与访问层,从而支持更灵活的组件部署模式
	 *******************************************************/
	for _, c := range _globalRegister {
		svr.RegisterService(c.Registry, c.Implement, c.Aspects...)
	}

	/********************************************************
	* 处理ServiceSetting的插件(统计/修改/添加/删除)等
	 ********************************************************/
	for _, p := range svr._servicePlugin {
		p(&svr._serviceSetting)
	}

	/********************************************************
	* 添加MethodSetting相应的AccessSetting以及MethodGroup
	 ********************************************************/
	for _, ss := range svr._serviceSetting {
		for _, ms := range ss.Methods {
			svr.settings[FullMethod(ms.Meta)] = ms
			// 添加method相应的RequestSetting
			if ms.Meta.Http.Get != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodGet,
					Path:        ms.Meta.Http.Get,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Put != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodPut,
					Path:        ms.Meta.Http.Put,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Post != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodPost,
					Path:        ms.Meta.Http.Post,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Delete != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodDelete,
					Path:        ms.Meta.Http.Delete,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Options != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodOptions,
					Path:        ms.Meta.Http.Options,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Head != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodHead,
					Path:        ms.Meta.Http.Head,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Patch != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodPatch,
					Path:        ms.Meta.Http.Patch,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Trace != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodTrace,
					Path:        ms.Meta.Http.Trace,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Connect != "" {
				svr._group.Handle(&Handler{
					Setting:     ms,
					Method:      http.MethodConnect,
					Path:        ms.Meta.Http.Connect,
					Status:      ms.Meta.Http.Status,
					Result:      ms.Meta.Http.Result,
					HandleChain: []HandleFunc{RestfulHandleFunc},
				})
			}
			if ms.Meta.Http.Websocket != "" {
				if svr.mux.upgrader == nil {
					svr.mux.upgrader = newWebsocketUpgrader(svr.config)
				}
				svr._group.Handle(&Handler{
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
	_requestSetting := svr._group.Flatten()

	/********************************************************
	* 处理RequestSetting的插件(缓存,日志,代理)等
	 ********************************************************/
	for _, p := range svr._requestPlugin {
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
			setting.Handler.Status = NvlI(setting.Handler.Status, profile.DefaultApplyStatus, http.StatusOK) // 很关键: 不能为0!
			setting.Handler.BodyMaxBytes = NvlI(setting.Handler.BodyMaxBytes, svr.config.HttpBodyMaxBytes, profile.HttpBodyMaxBytes, MAX_MEM)
			setting.Handler.FormMaxMemory = NvlI(setting.Handler.FormMaxMemory, svr.config.HttpFormMaxMemory, profile.HttpFormMaxMemory, MAX_MEM)
			setting.Handler.HandleChain = Join(svr._httpServerOption, setting.Plugins, setting.Filters, setting.Handler.HandleChain)
			svr.mux.route(method, path, setting.Handler) // 正式注册到mux
		}
	}
	svr.mux.initServeHTTP()

	/********************************************************
	* 初始与退出钩子机制.
	 ********************************************************/
	if len(svr.onInit) > 0 {
		for _, f := range svr.onInit {
			protect(f)
		}
	}
	if len(svr.onExit) > 0 {
		// 此处defer没用闭包
		defer func(exit []func()) {
			for _, f := range exit {
				protect(f)
			}
		}(svr.onExit)
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
	if svr.config.GrpcAddr != "" {
		var opts []grpc.ServerOption
		if len(svr._grpcServerOption) > 0 {
			opts = append(opts, svr._grpcServerOption...)
		}
		if svr.config.GrpcKeepAlive > 0 {
			opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{Time: svr.config.GrpcKeepAlive}))
		}
		if svr.config.GrpcKeepAlivePolicy > 0 {
			opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{MinTime: svr.config.GrpcKeepAlivePolicy, PermitWithoutStream: true}))
		}
		// 很关键: 必须确保"boostrap interceptor"位于最后位置, 即必须是grpc server最后添加的ServerOption
		grpcServer = grpc.NewServer(append(opts,
			grpc.ChainUnaryInterceptor(svr.bootstrapUnaryInterceptor),
			grpc.ChainStreamInterceptor(svr.bootstrapStreamInterceptor),
		)...)
		for _, ps := range svr._serviceSetting {
			// v0.9.9+支持service仅用于http
			if !ps.HttpOnly {
				grpcServer.RegisterService(ps.Desc, ps.Impl)
			}
		}
		for _, invoke := range svr._grpcServerInvoke {
			invoke(grpcServer)
		}
		// 注册grpc服务.暂时没有保护机制,由回调函数确保panic安全.
		if svr.onRegisterGrpcService != nil && svr.config.Name != "" {
			svr.onRegisterGrpcService(svr.config.Name, svr.config.GrpcAddr, false)
			defer svr.onRegisterGrpcService(svr.config.Name, svr.config.GrpcAddr, true)
		}
	}

	if svr.config.HttpAddr != "" {
		// 如果还要其他http server配置,请在Config添加
		httpServer = &http.Server{
			Addr:         svr.config.HttpAddr,
			Handler:      &svr.mux,
			ReadTimeout:  svr.config.HttpReadTimeout,
			WriteTimeout: svr.config.HttpWriteTimeout,
			IdleTimeout:  svr.config.HttpIdleTimeout,
		}
		if svr.config.HttpKeepAlive > 0 {
			httpServer.SetKeepAlivesEnabled(true)
		} else if svr.config.HttpKeepAlive < 0 {
			httpServer.SetKeepAlivesEnabled(false)
		}

		// 注册http服务.暂时没有保护机制,由回调函数确保panic安全.
		if svr.onRegisterHttpService != nil && svr.config.Name != "" {
			svr.onRegisterHttpService(svr.config.Name, svr.config.HttpAddr, false)
			defer svr.onRegisterHttpService(svr.config.Name, svr.config.HttpAddr, true)
		}
	}

	// 启动GRPC服务器
	if grpcServer != nil {
		grpcListener, err = graceGrpcListener(svr.config.GrpcAddr)
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
		httpListener, err = graceHttpListener(svr.config.HttpAddr, svr.config.HttpKeepAlive)
		if err != nil {
			log.Error("grace http listener error: %v", err)
			return
		}
		// 异步避免阻塞
		go func(httpServer *http.Server, httpListener net.Listener) {
			if svr.config.HttpCertFile != "" {
				if xrr := httpServer.ServeTLS(httpListener, svr.config.HttpCertFile, svr.config.HttpKeyFile); xrr != nil {
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
	svr.clean()

	// 正常启动回调机制
	if len(svr.onReady) > 0 {
		for _, f := range svr.onReady {
			protect(f)
		}
	}
	// 等待信号,优雅关闭或重启服务. 明确关闭HTTP服务器,返回response添加Connection:closed
	graceShutdownOrRestart(grpcServer, grpcListener, httpServer, httpListener, &svr.mux.closed)

	return
}

/**********************************************
* 流式拦截链尾部控制整体执行流程, 包括ErrorResult及Localize处理.
* 流式拦截链尾部必须位于拦截链尾位置(即通过grpc.ChainStreamInterceptor设置).
 **********************************************/
func (svr *Server) bootstrapStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	set := svr.settings[info.FullMethod]
	ctx := ss.Context()
	if set == nil {
		return StatusError(http.StatusNotFound, uint32(codes.NotFound), "Meta not found: %v", info.FullMethod)
	}
	defer func(set *MethodSetting, ctx context.Context, grpcPanicFunc GrpcPanicFunc) {
		if p := recover(); p != nil {
			err = grpcPanicFunc(set, ctx, p)
		}
	}(set, ctx, svr.grpcPanicFunc)

	// 前置处理
	idx, ctx, err := BeforeAspect(set, ctx, nil)

	// 业务调用
	if err == nil {
		err = handler(svr, ss)
	}

	// 后置处理
	_, err = AfterAspect(set, idx, ctx, nil, nil, err)

	// 错误转换(grpc默认关闭i18n追求更快性能)
	if err != nil && hasResMap && svr.config.GrpcI18nError {
		err = i18nGrpcError(ctx, err)
	}

	// 结果返回
	return
}

/**********************************************
* 一元拦截链尾部控制整体执行流程, 包括ErrorResult及Localize处理.
* 一元拦截链尾部必须位于拦截链尾位置(即通过grpc.ChainUnaryInterceptor设置).
 **********************************************/
func (svr *Server) bootstrapUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {

	set := svr.settings[info.FullMethod]
	if set == nil {
		return nil, StatusError(http.StatusNotFound, uint32(codes.NotFound), "Meta not found: %v", info.FullMethod)
	}
	defer func(set *MethodSetting, ctx context.Context, grpcPanicFunc GrpcPanicFunc) {
		if p := recover(); p != nil {
			err = grpcPanicFunc(set, ctx, p)
		}
	}(set, ctx, svr.grpcPanicFunc)

	// 前置处理
	idx, ctx, err := BeforeAspect(set, ctx, req)

	// 业务调用
	if err == nil {
		rsp, err = handler(ctx, req)
	}

	// 后置处理
	rsp, err = AfterAspect(set, idx, ctx, req, rsp, err)

	// 错误转换(grpc默认关闭i18n追求更快性能)
	if err != nil && hasResMap && svr.config.GrpcI18nError {
		err = i18nGrpcError(ctx, err)
	}

	// 结果返回
	return
}

func i18nGrpcError(c context.Context, err error) error {

	if result, ok := err.(StatusResult); ok {
		if md, ok := metadata.FromIncomingContext(c); ok {
			if vs, ok := md["accept-language"]; ok {
				if resMap := fastGetResMapByAcceptLanguage(vs[0]); resMap != nil {
					if rs, ok := resMap[result.GetCode()]; ok {
						result.SetMessage(rs.Message)
					}
				}
			}
		}
	} else if sta, ok := status.FromError(err); ok {
		if md, ok := metadata.FromIncomingContext(c); ok {
			if vs, ok := md["accept-language"]; ok {
				if resMap := fastGetResMapByAcceptLanguage(vs[0]); resMap != nil {
					if rs, ok := resMap[result.GetCode()]; ok {
						err = status.Error(sta.Code(), rs.Message)
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
			_ = ctx.WriteErrorResult(StatusError(http.StatusInternalServerError, http.StatusInternalServerError, fmt.Sprintf("internal server error: %+v", ctx.panic)))
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
type GrpcPanicFunc func(set *MethodSetting, ctx context.Context, p interface{}) error

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
