package interfaces

import "time"

// NonceStoreInterface issues single-use login nonces with a TTL.
// Consume returns false when the nonce is unknown or already used.
type NonceStoreInterface interface {
	Generate() (string, time.Time, error)
	Consume(nonce string) bool
}
