package mockserver

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Method string
	Path   string
	Handle httprouter.Handle
}

type Router struct {
	routes []Route
}

func WithRoutes(rtes []Route) *httptest.Server {
	router := httprouter.New()
	for _, route := range rtes {
		router.Handle(strings.ToUpper(route.Method), route.Path, route.Handle)
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}))
}
