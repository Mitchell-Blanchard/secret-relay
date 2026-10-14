package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mitchell.blanchard/secret-relay/internal/datastore"
	"github.com/mitchell.blanchard/secret-relay/internal/secret"
)

type store interface {
	Save(ctx context.Context, s secret.Secret) error
	Take(ctx context.Context, key string) (secret.Secret, error)
}

type app struct {
	logger *slog.Logger
	store  store
}

type config struct {
	dbHost     string
	dbPort     string
	dbPass     string
	serverHost string
	serverPort string
}

func getEnv(envVarName string, fallback string) string {
	if setVar := os.Getenv(envVarName); setVar != "" {
		return setVar
	}
	return fallback
}

func newConfig() config {
	return config{
		dbHost:     getEnv("DATABASE_HOST", "localhost"),
		dbPort:     getEnv("DATABASE_PORT", "6379"),
		dbPass:     getEnv("DATABASE_PASS", ""),
		serverHost: getEnv("SERVER_HOST", ""),
		serverPort: getEnv("SERVER_PORT", "8080"),
	}
}

func newStore(ctx context.Context, cfg config) (store, error) {
	dbAddr := cfg.dbHost + ":" + cfg.dbPort

	ds, err := datastore.NewRedis(ctx, dbAddr, cfg.dbPass)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func main() {
	ctx := context.Background()
	slogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := newConfig()

	st, err := newStore(ctx, cfg)
	if err != nil {
		slogger.Error("unable to create store", "error", err)
		os.Exit(1)
	}

	app := app{
		logger: slogger,
		store:  st,
	}

	err = app.start(cfg)
	slogger.Error("server stopped", "error", err)
	os.Exit(1)
}

func (app *app) start(cfg config) error {
	address := cfg.serverHost + ":" + cfg.serverPort

	srv := &http.Server{
		Addr:              address,
		Handler:           app.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return srv.ListenAndServe()
}
