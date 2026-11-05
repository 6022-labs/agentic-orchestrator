package interfaces

import (
	"time"
)

type JwtManagerInterface interface {
	VerifyJwtToken(token string) (string, error)
	CreateJwtToken(address string) (string, time.Time, error)
}
