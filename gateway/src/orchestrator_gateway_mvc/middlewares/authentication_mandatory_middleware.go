package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthenticationMandatoryMiddleware struct {
	logger *zap.Logger
}

func NewAuthenticationMandatoryMiddleware(
	logger *zap.Logger,
) *AuthenticationMandatoryMiddleware {
	return &AuthenticationMandatoryMiddleware{
		logger: logger,
	}
}

func (m *AuthenticationMandatoryMiddleware) Handle(c *fiber.Ctx) error {
	alreadyInjectedAuthenticatedWallet := c.Locals("authenticated-wallet")
	if alreadyInjectedAuthenticatedWallet == nil {
		return c.SendStatus(401)
	}

	return c.Next()
}
