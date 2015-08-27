package mockserver

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Route is common structure for setting up routes with the mockserver
type Route struct {
	Method string
	Path   string
	Handle httprouter.Handle
}

// WithRoutes creates a new mockserver that takes an array of routes and puts
// them into a httptest.Server
func WithRoutes(rtes []Route) *httptest.Server {
	router := httprouter.New()
	for _, route := range rtes {
		router.Handle(strings.ToUpper(route.Method), route.Path, route.Handle)
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}))
}
