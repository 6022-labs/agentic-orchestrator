package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	GET_AGENT_PROMPT_TEMPLATE_ENDPOINT  = "/agents/:collectionAddress/:collectionTokenId/prompt-templates/:key"
	SAVE_AGENT_PROMPT_TEMPLATE_ENDPOINT = "/agents/:collectionAddress/:collectionTokenId/prompt-templates/:key"
)

type PromptTemplatesController struct {
	logger                  *zap.Logger
	getAgentPromptTemplate  use_cases.GetAgentPromptTemplateInterface
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface
}

func NewPromptTemplatesController(
	logger *zap.Logger,
	getAgentPromptTemplate use_cases.GetAgentPromptTemplateInterface,
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface,
) *PromptTemplatesController {
	return &PromptTemplatesController{
		logger:                  logger,
		getAgentPromptTemplate:  getAgentPromptTemplate,
		saveAgentPromptTemplate: saveAgentPromptTemplate,
	}
}

func (c *PromptTemplatesController) InitRoutes(app fiber.Router) {
	app.Get(GET_AGENT_PROMPT_TEMPLATE_ENDPOINT, c.Get)
	app.Post(SAVE_AGENT_PROMPT_TEMPLATE_ENDPOINT, c.Save)
}

// @Summary Get Agent Prompt Template
// @Description Returns the agent-scoped prompt template for the given key.
// @Tags Prompt Templates
// @Produce json
// @Param collectionAddress path string true "Agent collection address"
// @Param collectionTokenId path string true "Agent collection token id"
// @Param key path string true "Prompt template key"
// @Success 200 {object} responses.AgentPromptTemplateResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agents/{collectionAddress}/{collectionTokenId}/prompt-templates/{key} [get]
func (c *PromptTemplatesController) Get(ctx fiber.Ctx) error {
	req := &requests.GetAgentPromptTemplateRequest{
		CollectionAddress: ctx.Params("collectionAddress"),
		CollectionTokenId: ctx.Params("collectionTokenId"),
		Key:               ctx.Params("key"),
	}
	template, err := c.getAgentPromptTemplate.Get(ctx.Context(), req)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.JSON(template)
}

// @Summary Save Agent Prompt Template
// @Description Creates or overwrites the agent-scoped prompt template for the given key.
// @Tags Prompt Templates
// @Accept json
// @Produce json
// @Param collectionAddress path string true "Agent collection address"
// @Param collectionTokenId path string true "Agent collection token id"
// @Param key path string true "Prompt template key"
// @Param request body requests.SaveAgentPromptTemplateRequest true "Save Agent Prompt Template Request"
// @Success 200 {object} responses.AgentPromptTemplateResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /agents/{collectionAddress}/{collectionTokenId}/prompt-templates/{key} [post]
func (c *PromptTemplatesController) Save(ctx fiber.Ctx) error {
	req := new(requests.SaveAgentPromptTemplateRequest)
	if err := ctx.Bind().Body(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req.CollectionAddress = ctx.Params("collectionAddress")
	req.CollectionTokenId = ctx.Params("collectionTokenId")
	req.Key = ctx.Params("key")

	template, err := c.saveAgentPromptTemplate.Save(ctx.Context(), req)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.JSON(template)
}
