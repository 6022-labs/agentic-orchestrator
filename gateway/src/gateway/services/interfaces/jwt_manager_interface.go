package interfaces

import "github.com/golang-jwt/jwt/v5"

type JwtManagerInterface interface {
	VerifyJwtToken(token string) (string, error)
	CreateJwtToken(address string, chainId uint64) (string, error)
	ExtractSubFromToken(token *jwt.Token, hash string) (*string, error)
}
