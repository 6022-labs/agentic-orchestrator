package database_connector_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	CREATE_AGENT_CONTEXT_ENDPOINT     = "/agent/context"
	GET_AGENT_CONTEXTS_ENDPOINT       = "/agent/:agentId/contexts"
	GET_AGENT_CONTEXT_BY_KEY_ENDPOINT = "/agent/:agentId/context/:key"
)

type AgentContextController struct {
	logger               *zap.Logger
	getAgentContexts     use_cases.GetAgentContextsInterface
	createAgentContext   use_cases.CreateAgentContextInterface
	getAgentContextByKey use_cases.GetAgentContextByKeyInterface
}

func NewAgentContextController(
	logger *zap.Logger,
	getAgentContexts use_cases.GetAgentContextsInterface,
	createAgentContext use_cases.CreateAgentContextInterface,
	getAgentContextByKey use_cases.GetAgentContextByKeyInterface,
) *AgentContextController {
	return &AgentContextController{
		logger:               logger,
		getAgentContexts:     getAgentContexts,
		createAgentContext:   createAgentContext,
		getAgentContextByKey: getAgentContextByKey,
	}
}

func (controller *AgentContextController) InitRoutes(app fiber.Router) {
	app.Post(CREATE_AGENT_CONTEXT_ENDPOINT, controller.CreateAgentContext)
	app.Get(GET_AGENT_CONTEXTS_ENDPOINT, controller.GetAgentContexts)
	app.Get(GET_AGENT_CONTEXT_BY_KEY_ENDPOINT, controller.GetAgentContextByKey)
}

// CreateAgentContext creates a new agent context.
// @Summary Create agent context
// @Description Creates a new agent context for the specified agent
// @Tags AgentContext
// @Accept json
// @Produce json
// @Param request body requests.CreateAgentContextRequest true "Create Agent Context Request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/context [post]
func (controller *AgentContextController) CreateAgentContext(c *fiber.Ctx) error {
	controller.logger.Info("Received request to create agent context")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.CreateAgentContextRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.createAgentContext.Create(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to create agent context", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully created agent context")

	return c.SendStatus(fiber.StatusCreated)
}

// GetAgentContextKeys retrieves the keys of the agent context.
// @Summary Get agent context keys
// @Description Retrieves the keys of the agent context for the specified agent
// @Tags AgentContext
// @Accept json
// @Produce json
// @Param agentId path string true "Agent ID"
// @Param request query requests.GetAgentContextsRequest true "Get Agent Context Keys Request"
// @Success 200 {object} []responses.GetAgentContextResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/{agentId}/contexts [get]
func (controller *AgentContextController) GetAgentContexts(c *fiber.Ctx) error {
	controller.logger.Info("Received request to get agent contexts")
	controller.logger.Debug(
		"Request",
		zap.Any("params", c.AllParams()),
		zap.Any("query", c.Queries()),
	)

	request := new(requests.GetAgentContextsRequest)
	if err := c.ParamsParser(request); err != nil {
		return err
	}

	if err := c.QueryParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.getAgentContexts.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get agent contexts", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully retrieved agent contexts")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// GetAgentContextByKey retrieves a specific agent context by key.
// @Summary Get agent context by key
// @Description Retrieves a specific agent context by key for the specified agent
// @Tags AgentContext
// @Accept json
// @Produce json
// @Param agentId path string true "Agent ID"
// @Param key path string true "Context Key"
// @Success 200 {object} responses.GetAgentContextResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/{agentId}/context/{key} [get]
func (controller *AgentContextController) GetAgentContextByKey(c *fiber.Ctx) error {
	controller.logger.Info("Received request to get agent context by key")
	controller.logger.Debug("Request", zap.Any("params", c.AllParams()))

	request := new(requests.GetAgentContextByKeyRequest)
	if err := c.ParamsParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.getAgentContextByKey.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get agent context by key", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully retrieved agent context by key")
	controller.logger.Debug("Response", zap.Any("response", response))

	if response == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(response)
}
