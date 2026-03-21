package middlewares

import (
	"encoding/base64"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SignatureAuthenticationInjectionMiddleware struct {
	logger         *zap.Logger
	ecdsaSignature interfaces.ECDSASignerInterface
}

func NewSignatureAuthenticationInjectionMiddleware(
	logger *zap.Logger,
	ecdsaSignature interfaces.ECDSASignerInterface,
) *SignatureAuthenticationInjectionMiddleware {
	return &SignatureAuthenticationInjectionMiddleware{
		logger:         logger,
		ecdsaSignature: ecdsaSignature,
	}
}

func (m *SignatureAuthenticationInjectionMiddleware) Handle(c *fiber.Ctx) error {
	alreadyInjectedAuthenticatedWallet := c.Locals("authenticated-wallet")
	if alreadyInjectedAuthenticatedWallet != nil {
		return c.Next()
	}

	requestQuery := c.Context().QueryArgs().QueryString()
	requestBody := c.Body()

	payload := append(requestQuery, requestBody...)

	signatureHeaders := c.GetReqHeaders()["X-Agentic-Orchestrator-Signature"]
	// Check if there is a item in the array
	if len(signatureHeaders) == 0 {
		m.logger.Warn("No signature found in headers")
		return c.Next()
	}

	signature, err := base64.StdEncoding.DecodeString(signatureHeaders[0])
	if err != nil {
		m.logger.Error("Cannot decode signature", zap.Error(err))
		return c.Next()
	}

	authenticatedWalletStr, err := m.ecdsaSignature.RecoverPublicKey(payload, signature)
	if err != nil {
		m.logger.Error("Cannot recover public key from signature", zap.Error(err))
		return c.Next()
	}

	authenticatedWallet := common.HexToAddress(*authenticatedWalletStr)

	c.Locals("authenticated-wallet", authenticatedWallet)

	return c.Next()
}
