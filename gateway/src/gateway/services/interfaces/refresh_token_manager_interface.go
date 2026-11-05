package interfaces

import "time"

type RefreshClaims struct {
	Address          string
	SessionId        string
	TokenId          string
	ExpiresAt        time.Time
	SessionExpiresAt time.Time
}

// RefreshTokenManagerInterface signs refresh tokens with a key derived from the access key,
// so neither kind ever validates as the other.
type RefreshTokenManagerInterface interface {
	// Issue starts a session capped at now + jwt.refresh_ttl. Each token expires jwt.refresh_idle_ttl after
	// it is signed, never past that cap.
	Issue(address string) (string, *RefreshClaims, error)
	Reissue(claims RefreshClaims, tokenId string) (string, *RefreshClaims, error)
	NewTokenId() (string, error)
	Parse(token string) (*RefreshClaims, error)
}
