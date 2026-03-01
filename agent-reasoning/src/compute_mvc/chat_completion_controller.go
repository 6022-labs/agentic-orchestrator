package compute_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	CHAT_COMPLETION_ENDPOINT = "/chat/completion"
)

type ChatCompletionController struct {
	logger                      *zap.Logger
	llmConnectionChatCompletion use_cases.LLMConnectionChatCompletionInterface
}

func NewChatCompletionController(
	logger *zap.Logger,
	llmConnectionChatCompletion use_cases.LLMConnectionChatCompletionInterface,
) *ChatCompletionController {
	return &ChatCompletionController{
		logger:                      logger,
		llmConnectionChatCompletion: llmConnectionChatCompletion,
	}
}

func (controller *ChatCompletionController) InitRoutes(app fiber.Router) {
	app.Post(CHAT_COMPLETION_ENDPOINT, controller.ChatCompletion)
}

// @Summary Chat Completion
// @Description Handles chat completion requests
// @Tags Chat Completion
// @Accept json
// @Produce json
// @Param request body requests.LLMConnectionChatCompletionRequest true "Chat Completion Request"
// @Success 200 {object} string "Chat Completion Response"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /chat/completion [post]
func (controller *ChatCompletionController) ChatCompletion(c *fiber.Ctx) error {
	controller.logger.Debug("Request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
	)

	request := new(requests.LLMConnectionChatCompletionRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request",
		zap.String("agentId", request.AgentId.String()),
	)

	response, err := controller.llmConnectionChatCompletion.Complete(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to complete chat",
			zap.String("agentId", request.AgentId.String()),
			zap.Error(err),
		)

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Chat completion successful")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}
