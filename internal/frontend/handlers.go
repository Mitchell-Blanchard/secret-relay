package frontend

import (
	"bytes"
	"net/http"
)

func (f *Frontend) handleHome(wr http.ResponseWriter, r *http.Request) {
	name := "home.html.tmpl"

	tmpl, ok := f.templates[name]
	if !ok {
		f.logger.Error("template not found", "path", r.URL.Path)
		http.Error(wr, "internal server error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, nil); err != nil {
		f.logger.Error("executing template", "error", err, "path", r.URL.Path)
		http.Error(wr, "internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := buf.WriteTo(wr); err != nil {
		f.logger.Warn("writing response", "error", err, "path", r.URL.Path)
	}
}
