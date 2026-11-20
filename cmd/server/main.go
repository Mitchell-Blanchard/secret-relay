package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mitchell.blanchard/secret-relay/internal/datastore"
	"github.com/mitchell.blanchard/secret-relay/internal/secret"
)

type store interface {
	Save(ctx context.Context, s secret.Secret) error
	Take(ctx context.Context, key string) (secret.Secret, error)
	Close()
}

type application struct {
	logger *slog.Logger
	store  store
	cfg    config
}

type config struct {
	dbHost         string
	dbPort         string
	dbPass         string
	address        string
	enableFrontend bool
}

func getEnvStr(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) (bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback, fmt.Errorf("parsing %s: %w", key, err)
	}
	return b, nil
}

func newConfig() (config, error) {
	enableFrontend, err := getEnvBool("ENABLE_FRONTEND", true)
	if err != nil {
		return config{}, err
	}

	address := getEnvStr("SERVER_HOST", "") + ":" + getEnvStr("SERVER_PORT", "8080")

	return config{
		dbHost:         getEnvStr("DATABASE_HOST", "localhost"),
		dbPort:         getEnvStr("DATABASE_PORT", "6379"),
		dbPass:         getEnvStr("DATABASE_PASS", ""),
		address:        address,
		enableFrontend: enableFrontend,
	}, nil
}

func (app *application) newStore(ctx context.Context) (store, error) {
	dbAddr := app.cfg.dbHost + ":" + app.cfg.dbPort

	ds, err := datastore.NewRedis(ctx, dbAddr, app.cfg.dbPass)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func run(logger *slog.Logger) error {
	var (
		ctx = context.Background()
		err error
		app = application{
			logger: logger,
		}
	)

	app.cfg, err = newConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	app.store, err = app.newStore(ctx)
	if err != nil {
		return fmt.Errorf("creating store: %w", err)
	}
	defer app.store.Close()

	routes, err := app.routes()
	if err != nil {
		return fmt.Errorf("setting routes: %w", err)
	}

	app.logger.Info("starting server", "address", app.cfg.address)
	srv := &http.Server{
		Addr:              app.cfg.address,
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return srv.ListenAndServe()
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("server failure", "error", err)
		os.Exit(1)
	}
}
