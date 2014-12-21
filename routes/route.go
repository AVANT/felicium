package routes

import (
	"fmt"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		HomeHandler,
	},
}
