package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
)

const (
	nonceCachePathTemplate = "auth/nonces/%s"
	nonceTtl               = 5 * time.Minute
)

// NonceStore keeps single-use login nonces in the shared cache storage.
type NonceStore struct {
	cacheStorage interfaces.CacheStorageInterface
}

func NewNonceStore(cacheStorage interfaces.CacheStorageInterface) *NonceStore {
	return &NonceStore{cacheStorage: cacheStorage}
}

func (s *NonceStore) Generate() (string, time.Time, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate nonce: %w", err)
	}
	nonce := hex.EncodeToString(buf)

	ttl := nonceTtl
	if err := s.cacheStorage.Store(fmt.Sprintf(nonceCachePathTemplate, nonce), "1", &ttl); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to store nonce: %w", err)
	}

	return nonce, time.Now().Add(nonceTtl), nil
}

func (s *NonceStore) Consume(nonce string) bool {
	_, err := s.cacheStorage.ReadAndDelete(fmt.Sprintf(nonceCachePathTemplate, nonce))
	return err == nil
}
