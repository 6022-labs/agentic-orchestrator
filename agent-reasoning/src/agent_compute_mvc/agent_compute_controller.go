package agent_compute_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	AGENT_COMPUTE_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/compute"
)

type AgentComputeController struct {
	logger       *zap.Logger
	agentCompute use_cases.AgentComputeInterface
}

func NewAgentComputeController(
	logger *zap.Logger,
	agentCompute use_cases.AgentComputeInterface,
) *AgentComputeController {
	return &AgentComputeController{
		logger:       logger,
		agentCompute: agentCompute,
	}
}

func (controller *AgentComputeController) InitRoutes(app fiber.Router) {
	app.Post(AGENT_COMPUTE_ENDPOINT, controller.Compute)
}

func (controller *AgentComputeController) Compute(c *fiber.Ctx) error {
	request := new(requests.AgentComputeRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		return err
	}

	response, err := controller.agentCompute.Compute(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to compute agent",
			zap.String("agentCollectionAddress", request.AgentCollectionAddress),
			zap.Uint64("agentCollectionTokenId", request.AgentCollectionTokenId),
			zap.Error(err),
		)

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}
