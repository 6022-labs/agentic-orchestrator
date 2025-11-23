package interfaces

import "time"

type CacheStorageInterface interface {
	Store(key string, value string, expiration *time.Duration) error
	Read(key string) (string, error)
	Delete(key string) error
}
