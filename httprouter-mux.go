package protoapi

import (
	"net/http"
	"strconv"
	"sync"
)

const (
	HandleChainCapacity = 128 // HandleChain容许的最大长度上限
	InsensitiveCapacity = 256 // 用于忽略大小写的查找过程
)

/*
请在调用route()完毕后务必调用done()对pool进行初始化
*/
type mux struct {
	nodeGet      *node
	nodeHead     *node
	nodePost     *node
	nodePut      *node
	nodePatch    *node
	nodeDelete   *node
	nodeConnect  *node
	nodeOptions  *node
	nodeTrace    *node
	httpPanic    *Handler
	httpNotFound *Handler
	upgrader     *Upgrader
	contexts     sync.Pool
	maxParams    uint16
	maxSections  uint16
	closed       uint32 // HTTP服务器关闭后返回的response添加Connection:closed.防止keep-alive影响!
}

var _ http.Handler = (*mux)(nil)

func (m *mux) root(method string) *node {
	switch method {
	case http.MethodPost:
		return m.nodePost
	case http.MethodGet:
		return m.nodeGet
	case http.MethodPut:
		return m.nodePut
	case http.MethodDelete:
		return m.nodeDelete
	case http.MethodHead:
		return m.nodeHead
	case http.MethodConnect:
		return m.nodeConnect
	case http.MethodOptions:
		return m.nodeOptions
	case http.MethodPatch:
		return m.nodePatch
	case http.MethodTrace:
		return m.nodeTrace
	}
	return nil
}

func (m *mux) must(method string) *node {
	switch method {
	case http.MethodPost:
		if m.nodePost == nil {
			m.nodePost = new(node)
		}
		return m.nodePost
	case http.MethodGet:
		if m.nodeGet == nil {
			m.nodeGet = new(node)
		}
		return m.nodeGet
	case http.MethodPut:
		if m.nodePut == nil {
			m.nodePut = new(node)
		}
		return m.nodePut
	case http.MethodDelete:
		if m.nodeDelete == nil {
			m.nodeDelete = new(node)
		}
		return m.nodeDelete
	case http.MethodHead:
		if m.nodeHead == nil {
			m.nodeHead = new(node)
		}
		return m.nodeHead
	case http.MethodConnect:
		if m.nodeConnect == nil {
			m.nodeConnect = new(node)
		}
		return m.nodeConnect
	case http.MethodOptions:
		if m.nodeOptions == nil {
			m.nodeOptions = new(node)
		}
		return m.nodeOptions
	case http.MethodPatch:
		if m.nodePatch == nil {
			m.nodePatch = new(node)
		}
		return m.nodePatch
	case http.MethodTrace:
		if m.nodeTrace == nil {
			m.nodeTrace = new(node)
		}
		return m.nodeTrace
	}
	panic("invalid method: " + method)
}

// 除panic与notfound外, 其他必须通过RouterConfig的相关方法添加, 实现conf配置化预处理
func (m *mux) route(method string, path string, node *Handler) {
	// !关键: 保证Context.Next()可以正常执行的必要条件!
	if len(node.HandleChain) >= HandleChainCapacity {
		panic("handle chain exceed maximum length: " + strconv.Itoa(HandleChainCapacity))
	}
	m.must(method).addRoute(path, node)
	if maxParams := countParams(path); maxParams > m.maxParams {
		m.maxParams = maxParams
	}
	if maxSections := countSections(path); maxSections > m.maxSections {
		m.maxSections = maxSections
	}
}

// 执行ServeHTTP()前必须调用serveINIT()初始化context pool!
func (m *mux) initServeHTTP() {
	m.contexts.New = func() interface{} {
		c := &Context{
			mux:            m,
			params:         make([]Param, 0, m.maxParams),
			skippedNodes:   make([]skippedNode, 0, m.maxSections),
			cipath:         make([]byte, 0, InsensitiveCapacity),
			cibuff:         [4]byte{},
			ResponseWriter: new(proxyResponseWriter),
		}
		c.stream.c = c //相互引用
		return c
	}
}

func (m *mux) HttpPanicFunc(f HandleFunc) {
	m.HttpPanic(&Handler{
		Meta:        Meta(http.StatusInternalServerError, Http_simple),
		HandleChain: []HandleFunc{f},
	})
}

func (m *mux) HttpPanic(h *Handler) {
	// 避免空指针
	if h != nil {
		m.httpPanic = h
	}
}

func (m *mux) HttpNotFoundFunc(f HandleFunc) {
	m.HttpNotFound(&Handler{
		Meta:        Meta(http.StatusNotFound, Http_simple),
		HandleChain: []HandleFunc{f},
	})
}

func (m *mux) HttpNotFound(h *Handler) {
	if h != nil {
		m.httpNotFound = h
	}
}

// 明确关闭HTTP服务器,返回response添加Connection:closed
func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 在router查找handler, 获取params及handlers, 设置到context执行, 最后将context归还
	c := m.contexts.Get().(*Context)
	defer func() {
		// 异常处理
		if p := recover(); p != nil {
			c.panicNext(p, m.httpPanic)
		}
		// 缓存前清理外部引用保证及时释放内存
		m.contexts.Put(c.clean())
	}()

	// 查找handle. 需要处理RedirectTrailingSlash,RedirectFixedPath,HandleMethodNotAllowed,
	if t := m.root(r.Method); t != nil {
		path := r.URL.Path
		c.params = c.params[0:0]
		c.skippedNodes = c.skippedNodes[0:0]
		_ = t.getValue(path, c)
		if c.Handler == nil && path != "/" {
			// 尝试忽略大小写的情况
			// 临时变量初始化. 缓存在context减少不必要GC!
			c.cipath = c.cipath[0:0]
			c.cibuff[0], c.cibuff[1], c.cibuff[2], c.cibuff[3] = 0, 0, 0, 0
			// 保证足够大小
			if mx := len(path); mx >= cap(c.cipath) {
				c.cipath = make([]byte, 0, mx+1)
			}
			// 尝试忽略大小写等情况
			if fixedPath := UnsafeString(t.findCaseInsensitivePathRec(
				cleanPath(path),
				c.cipath,
				c.cibuff,
				true,
			)); len(fixedPath) > 0 && fixedPath != path {
				c.params = c.params[0:0]
				c.skippedNodes = c.skippedNodes[0:0]
				_ = t.getValue(fixedPath, c) // 此处可以使用非安全字串
			}
		}
	}

	if c.Handler != nil {
		if c.Handler.BodyMaxBytes > 0 {
			// 根据MaxBodyBytes限制请求体大小. 保护request解析过程被恶意攻击!
			r.Body = http.MaxBytesReader(w, r.Body, c.Handler.BodyMaxBytes)
			r.ContentLength = c.Handler.BodyMaxBytes
		}
	} else {
		c.Handler = m.httpNotFound
	}

	// 处理请求
	c.serveNext(w, r, c.Handler)

}
