# mockserver

[![Go Doc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/apriendeau/mockserver)

When do you write a library? Its either:

A. You need something that noone has done
B. Something is just urking you because you need it
C. All of the Above

It this case, its C.

I have seen no good example or explanation on mocking when it comes to doing
integration testing. You don't want your tests to always all a third party API
or you need to create a mockserver of your API. Either way, the key feature
about this mockserver is *ROUTING* which I went with [httprouter](github.com/julienschmidt/httprouter) because I like
speed. Yes, there is helper functions to generate one off `httptest.Server` for
JSON and XML as well a Simple one.

This merely is exposing the unmodified power of routing from httprouter to
testing.


## Let's get this party started

1. Routes (borrowed from the routing tests)

```golang
package your_test

import (
	"testing"

	"github.com/apriendeau/mockserver"
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

func TestAllTheThings(t *testing) {
	assert := assert.New(t)
	routes := sampleRoutes()
	server := mockserver.WithRoutes(rtes)
	defer server.Close()

	resp, err := http.Get(server.URL + "/random")
	assert.NoError(err)
	// This will pass
	assert.Equal(404, resp.StatusCode)

	// How about some parameters
	resp, err = http.Post(server.URL+"/welcome/tony", "plain/text", bytes.NewBuffer([]byte{}))
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)
	body, err := parseResp(resp.Body)
	assert.NoError(err)
	assert.Equal(body, "Welcome tony") // Booyah biz-nitches

	// How about the same route but an umplemented method?
	resp, err = http.Get(server.URL + "/welcome/tony")
	assert.NoError(err)
	assert.Equal(405, resp.StatusCode) // again, booyah.
}
```
