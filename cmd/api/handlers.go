package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/mitchell.blanchard/secret-relay/internal/secret"
)

type createSecretRequest struct {
	Secret     string `json:"secret"`
	TTLSeconds int    `json:"ttl_seconds"`
}

func (app *application) handleCreateSecret(w http.ResponseWriter, r *http.Request) {
	var req createSecretRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		msg := "invalid request body"
		if errors.Is(err, io.EOF) {
			msg = "request body must not be empty"
		}
		app.clientError(w, r, msg, err)
		return
	}

	sec, err := secret.New(req.Secret, time.Duration(req.TTLSeconds)*time.Second)
	switch {
	case err == nil: // ok
	case errors.Is(err, secret.ErrEmptyValue):
		app.clientError(w, r, "invalid secret", err)
		return
	case errors.Is(err, secret.ErrInvalidTTL):
		app.clientError(w, r, "invalid expiry", err)
		return
	default:
		app.serverError(w, r, "unable to create secret", err)
		return
	}

	if err := app.store.Save(r.Context(), sec); err != nil {
		app.serverError(w, r, "unable to store input", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	key := map[string]string{"key": sec.Key}
	if err := json.NewEncoder(w).Encode(key); err != nil {
		app.logger.Warn("writing response", "error", err)
	}
}
