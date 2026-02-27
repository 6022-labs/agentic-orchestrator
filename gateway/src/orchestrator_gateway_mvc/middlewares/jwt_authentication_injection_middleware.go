package middlewares

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JWTAuthenticationInjectionMiddleware struct {
	logger      *zap.Logger
	jwtSettings *settings.JwtSettings
	jwtManager  interfaces.JwtManagerInterface
}

func NewJWTAuthenticationInjectionMiddleware(
	logger *zap.Logger,
	jwtSettings *settings.JwtSettings,
	jwtManager interfaces.JwtManagerInterface,
) *JWTAuthenticationInjectionMiddleware {
	return &JWTAuthenticationInjectionMiddleware{
		logger:      logger,
		jwtManager:  jwtManager,
		jwtSettings: jwtSettings,
	}
}

func (m *JWTAuthenticationInjectionMiddleware) Handle(c *fiber.Ctx) error {
	alreadyInjectedAuthenticatedWallet := c.Locals("authenticated-wallet")
	if alreadyInjectedAuthenticatedWallet != nil {
		return c.Next()
	}

	authorizationHeaders := c.GetReqHeaders()["Authorization"]
	if len(authorizationHeaders) == 0 {
		m.logger.Info("No JWT found in request headers")
		return c.Next()
	}
	authorizationHeader := authorizationHeaders[0]
	if len(strings.TrimSpace(authorizationHeader)) == 0 {
		m.logger.Info("No JWT found in request headers")
		return c.Next()
	}

	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderParts) != 2 {
		m.logger.Info("Invalid JWT format in request headers")
		return c.Next()
	}

	if strings.ToLower(authorizationHeaderParts[0]) != "bearer" {
		m.logger.Info("Invalid JWT format in request headers")
		return c.Next()
	}

	jwtStr := authorizationHeaderParts[1]
	if len(strings.TrimSpace(jwtStr)) == 0 {
		m.logger.Info("No JWT found in request headers")
		return c.Next()
	}

	jwtVal, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			m.logger.Info("Invalid JWT signing method")
			return nil, nil
		}
		return []byte(m.jwtSettings.Hash), nil
	})
	if err != nil {
		m.logger.Error("Invalid JWT", zap.Error(err))
		return c.Next()
	}

	authenticatedWalletStr, err := m.jwtManager.ExtractSubFromToken(jwtVal, m.jwtSettings.Hash)
	if err != nil {
		m.logger.Error("Cannot extract user id from JWT", zap.Error(err))
		return c.Next()
	}

	authenticatedWallet := common.HexToAddress(*authenticatedWalletStr)

	c.Locals("authenticated-wallet", authenticatedWallet)

	return c.Next()
}
