package slack_bridge_internal_mvc

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	SET_WALLET_OWNERSHIP_TO_SLACK_USER_ENDPOINT = "/wallet/set-ownership-to-slack-user"
)

type WalletController struct {
	logger                          *zap.Logger
	setWalletOwnershipToSlackPlayer use_cases.SetWalletOwnershipToSlackUserInterface
}

func NewWalletController(
	logger *zap.Logger,
	setWalletOwnershipToSlackPlayer use_cases.SetWalletOwnershipToSlackUserInterface,
) *WalletController {
	return &WalletController{
		logger:                          logger,
		setWalletOwnershipToSlackPlayer: setWalletOwnershipToSlackPlayer,
	}
}

func (controller *WalletController) InitRoutes(app fiber.Router) {
	app.Post(SET_WALLET_OWNERSHIP_TO_SLACK_USER_ENDPOINT, controller.SetWalletOwnershipToSlackUser)
}

// SetWalletOwnershipToSlackUser godoc
// @Summary      Set wallet ownership to Slack user
// @Description  Confirms wallet ownership for a Slack user by providing a signed challenge
// @Tags         Wallet
// @Accept       json
// @Produce      plain
// @Param        request body requests.SetWalletOwnershipToSlackUserRequest true "Wallet ownership verification"
// @Success      200 {string} string "OK"
// @Failure      400 {string} string "Bad Request"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /wallet/set-ownership-to-slack-user [post]
func (controller *WalletController) SetWalletOwnershipToSlackUser(c *fiber.Ctx) error {
	controller.logger.Info("Set wallet ownership to Slack user request received")
	controller.logger.Debug("Request", zap.Any("params", c.AllParams()))

	request := new(requests.SetWalletOwnershipToSlackUserRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.setWalletOwnershipToSlackPlayer.Set(request)
	if err != nil {
		if errors.Is(err, requests.ErrValidation) {
			controller.logger.Warn("Set wallet ownership failed validation", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		controller.logger.Error("Failed to set wallet ownership", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Set wallet ownership to Slack user request processed successfully")

	return c.SendStatus(fiber.StatusOK)
}
