package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	_ "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	EXPORT_AGENT_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/export"
)

type ExportAgentController struct {
	logger      *zap.Logger
	exportAgent use_cases.ExportAgentInterface
}

func NewExportAgentController(
	logger *zap.Logger,
	exportAgent use_cases.ExportAgentInterface,
) *ExportAgentController {
	return &ExportAgentController{
		logger:      logger,
		exportAgent: exportAgent,
	}
}

func (controller *ExportAgentController) InitAuthenticatedRoutes(app fiber.Router) {
	app.Post(EXPORT_AGENT_ENDPOINT, controller.ExportAgent)
}

func (controller *ExportAgentController) InitUnauthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

// ExportAgent godoc
// @Summary      Export agent data
// @Description  Exports the agent in a portable format
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.ExportAgentResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/export [post]
func (controller *ExportAgentController) ExportAgent(c *fiber.Ctx) error {
	controller.logger.Info("Received request to export agent")
	controller.logger.Debug("Request", zap.Any("request", c.Body()))

	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	response, err := controller.exportAgent.Export(request)
	if err != nil {
		controller.logger.Error("Failed to export agent", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully exported agent")

	return c.JSON(response)
}
