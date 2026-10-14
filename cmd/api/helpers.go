package main

import (
	"net/http"
)

func (app *application) clientError(w http.ResponseWriter, r *http.Request, statusCode int, msg string, err error) {
	app.logger.Warn(msg, "error", err, "method", r.Method, "path", r.URL.Path, "status", statusCode)
	http.Error(w, msg, statusCode)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	app.logger.Error(msg, "error", err, "method", r.Method, "path", r.URL.Path, "status", http.StatusInternalServerError)
	http.Error(w, msg, http.StatusInternalServerError)
}
