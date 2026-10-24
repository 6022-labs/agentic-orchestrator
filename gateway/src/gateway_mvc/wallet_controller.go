package orchestrator_gateway_mvc

import (
	_ "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v3"
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
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.verifyWalletForBridge.Verify(c.Context(), request)
	if err != nil {
		return mvc.WriteError(c, err)
	}

	controller.logger.Info("Verify wallet for bridge request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}
