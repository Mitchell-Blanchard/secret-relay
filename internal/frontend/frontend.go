package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
)

//go:embed static
var staticFS embed.FS

type Frontend struct {
	logger    *slog.Logger
	templates map[string]*template.Template
}

func New(l *slog.Logger) (*Frontend, error) {
	ts := map[string]*template.Template{}

	pages, err := fs.Glob(staticFS, "static/html/pages/*.html.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		t, err := template.ParseFS(staticFS, "static/html/base.html.tmpl", page)
		if err != nil {
			return nil, err
		}
		ts[path.Base(page)] = t
	}

	return &Frontend{
		logger:    l.With("component", "frontend"),
		templates: ts,
	}, nil
}

func (f *Frontend) Routes(mux *http.ServeMux) {
	mux.Handle("GET /static/css/", http.FileServerFS(staticFS))
	mux.HandleFunc("GET /{$}", f.handleHome)
}
