package interfaces

import "time"

// CacheStorageInterface is the broker's cache port; a nil expiration stores
// the value without a TTL.
type CacheStorageInterface interface {
	Store(key string, value string, expiration *time.Duration) error
	Read(key string) (string, error)
	Delete(key string) error
}
