package main

import (
	"fmt"
	"net/http"

	"github.com/mitchell.blanchard/secret-relay/internal/frontend"
)

func (app *application) routes() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /secrets", app.handleCreateSecret)
	mux.HandleFunc("GET /secrets/{key}", app.handleTakeSecret)

	if !app.cfg.enableFrontend {
		return mux, nil
	}

	fe, err := frontend.New(app.logger)
	if err != nil {
		return nil, fmt.Errorf("enabling frontend: %w", err)
	}

	fe.Routes(mux)
	return mux, nil
}
