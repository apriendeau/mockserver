package mockserver

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"text/template"
)

// Template is exported for the sake of name collision and name spacing.
type Template struct {
	tmpl *template.Template
}

// NewTemplate sets up test/template. It is if you want to dynamically
// populate your response. The base just is the text/template syntax
func NewTemplate(name, base string) (*Template, error) {
	tmpl, err := template.New(name).Parse(base)
	if err != nil {
		return nil, err
	}
	return &Template{
		tmpl: tmpl,
	}, nil
}

// Server creates an httptest.Server for you. If your template fails with the
// provided object, the http server will return a 500, set the content-type to
// plain/text and the error will be written since it happens at the time of
// execution
func (t *Template) Server(code int, contentType string, obj interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg []byte
		buf := bytes.NewBuffer(msg)
		if err := t.tmpl.Execute(buf, obj); err != nil {
			code = 500
			contentType = "plain/text"
			fmt.Fprint(w, "Template server error:", err.Error())
			return
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", contentType)
		fmt.Fprint(w, buf.String())
	}))
}
