package protoapi

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestSetting struct {
	Filters []HandleFunc // 过滤规则, 通过server api注册的间接HandleFunc
	Plugins []HandleFunc // 插件规则, 通过Config注册的HandleFunc
	Handler *Handler     // 处理句柄
}

type routerGroup struct {
	path     string
	filters  []HandleFunc
	settings map[string]map[string]*RequestSetting
	children []*routerGroup
}

func newRouterGroup(path string, filters ...HandleFunc) *routerGroup {
	return &routerGroup{
		path:     path,
		filters:  filters,
		settings: make(map[string]map[string]*RequestSetting),
	}
}

func (rc *routerGroup) Group(path string, hs ...HandleFunc) *routerGroup {
	child := newRouterGroup(path, hs...)
	rc.children = append(rc.children, child)
	return child
}

func (rc *routerGroup) Use(filters ...HandleFunc) *routerGroup {
	rc.filters = append(rc.filters, filters...)
	return rc
}

// Deprecated: 请改用HandleFunc()!
func (rc *routerGroup) Handle(method string, path string, hs ...HandleFunc) *routerGroup {
	return rc.HandleFunc(method, path, hs...)
}

func (rc *routerGroup) HandleFunc(method string, path string, hs ...HandleFunc) *routerGroup {
	return rc.Handler(method, &Handler{
		Path:        path,
		HandleChain: hs,
	})
}

func (rc *routerGroup) Handler(method string, hd *Handler) *routerGroup {
	nodes, ok := rc.settings[method]
	if !ok {
		nodes = make(map[string]*RequestSetting)
		rc.settings[method] = nodes
	} else if v, ok := nodes[hd.Path]; ok {
		panic(fmt.Sprintf("handle duplicate path: %v %+v, new: %+v, old: %+v", method, hd.Path, hd, v))
	}
	nodes[hd.Path] = &RequestSetting{
		Handler: hd,
	}
	return rc

}

func (rc *routerGroup) Any(path string, f ...HandleFunc) *routerGroup {
	rc.HandleFunc(http.MethodGet, path, f...)
	rc.HandleFunc(http.MethodPost, path, f...)
	rc.HandleFunc(http.MethodPut, path, f...)
	rc.HandleFunc(http.MethodDelete, path, f...)
	rc.HandleFunc(http.MethodHead, path, f...)
	rc.HandleFunc(http.MethodPatch, path, f...)
	rc.HandleFunc(http.MethodOptions, path, f...)
	rc.HandleFunc(http.MethodConnect, path, f...)
	rc.HandleFunc(http.MethodTrace, path, f...)
	return rc
}
func (rc *routerGroup) GET(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodGet, path, f...)
}
func (rc *routerGroup) POST(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodPost, path, f...)
}
func (rc *routerGroup) PUT(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodPut, path, f...)
}
func (rc *routerGroup) DELETE(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodDelete, path, f...)
}
func (rc *routerGroup) HEAD(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodHead, path, f...)
}
func (rc *routerGroup) PATCH(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodPatch, path, f...)
}
func (rc *routerGroup) OPTIONS(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodOptions, path, f...)
}
func (rc *routerGroup) CONNECT(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodConnect, path, f...)
}
func (rc *routerGroup) TRACE(path string, f ...HandleFunc) *routerGroup {
	return rc.HandleFunc(http.MethodTrace, path, f...)
}

func (rc *routerGroup) StaticFile(path string, file string) *routerGroup {
	if strings.Contains(path, ":") || strings.Contains(path, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := &Handler{
		Path: path,
		HandleChain: []HandleFunc{
			func(c *Context) {
				// 底层直接使用http api适配
				http.ServeFile(c.ResponseWriter, c.Request, file)
			},
		},
	}
	rc.Handler(http.MethodGet, handler)
	rc.Handler(http.MethodHead, handler)
	return rc
}

func (rc *routerGroup) StaticDir(prefix string, dir string) *routerGroup {
	return rc.StaticFS(prefix, http.Dir(dir))
}

func (rc *routerGroup) StaticFS(prefix string, fs http.FileSystem) *routerGroup {
	if strings.Contains(prefix, ":") || strings.Contains(prefix, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	// 清除最后的"/", 避免影响前缀剔除操作
	plen := len(prefix)
	if last := plen - 1; last >= 0 && prefix[last] == '/' {
		prefix = prefix[:last]
		plen--
	}

	staticFs := http.FileServer(fs)
	handler := &Handler{
		Path: prefix + "/*file",
		HandleChain: []HandleFunc{
			func(c *Context) {
				// 必须是path前缀，否则404
				ulen := len(c.Request.URL.Path)
				if ulen == plen {
					// 重定向到带"/"结尾的路径
					c.ResponseWriter.Header().Set("Location", prefix+"/")
					c.ResponseWriter.WriteHeader(http.StatusMovedPermanently)
				} else if ulen > plen && c.Request.URL.Path[plen] == '/' {
					c.Request.URL.Path = c.Request.URL.Path[plen:]
					staticFs.ServeHTTP(c.ResponseWriter, c.Request)
				} else {
					http.NotFound(c.ResponseWriter, c.Request)
				}
			},
		},
	}

	rc.Handler(http.MethodGet, handler)
	rc.Handler(http.MethodHead, handler)

	return rc
}

// Flatten 返回{httpMethod:{path:setting}}
func (rc *routerGroup) Flatten() (ret map[string]map[string]*RequestSetting) {
	ret = make(map[string]map[string]*RequestSetting)
	flatten(ret, rc, "", nil)
	return
}

// flatten 递归展开gro
func flatten(ret map[string]map[string]*RequestSetting, group *routerGroup, prefix string, filters []HandleFunc) {
	if group.path != "" {
		prefix = prefix + group.path
	}
	if len(group.filters) > 0 {
		filters = append(filters, group.filters...)
	}
	for method, pathNodeMap := range group.settings {
		nodes := ret[method]
		if nodes == nil {
			nodes = make(map[string]*RequestSetting)
		}
		for path, node := range pathNodeMap {
			node.Filters = filters
			nodes[prefix+path] = node
		}
		ret[method] = nodes
	}
	// 递归遍历子结点
	for _, child := range group.children {
		flatten(ret, child, prefix, filters)
	}
}
