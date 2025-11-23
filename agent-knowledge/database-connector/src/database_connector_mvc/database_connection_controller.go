package database_connector_mvc

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	CREATE_DATABASE_CONNECTION_ENDPOINT           = "/database/connection"
	REMOVE_DATABASE_CONNECTION_ENDPOINT           = "/database/connection"
	SETUP_AGENT_DATABASE_ENDPOINT                 = "/database/connection/setup"
	GET_DATABASE_CONNECTION_SETUP_STATUS_ENDPOINT = "/database/connection/setup-status"
)

type DatabaseConnectionController struct {
	logger                           *zap.Logger
	setupAgentDatabase               use_cases.SetupAgentDatabaseInterface
	createDatabaseConnection         use_cases.CreateDatabaseConnectionInterface
	removeDatabaseConnection         use_cases.RemoveDatabaseConnectionInterface
	getDatabaseConnectionSetupStatus use_cases.GetDatabaseConnectionSetupStatusInterface
}

func NewDatabaseConnectionController(
	logger *zap.Logger,
	setupAgentDatabase use_cases.SetupAgentDatabaseInterface,
	createDatabaseConnection use_cases.CreateDatabaseConnectionInterface,
	removeDatabaseConnection use_cases.RemoveDatabaseConnectionInterface,
	getDatabaseConnectionSetupStatus use_cases.GetDatabaseConnectionSetupStatusInterface,
) *DatabaseConnectionController {
	return &DatabaseConnectionController{
		logger:                           logger,
		setupAgentDatabase:               setupAgentDatabase,
		createDatabaseConnection:         createDatabaseConnection,
		removeDatabaseConnection:         removeDatabaseConnection,
		getDatabaseConnectionSetupStatus: getDatabaseConnectionSetupStatus,
	}
}

func (controller *DatabaseConnectionController) InitRoutes(app fiber.Router) {
	app.Post(SETUP_AGENT_DATABASE_ENDPOINT, controller.SetupAgentDatabase)
	app.Post(CREATE_DATABASE_CONNECTION_ENDPOINT, controller.CreateDatabaseConnection)
	app.Delete(REMOVE_DATABASE_CONNECTION_ENDPOINT, controller.RemoveDatabaseConnection)
	app.Get(GET_DATABASE_CONNECTION_SETUP_STATUS_ENDPOINT, controller.GetDatabaseConnectionSetupStatus)
}

// CreateDatabaseConnection handles the creation of a new database connection.
// @Summary Create a new database connection
// @Description Creates a new database connection for the specified agent
// @Tags DatabaseConnection
// @Accept json
// @Produce json
// @Param request body requests.CreateDatabaseConnectionRequest true "Create Database Connection Request"
// @Success 201 {object} responses.CreateDatabaseConnectionResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /database/connection [post]
func (controller *DatabaseConnectionController) CreateDatabaseConnection(c *fiber.Ctx) error {
	controller.logger.Info("Received request to create database connection")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.CreateDatabaseConnectionRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.createDatabaseConnection.Create(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("LLM connection created successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.Status(fiber.StatusCreated).JSON(response)
}

// RemoveDatabaseConnection handles the removal of an existing database connection.
// @Summary Remove a database connection
// @Description Deletes an existing database connection for the specified agent
// @Tags DatabaseConnection
// @Accept json
// @Produce json
// @Param request body requests.AgentIdRequest true "Remove Database Connection Request"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /database/connection [delete]
func (controller *DatabaseConnectionController) RemoveDatabaseConnection(c *fiber.Ctx) error {
	controller.logger.Info("Received request to delete database connection")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.AgentIdRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.removeDatabaseConnection.Remove(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("LLM connection deleted successfully")

	return c.SendStatus(fiber.StatusNoContent)
}

// SetupAgentDatabase handles the setup of the agent database.
// @Summary Setup agent database
// @Description Sets up the database for the specified agent
// @Tags DatabaseConnection
// @Accept json
// @Produce json
// @Param request body requests.SetupAgentDatabaseRequest true "Setup Agent Database Request"
// @Success 200 {string} string "OK"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /database/connection/setup [post]
func (controller *DatabaseConnectionController) SetupAgentDatabase(c *fiber.Ctx) error {
	controller.logger.Info("Setting up agent database")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.SetupAgentDatabaseRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.setupAgentDatabase.Setup(request)

	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to setup agent database", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Agent database setup completed successfully")

	return c.SendStatus(fiber.StatusOK)
}

// GetDatabaseConnectionSetupStatus handles retrieving the database connection setup status.
// @Summary Get database connection setup status
// @Description Retrieves the setup status of the database connection for the specified agent
// @Tags DatabaseConnection
// @Accept json
// @Produce json
// @Param request query requests.AgentIdRequest true "Get Database Connection Setup Status Request"
// @Success 200 {object} responses.GetDatabaseConnectionSetupStatusResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /database/connection/setup-status [get]
func (controller *DatabaseConnectionController) GetDatabaseConnectionSetupStatus(c *fiber.Ctx) error {
	controller.logger.Info("Getting database connection setup status")
	controller.logger.Debug("Request", zap.Any("query", c.Queries()))

	request := new(requests.AgentIdRequest)
	if err := c.QueryParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.getDatabaseConnectionSetupStatus.Get(request)

	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get database connection setup status", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Database connection setup status retrieved successfully")

	return c.JSON(response)
}
