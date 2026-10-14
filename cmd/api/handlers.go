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

	err := dec.Decode(&req)
	if errors.Is(err, io.EOF) {
		app.clientError(w, r, http.StatusBadRequest, "request body must not be empty", err)
		return
	}
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest, "invalid request body", err)
		return
	}

	sec, err := secret.New(req.Secret, time.Duration(req.TTLSeconds)*time.Second)
	if err != nil {
		switch {
		case errors.Is(err, secret.ErrEmptyValue):
			app.clientError(w, r, http.StatusBadRequest, "secret must not be empty", err)
		case errors.Is(err, secret.ErrInvalidTTL):
			app.clientError(w, r, http.StatusBadRequest, "invalid expiry", err)
		default:
			app.serverError(w, r, "unable to create secret", err)
		}
		return
	}

	if err := app.store.Save(r.Context(), sec); err != nil {
		app.serverError(w, r, "unable to store secret", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]string{"key": sec.Key}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		app.logger.Warn("writing response", "error", err)
	}
}

func (app *application) handleTakeSecret(w http.ResponseWriter, r *http.Request) {
	sec, err := app.store.Take(r.Context(), r.PathValue("key"))
	if errors.Is(err, secret.ErrNotFound) {
		app.clientError(w, r, http.StatusNotFound, "secret not found", err)
		return
	}
	if err != nil {
		app.serverError(w, r, "unable to retrieve secret", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{"secret": sec.Value}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		app.logger.Warn("writing response", "error", err)
	}
}
