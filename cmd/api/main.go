package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mitchell.blanchard/secret-relay/internal/database"
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
	dbHost string
	dbPort string
	dbPass string
}

func getEnv(env string, d string) string {
	if envVar := os.Getenv(env); envVar != "" {
		return envVar
	}
	return d
}

func newConfig() config {
	return config{
		dbHost: getEnv("DATABASE_HOST", "localhost"),
		dbPort: getEnv("DATABASE_PORT", "6379"),
		dbPass: getEnv("DATABASE_PASS", ""),
	}
}

func newStore(ctx context.Context) (store, error) {
	cfg := newConfig()
	dbAddr := cfg.dbHost + ":" + cfg.dbPort

	conn, err := database.NewRedis(ctx, dbAddr, cfg.dbPass)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	ctx := context.Background()
	slogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dataStore, err := newStore(ctx)
	if err != nil {
		slogger.Error("unable to create store: ", "error", err)
		os.Exit(1)
	}

	app := app{
		logger: slogger,
		store:  dataStore,
	}

	fmt.Println(app)
}
