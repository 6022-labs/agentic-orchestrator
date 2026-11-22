package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/golang-jwt/jwt/v5"
)

const (
	refreshTokenType    = "refresh"
	refreshKeyDomain    = ":refresh"
	refreshIdByteLength = 32
	refreshTypeClaim    = "typ"
	refreshTokenIdClaim = "jti"
	sessionExpiryClaim  = "sxp"
)

var ErrRefreshTokenMalformed = errors.New("refresh token is malformed")

type RefreshTokenManager struct {
	jwtSettings *settings.JwtSettings
	signingKey  []byte
	now         func() time.Time
}

func NewRefreshTokenManager(jwtSettings *settings.JwtSettings) *RefreshTokenManager {
	key := sha256.Sum256([]byte(jwtSettings.Hash + refreshKeyDomain))
	return &RefreshTokenManager{
		jwtSettings: jwtSettings,
		signingKey:  key[:],
		now:         time.Now,
	}
}

func (m *RefreshTokenManager) Issue(address string) (string, *interfaces.RefreshClaims, error) {
	sessionId, err := m.NewTokenId()
	if err != nil {
		return "", nil, err
	}
	tokenId, err := m.NewTokenId()
	if err != nil {
		return "", nil, err
	}

	return m.Reissue(interfaces.RefreshClaims{
		Address:          address,
		SessionId:        sessionId,
		SessionExpiresAt: m.now().Add(time.Duration(m.jwtSettings.RefreshTtl) * time.Second),
	}, tokenId)
}

func (m *RefreshTokenManager) Reissue(
	claims interfaces.RefreshClaims,
	tokenId string,
) (string, *interfaces.RefreshClaims, error) {
	now := m.now()
	claims.TokenId = tokenId
	claims.ExpiresAt = now.Add(time.Duration(m.jwtSettings.RefreshIdleTtl) * time.Second)
	if claims.SessionExpiresAt.Before(claims.ExpiresAt) {
		claims.ExpiresAt = claims.SessionExpiresAt
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		refreshTypeClaim:    refreshTokenType,
		subjectClaim:        claims.Address,
		sessionIdClaim:      claims.SessionId,
		refreshTokenIdClaim: claims.TokenId,
		issuedAtClaim:       now.Unix(),
		expiresAtClaim:      claims.ExpiresAt.Unix(),
		sessionExpiryClaim:  claims.SessionExpiresAt.Unix(),
	})

	signed, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", nil, err
	}
	return signed, &claims, nil
}

func (m *RefreshTokenManager) NewTokenId() (string, error) {
	buf := make([]byte, refreshIdByteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func (m *RefreshTokenManager) Parse(token string) (*interfaces.RefreshClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		claims,
		func(*jwt.Token) (any, error) { return m.signingKey, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(m.now),
	)
	if err != nil {
		return nil, err
	}

	typ, _ := claims[refreshTypeClaim].(string)
	sub, _ := claims[subjectClaim].(string)
	sid, _ := claims[sessionIdClaim].(string)
	jti, _ := claims[refreshTokenIdClaim].(string)
	sxp, _ := claims[sessionExpiryClaim].(float64)
	if typ != refreshTokenType || sub == "" || sid == "" || jti == "" || sxp <= 0 {
		return nil, ErrRefreshTokenMalformed
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	return &interfaces.RefreshClaims{
		Address:          sub,
		SessionId:        sid,
		TokenId:          jti,
		ExpiresAt:        expiresAt.Time,
		SessionExpiresAt: time.Unix(int64(sxp), 0),
	}, nil
}
