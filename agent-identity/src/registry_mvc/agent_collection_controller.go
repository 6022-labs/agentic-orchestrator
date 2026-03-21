package registry_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	LIST_AGENT_COLLECTIONS_ENDPOINT = "/agent-collections"
)

type AgentCollectionController struct {
	logger               *zap.Logger
	listAgentCollections use_cases.ListAgentCollectionsInterface
}

func NewAgentCollectionController(
	logger *zap.Logger,
	listAgentCollections use_cases.ListAgentCollectionsInterface,
) *AgentCollectionController {
	return &AgentCollectionController{
		logger:               logger,
		listAgentCollections: listAgentCollections,
	}
}

func (controller *AgentCollectionController) RegisterRoutes(app fiber.Router) {
	app.Get(LIST_AGENT_COLLECTIONS_ENDPOINT, controller.ListAgentCollections)
}

func (controller *AgentCollectionController) ListAgentCollections(c *fiber.Ctx) error {
	collections, err := controller.listAgentCollections.List()
	if err != nil {
		controller.logger.Error("Failed to list agent collections", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(collections)
}
