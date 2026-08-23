package interfaces

import (
	"context"

	"time"
)

type CacheStorageInterface interface {
	Store(ctx context.Context, key string, value string, expiration *time.Duration) error
	// ReadAndDelete atomically reads and drops a key; a nil value with a nil error is a miss.
	ReadAndDelete(ctx context.Context, key string) (*string, error)
}
