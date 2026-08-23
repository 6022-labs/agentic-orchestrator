package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/zapext"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	LIST_AGENTS_ENDPOINT    = "/agents"
	GET_AGENT_ENDPOINT      = "/agent/:agentCollectionAddress/:agentCollectionTokenId"
	GET_AGENT_CARD_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/card"
	GET_AGENT_URL_ENDPOINT  = "/agent/:agentCollectionAddress/:agentCollectionTokenId/url"

	GET_AGENT_PROMPT_TEMPLATE_ENDPOINT  = "/agent/:agentCollectionAddress/:agentCollectionTokenId/prompt-template/:key"
	SAVE_AGENT_PROMPT_TEMPLATE_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/prompt-template/:key"
)

type AgentController struct {
	logger                  *zap.Logger
	getAgent                use_cases.GetAgentInterface
	listAgents              use_cases.ListAgentsInterface
	getAgentCard            use_cases.GetAgentCardInterface
	getAgentUrl             use_cases.GetAgentUrlInterface
	getAgentPromptTemplate  use_cases.GetAgentPromptTemplateInterface
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface
}

func NewAgentController(
	logger *zap.Logger,
	getAgent use_cases.GetAgentInterface,
	listAgents use_cases.ListAgentsInterface,
	getAgentCard use_cases.GetAgentCardInterface,
	getAgentUrl use_cases.GetAgentUrlInterface,
	getAgentPromptTemplate use_cases.GetAgentPromptTemplateInterface,
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface,
) *AgentController {
	return &AgentController{
		logger:                  logger,
		getAgent:                getAgent,
		listAgents:              listAgents,
		getAgentCard:            getAgentCard,
		getAgentUrl:             getAgentUrl,
		getAgentPromptTemplate:  getAgentPromptTemplate,
		saveAgentPromptTemplate: saveAgentPromptTemplate,
	}
}

func (controller *AgentController) InitAuthenticatedRoutes(app fiber.Router) {
	app.Get(GET_AGENT_CARD_ENDPOINT, controller.GetAgentCard)
	app.Get(GET_AGENT_URL_ENDPOINT, controller.GetAgentUrl)
	app.Get(GET_AGENT_PROMPT_TEMPLATE_ENDPOINT, controller.GetAgentPromptTemplate)
	app.Post(SAVE_AGENT_PROMPT_TEMPLATE_ENDPOINT, controller.SaveAgentPromptTemplate)
}

func (controller *AgentController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Get(GET_AGENT_ENDPOINT, controller.GetAgent)
	app.Get(LIST_AGENTS_ENDPOINT, controller.ListAgents)
}

// GetAgent godoc
// @Summary      Get agent information
// @Description  Retrieves details of a specific agent
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.AgentResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId} [get]
func (controller *AgentController) GetAgent(c fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	response, err := controller.getAgent.Get(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to get agent", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	if response == nil {
		return c.Status(fiber.StatusNotFound).SendString("agent not found")
	}

	return c.JSON(response)
}

// ListAgents godoc
// @Summary      List agents
// @Description  Returns all agents synced from the chain
// @Tags         agent
// @Produce      json
// @Param        page query int false "Page number"
// @Param        elementsPerPage query int false "Elements per page"
// @Param        sortDirection query string false "Sort direction (asc, desc)"
// @Param        sortBy query string false "Sort by field"
// @Param        isMinted query bool false "Filter on mint state (true = minted only, false = proposals only)"
// @Success      200 {object} responses.ListAgentsResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agents [get]
func (controller *AgentController) ListAgents(c fiber.Ctx) error {
	request := new(requests.ListAgentsRequest)
	if err := c.Bind().Query(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	resp, err := controller.listAgents.List(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to list agents", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Debug(
		"List agents response",
		zap.Any("page", resp.Page),
		zap.Any("totalElements", resp.TotalElements),
		zapext.ArrayTruncateEnd("agents", resp.Agents, 20),
	)

	return c.JSON(resp)
}

// GetAgentCard godoc
// @Summary      Get agent node card
// @Description  Returns the live /.well-known/6022 card of the agent's node
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.AgentCardResponse
// @Failure      404 {object} errors.NotFoundError
// @Failure      502 {string} string "Bad Gateway"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/card [get]
func (controller *AgentController) GetAgentCard(c fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	response, err := controller.getAgentCard.Get(c.Context(), request)
	if err != nil {
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return c.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		controller.logger.Warn("Failed to get agent card", zap.Error(err))
		return c.Status(fiber.StatusBadGateway).SendString("failed to reach agent node")
	}

	return c.JSON(response)
}

// GetAgentUrl godoc
// @Summary      Get agent node url
// @Description  Returns the agent node's URL resolved from its ENS url text record
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.AgentUrlResponse
// @Failure      404 {object} errors.NotFoundError
// @Failure      502 {string} string "Bad Gateway"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/url [get]
func (controller *AgentController) GetAgentUrl(c fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	response, err := controller.getAgentUrl.Get(c.Context(), request)
	if err != nil {
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return c.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		controller.logger.Warn("Failed to get agent url", zap.Error(err))
		return c.Status(fiber.StatusBadGateway).SendString("failed to resolve agent url")
	}

	return c.JSON(response)
}

// GetAgentPromptTemplate godoc
// @Summary      Get agent prompt template
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        key path string true "template key"
// @Success      200 {object} responses.AgentPromptTemplateResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/prompt-template/{key} [get]
func (controller *AgentController) GetAgentPromptTemplate(c fiber.Ctx) error {
	request := new(requests.GetAgentPromptTemplateRequest)
	if err := c.Bind().URI(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	response, err := controller.getAgentPromptTemplate.Get(c.Context(), request)
	if err != nil {
		return controller.writeUseCaseError(c, err)
	}

	return c.JSON(response)
}

// SaveAgentPromptTemplate godoc
// @Summary      Save agent prompt template
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        key path string true "template key"
// @Param        request body requests.SaveAgentPromptTemplateRequest true "template content"
// @Success      200 {object} responses.AgentPromptTemplateResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/prompt-template/{key} [post]
func (controller *AgentController) SaveAgentPromptTemplate(c fiber.Ctx) error {
	request := new(requests.SaveAgentPromptTemplateRequest)
	if err := c.Bind().Body(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := c.Bind().URI(&request.GetAgentPromptTemplateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	response, err := controller.saveAgentPromptTemplate.Save(c.Context(), request)
	if err != nil {
		return controller.writeUseCaseError(c, err)
	}

	return c.JSON(response)
}

// writeUseCaseError maps domain error types to HTTP statuses.
func (controller *AgentController) writeUseCaseError(c fiber.Ctx, err error) error {
	if validationErr, ok := err.(*errors.ValidationError); ok {
		return c.Status(fiber.StatusBadRequest).JSON(validationErr)
	}
	if notFoundErr, ok := err.(*errors.NotFoundError); ok {
		return c.Status(fiber.StatusNotFound).JSON(notFoundErr)
	}
	controller.logger.Error("internal error", zap.Error(err))
	return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
}
