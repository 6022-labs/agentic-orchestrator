package agent_compute_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	CREATE_MCP_CONNECTION_ENDPOINT = "/mcp/connection"
	REMOVE_MCP_CONNECTION_ENDPOINT = "/mcp/connection"
	LIST_MCP_CONNECTIONS_ENDPOINT  = "/mcp/connections"
)

type MCPConnectionController struct {
	logger              *zap.Logger
	listMCPConnections  use_cases.ListMCPConnectionsInterface
	createMCPConnection use_cases.CreateMCPConnectionInterface
	removeMCPConnection use_cases.RemoveMCPConnectionInterface
}

func NewMCPConnectionController(
	logger *zap.Logger,
	listMCPConnections use_cases.ListMCPConnectionsInterface,
	createMCPConnection use_cases.CreateMCPConnectionInterface,
	removeMCPConnection use_cases.RemoveMCPConnectionInterface,
) *MCPConnectionController {
	return &MCPConnectionController{
		logger:              logger,
		listMCPConnections:  listMCPConnections,
		createMCPConnection: createMCPConnection,
		removeMCPConnection: removeMCPConnection,
	}
}

func (controller *MCPConnectionController) InitRoutes(app fiber.Router) {
	app.Get(LIST_MCP_CONNECTIONS_ENDPOINT, controller.ListMCPConnections)
	app.Post(CREATE_MCP_CONNECTION_ENDPOINT, controller.CreateMCPConnection)
	app.Delete(REMOVE_MCP_CONNECTION_ENDPOINT, controller.RemoveMCPConnection)
}

// @Summary List MCP Connections
// @Description Retrieves a list of MCP connections for the given agent
// @Tags MCP Connection
// @Accept json
// @Produce json
// @Param request query requests.ListMCPConnectionsRequest true "List MCP Connections Request"
// @Success 200 {object} responses.ListMCPConnectionsResponse "OK"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /mcp/connections [get]
func (controller *MCPConnectionController) ListMCPConnections(c *fiber.Ctx) error {
	controller.logger.Debug("Request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Any("query", c.Queries()),
	)

	request := new(requests.ListMCPConnectionsRequest)
	if err := c.QueryParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.listMCPConnections.List(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to list mcp connections", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("MCP connections listed successfully", zap.Any("response", response))

	return c.JSON(response)
}

// @Summary Create MCP Connection
// @Description Creates a new MCP connection
// @Tags MCP Connection
// @Accept json
// @Produce json
// @Param request body requests.CreateMCPConnectionRequest true "Create MCP Connection Request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /mcp/connection [post]
func (controller *MCPConnectionController) CreateMCPConnection(c *fiber.Ctx) error {
	controller.logger.Debug("Request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Any("query", c.Queries()),
	)

	request := new(requests.CreateMCPConnectionRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.createMCPConnection.Create(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to create mcp connection", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("MCP connection created successfully")

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Remove MCP Connection
// @Description Removes an existing MCP connection
// @Tags MCP Connection
// @Accept json
// @Produce json
// @Param request body requests.IdRequest true "Remove MCP Connection Request"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /mcp/connection [delete]
func (controller *MCPConnectionController) RemoveMCPConnection(c *fiber.Ctx) error {
	controller.logger.Debug("Request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Any("query", c.Queries()),
	)

	request := new(requests.IdRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.removeMCPConnection.Remove(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to remove mcp connection", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("MCP connection removed successfully")

	return c.SendStatus(fiber.StatusNoContent)
}
