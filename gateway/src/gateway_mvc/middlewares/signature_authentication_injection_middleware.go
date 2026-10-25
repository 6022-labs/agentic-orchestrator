package middlewares

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const signatureHeader = "X-Agentic-Orchestrator-Signature"

type SignatureAuthenticationInjectionMiddleware struct {
	logger                    *zap.Logger
	authenticateSignedRequest use_cases.AuthenticateSignedRequestInterface
}

func NewSignatureAuthenticationInjectionMiddleware(
	logger *zap.Logger,
	authenticateSignedRequest use_cases.AuthenticateSignedRequestInterface,
) *SignatureAuthenticationInjectionMiddleware {
	return &SignatureAuthenticationInjectionMiddleware{
		logger:                    logger,
		authenticateSignedRequest: authenticateSignedRequest,
	}
}

func (m *SignatureAuthenticationInjectionMiddleware) Handle(c fiber.Ctx) error {
	if _, ok := helpers.CallerIdFromContext(c.Context()); ok {
		return c.Next()
	}

	if c.Get(signatureHeader) == "" {
		return c.Next()
	}

	request := new(requests.AuthenticateSignedRequestRequest)
	if err := c.Bind().Header(request); err != nil {
		m.logger.Info("Cannot read signed request headers", zap.Error(err))
		return c.Next()
	}
	request.Method = c.Method()
	request.Path = c.Path()
	request.Query = string(c.RequestCtx().URI().QueryString())
	request.Body = c.Body()

	response, err := m.authenticateSignedRequest.Execute(c.Context(), request)
	if err != nil {
		m.logger.Info("Rejected signed request", zap.Error(err))
		return c.Next()
	}

	authenticatedWallet := common.HexToAddress(response.Wallet)

	c.SetContext(helpers.ContextWithCallerId(c.Context(), authenticatedWallet.Hex()))
	c.Locals(AuthMethodLocal, metrics_interfaces.AuthMethodSignature)

	return c.Next()
}
