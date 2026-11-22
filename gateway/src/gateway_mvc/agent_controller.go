package orchestrator_gateway_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/gofiber/fiber/v3"
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
	errorWriter             *mvc.ErrorWriter
	getAgent                use_cases.GetAgentInterface
	listAgents              use_cases.ListAgentsInterface
	getAgentCard            use_cases.GetAgentCardInterface
	getAgentUrl             use_cases.GetAgentUrlInterface
	getAgentPromptTemplate  use_cases.GetAgentPromptTemplateInterface
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface
}

func NewAgentController(
	errorWriter *mvc.ErrorWriter,
	getAgent use_cases.GetAgentInterface,
	listAgents use_cases.ListAgentsInterface,
	getAgentCard use_cases.GetAgentCardInterface,
	getAgentUrl use_cases.GetAgentUrlInterface,
	getAgentPromptTemplate use_cases.GetAgentPromptTemplateInterface,
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface,
) *AgentController {
	return &AgentController{
		errorWriter:             errorWriter,
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
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      404 {object} errors.NotFoundError "agent_not_found"
// @Failure      500 {object} errors.InternalError
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId} [get]
func (controller *AgentController) GetAgent(c fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	response, err := controller.getAgent.Get(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}
	if response == nil {
		return controller.errorWriter.Write(c, apperrors.NewNotFoundError("agent_not_found", "agent not found"))
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
// @Failure      400 {object} errors.BadRequestError "invalid_query"
// @Failure      500 {object} errors.InternalError
// @Router       /agents [get]
func (controller *AgentController) ListAgents(c fiber.Ctx) error {
	request := new(requests.ListAgentsRequest)
	if err := c.Bind().Query(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_query", err.Error()))
	}

	resp, err := controller.listAgents.List(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

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
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      404 {object} errors.NotFoundError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/card [get]
func (controller *AgentController) GetAgentCard(c fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	response, err := controller.getAgentCard.Get(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
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
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      404 {object} errors.NotFoundError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/url [get]
func (controller *AgentController) GetAgentUrl(c fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	response, err := controller.getAgentUrl.Get(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(response)
}

// GetAgentPromptTemplate godoc
// @Summary      Get agent prompt template
// @Description  The caller's wallet must hold the agent NFT or be the agent's own evm address
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        key path string true "template key"
// @Success      200 {object} responses.AgentPromptTemplateResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      401 {object} errors.UnauthorizedError
// @Failure      403 {object} errors.ForbiddenError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {object} errors.InternalError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/prompt-template/{key} [get]
func (controller *AgentController) GetAgentPromptTemplate(c fiber.Ctx) error {
	request := new(requests.GetAgentPromptTemplateRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	response, err := controller.getAgentPromptTemplate.Get(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(response)
}

// SaveAgentPromptTemplate godoc
// @Summary      Save agent prompt template
// @Description  The caller's wallet must hold the agent NFT or be the agent's own evm address
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        key path string true "template key"
// @Param        request body requests.SaveAgentPromptTemplateRequest true "template content"
// @Success      200 {object} responses.AgentPromptTemplateResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_body or invalid_params"
// @Failure      401 {object} errors.UnauthorizedError
// @Failure      403 {object} errors.ForbiddenError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {object} errors.InternalError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/prompt-template/{key} [post]
func (controller *AgentController) SaveAgentPromptTemplate(c fiber.Ctx) error {
	request := new(requests.SaveAgentPromptTemplateRequest)
	if err := c.Bind().Body(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_body", err.Error()))
	}
	if err := c.Bind().URI(&request.GetAgentPromptTemplateRequest); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	response, err := controller.saveAgentPromptTemplate.Save(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(response)
}
