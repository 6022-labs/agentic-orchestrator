package registry_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/gofiber/fiber/v3"
)

const (
	LIST_AGENTS_ENDPOINT             = "/agents"
	GET_AGENT_ENDPOINT               = "/agent/:agentCollectionAddress/:agentCollectionTokenId"
	GET_AGENT_BY_ENS_DOMAIN_ENDPOINT = "/agent/ens/:ensDomain"
)

type AgentController struct {
	errorWriter               *mvc.ErrorWriter
	getAgent                  use_cases.GetAgentInterface
	getMintedAgentByEnsDomain use_cases.GetMintedAgentByEnsDomainInterface
	listAgents                use_cases.ListAgentsInterface
}

func NewAgentController(
	errorWriter *mvc.ErrorWriter,
	getAgent use_cases.GetAgentInterface,
	getMintedAgentByEnsDomain use_cases.GetMintedAgentByEnsDomainInterface,
	listAgents use_cases.ListAgentsInterface,
) *AgentController {
	return &AgentController{
		errorWriter:               errorWriter,
		getAgent:                  getAgent,
		getMintedAgentByEnsDomain: getMintedAgentByEnsDomain,
		listAgents:                listAgents,
	}
}

func (controller *AgentController) RegisterRoutes(app fiber.Router) {
	app.Get(GET_AGENT_BY_ENS_DOMAIN_ENDPOINT, controller.GetMintedAgentByEnsDomain)
	app.Get(GET_AGENT_ENDPOINT, controller.GetAgent)
	app.Get(LIST_AGENTS_ENDPOINT, controller.ListAgents)
}

// GetAgent godoc
// @Summary      Get agent
// @Description  Retrieves details of a specific agent by collection address and token ID
// @Tags         Agent
// @Produce      json
// @Param        agentCollectionAddress path string true "Agent collection address"
// @Param        agentCollectionTokenId path int true "Agent collection token ID"
// @Success      200 {object} responses.AgentResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      404 {object} errors.NotFoundError "agent_not_found"
// @Failure      500 {object} errors.InternalError "internal_error"
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
		return controller.errorWriter.Write(c, apperrors.NewNotFoundError("agent_not_found", "Agent not found"))
	}

	return c.JSON(response)
}

// ListAgents godoc
// @Summary      List agents
// @Description  Retrieves a paginated list of agents with optional filtering
// @Tags         Agent
// @Produce      json
// @Param        page query int false "Page number"
// @Param        elementsPerPage query int false "Elements per page"
// @Param        sortDirection query string false "Sort direction (asc, desc)"
// @Param        sortBy query string false "Sort by field"
// @Param        isMinted query bool false "Filter on mint state (true = minted only, false = proposals only)"
// @Success      200 {object} responses.ListAgentsResponse
// @Failure      400 {object} errors.BadRequestError "invalid_query"
// @Failure      500 {object} errors.InternalError "internal_error"
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

// GetMintedAgentByEnsDomain godoc
// @Summary      Get agent by ENS domain
// @Description  Retrieves a minted agent by its ENS domain name
// @Tags         Agent
// @Produce      json
// @Param        ensDomain path string true "ENS domain name"
// @Success      200 {object} responses.AgentResponse
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      404 {object} errors.NotFoundError "agent_not_found"
// @Failure      500 {object} errors.InternalError "internal_error"
// @Router       /agent/ens/{ensDomain} [get]
func (controller *AgentController) GetMintedAgentByEnsDomain(c fiber.Ctx) error {
	request := new(requests.GetAgentByEnsDomainRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(c, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	response, err := controller.getMintedAgentByEnsDomain.Get(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}
	if response == nil {
		return controller.errorWriter.Write(c, apperrors.NewNotFoundError("agent_not_found", "Agent not found"))
	}

	return c.JSON(response)
}
