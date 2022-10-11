package ghostengine

import "net/http"

type HandlerFunc func (c Context)
type Context interface{
	Next()
}

type internalContext struct {
	middlewares []HandlerFunc
	index int
}

func newContext(w http.ResponseWriter, r *http.Request) *internalContext {
	return &internalContext{
		index: -1,
	}
}

func (i *internalContext) Next() {

}