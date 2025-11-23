package llm_connector_mvc

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	SAVE_LLM_CONNECTION_ENDPOINT             = "/llm/connection"
	REMOVE_LLM_CONNECTION_ENDPOINT           = "/llm/connection"
	GET_LLM_CONNECTION_SETTINGS_ENDPOINT     = "/llm/connection/settings"
	GET_LLM_CONNECTION_SETUP_STATUS_ENDPOINT = "/llm/connection/setup-status"
)

type LLMConnectionController struct {
	logger                      *zap.Logger
	saveLLMConnection           use_cases.SaveLLMConnectionInterface
	removeLLMConnection         use_cases.RemoveLLMConnectionInterface
	getLLMConnectionSettings    use_cases.GetLLMConnectionSettingsInterface
	getLLMConnectionSetupStatus use_cases.GetLLMConnectionSetupStatusInterface
}

func NewLLMConnectionController(
	logger *zap.Logger,
	saveLLMConnection use_cases.SaveLLMConnectionInterface,
	removeLLMConnection use_cases.RemoveLLMConnectionInterface,
	getLLMConnectionSettings use_cases.GetLLMConnectionSettingsInterface,
	getLLMConnectionSetupStatus use_cases.GetLLMConnectionSetupStatusInterface,
) *LLMConnectionController {
	return &LLMConnectionController{
		logger:                      logger,
		saveLLMConnection:           saveLLMConnection,
		removeLLMConnection:         removeLLMConnection,
		getLLMConnectionSettings:    getLLMConnectionSettings,
		getLLMConnectionSetupStatus: getLLMConnectionSetupStatus,
	}
}

func (controller *LLMConnectionController) InitRoutes(app fiber.Router) {
	app.Post(SAVE_LLM_CONNECTION_ENDPOINT, controller.SaveLLMConnection)
	app.Get(GET_LLM_CONNECTION_SETTINGS_ENDPOINT, controller.GetSettings)
	app.Delete(REMOVE_LLM_CONNECTION_ENDPOINT, controller.RemoveLLMConnection)
	app.Get(GET_LLM_CONNECTION_SETUP_STATUS_ENDPOINT, controller.GetSetupStatus)
}

// @Summary Save LLM Connection
// @Description Creates or update a LLM connection based on the agent id
// @Tags LLM Connection
// @Accept json
// @Produce json
// @Param request body requests.SaveLLMConnectionRequest true "Save LLM Connection Request"
// @Success 201 {string} string "Saved"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /llm/connection [post]
func (controller *LLMConnectionController) SaveLLMConnection(c *fiber.Ctx) error {
	controller.logger.Debug("Request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Any("body", c.Body()),
	)

	request := new(requests.SaveLLMConnectionRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.saveLLMConnection.Save(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to create llm connection", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("LLM connection saved successfully")

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Remove LLM Connection
// @Description Removes an existing LLM connection
// @Tags LLM Connection
// @Accept json
// @Produce json
// @Param request body requests.AgentIdRequest true "Remove LLM Connection Request"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /llm/connection [delete]
func (controller *LLMConnectionController) RemoveLLMConnection(c *fiber.Ctx) error {
	controller.logger.Info("Received request to remove LLM connection")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.AgentIdRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.removeLLMConnection.Remove(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to remove llm connection",
			zap.String("agentId", request.AgentId.String()),
			zap.Error(err),
		)

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("LLM connection created successfully")

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Get LLM Connection Settings
// @Description Retrieves settings for an LLM connection
// @Tags LLM Connection
// @Accept json
// @Produce json
// @Param request query requests.AgentIdRequest true "Get LLM Connection Settings Request"
// @Success 200 {object} responses.GetLLMConnectionSettingsResponse "OK"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /llm/connection/settings [get]
func (controller *LLMConnectionController) GetSettings(c *fiber.Ctx) error {
	controller.logger.Info("Received request to get LLM connection metadata")
	controller.logger.Debug("Request", zap.Any("query", c.Queries()))

	request := new(requests.AgentIdRequest)
	if err := c.QueryParser(request); err != nil {
		return err
	}

	metadata, err := controller.getLLMConnectionSettings.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get llm connection metadata",
			zap.Error(err),
			zap.String("agentId", request.AgentId.String()),
		)

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("LLM connection metadata retrieved successfully")
	controller.logger.Debug(
		"Response",
		zap.Any("response", metadata),
	)

	return c.JSON(metadata)
}

// @Summary Get LLM Connection Setup Status
// @Description Retrieves the setup status for an LLM connection
// @Tags LLM Connection
// @Accept json
// @Produce json
// @Param request query requests.AgentIdRequest true "Get LLM Connection Setup Status Request"
// @Success 200 {object} responses.GetLLMConnectionSetupStatusResponse "OK"
// @Failure 400 {object} errors.ValidationError "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /llm/connection/setup-status [get]
func (controller *LLMConnectionController) GetSetupStatus(c *fiber.Ctx) error {
	controller.logger.Info("Received request to get LLM connection setup status")
	controller.logger.Debug("Request", zap.Any("query", c.Queries()))

	request := new(requests.AgentIdRequest)
	if err := c.QueryParser(request); err != nil {
		return err
	}

	response, err := controller.getLLMConnectionSetupStatus.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get llm connection setup status",
			zap.Error(err),
			zap.String("agentId", request.AgentId.String()),
		)

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("LLM connection setup status retrieved successfully")
	controller.logger.Debug(
		"Response",
		zap.Any("response", response),
	)

	return c.JSON(response)
}
