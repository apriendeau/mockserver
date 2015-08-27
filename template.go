package mockserver

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"text/template"
)

type Template struct {
	object interface{}
	tmpl   *template.Template
}

func NewTemplate(name, base string, obj interface{}) (*Template, error) {
	tmpl, err := template.New(name).Parse(base)
	if err != nil {
		return nil, err
	}
	return &Template{
		object: obj,
		tmpl:   tmpl,
	}, nil
}

func (t *Template) Server(code int, contentType string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg []byte
		buf := bytes.NewBuffer(msg)
		if err := t.tmpl.Execute(buf, t.object); err != nil {
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
