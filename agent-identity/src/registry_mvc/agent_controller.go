package registry_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	LIST_AGENTS_ENDPOINT             = "/agents"
	GET_AGENT_ENDPOINT               = "/agent/:agentCollectionAddress/:agentCollectionTokenId"
	GET_AGENT_BY_ENS_DOMAIN_ENDPOINT = "/agent/ens/:ensDomain"
)

type AgentController struct {
	logger                    *zap.Logger
	getAgent                  use_cases.GetAgentInterface
	getMintedAgentByEnsDomain use_cases.GetMintedAgentByEnsDomainInterface
	listAgents                use_cases.ListAgentsInterface
}

func NewAgentController(
	logger *zap.Logger,
	getAgent use_cases.GetAgentInterface,
	getMintedAgentByEnsDomain use_cases.GetMintedAgentByEnsDomainInterface,
	listAgents use_cases.ListAgentsInterface,
) *AgentController {
	return &AgentController{
		logger:                    logger,
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
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId} [get]
func (controller *AgentController) GetAgent(c *fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	response, err := controller.getAgent.Get(request)
	if err != nil {
		controller.logger.Error("Failed to get agent", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
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
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agents [get]
func (controller *AgentController) ListAgents(c *fiber.Ctx) error {
	request := new(requests.ListAgentsRequest)
	if err := c.QueryParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	resp, err := controller.listAgents.List(request)
	if err != nil {
		controller.logger.Error("Failed to list agents", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
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
// @Failure      404 {string} string "Agent not found"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/ens/{ensDomain} [get]
func (controller *AgentController) GetMintedAgentByEnsDomain(c *fiber.Ctx) error {
	request := new(requests.GetAgentByEnsDomainRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request params", zap.Error(err))
		return err
	}

	response, err := controller.getMintedAgentByEnsDomain.Get(request)
	if err != nil {
		controller.logger.Error("Failed to get agent by ENS domain", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	if response == nil {
		return c.Status(fiber.StatusNotFound).SendString("Agent not found")
	}

	return c.JSON(response)
}
