package database_connector_mvc

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	LINK_CONVERSATION_TO_AGENT_CONTEXT       = "/agent/context/link-conversation"
	GET_LINKED_AGENT_CONTEXT_OF_CONVERSATION = "/agent/context/linked-conversation"
)

type AgentContextConversationMappingController struct {
	logger                              *zap.Logger
	linkConversationToAgentContext      use_cases.LinkConversationToAgentContextInterface
	getLinkedAgentContextOfConversation use_cases.GetLinkedAgentContextOfConversationInterface
}

func NewAgentContextConversationMappingController(
	logger *zap.Logger,
	linkConversationToAgentContext use_cases.LinkConversationToAgentContextInterface,
	getLinkedAgentContextOfConversation use_cases.GetLinkedAgentContextOfConversationInterface,
) *AgentContextConversationMappingController {
	return &AgentContextConversationMappingController{
		logger:                              logger,
		linkConversationToAgentContext:      linkConversationToAgentContext,
		getLinkedAgentContextOfConversation: getLinkedAgentContextOfConversation,
	}
}

func (controller *AgentContextConversationMappingController) InitRoutes(app fiber.Router) {
	app.Post(LINK_CONVERSATION_TO_AGENT_CONTEXT, controller.LinkConversationToAgentContext)
	app.Get(GET_LINKED_AGENT_CONTEXT_OF_CONVERSATION, controller.GetLinkedAgentContextOfConversation)
}

// LinkConversationToAgentContext links a conversation to an agent context.
// @Summary Link conversation to agent context
// @Description Links a conversation to an agent context
// @Tags AgentContext
// @Accept json
// @Produce json
// @Param request body requests.LinkConversationToAgentContextRequest true "Link Conversation to Agent Context Request"
// @Success 200 {string} string "OK"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/context/link-conversation [post]
func (controller *AgentContextConversationMappingController) LinkConversationToAgentContext(c *fiber.Ctx) error {
	controller.logger.Info("Received request to link conversation to agent context")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.LinkConversationToAgentContextRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.linkConversationToAgentContext.Link(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully linked conversation to agent context")

	return c.SendStatus(fiber.StatusOK)
}

// GetLinkedAgentContextOfConversation retrieves the linked agent context of a conversation.
// @Summary Get linked agent context of conversation
// @Description Retrieves the linked agent context of a conversation
// @Tags AgentContext
// @Accept json
// @Produce json
// @Param request query requests.GetLinkedAgentContextOfConversationRequest true "Get Linked Agent Context of Conversation Request"
// @Success 200 {object} responses.GetLinkedAgentContextOfConversationResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/context/linked-conversation [get]
func (controller *AgentContextConversationMappingController) GetLinkedAgentContextOfConversation(c *fiber.Ctx) error {
	controller.logger.Info("Received request to get linked agent context of conversation")
	controller.logger.Debug("Request", zap.Any("query", c.Queries()))

	request := new(requests.GetLinkedAgentContextOfConversationRequest)
	if err := c.QueryParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.getLinkedAgentContextOfConversation.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get linked agent context of conversation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully retrieved linked agent context of conversation", zap.Any("response", response))
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}
