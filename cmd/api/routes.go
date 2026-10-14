package main

import (
	"net/http"
)

func (app *app) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /secrets", app.handleCreateSecret)

	return mux
}
