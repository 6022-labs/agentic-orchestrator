package interfaces

import (
	"context"

	"time"
)

type CacheStorageInterface interface {
	Store(ctx context.Context, key string, value string, expiration *time.Duration) error
	ReadAndDelete(ctx context.Context, key string) (string, error)
}
