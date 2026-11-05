package jwt

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	jwtSettings *settings.JwtSettings
}

func NewJwtManager(jwtSettings *settings.JwtSettings) *JwtManager {
	return &JwtManager{
		jwtSettings: jwtSettings,
	}
}

func (j *JwtManager) CreateJwtToken(address string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(j.jwtSettings.LivingTime) * time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": address,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": expiresAt.Unix(),
	})

	signed, err := token.SignedString([]byte(j.jwtSettings.Hash))
	return signed, expiresAt, err
}

func (j *JwtManager) VerifyJwtToken(token string) (string, error) {
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

	return sub, nil
}
