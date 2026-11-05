package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
)

const (
	spentRefreshTokenCachePathTemplate = "auth/refresh/%s/spent/%s"
	revokedSessionCachePathTemplate    = "auth/refresh/%s/revoked"
	revokedSessionMarker               = "1"
	minRefreshEntryTtl                 = time.Second
)

type RefreshSessionStore struct {
	cacheStorage interfaces.CacheStorageInterface
	now          func() time.Time
}

func NewRefreshSessionStore(cacheStorage interfaces.CacheStorageInterface) *RefreshSessionStore {
	return &RefreshSessionStore{cacheStorage: cacheStorage, now: time.Now}
}

func (s *RefreshSessionStore) Spend(
	ctx context.Context,
	sessionId string,
	tokenId string,
	spent interfaces.SpentRefreshToken,
	expiresAt time.Time,
) (*interfaces.SpentRefreshToken, bool, error) {
	payload, err := json.Marshal(spent)
	if err != nil {
		return nil, false, fmt.Errorf("failed to encode spent refresh token: %w", err)
	}

	path := fmt.Sprintf(spentRefreshTokenCachePathTemplate, sessionId, tokenId)
	stored, err := s.cacheStorage.StoreIfAbsent(ctx, path, string(payload), s.ttlUntil(expiresAt))
	if err != nil {
		return nil, false, fmt.Errorf("failed to spend refresh token: %w", err)
	}
	if stored {
		return &spent, true, nil
	}

	prior, err := s.cacheStorage.Read(ctx, path)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read spent refresh token: %w", err)
	}
	if prior == nil {
		return nil, false, nil
	}

	var record interfaces.SpentRefreshToken
	if err := json.Unmarshal([]byte(*prior), &record); err != nil {
		return nil, false, fmt.Errorf("failed to decode spent refresh token: %w", err)
	}
	return &record, false, nil
}

func (s *RefreshSessionStore) Revoke(ctx context.Context, sessionId string, expiresAt time.Time) error {
	ttl := s.ttlUntil(expiresAt)
	path := fmt.Sprintf(revokedSessionCachePathTemplate, sessionId)
	if err := s.cacheStorage.Store(ctx, path, revokedSessionMarker, &ttl); err != nil {
		return fmt.Errorf("failed to revoke refresh session: %w", err)
	}
	return nil
}

func (s *RefreshSessionStore) IsRevoked(ctx context.Context, sessionId string) (bool, error) {
	marker, err := s.cacheStorage.Read(ctx, fmt.Sprintf(revokedSessionCachePathTemplate, sessionId))
	if err != nil {
		return false, fmt.Errorf("failed to read refresh session: %w", err)
	}
	return marker != nil, nil
}

// ttlUntil floors at a second: a zero TTL would store the key forever.
func (s *RefreshSessionStore) ttlUntil(expiresAt time.Time) time.Duration {
	return max(expiresAt.Sub(s.now()), minRefreshEntryTtl)
}
