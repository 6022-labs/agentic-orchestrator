package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/common"
)

const (
	challengeCachePathTemplate = "auth/challenges/%s"
	nonceByteLength            = 32
)

// NonceStore keeps single-use login challenges in the shared cache storage.
type NonceStore struct {
	cacheStorage interfaces.CacheStorageInterface
	ttl          time.Duration
}

func NewNonceStore(cacheStorage interfaces.CacheStorageInterface, jwtSettings *settings.JwtSettings) *NonceStore {
	return &NonceStore{
		cacheStorage: cacheStorage,
		ttl:          time.Duration(jwtSettings.NonceTtl) * time.Second,
	}
}

func (s *NonceStore) Generate(address common.Address, chainId uint64) (*interfaces.ChallengeEntry, error) {
	buf := make([]byte, nonceByteLength)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	now := time.Now()
	entry := interfaces.ChallengeEntry{
		Nonce:     hex.EncodeToString(buf),
		ChainId:   chainId,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.ttl),
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to encode challenge: %w", err)
	}

	ttl := s.ttl
	if err := s.cacheStorage.Store(s.cachePath(address), string(payload), &ttl); err != nil {
		return nil, fmt.Errorf("failed to store challenge: %w", err)
	}

	return &entry, nil
}

func (s *NonceStore) Consume(address common.Address) *interfaces.ChallengeEntry {
	payload, err := s.cacheStorage.ReadAndDelete(s.cachePath(address))
	if err != nil {
		return nil
	}
	var entry interfaces.ChallengeEntry
	if err := json.Unmarshal([]byte(payload), &entry); err != nil {
		return nil
	}
	return &entry
}

func (s *NonceStore) cachePath(address common.Address) string {
	return fmt.Sprintf(challengeCachePathTemplate, strings.ToLower(address.Hex()))
}
