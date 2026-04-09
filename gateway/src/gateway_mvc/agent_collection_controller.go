package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	_ "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
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

func (controller *AgentCollectionController) InitAuthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

func (controller *AgentCollectionController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Get(LIST_AGENT_COLLECTIONS_ENDPOINT, controller.ListAgentCollections)
}

// ListAgentCollections godoc
// @Summary      List agent collections
// @Description  Retrieves all available agent collections
// @Tags         agent-collections
// @Produce      json
// @Success      200 {array} responses.AgentCollectionResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent-collections [get]
func (controller *AgentCollectionController) ListAgentCollections(c *fiber.Ctx) error {
	controller.logger.Info("List agent collections request received")

	collections, err := controller.listAgentCollections.List()
	if err != nil {
		controller.logger.Error("Failed to list agent collections", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("List agent collections request completed successfully", zap.Int("collections_count", len(collections)))
	controller.logger.Debug("List agent collections response", zap.Any("collections", collections))

	return c.JSON(collections)
}
