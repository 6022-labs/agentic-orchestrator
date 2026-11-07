package orchestrator_gateway_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	VERIFY_FOR_BRIDGE_WALLET_ENDPOINT = "/wallet/verify-for-bridge"
)

type WalletController struct {
	logger                *zap.Logger
	errorWriter           *mvc.ErrorWriter
	verifyWalletForBridge use_cases.VerifyWalletForBridgeInterface
}

func NewWalletController(
	logger *zap.Logger,
	errorWriter *mvc.ErrorWriter,
	verifyWalletForBridge use_cases.VerifyWalletForBridgeInterface,
) *WalletController {
	return &WalletController{
		logger:                logger,
		errorWriter:           errorWriter,
		verifyWalletForBridge: verifyWalletForBridge,
	}
}

func (controller *WalletController) InitAuthenticatedRoutes(app fiber.Router) {
	app.Post(VERIFY_FOR_BRIDGE_WALLET_ENDPOINT, controller.VerifyWalletForBridge)
}

func (controller *WalletController) InitUnauthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

// VerifyWalletForBridge godoc
// @Summary      verify wallet for bridge
// @Description  Verify wallet ownership for platform bridge
// @Tags         wallet
// @Accept       multipart/form-data
// @Produce      json
// @Param        request body requests.VerifyWalletForBridgeRequest true "wallet verification request"
// @Success      204 "No Content"
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_body"
// @Failure      401 {object} errors.UnauthorizedError "unauthorized"
// @Failure      500 {object} errors.InternalError
// @Security     BearerAuth
// @Router       /wallet/verify-for-bridge [post]
func (controller *WalletController) VerifyWalletForBridge(c fiber.Ctx) error {
	controller.logger.Info("Verify wallet for bridge request received")
	params := make(map[string]string, len(c.Route().Params))
	for _, name := range c.Route().Params {
		params[name] = c.Params(name)
	}
	controller.logger.Debug("Request", zap.Any("params", params))

	request := new(requests.VerifyWalletForBridgeRequest)
	if err := c.Bind().Body(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_body", err.Error()))
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.verifyWalletForBridge.Verify(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	controller.logger.Info("Verify wallet for bridge request completed successfully")

	return c.SendStatus(fiber.StatusNoContent)
}
