package secret

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	Key   string
	Value string
	TTL   time.Duration
}

var (
	ErrNotFound   = errors.New("secret not found")
	ErrEmptyValue = errors.New("secret value is empty")
	ErrInvalidTTL = errors.New("ttl must be between 1 minute and 7 days")
)

const (
	minTTL = time.Minute
	maxTTL = 7 * 24 * time.Hour
)

func New(value string, ttl time.Duration) (Secret, error) {
	if value == "" {
		return Secret{}, ErrEmptyValue
	}
	if ttl < minTTL || ttl > maxTTL {
		return Secret{}, ErrInvalidTTL
	}

	return Secret{
		Key:   uuid.NewString(),
		Value: value,
		TTL:   ttl,
	}, nil
}
