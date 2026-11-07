package conversation_broker_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v3"
)

const (
	GET_AGENT_PROMPT_TEMPLATE_ENDPOINT  = "/agents/:collectionAddress/:collectionTokenId/prompt-templates/:key"
	SAVE_AGENT_PROMPT_TEMPLATE_ENDPOINT = "/agents/:collectionAddress/:collectionTokenId/prompt-templates/:key"
)

type PromptTemplatesController struct {
	errorWriter             *mvc.ErrorWriter
	getAgentPromptTemplate  use_cases.GetAgentPromptTemplateInterface
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface
}

func NewPromptTemplatesController(
	errorWriter *mvc.ErrorWriter,
	getAgentPromptTemplate use_cases.GetAgentPromptTemplateInterface,
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface,
) *PromptTemplatesController {
	return &PromptTemplatesController{
		errorWriter:             errorWriter,
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
// @Failure 400 {object} errors.BadRequestError "invalid_params"
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {object} errors.InternalError "internal_error"
// @Router /agents/{collectionAddress}/{collectionTokenId}/prompt-templates/{key} [get]
func (c *PromptTemplatesController) Get(ctx fiber.Ctx) error {
	req := new(requests.GetAgentPromptTemplateRequest)
	if err := ctx.Bind().URI(req); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}
	template, err := c.getAgentPromptTemplate.Get(ctx.Context(), req)
	if err != nil {
		return c.errorWriter.Write(ctx, err)
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
// @Failure 400 {object} errors.BadRequestError "invalid_body or invalid_params"
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {object} errors.InternalError "internal_error"
// @Router /agents/{collectionAddress}/{collectionTokenId}/prompt-templates/{key} [post]
func (c *PromptTemplatesController) Save(ctx fiber.Ctx) error {
	req := new(requests.SaveAgentPromptTemplateRequest)
	if err := ctx.Bind().Body(req); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_body", err.Error()))
	}
	if err := ctx.Bind().URI(req); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	template, err := c.saveAgentPromptTemplate.Save(ctx.Context(), req)
	if err != nil {
		return c.errorWriter.Write(ctx, err)
	}
	return ctx.JSON(template)
}
