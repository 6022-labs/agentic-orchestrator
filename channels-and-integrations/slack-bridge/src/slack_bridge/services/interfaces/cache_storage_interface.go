package interfaces

import (
	"context"

	"time"
)

type CacheStorageInterface interface {
	Store(ctx context.Context, key string, value any, expiration *time.Duration) error
	ReadBool(ctx context.Context, key string) (*bool, error)
	ReadString(ctx context.Context, key string) (*string, error)
}
