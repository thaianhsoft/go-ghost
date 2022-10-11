package ghostengine

import "net/http"

type engine struct {
	*RouterGroup
	router *internalRouter
}

func Engine() *engine {
	e := &engine{
		RouterGroup: &RouterGroup{},
		router: newRouter(),
	}
	e.RouterGroup.engine = e
	return e
}


func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

func (e *engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}