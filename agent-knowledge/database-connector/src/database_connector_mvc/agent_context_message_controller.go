package database_connector_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/zapext"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	ADD_MESSAGES_TO_AGENT_CONTEXT = "/agent/context/messages"
	GET_AGENT_CONTEXT_MESSAGES    = "/agent/context/messages"
)

type AgentContextMessageController struct {
	logger                    *zap.Logger
	getAgentContextMessages   use_cases.GetAgentContextMessagesInterface
	addMessagesToAgentContext use_cases.AddMessagesToAgentContextInterface
}

func NewAgentContextMessageController(
	logger *zap.Logger,
	getAgentContextMessages use_cases.GetAgentContextMessagesInterface,
	addMessagesToAgentContext use_cases.AddMessagesToAgentContextInterface,
) *AgentContextMessageController {
	return &AgentContextMessageController{
		logger:                    logger,
		getAgentContextMessages:   getAgentContextMessages,
		addMessagesToAgentContext: addMessagesToAgentContext,
	}
}

func (controller *AgentContextMessageController) InitRoutes(app fiber.Router) {
	app.Get(GET_AGENT_CONTEXT_MESSAGES, controller.GetAgentContextMessages)
	app.Post(ADD_MESSAGES_TO_AGENT_CONTEXT, controller.AddMessagesToAgentContext)
}

// AddMessageToAgentContext adds a message to the agent context.
// @Summary Add message to agent context
// @Description Adds a message to the agent context for the specified agent
// @Tags AgentContextMessage
// @Accept json
// @Produce json
// @Param request body requests.AddMessagesToAgentContextRequest true "Add Message To Agent Context Request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/context/messages [post]
func (controller *AgentContextMessageController) AddMessagesToAgentContext(c *fiber.Ctx) error {
	controller.logger.Info("Received request to add message to agent context")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.AddMessagesToAgentContextRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.addMessagesToAgentContext.Add(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to add message to agent context", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully added message to agent context")

	return c.SendStatus(fiber.StatusCreated)
}

// GetAgentContextMessages retrieves messages from the agent context.
// @Summary Get agent context messages
// @Description Retrieves messages from the agent context for the specified agent
// @Tags AgentContextMessage
// @Accept json
// @Produce json
// @Param request query requests.GetAgentContextMessagesRequest true "Get Agent Context Messages Request"
// @Success 200 {object} []responses.GetAgentContextMessageResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/context/messages [get]
func (controller *AgentContextMessageController) GetAgentContextMessages(c *fiber.Ctx) error {
	controller.logger.Info("Received request to get agent context messages")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.GetAgentContextMessagesRequest)
	if err := c.QueryParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.getAgentContextMessages.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get agent context messages", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully retrieved agent context messages")
	controller.logger.Debug(
		"Response",
		zapext.ArrayTruncateStart("response", response, 5),
	)

	return c.JSON(response)
}
