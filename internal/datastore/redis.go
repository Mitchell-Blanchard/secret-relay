package datastore

import (
	"context"
	"errors"
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
	sec, err := r.client.GetDel(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return secret.Secret{}, secret.ErrNotFound
	}
	if err != nil {
		return secret.Secret{}, fmt.Errorf("taking secret: %w", err)
	}

	return secret.Secret{Key: key, Value: sec}, nil
}
