package interfaces

import (
	"context"
	"time"
)

// JwtManagerInterface mints and verifies access JWTs tied to a refresh session through their sid claim;
// an access JWT never outlives its session.
type JwtManagerInterface interface {
	// VerifyJwtToken returns the subject address; it rejects a token whose session is missing or revoked,
	// and fails closed when the session store errors.
	VerifyJwtToken(ctx context.Context, token string) (string, error)
	CreateJwtToken(address string, sessionId string, sessionExpiresAt time.Time) (string, time.Time, error)
}
