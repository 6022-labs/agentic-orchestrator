package database_connector_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	SAVE_AGENT_PROMPT_TEMPLATE = "/agent/prompt/template"
	GET_AGENT_PROMPT_TEMPLATES = "/agent/prompt/template"
)

type AgentPromptTemplateController struct {
	logger                  *zap.Logger
	getAgentPromptTemplate  use_cases.GetAgentPromptTemplateInterface
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface
}

func NewAgentPromptTemplateController(
	logger *zap.Logger,
	getAgentPromptTemplate use_cases.GetAgentPromptTemplateInterface,
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface,
) *AgentPromptTemplateController {
	return &AgentPromptTemplateController{
		logger:                  logger,
		getAgentPromptTemplate:  getAgentPromptTemplate,
		saveAgentPromptTemplate: saveAgentPromptTemplate,
	}
}

func (controller *AgentPromptTemplateController) InitRoutes(app fiber.Router) {
	app.Post(SAVE_AGENT_PROMPT_TEMPLATE, controller.SaveAgentPromptTemplate)
	app.Get(GET_AGENT_PROMPT_TEMPLATES, controller.GetAgentPromptTemplates)
}

// SaveAgentPromptTemplate saves a new agent prompt template.
// @Summary Save agent prompt template
// @Description Saves a new agent prompt template for the specified agent
// @Tags AgentPromptTemplate
// @Accept json
// @Produce json
// @Param request body requests.SaveAgentPromptTemplateRequest true "Save Agent Prompt Template Request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/prompt/template [post]
func (controller *AgentPromptTemplateController) SaveAgentPromptTemplate(c *fiber.Ctx) error {
	controller.logger.Info("Received request to save agent prompt template")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.SaveAgentPromptTemplateRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.saveAgentPromptTemplate.Save(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to save agent prompt template", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully saved agent prompt template")

	return c.SendStatus(fiber.StatusCreated)
}

// GetAgentPromptTemplates retrieves agent prompt templates.
// @Summary Get agent prompt templates
// @Description Retrieves agent prompt templates for the specified agent
// @Tags AgentPromptTemplate
// @Accept json
// @Produce json
// @Param agentId query string true "Agent ID"
// @Param templateKey query string false "Template Key"
// @Success 200 {object} interface{}
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agent/prompt/template [get]
func (controller *AgentPromptTemplateController) GetAgentPromptTemplates(c *fiber.Ctx) error {
	controller.logger.Info("Received request to get agent prompt templates")
	controller.logger.Debug("Request query", zap.Any("query", c.Queries()))

	request := new(requests.GetAgentPromptTemplateRequest)
	if err := c.QueryParser(request); err != nil {
		return err
	}

	response, err := controller.getAgentPromptTemplate.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get agent prompt templates", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully retrieved agent prompt templates")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}
