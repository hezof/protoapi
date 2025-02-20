package protoapi

import (
	"net/http"
	"net/url"
	"strings"
)

// 对于in为body/form, path, query, header, cookie的参数, 按照OAS 3.0规范解析:
// body/form: style=form, explode=?{false|true}
// path: style=simple, explode=?{false|true}
// query: style=form, explode=?{false|true}
// header: style=simple, explode=?{false|true}
// cookie: style=form, explode=?{false|true}
// 详细内容参考: https://swagger.io/docs/specification/v3_0/serialization

/****************************************
*	form: style=form, explode=?
****************************************/

// parsePostForm 同时适配multipart/form-data 与 application/x-www-form-urlencoded
func parsePostForm(request *http.Request, maxMemory int64) error {
	if strings.HasPrefix(request.Header.Get("Content-Type"), "multipart/form-data") {
		return request.ParseMultipartForm(maxMemory)
	}
	return request.ParseForm()
}

// FormValue 获取首个form参数
func (ctx *Context) FormValue(name string) (string, error) {
	if ctx.Request.PostForm == nil {
		err := parsePostForm(ctx.Request, ctx.Handler.FormMaxMemory)
		if err != nil {
			return "", err
		}
	}
	if vs := ctx.Request.PostForm[name]; len(vs) > 0 {
		return vs[0], nil
	}
	return "", nil
}

// FormValueRepeated 根据explode解析form array参数
func (ctx *Context) FormValueRepeated(name string, explode bool) ([]string, error) {
	if ctx.Request.PostForm == nil {
		err := parsePostForm(ctx.Request, ctx.Handler.FormMaxMemory)
		if err != nil {
			return nil, err
		}
	}
	if explode {
		// collectFormat=multi(style=form, explode=true)
		return ctx.Request.PostForm[name], nil
	} else {
		// collectFormat=csv(style=form, explode=false)
		if vs := ctx.Request.PostForm[name]; len(vs) > 0 {
			return strings.Split(vs[0], ","), nil
		}
		return nil, nil
	}
}

// FormValueMap 根据explode解析form object参数
func (ctx *Context) FormValueMap(name string, explode bool) (map[string]string, error) {
	if ctx.Request.PostForm == nil {
		err := parsePostForm(ctx.Request, ctx.Handler.FormMaxMemory)
		if err != nil {
			return nil, err
		}
	}
	if explode {
		// collectFormat=multi(style=form, explode=true)
		rt := make(map[string]string)
		for _, kv := range ctx.Request.PostForm[name] {
			if ps := strings.IndexByte(kv, '='); ps != -1 {
				rt[kv[:ps]] = kv[ps+1:]
			}
		}
		return rt, nil
	} else {
		// collectFormat=csv(style=form, explode=false)
		if vs := ctx.Request.PostForm[name]; len(vs) > 0 {
			rt := make(map[string]string)
			ss := strings.Split(vs[0], ",")
			for i, n := 1, len(ss); i < n; i += 2 {
				rt[ss[i-1]] = ss[i]
			}
			return rt, nil
		}
		return nil, nil
	}
}

/****************************************
*	path: style=simple, explode=?
****************************************/

// PathValue 获取首个path参数
func (ctx *Context) PathValue(name string) (string, error) {
	for _, pm := range ctx.params {
		if pm.Key == name {
			return pm.Value, nil
		}
	}
	return "", nil
}

// PathValueRepeated 根据explode解析path array参数
func (ctx *Context) PathValueRepeated(name string, explode bool) ([]string, error) {
	for _, pm := range ctx.params {
		if pm.Key == name {
			// path的explode情形都相同, 所以可以忽略
			return strings.Split(pm.Value, ","), nil
		}
	}
	return nil, nil
}

// PathValueMap 根据explode解析path object参数
func (ctx *Context) PathValueMap(name string, explode bool) (map[string]string, error) {
	for _, pm := range ctx.params {
		if pm.Key == name {
			if explode {
				rt := make(map[string]string)
				for _, kv := range strings.Split(pm.Value, ",") {
					if ps := strings.IndexByte(kv, '='); ps != -1 {
						rt[kv[:ps]] = kv[ps+1:]
					}
				}
				return rt, nil
			} else {
				rt := make(map[string]string)
				ss := strings.Split(pm.Value, ",")
				for i, n := 1, len(ss); i < n; i += 2 {
					rt[ss[i-1]] = ss[i]
				}
				return rt, nil
			}
		}
	}
	return nil, nil
}

/****************************************
*	query: style=form, explode=?
****************************************/

func (ctx *Context) QueryValue(name string) (string, error) {
	if ctx.query == nil {
		var err error
		ctx.query, err = url.ParseQuery(ctx.Request.URL.RawQuery)
		if err != nil {
			return "", err
		}
	}
	if vs, _ := ctx.query[name]; len(vs) > 0 {
		return vs[0], nil
	}
	return "", nil
}
func (ctx *Context) QueryValueRepeated(name string, explode bool) ([]string, error) {
	if ctx.query == nil {
		var err error
		ctx.query, err = url.ParseQuery(ctx.Request.URL.RawQuery)
		if err != nil {
			return nil, err
		}
	}
	if explode {
		return ctx.query[name], nil
	} else {
		if vs := ctx.query[name]; len(vs) > 0 {
			return strings.Split(vs[0], ","), nil
		}
		return nil, nil
	}
}
func (ctx *Context) QueryValueMap(name string, explode bool) (map[string]string, error) {
	if ctx.query == nil {
		var err error
		ctx.query, err = url.ParseQuery(ctx.Request.URL.RawQuery)
		if err != nil {
			return nil, err
		}
	}
	if explode {
		// collectFormat=multi(style=form, explode=true)
		rt := make(map[string]string)
		for _, kv := range ctx.query[name] {
			if ps := strings.IndexByte(kv, '='); ps != -1 {
				rt[kv[:ps]] = kv[ps+1:]
			}
		}
		return rt, nil
	} else {
		// collectFormat=csv(style=form, explode=false)
		if vs := ctx.query[name]; len(vs) > 0 {
			rt := make(map[string]string)
			ss := strings.Split(vs[0], ",")
			for i, n := 1, len(ss); i < n; i += 2 {
				rt[ss[i-1]] = ss[i]
			}
			return rt, nil
		}
		return nil, nil
	}
}

/****************************************
*	header: style=simple, explode=?
****************************************/

func (ctx *Context) HeaderValue(name string) (string, error) {
	return ctx.Request.Header.Get(name), nil
}
func (ctx *Context) HeaderValueRepeated(name string, explode bool) ([]string, error) {
	if explode {
		return ctx.Request.Header.Values(name), nil
	} else {
		return strings.Split(ctx.Request.Header.Get(name), ","), nil
	}
}
func (ctx *Context) HeaderValueMap(name string, explode bool) (map[string]string, error) {
	if explode {
		rt := make(map[string]string)
		for _, kv := range ctx.Request.Header.Values(name) {
			if ps := strings.IndexByte(kv, '='); ps != -1 {
				rt[kv[:ps]] = kv[ps+1:]
			}
		}
		return rt, nil
	} else {
		rt := make(map[string]string)
		ss := strings.Split(ctx.Request.Header.Get(name), ",")
		for i, n := 1, len(ss); i < n; i += 2 {
			rt[ss[i-1]] = ss[i]
		}
		return rt, nil
	}
}

/****************************************
*	cookie
****************************************/

func (ctx *Context) CookieValue(name string) (string, error) {
	for _, ck := range ctx.Request.Cookies() {
		if ck.Name == name {
			return ck.Value, nil
		}
	}
	return "", nil
}
func (ctx *Context) CookieValueRepeated(name string, explode bool) ([]string, error) {
	if explode {
		rt := make([]string, 0, 4)
		for _, ck := range ctx.Request.Cookies() {
			if ck.Name == name {
				rt = append(rt, ck.Value)
			}
		}
		return rt, nil
	} else {
		for _, ck := range ctx.Request.Cookies() {
			if ck.Name == name {
				return strings.Split(ck.Value, ","), nil
			}
		}
		return nil, nil
	}
}
func (ctx *Context) CookieValueMap(name string, explode bool) (map[string]string, error) {
	if explode {
		rt := make(map[string]string)
		for _, ck := range ctx.Request.Cookies() {
			if ck.Name == name {
				if ps := strings.IndexByte(ck.Value, '='); ps != -1 {
					rt[ck.Value[:ps]] = ck.Value[ps+1:]
				}
			}
		}
		return rt, nil
	} else {
		for _, ck := range ctx.Request.Cookies() {
			if ck.Name == name {
				rt := make(map[string]string)
				vs := strings.Split(ck.Value, ",")
				for i, n := 1, len(vs); i < n; i += 2 {
					rt[vs[i-1]] = vs[i]
				}
				return rt, nil
			}
		}
		return nil, nil
	}
}
