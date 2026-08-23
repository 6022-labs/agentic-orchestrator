package middlewares

import (
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/metrics/interfaces"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AuthMethodLocal carries which injector authenticated the request, for the auth metric label.
const AuthMethodLocal = "auth-method"

type AuthenticationMandatoryMiddleware struct {
	logger         *zap.Logger
	gatewayMetrics metrics_interfaces.GatewayMetricsInterface
}

func NewAuthenticationMandatoryMiddleware(
	logger *zap.Logger,
	gatewayMetrics metrics_interfaces.GatewayMetricsInterface,
) *AuthenticationMandatoryMiddleware {
	return &AuthenticationMandatoryMiddleware{
		logger:         logger,
		gatewayMetrics: gatewayMetrics,
	}
}

func (m *AuthenticationMandatoryMiddleware) Handle(c *fiber.Ctx) error {
	alreadyInjectedAuthenticatedWallet := c.Locals("authenticated-wallet")
	if alreadyInjectedAuthenticatedWallet == nil {
		m.gatewayMetrics.RecordAuth(c.UserContext(), metrics_interfaces.AuthMethodNone, false)
		return c.SendStatus(401)
	}

	method, ok := c.Locals(AuthMethodLocal).(string)
	if !ok {
		method = metrics_interfaces.AuthMethodNone
	}
	m.gatewayMetrics.RecordAuth(c.UserContext(), method, true)

	return c.Next()
}
