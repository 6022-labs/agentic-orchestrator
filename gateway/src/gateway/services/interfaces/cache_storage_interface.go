package interfaces

import "time"

type CacheStorageInterface interface {
	Store(key string, value string, expiration *time.Duration) error
	ReadAndDelete(key string) (string, error)
}
