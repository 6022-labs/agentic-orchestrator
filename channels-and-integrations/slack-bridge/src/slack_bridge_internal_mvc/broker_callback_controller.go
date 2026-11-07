package slack_bridge_internal_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	BROKER_CALLBACK_ENDPOINT = "/broker/callback"
)

type BrokerCallbackController struct {
	logger                      *zap.Logger
	errorWriter                 *mvc.ErrorWriter
	postBrokerCompletionToSlack use_cases.PostBrokerCompletionToSlackInterface
}

func NewBrokerCallbackController(
	logger *zap.Logger,
	errorWriter *mvc.ErrorWriter,
	postBrokerCompletionToSlack use_cases.PostBrokerCompletionToSlackInterface,
) *BrokerCallbackController {
	return &BrokerCallbackController{
		logger:                      logger,
		errorWriter:                 errorWriter,
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
// @Produce      json
// @Param        request body requests.PostBrokerCompletionToSlackRequest true "Broker callback payload"
// @Success      204 {string} string ""
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_body"
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {object} errors.InternalError
// @Router       /broker/callback [post]
func (c *BrokerCallbackController) Handle(ctx fiber.Ctx) error {
	request := new(requests.PostBrokerCompletionToSlackRequest)
	if err := ctx.Bind().Body(request); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_body", err.Error()))
	}

	if err := c.postBrokerCompletionToSlack.Post(ctx.Context(), request); err != nil {
		return c.errorWriter.Write(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
