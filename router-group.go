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

type _group struct {
	path     string
	filters  []HandleFunc
	settings map[string]map[string]*RequestSetting
	children []*_group
}

func newGroup(path string, filters ...HandleFunc) *_group {
	return &_group{
		path:     path,
		filters:  filters,
		settings: make(map[string]map[string]*RequestSetting),
	}
}
func (rc *_group) BasePath() string {
	return rc.path
}

func (rc *_group) Group(path string, hs ...HandleFunc) GroupRouter {
	child := newGroup(path, hs...)
	rc.children = append(rc.children, child)
	return child
}

func (rc *_group) Use(filters ...HandleFunc) GroupRouter {
	rc.filters = append(rc.filters, filters...)
	return rc
}

func (rc *_group) HandleFunc(method string, path string, hs ...HandleFunc) GroupRouter {
	return rc.Handle(&Handler{
		Method:      method,
		Path:        path,
		HandleChain: hs,
	})
}

func (rc *_group) Handle(hd *Handler) GroupRouter {

	setting, ok := rc.settings[hd.Method]
	if !ok {
		setting = make(map[string]*RequestSetting)
		rc.settings[hd.Method] = setting
	} else if set, ok := setting[hd.Path]; ok {
		panic(fmt.Sprintf("handle duplicate path: %v %+v, new: %+v, old: %+v", hd.Method, hd.Path, hd, set))
	}
	setting[hd.Path] = &RequestSetting{
		Handler: hd,
	}
	return rc

}

func (rc *_group) Any(path string, f ...HandleFunc) GroupRouter {
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
func (rc *_group) GET(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodGet, path, f...)
}
func (rc *_group) POST(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodPost, path, f...)
}
func (rc *_group) PUT(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodPut, path, f...)
}
func (rc *_group) DELETE(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodDelete, path, f...)
}
func (rc *_group) HEAD(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodHead, path, f...)
}
func (rc *_group) PATCH(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodPatch, path, f...)
}
func (rc *_group) OPTIONS(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodOptions, path, f...)
}
func (rc *_group) CONNECT(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodConnect, path, f...)
}
func (rc *_group) TRACE(path string, f ...HandleFunc) GroupRouter {
	return rc.HandleFunc(http.MethodTrace, path, f...)
}

func (rc *_group) StaticFile(path string, file string) GroupRouter {
	if strings.Contains(path, ":") || strings.Contains(path, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	call := func(c *Context) {
		http.ServeFile(c.ResponseWriter.ResponseWriter, c.Request, file)
	}
	rc.HandleFunc(http.MethodGet, path, call)
	rc.HandleFunc(http.MethodHead, path, call)
	return rc
}

func (rc *_group) Static(prefix string, root string) GroupRouter {
	return rc.StaticFS(prefix, http.Dir(root))
}

func (rc *_group) StaticFS(prefix string, fs http.FileSystem) GroupRouter {
	if strings.Contains(prefix, ":") || strings.Contains(prefix, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	// 清除最后的"/", 避免影响前缀剔除操作
	plen := len(prefix)
	if last := plen - 1; last >= 0 && prefix[last] == '/' {
		prefix = prefix[:last]
		plen--
	}
	serv := http.FileServer(fs)
	path := prefix + "/*file"
	call := func(c *Context) {
		// 必须是path前缀，否则404
		ulen := len(c.Request.URL.Path)
		if ulen == plen {
			// 重定向到带"/"结尾的路径
			c.ResponseWriter.Header().Set("Location", prefix+"/")
			c.ResponseWriter.WriteHeader(http.StatusMovedPermanently)
		} else if ulen > plen && c.Request.URL.Path[plen] == '/' {
			c.Request.URL.Path = c.Request.URL.Path[plen:]
			serv.ServeHTTP(c.ResponseWriter.ResponseWriter, c.Request)
		} else {
			http.NotFound(c.ResponseWriter.ResponseWriter, c.Request)
		}
	}
	rc.HandleFunc(http.MethodGet, path, call)
	rc.HandleFunc(http.MethodHead, path, call)

	return rc
}

// Flatten 返回{httpMethod:{path:setting}}
func (rc *_group) Flatten() (ret map[string]map[string]*RequestSetting) {
	ret = make(map[string]map[string]*RequestSetting)
	flatten(ret, rc, "", nil)
	return
}

// flatten 递归展开gro
func flatten(ret map[string]map[string]*RequestSetting, group *_group, prefix string, filters []HandleFunc) {
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
