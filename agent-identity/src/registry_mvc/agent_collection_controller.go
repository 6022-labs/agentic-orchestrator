package registry_mvc

import (
	_ "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/gofiber/fiber/v3"
)

const (
	LIST_AGENT_COLLECTIONS_ENDPOINT = "/agent-collections"
)

type AgentCollectionController struct {
	errorWriter          *mvc.ErrorWriter
	listAgentCollections use_cases.ListAgentCollectionsInterface
}

func NewAgentCollectionController(
	errorWriter *mvc.ErrorWriter,
	listAgentCollections use_cases.ListAgentCollectionsInterface,
) *AgentCollectionController {
	return &AgentCollectionController{
		errorWriter:          errorWriter,
		listAgentCollections: listAgentCollections,
	}
}

func (controller *AgentCollectionController) RegisterRoutes(app fiber.Router) {
	app.Get(LIST_AGENT_COLLECTIONS_ENDPOINT, controller.ListAgentCollections)
}

// ListAgentCollections godoc
// @Summary      List agent collections
// @Description  Retrieves all available agent collections
// @Tags         AgentCollection
// @Produce      json
// @Success      200 {array} responses.AgentCollectionResponse
// @Failure      500 {object} errors.InternalError "internal_error"
// @Router       /agent-collections [get]
func (controller *AgentCollectionController) ListAgentCollections(c fiber.Ctx) error {
	collections, err := controller.listAgentCollections.List(c.Context())
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(collections)
}
