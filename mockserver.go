package mockserver

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Simple does little to no setup for you just takes status code, content-type
// and the body
func Simple(code int, contentType, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", contentType)
		fmt.Fprint(w, body)
	}))
}

// JSON does a little more magic, where it takes a code and interface{}. It
// will set the content-type to application/json for you and marshal
// the object for the test server
func JSON(code int, body interface{}) (*httptest.Server, error) {
	b, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	return Simple(code, "application/json", string(b)), nil
}

// XML  does a little more magic, where it takes a code and interface{}. It
// will set the content-type to text/xml for you and marshal
// the object for the test server
func XML(code int, body interface{}) (*httptest.Server, error) {
	b, err := xml.Marshal(body)
	if err != nil {
		return nil, err
	}
	return Simple(code, "text/xml", string(b)), nil
}
