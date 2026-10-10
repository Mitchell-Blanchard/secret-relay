package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type config struct {
	dbHost   string
	dbPort   string
	dbPass   string
	database int
}

func getEnv(e string, d string) string {
	if env := os.Getenv(e); env != "" {
		return env
	}
	return d
}

func initConfig() (config, error) {
	var cfg config

	var err error
	cfg.dbHost = getEnv("DATABASE_HOST", "localhost")
	cfg.dbPort = getEnv("DATABASE_PORT", "6379")
	cfg.dbPass = getEnv("DATABASE_PASS", "")
	cfg.database, err = strconv.Atoi(getEnv("DATABASE_DB", "0"))
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}

func main() {
	ctx := context.Background()
	slogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := initConfig()
	if err != nil {
		slogger.Error("failed to initialise config", "error", err)
		os.Exit(1)
	}

	databaseEndpoint := cfg.dbHost + ":" + cfg.dbPort
	slogger.Info("connecting to database", "endpoint", databaseEndpoint)

	rdb := redis.NewClient(&redis.Options{
		Addr:     databaseEndpoint,
		Password: cfg.dbPass,
		DB:       cfg.database,
	})
	defer rdb.Close()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		slogger.Error("error connecting to database", "error", err)
		os.Exit(1)
	}

	slogger.Info("successfully connected to database", "endpoint", databaseEndpoint)
}
