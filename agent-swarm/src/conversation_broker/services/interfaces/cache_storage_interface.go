package interfaces

import (
	"context"

	"time"
)

// CacheStorageInterface is the broker's cache port; a nil expiration stores
// the value without a TTL.
type CacheStorageInterface interface {
	Store(ctx context.Context, key string, value string, expiration *time.Duration) error
	Read(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
