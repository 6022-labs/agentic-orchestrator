package interfaces

import (
	"context"

	"time"
)

type CacheStorageInterface interface {
	Store(ctx context.Context, key string, value string, expiration *time.Duration) error
	// ReadAndDelete atomically reads and drops a key; a nil value with a nil error is a miss.
	ReadAndDelete(ctx context.Context, key string) (*string, error)
	// StoreIfAbsent stores a key only if it doesn't exist yet, and reports whether it stored it.
	StoreIfAbsent(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
}
