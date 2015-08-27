package mockserver_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/apriendeau/mockserver"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func sampleRoutes() []mockserver.Route {
	first := mockserver.Route{
		Method: "GET",
		Path:   "/hello/world",
		Handle: httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "plain/text")
			fmt.Fprint(w, "Hello world")
		}),
	}
	second := mockserver.Route{
		Method: "POST",
		Path:   "/welcome/:name",
		Handle: httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "plain/text")
			fmt.Fprint(w, "Welcome ", p.ByName("name"))
		}),
	}
	return []mockserver.Route{first, second}

}

func TestWithRoutes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	rtes := sampleRoutes()
	server := mockserver.WithRoutes(rtes)
	defer server.Close()
	// Check 404 route
	resp, err := http.Get(server.URL + "/random")
	assert.NoError(err)
	assert.Equal(404, resp.StatusCode)
	// check POST
	resp, err = http.Post(server.URL+"/welcome/tony", "plain/text", bytes.NewBuffer([]byte{}))
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)
	// check body response
	body, err := parseResp(resp.Body)
	assert.NoError(err)
	assert.Equal(body, "Welcome tony")
	// check 405
	resp, err = http.Get(server.URL + "/welcome/tony")
	assert.NoError(err)
	assert.Equal(405, resp.StatusCode)

	// check second route
	resp, err = http.Get(server.URL + "/hello/world")
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)
	// check hello world
	body, err = parseResp(resp.Body)
	assert.NoError(err)
	assert.Equal("Hello world", body)
}
