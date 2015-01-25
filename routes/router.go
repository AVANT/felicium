package routes

import (
	"github.com/avant/felicium/Godeps/_workspace/src/github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
