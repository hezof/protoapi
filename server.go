package protoapi

import (
	"context"
	"fmt"
	"github.com/hezof/protoapi/internal/websocket"
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
	onInit                []func()                           // 初始前回调
	onReady               []func()                           // 启动前回调
	onExit                []func()                           // 退出前回调
	onRegisterGrpcService func(name, addr string, undo bool) // grpc服务注册回调
	onRegisterHttpService func(name, addr string, undo bool) // http服务注册回调
	/*
		下述过程数据在服务最终启动前应该清除释放内存!
	*/
	*routerGroup                                            // 分组路由(过程数据)
	_serviceSetting   []*ServiceSetting                     // 服务设置(过程数据)
	_methodSetting    map[string]*MethodSetting             // 方法设置(过程数据)
	_requestSetting   map[string]map[string]*RequestSetting // 平面后的规则(过程数据)
	_serviceAspect    []ServiceAspect                       // 服务校验器(过程参数)
	_servicePlugin    []ServicePlugin                       // 服务插件(过程数据)
	_requestPlugin    []RequestPlugin                       // 请求插件(过程数据)
	_grpcServerOption []grpc.ServerOption                   // grpc拦截器(过程数据)
	_httpServerOption []HandleFunc                          // http拦截器(过程数据)
	_grpcServerInvoke []func(server *grpc.Server)           // grpc服务器回调(过程数据)
	_grpcPanicFunc    GrpcPanicFunc                         // Grpc panic函数
}

// 整理不再使用的内存变量
func (s *Server) clean() {
	s.routerGroup = nil
	s._serviceSetting = nil
	s._methodSetting = nil
	s._requestSetting = nil
	s._serviceAspect = nil
	s._servicePlugin = nil
	s._requestPlugin = nil
	s._grpcServerOption = nil
	s._httpServerOption = nil
	s._grpcServerInvoke = nil
	s._grpcPanicFunc = nil
}

func NewServer(c *Config) *Server {
	return &Server{
		mux: mux{
			httpPanic:    defaultHttpPanicHandler,
			httpNotFound: defaultHttpNotFoundHandler,
		},
		config:         mergeConfig(c),
		routerGroup:    newRouterGroup(""),
		_grpcPanicFunc: defaultGrpcPanicFunc,
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
		s._grpcPanicFunc = f
	}
}

/**********************************************
	服务注册
 **********************************************/

func (s *Server) RegisterService(registry ServiceRegistry, implement interface{}, aspects ...ServiceAspect) *Server {
	aspects = orderServiceAspects(s._serviceAspect, aspects)
	serviceSetting := registry(implement, aspects)
	// 设置parent指向. 后续处理依赖parent.
	for _, methodSetting := range serviceSetting.Methods {
		methodSetting.parent = serviceSetting
	}
	serviceSetting.Aspects = aspects
	s._serviceSetting = append(s._serviceSetting, serviceSetting)
	return s
}

/**********************************************
	启动监听
 **********************************************/

func (s *Server) protect(f func()) {
	defer func() {
		if per := recover(); per != nil {
			log.Error("panic: %+v\n%v", per, StackTrace(1, "\n"))
		}
	}()
	f()
}

func (s *Server) ListenAndServe() (err error) {

	/********************************************************
	* 如果没有配置grpc或http地址则自动结束服务初始化流程!
	 ********************************************************/
	if s.config.GrpcAddr == "" && s.config.HttpAddr == "" {
		return
	}

	/********************************************************
	* 初始与退出钩子机制.
	 ********************************************************/
	if len(s.onInit) > 0 {
		for _, f := range s.onInit {
			s.protect(f)
		}
	}
	if len(s.onExit) > 0 {
		// 此处defer没用闭包
		defer func(exit []func()) {
			for _, f := range exit {
				s.protect(f)
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
		upgrader     *websocket.Upgrader
	)
	defer func() {
		if grpcListener != nil {
			grpcListener.Close()
		}
		if httpListener != nil {
			httpListener.Close()
		}
	}()

	/*******************************************************
	* 全局服务组件Visitor
	* 该步骤主要用于分离服务层与访问层,从而支持更灵活的组件部署模式
	 *******************************************************/
	globalServiceVisitor(s)

	/********************************************************
	* 处理ServiceSetting的插件(统计/修改/添加/删除)等
	 ********************************************************/
	for _, p := range s._servicePlugin {
		p(&s._serviceSetting)
	}

	/********************************************************
	* 添加MethodSetting相应的AccessSetting以及MethodGroup
	 ********************************************************/
	s._methodSetting = make(map[string]*MethodSetting)
	for _, ss := range s._serviceSetting {
		for _, ms := range ss.Methods {
			// 计算全方法
			ms.FullMethod = fullMethod(ms.Package, ms.Service, ms.Method)
			// 添加方法设置
			s._methodSetting[ms.FullMethod] = ms

			// 添加method相应的RequestSetting
			if ms.Websocket != "" {
				if upgrader == nil {
					upgrader = newWebsocketUpgrader(s.config)
				}
				s.routerGroup.Handler(http.MethodGet, &Handler{
					Meta:         ms,
					Method:       http.MethodGet,
					Path:         ms.Websocket,
					HandleChain:  []HandleFunc{WebsocketHandleFunc(ms.Call, upgrader)},
					BodyMaxBytes: -1, // 如果是Websocket长链接则自动忽略BodyMaxBytes参数
				})
			}
			if ms.Post != "" {
				s.routerGroup.Handler(http.MethodPost, &Handler{
					Meta:        ms,
					Method:      http.MethodPost,
					Path:        ms.Post,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Get != "" {
				s.routerGroup.Handler(http.MethodGet, &Handler{
					Meta:        ms,
					Method:      http.MethodGet,
					Path:        ms.Get,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Put != "" {
				s.routerGroup.Handler(http.MethodPut, &Handler{
					Meta:        ms,
					Method:      http.MethodPut,
					Path:        ms.Put,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Delete != "" {
				s.routerGroup.Handler(http.MethodDelete, &Handler{
					Meta:        ms,
					Method:      http.MethodDelete,
					Path:        ms.Delete,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Options != "" {
				s.routerGroup.Handler(http.MethodOptions, &Handler{
					Meta:        ms,
					Method:      http.MethodOptions,
					Path:        ms.Options,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Head != "" {
				s.routerGroup.Handler(http.MethodHead, &Handler{
					Meta:        ms,
					Method:      http.MethodHead,
					Path:        ms.Head,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Patch != "" {
				s.routerGroup.Handler(http.MethodPatch, &Handler{
					Meta:        ms,
					Method:      http.MethodPatch,
					Path:        ms.Patch,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Trace != "" {
				s.routerGroup.Handler(http.MethodTrace, &Handler{
					Meta:        ms,
					Method:      http.MethodTrace,
					Path:        ms.Trace,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
			if ms.Connect != "" {
				s.routerGroup.Handler(http.MethodConnect, &Handler{
					Meta:        ms,
					Method:      http.MethodConnect,
					Path:        ms.Connect,
					HandleChain: []HandleFunc{RestfulHandleFunc(ms.Call)},
				})
			}
		}
	}

	/********************************************************
	* 扁平化router的AccessSetting: map{httpMethod: map{httpPath:requestSetting}}
	* 后面应用AccessOption时实现proxy, access, cache等的拦截逻辑
	 ********************************************************/
	s._requestSetting = s.routerGroup.Flatten()

	/********************************************************
	* 处理RequestSetting的插件(缓存,日志,代理)等
	 ********************************************************/
	for _, p := range s._requestPlugin {
		p(s._requestSetting)
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
	for method, pathSetting := range s._requestSetting {
		for path, setting := range pathSetting {
			handler := setting.Handler
			handler.Method = method
			handler.Path = path // NOTE: 通过martin-like注册的Handler可能没有path
			handler.Meta.Status = NvlI(handler.Meta.Status, profile.DefaultApplyStatus)
			handler.BodyMaxBytes = NvlI(setting.Handler.BodyMaxBytes, s.config.HttpBodyMaxBytes)
			handler.FormMaxMemory = NvlI(setting.Handler.FormMaxMemory, s.config.HttpFormMaxMemory)
			handler.HandleChain = Join(s._httpServerOption, setting.Plugins, setting.Filters, setting.Handler.HandleChain)
			s.mux.route(method, path, handler) // 正式注册到mux
		}
	}
	/********************************************************
	* 注意: 启用mux.ServeHTTP()之前必须调用mux.initServeHTTP()初始化
	 ********************************************************/
	s.mux.initServeHTTP()

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
			generateBootstrapUnaryInterceptor(s._methodSetting, s._grpcPanicFunc, s.config.GrpcI18nError),
			generateBootstrapStreamInterceptor(s._methodSetting, s._grpcPanicFunc, s.config.GrpcI18nError),
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
			Addr:    s.config.HttpAddr,
			Handler: &s.mux,
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
		go func() {
			var err error
			if err = grpcServer.Serve(grpcListener); err != nil {
				log.Error("grace grpc server error: %v", err)
			}
		}()
	}

	// 启动HTTP服务器
	if httpServer != nil {
		httpListener, err = graceHttpListener(s.config.HttpAddr, s.config.HttpKeepAlive)
		if err != nil {
			log.Error("grace http listener error: %v", err)
			return
		}
		// 异步避免阻塞
		go func() {
			var err error
			if s.config.HttpCertFile != "" {
				err = httpServer.ServeTLS(httpListener, s.config.HttpCertFile, s.config.HttpKeyFile)
			} else {
				err = httpServer.Serve(httpListener)
			}
			if err != nil {
				log.Error("grace http server error: %v", err)
			}
		}()
	}

	// 整理不再使用的内存变量
	s.clean()

	// 正常启动回调机制
	if len(s.onReady) > 0 {
		for _, f := range s.onReady {
			s.protect(f)
		}
	}
	// 等待信号,优雅关闭或重启服务. 明确关闭HTTP服务器,返回response添加Connection:closed
	graceShutdownOrRestart(grpcServer, grpcListener, httpServer, httpListener, &s.mux.closed)

	return
}

/**********************************************
	默认实现
 **********************************************/

var defaultHttpPanicHandler = &Handler{
	HandleChain: []HandleFunc{
		func(ctx *Context) {
			log.Error("panic: %+v\n%v", ctx.panic, StackTrace(2, "\n"))
			ctx.WriteErrorResult(StatusError(http.StatusInternalServerError, http.StatusInternalServerError, "internal server error: %+v", ctx.panic))
		},
	},
}

var defaultHttpNotFoundHandler = &Handler{
	HandleChain: []HandleFunc{
		func(ctx *Context) {
			ctx.WriteErrorResult(StatusError(http.StatusNotFound, http.StatusNotFound, "not found"))
		},
	},
}

// GrpcPanicFunc grpc panic处理函数
type GrpcPanicFunc func(meta *MethodSetting, ctx context.Context, p interface{}) error

func defaultGrpcPanicFunc(meta *MethodSetting, ctx context.Context, p interface{}) error {
	log.Error("panic: %+v\n%v", p, StackTrace(2, "\n"))
	return status.Error(codes.Internal, fmt.Sprintf("panic: %v", p))
}

/**********************************************
* 流式拦截链尾部控制整体执行流程, 包括ErrorResult及Localize处理.
* 流式拦截链尾部必须位于拦截链尾位置(即通过grpc.ChainStreamInterceptor设置).
 **********************************************/
func generateBootstrapStreamInterceptor(metas map[string]*MethodSetting, grpcPanicFunc GrpcPanicFunc, grpcI18nError bool) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var ctx = ss.Context()
		var meta = metas[info.FullMethod]
		if meta == nil {
			return StatusError(http.StatusNotFound, uint32(codes.NotFound), "meta not found: %v", info.FullMethod)
		}
		defer func(meta *MethodSetting, ctx context.Context, grpcPanicFunc GrpcPanicFunc) {
			if p := recover(); p != nil {
				err = grpcPanicFunc(meta, ctx, p)
			}
		}(meta, ctx, grpcPanicFunc)
		var aspects = meta.parent.Aspects

		var idx = -1
		for _, s := range aspects {
			idx++
			if ctx, err = s.Before(meta, ctx, nil); err != nil {
				goto __AFTER__
			}
		}

		// 业务调用
		err = handler(srv, ss)

	__AFTER__:
		for idx >= 0 {
			ctx, _, err = aspects[idx].After(meta, ctx, nil, nil, err)
			idx--
		}

		// 错误转换(grpc默认关闭i18n追求更快性能)
		if grpcI18nError && lenResMap > 0 && err != nil {
			err = i18nGrpcError(ctx, err)
		}

		// 结果返回
		return
	})
}

/**********************************************
* 一元拦截链尾部控制整体执行流程, 包括ErrorResult及Localize处理.
* 一元拦截链尾部必须位于拦截链尾位置(即通过grpc.ChainUnaryInterceptor设置).
 **********************************************/
func generateBootstrapUnaryInterceptor(metas map[string]*MethodSetting, grpcPanicFunc GrpcPanicFunc, grpcI18nError bool) grpc.ServerOption {

	// 实现采用goto确保控制流程清晰!
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {

		var meta = metas[info.FullMethod]
		if meta == nil {
			return nil, StatusError(http.StatusNotFound, uint32(codes.NotFound), "meta not found: %v", info.FullMethod)
		}
		defer func(meta *MethodSetting, ctx context.Context, grpcPanicFunc GrpcPanicFunc) {
			if p := recover(); p != nil {
				err = grpcPanicFunc(meta, ctx, p)
			}
		}(meta, ctx, grpcPanicFunc)
		var aspects = meta.parent.Aspects

		var idx = -1
		for _, s := range aspects {
			idx++
			if ctx, err = s.Before(meta, ctx, req); err != nil {
				goto __AFTER__
			}
		}

		// 语法校验
		if vd, ok := req.(MessageValidator); ok {
			if err = vd.Validate(ctx); err != nil {
				goto __AFTER__
			}
		}
		// 业务调用
		rsp, err = handler(ctx, req)

	__AFTER__:
		for idx >= 0 {
			ctx, rsp, err = aspects[idx].After(meta, ctx, req, rsp, err)
			idx--
		}

		// 错误转换(grpc默认关闭i18n追求更快性能)
		if grpcI18nError && lenResMap > 0 && err != nil {
			err = i18nGrpcError(ctx, err)
		}

		// 结果返回
		return
	})
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
