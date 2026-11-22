package middlewares

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type JWTAuthenticationInjectionMiddleware struct {
	logger     *zap.Logger
	jwtManager interfaces.JwtManagerInterface
}

func NewJWTAuthenticationInjectionMiddleware(
	logger *zap.Logger,
	jwtManager interfaces.JwtManagerInterface,
) *JWTAuthenticationInjectionMiddleware {
	return &JWTAuthenticationInjectionMiddleware{
		logger:     logger,
		jwtManager: jwtManager,
	}
}

func (m *JWTAuthenticationInjectionMiddleware) Handle(c fiber.Ctx) error {
	if _, ok := helpers.CallerIdFromContext(c.Context()); ok {
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

	authenticatedWalletStr, err := m.jwtManager.VerifyJwtToken(c.Context(), jwtStr)
	if err != nil {
		m.logger.Error("Invalid JWT", zap.Error(err))
		return c.Next()
	}

	authenticatedWallet := common.HexToAddress(authenticatedWalletStr)

	c.SetContext(helpers.ContextWithCallerId(c.Context(), authenticatedWallet.Hex()))
	c.Locals(AuthMethodLocal, metrics_interfaces.AuthMethodJwt)

	return c.Next()
}
