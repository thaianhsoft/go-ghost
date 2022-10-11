package ghostengine

type RouterGroup struct {
	*engine
	prefix string
	middlewares []HandlerFunc
	childGroups []*RouterGroup
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: r.prefix + prefix,
		engine: r.engine,
	}
	return newGroup
}

func (r *RouterGroup) GET(url string, handler HandlerFunc) {
	r.addRoute("GET", url, handler)
}

func (r *RouterGroup) POST(url string, handler HandlerFunc) {
	r.addRoute("POST", url, handler)
}

func (r *RouterGroup) addRoute(method string, url string, handler HandlerFunc) {
	pattern := r.prefix + url
	r.engine.router.insertRoute(method, pattern, handler)
}






