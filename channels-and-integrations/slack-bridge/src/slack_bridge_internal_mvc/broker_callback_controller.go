package slack_bridge_internal_mvc

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	BROKER_CALLBACK_ENDPOINT = "/internal/broker/callback"
)

type BrokerCallbackController struct {
	logger                      *zap.Logger
	postBrokerCompletionToSlack use_cases.PostBrokerCompletionToSlackInterface
}

func NewBrokerCallbackController(
	logger *zap.Logger,
	postBrokerCompletionToSlack use_cases.PostBrokerCompletionToSlackInterface,
) *BrokerCallbackController {
	return &BrokerCallbackController{
		logger:                      logger,
		postBrokerCompletionToSlack: postBrokerCompletionToSlack,
	}
}

func (c *BrokerCallbackController) InitRoutes(app fiber.Router) {
	app.Post(BROKER_CALLBACK_ENDPOINT, c.Handle)
}

// Handle godoc
// @Summary      Conversation broker completion callback
// @Description  Receives an agent completion from the conversation broker and posts it to the originating Slack thread.
// @Tags         BrokerCallback
// @Accept       json
// @Produce      plain
// @Param        request body requests.PostBrokerCompletionToSlackRequest true "Broker callback payload"
// @Success      204 {string} string ""
// @Failure      400 {string} string "Bad Request"
// @Failure      404 {string} string "Not Found"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /internal/broker/callback [post]
func (c *BrokerCallbackController) Handle(ctx *fiber.Ctx) error {
	request := new(requests.PostBrokerCompletionToSlackRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.logger.Warn("Failed to parse broker callback payload", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := c.postBrokerCompletionToSlack.Post(request); err != nil {
		if errors.Is(err, use_cases.ErrConversationLinkNotFound) || errors.Is(err, use_cases.ErrSlackAgentNotFound) {
			c.logger.Warn("Broker callback referred to unknown entity", zap.Error(err))
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		c.logger.Error("Failed to post broker completion to slack", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
