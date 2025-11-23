package interfaces

import (
	"time"
)

type CacheStorageInterface interface {
	Store(key string, value any, expiration *time.Duration) error
	ReadBool(key string) (*bool, error)
	ReadString(key string) (*string, error)
}
