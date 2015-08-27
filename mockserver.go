package mockserver

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func Clean(code int, contentType, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", contentType)
		fmt.Fprint(w, body)
	}))
}

func JSON(code int, body interface{}) (*httptest.Server, error) {
	b, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	return Clean(code, "application/json", string(b)), nil
}

func XML(code int, body interface{}) (*httptest.Server, error) {
	b, err := xml.Marshal(body)
	if err != nil {
		return nil, err
	}
	return Clean(code, "text/xml", string(b)), nil
}
