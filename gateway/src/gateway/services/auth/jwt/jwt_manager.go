package jwt

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
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

func (j *JwtManager) CreateJwtToken(address string, chainId uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": address,
		"cId": chainId,
		"iat": time.Now().Unix(),                                   // Generation time
		"nbf": time.Now().Unix(),                                   // Cannot use before
		"exp": time.Now().Unix() + int64(j.jwtSettings.LivingTime), // Expires at
	})

	return token.SignedString([]byte(j.jwtSettings.Hash))
}

func (j *JwtManager) VerifyJwtToken(token string) (string, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.jwtSettings.Hash), nil
	})
	if err != nil {
		return "", err
	}

	return claims["sub"].(string), nil
}

func (j *JwtManager) ExtractSubFromToken(token *jwt.Token, hash string) (*string, error) {
	claims := token.Claims.(jwt.MapClaims)
	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("error while parsing sub from token")
	}

	return &sub, nil
}
