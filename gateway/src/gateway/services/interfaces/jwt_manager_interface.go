package interfaces

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManagerInterface interface {
	VerifyJwtToken(token string) (string, error)
	CreateJwtToken(address string) (string, time.Time, error)
	ExtractSubFromToken(token *jwt.Token, hash string) (*string, error)
}
