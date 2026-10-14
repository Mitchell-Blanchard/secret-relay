package datastore

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
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	return &Redis{client: client}, nil
}

func (r *Redis) Save(ctx context.Context, s secret.Secret) error {
	if err := r.client.Set(ctx, s.Key, s.Value, s.TTL).Err(); err != nil {
		return fmt.Errorf("saving secret: %w", err)
	}
	return nil
}

func (r *Redis) Take(ctx context.Context, key string) (secret.Secret, error) {
	return secret.Secret{}, nil
}
