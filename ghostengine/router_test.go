package ghostengine

import (
	"fmt"
	"testing"
)

func TestRouter(t *testing.T) {
	router := newRouter()
	router.insertRoute("GET", "/api/:role/:id/notes", func(c Context) {

	})
	router.insertRoute("GET", "/api/user/:id/books", func(c Context) {
		fmt.Println("/api/user/:id/notes is pattern one")
	})

	fn, params := router.getRoute("GET", "/api/user/3/books")
	fmt.Println(fn, params)

}