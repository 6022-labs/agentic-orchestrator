package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/zapext"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	LIST_AGENTS_ENDPOINT = "/agents"
	GET_AGENT_ENDPOINT   = "/agent/:agentCollectionAddress/:agentCollectionTokenId"
)

type AgentController struct {
	logger     *zap.Logger
	getAgent   use_cases.GetAgentInterface
	listAgents use_cases.ListAgentsInterface
}

func NewAgentController(
	logger *zap.Logger,
	getAgent use_cases.GetAgentInterface,
	listAgents use_cases.ListAgentsInterface,
) *AgentController {
	return &AgentController{
		logger:     logger,
		getAgent:   getAgent,
		listAgents: listAgents,
	}
}

func (controller *AgentController) InitAuthenticatedRoutes(app fiber.Router) {
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
func (controller *AgentController) GetAgent(c *fiber.Ctx) error {
	controller.logger.Info("Get agent request received")
	controller.logger.Debug("Request", zap.Any("params", c.AllParams()))

	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.getAgent.Get(request)
	if err != nil {
		controller.logger.Error("Failed to get agent", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Get agent request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// ListAgents godoc
// @Summary      List agents
// @Description  Returns all agents managed by the orchestrator
// @Tags         agent
// @Produce      json
// @Param        request query requests.ListAgentsRequest true "list agents request"
// @Success      200 {object} responses.ListAgentsResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agents [get]
func (controller *AgentController) ListAgents(c *fiber.Ctx) error {
	controller.logger.Info("List agents request received")
	controller.logger.Debug("Request", zap.Any("queries", c.Queries()))

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

	controller.logger.Info("List agents request completed successfully", zap.Int("agents_count", len(resp.Agents)))
	controller.logger.Debug(
		"List agents response",
		zap.Any("page", resp.Page),
		zap.Any("totalElements", resp.TotalElements),
		zap.Any("elementsPerPage", resp.ElementsPerPage),
		zapext.ArrayTruncateEnd("agents", resp.Agents, 20),
	)

	return c.JSON(resp)
}
