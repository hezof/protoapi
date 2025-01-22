package protoapi

/****************************************
*	form
****************************************/

func (ctx *Context) FormValue(key string) string {
	return ""
}
func (ctx *Context) FormValueSlice(key string, explode bool) []string {
	return nil
}
func (ctx *Context) FormValueMap(key string, explode bool) map[string]string {
	return nil
}

/****************************************
*	path
****************************************/

func (ctx *Context) PathValue(key string) string {
	return ""
}
func (ctx *Context) PathValueSlice(key string, explode bool) []string {
	return nil
}
func (ctx *Context) PathValueMap(key string, explode bool) map[string]string {
	return nil
}

/****************************************
*	query
****************************************/

func (ctx *Context) QueryValue(key string) string {
	return ""
}
func (ctx *Context) QueryValueSlice(key string, explode bool) []string {
	return nil
}
func (ctx *Context) QueryValueMap(key string, explode bool) map[string]string {
	return nil
}

/****************************************
*	header
****************************************/

func (ctx *Context) HeaderValue(key string) string {
	return ""
}
func (ctx *Context) HeaderValueSlice(key string, explode bool) []string {
	return nil
}
func (ctx *Context) HeaderValueMap(key string, explode bool) map[string]string {
	return nil
}

/****************************************
*	cookie
****************************************/

func (ctx *Context) CookieValue(key string) string {
	return ""
}
func (ctx *Context) CookieValueSlice(key string, explode bool) []string {
	return nil
}
func (ctx *Context) CookieValueMap(key string, explode bool) map[string]string {
	return nil
}
