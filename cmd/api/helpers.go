package main

import (
	"net/http"
)

func (app *app) clientError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	app.logger.Warn(msg, "error", err, "method", r.Method, "path", r.URL.Path)
	http.Error(w, msg, http.StatusBadRequest)
}

func (app *app) serverError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	app.logger.Error(msg, "error", err, "method", r.Method, "path", r.URL.Path)
	http.Error(w, msg, http.StatusInternalServerError)
}
