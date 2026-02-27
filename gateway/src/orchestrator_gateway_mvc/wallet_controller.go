package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	VERIFY_FOR_BRIDGE_WALLET_ENDPOINT = "/wallet/verify-for-bridge"
)

type WalletController struct {
	logger                *zap.Logger
	verifyWalletForBridge use_cases.VerifyWalletForBridgeInterface
}

func NewWalletController(
	logger *zap.Logger,
	verifyWalletForBridge use_cases.VerifyWalletForBridgeInterface,
) *WalletController {
	return &WalletController{
		logger:                logger,
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
// @Success      200 {string} string "OK"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /wallet/verify-for-bridge [post]
func (controller *WalletController) VerifyWalletForBridge(c *fiber.Ctx) error {
	controller.logger.Info("Verify wallet for bridge request received")
	controller.logger.Debug("Request", zap.Any("params", c.AllParams()))

	request := new(requests.VerifyWalletForBridgeRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.verifyWalletForBridge.Verify(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to verify wallet for bridge", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Verify wallet for bridge request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}
