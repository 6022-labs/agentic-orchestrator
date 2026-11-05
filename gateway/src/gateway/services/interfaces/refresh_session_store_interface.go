package interfaces

import (
	"context"
	"time"
)

type SpentRefreshToken struct {
	SuccessorId string    `json:"successorId"`
	SpentAt     time.Time `json:"spentAt"`
}

// RefreshSessionStoreInterface tracks spent refresh tokens and revoked sessions until the session expires.
type RefreshSessionStoreInterface interface {
	// Spend atomically marks tokenId spent; if it already was, it returns the first record (nil if gone) and false.
	Spend(
		ctx context.Context,
		sessionId string,
		tokenId string,
		spent SpentRefreshToken,
		expiresAt time.Time,
	) (*SpentRefreshToken, bool, error)
	Revoke(ctx context.Context, sessionId string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, sessionId string) (bool, error)
}
