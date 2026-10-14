package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /secrets", app.handleCreateSecret)
	mux.HandleFunc("GET /secrets/{key}", app.handleTakeSecret)

	return mux
}
