package database

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchell.blanchard/secret-relay/internal/secret"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, addr string, pass string) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("connecting to redis at %s: %w", addr, err)
	}

	return &Redis{client: rdb}, nil
}

func (r *Redis) Save(ctx context.Context, s secret.Secret) error {
	if err := r.client.Set(ctx, s.Key, s.Value, s.TTL).Err(); err != nil {
		return fmt.Errorf("saving secret %s: %w", s.Key, err)
	}
	return nil
}

func (r *Redis) Take(ctx context.Context, key string) (secret.Secret, error) {
	return secret.Secret{}, nil
}
