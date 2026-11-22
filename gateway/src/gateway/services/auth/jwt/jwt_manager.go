package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrAccessTokenSessionMissing = errors.New("access token has no session")
	ErrAccessTokenSessionRevoked = errors.New("access token session revoked")
)

type JwtManager struct {
	jwtSettings         *settings.JwtSettings
	refreshSessionStore interfaces.RefreshSessionStoreInterface
}

func NewJwtManager(
	jwtSettings *settings.JwtSettings,
	refreshSessionStore interfaces.RefreshSessionStoreInterface,
) *JwtManager {
	return &JwtManager{
		jwtSettings:         jwtSettings,
		refreshSessionStore: refreshSessionStore,
	}
}

func (j *JwtManager) CreateJwtToken(
	address string,
	sessionId string,
	sessionExpiresAt time.Time,
) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(j.jwtSettings.LivingTime) * time.Second)
	if sessionExpiresAt.Before(expiresAt) {
		expiresAt = sessionExpiresAt
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		subjectClaim:   address,
		issuedAtClaim:  time.Now().Unix(),
		notBeforeClaim: time.Now().Unix(),
		expiresAtClaim: expiresAt.Unix(),
		sessionIdClaim: sessionId,
	})

	signed, err := token.SignedString([]byte(j.jwtSettings.Hash))
	return signed, expiresAt, err
}

func (j *JwtManager) VerifyJwtToken(ctx context.Context, token string) (string, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.jwtSettings.Hash), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}

	sub, err := claims.GetSubject()
	if err != nil || !common.IsHexAddress(sub) {
		return "", jwt.ErrTokenInvalidSubject
	}

	sessionId, _ := claims[sessionIdClaim].(string)
	if sessionId == "" {
		return "", ErrAccessTokenSessionMissing
	}
	revoked, err := j.refreshSessionStore.IsRevoked(ctx, sessionId)
	if err != nil {
		return "", fmt.Errorf("failed to check access token session: %w", err)
	}
	if revoked {
		return "", ErrAccessTokenSessionRevoked
	}

	return sub, nil
}
