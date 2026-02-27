package slack_bridge_internal_mvc

import (
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
		controller.logger.Error("Failed to set wallet ownership", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Set wallet ownership to Slack user request processed successfully")

	return c.SendStatus(fiber.StatusOK)
}
