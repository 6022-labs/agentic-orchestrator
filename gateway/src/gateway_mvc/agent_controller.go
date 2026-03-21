package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/zapext"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	LIST_AGENTS_ENDPOINT              = "/agents"
	CREATE_AGENT_FIRST_STEP_ENDPOINT  = "/agent/create/first-step"
	CREATE_AGENT_SECOND_STEP_ENDPOINT = "/agent/create/second-step"
	GET_AGENT_ENDPOINT                = "/agent/:agentCollectionAddress/:agentCollectionTokenId"
	AGENT_IS_INTERNAL_ENDPOINT        = "/agent/:agentCollectionAddress/:agentCollectionTokenId/is-internal"

	AGENT_COMPUTE_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/compute"

	GET_AGENT_PROMPT_TEMPLATE_ENDPOINT  = "/agent/:agentCollectionAddress/:agentCollectionTokenId/prompt-template/:key"
	SAVE_AGENT_PROMPT_TEMPLATE_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/prompt-template/:key"

	CREATE_AGENT_CONTEXT_ENDPOINT         = "/agent/:agentCollectionAddress/:agentCollectionTokenId/context"
	LIST_AGENT_CONTEXT_KEYS_ENDPOINT      = "/agent/:agentCollectionAddress/:agentCollectionTokenId/contexts/keys"
	ADD_MESSAGE_TO_AGENT_CONTEXT_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/context/:contextKey/message"

	SETUP_AGENT_DATABASE_SETTINGS_FIRST_STEP_ENDPOINT  = "/agent/:agentCollectionAddress/:agentCollectionTokenId/database/settings/first-step"
	SETUP_AGENT_DATABASE_SETTINGS_SECOND_STEP_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/database/settings/second-step"

	SETUP_AGENT_LLM_SETTINGS_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/llm/settings"
	GET_AGENT_LLM_SETTINGS_ENDPOINT   = "/agent/:agentCollectionAddress/:agentCollectionTokenId/llm/settings"

	LIST_AGENT_MCP_CONNECTIONS_ENDPOINT  = "/agent/:agentCollectionAddress/:agentCollectionTokenId/mcp/connections"
	CREATE_AGENT_MCP_CONNECTION_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/mcp/connection"
	REMOVE_AGENT_MCP_CONNECTION_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/mcp/connection/:mcpConnectionId"
)

type AgentController struct {
	logger                               *zap.Logger
	getAgent                             use_cases.GetAgentInterface
	listAgents                           use_cases.ListAgentsInterface
	agentCompute                         use_cases.AgentComputeInterface
	isInternalAgent                      use_cases.IsInternalAgentInterface
	createAgentContext                   use_cases.CreateAgentContextInterface
	getAgentLLMSettings                  use_cases.GetAgentLLMSettingsInterface
	listAgentContextKeys                 use_cases.ListAgentContextKeysInterface
	createAgentFirstStep                 use_cases.CreateAgentFirstStepInterface
	createAgentSecondStep                use_cases.CreateAgentSecondStepInterface
	setupAgentLLMSettings                use_cases.SetupAgentLLMSettingsInterface
	getAgentPromptTemplate               use_cases.GetAgentPromptTemplateInterface
	listAgentMCPConnections              use_cases.ListAgentMCPConnectionsInterface
	saveAgentPromptTemplate              use_cases.SaveAgentPromptTemplateInterface
	addMessageToAgentContext             use_cases.AddMessageToAgentContextInterface
	removeAgentMCPConnection             use_cases.RemoveAgentMCPConnectionInterface
	createAgentMCPConnection             use_cases.CreateAgentMCPConnectionInterface
	setupAgentDatabaseSettingsFirstStep  use_cases.SetupAgentDatabaseSettingsFirstStepInterface
	setupAgentDatabaseSettingsSecondStep use_cases.SetupAgentDatabaseSettingsSecondStepInterface
}

func NewAgentController(
	logger *zap.Logger,
	getAgent use_cases.GetAgentInterface,
	listAgents use_cases.ListAgentsInterface,
	agentCompute use_cases.AgentComputeInterface,
	isInternalAgent use_cases.IsInternalAgentInterface,
	createAgentContext use_cases.CreateAgentContextInterface,
	getAgentLLMSettings use_cases.GetAgentLLMSettingsInterface,
	listAgentContextKeys use_cases.ListAgentContextKeysInterface,
	createAgentFirstStep use_cases.CreateAgentFirstStepInterface,
	createAgentSecondStep use_cases.CreateAgentSecondStepInterface,
	setupAgentLLMSettings use_cases.SetupAgentLLMSettingsInterface,
	getAgentPromptTemplate use_cases.GetAgentPromptTemplateInterface,
	listAgentMCPConnections use_cases.ListAgentMCPConnectionsInterface,
	saveAgentPromptTemplate use_cases.SaveAgentPromptTemplateInterface,
	addMessageToAgentContext use_cases.AddMessageToAgentContextInterface,
	removeAgentMCPConnection use_cases.RemoveAgentMCPConnectionInterface,
	createAgentMCPConnection use_cases.CreateAgentMCPConnectionInterface,
	setupAgentDatabaseSettingsFirstStep use_cases.SetupAgentDatabaseSettingsFirstStepInterface,
	setupAgentDatabaseSettingsSecondStep use_cases.SetupAgentDatabaseSettingsSecondStepInterface,
) *AgentController {
	return &AgentController{
		logger:                               logger,
		getAgent:                             getAgent,
		listAgents:                           listAgents,
		agentCompute:                         agentCompute,
		isInternalAgent:                      isInternalAgent,
		createAgentContext:                   createAgentContext,
		getAgentLLMSettings:                  getAgentLLMSettings,
		listAgentContextKeys:                 listAgentContextKeys,
		createAgentFirstStep:                 createAgentFirstStep,
		createAgentSecondStep:                createAgentSecondStep,
		setupAgentLLMSettings:                setupAgentLLMSettings,
		getAgentPromptTemplate:               getAgentPromptTemplate,
		listAgentMCPConnections:              listAgentMCPConnections,
		saveAgentPromptTemplate:              saveAgentPromptTemplate,
		addMessageToAgentContext:             addMessageToAgentContext,
		removeAgentMCPConnection:             removeAgentMCPConnection,
		createAgentMCPConnection:             createAgentMCPConnection,
		setupAgentDatabaseSettingsFirstStep:  setupAgentDatabaseSettingsFirstStep,
		setupAgentDatabaseSettingsSecondStep: setupAgentDatabaseSettingsSecondStep,
	}
}

func (controller *AgentController) InitAuthenticatedRoutes(app fiber.Router) {
	app.Post(CREATE_AGENT_FIRST_STEP_ENDPOINT, controller.CreateAgentFirstStep)
	app.Post(CREATE_AGENT_SECOND_STEP_ENDPOINT, controller.CreateAgentSecondStep)

	app.Get(AGENT_IS_INTERNAL_ENDPOINT, controller.IsInternalAgent)
	app.Post(AGENT_COMPUTE_ENDPOINT, controller.Compute)

	app.Get(GET_AGENT_PROMPT_TEMPLATE_ENDPOINT, controller.GetAgentPromptTemplate)
	app.Post(SAVE_AGENT_PROMPT_TEMPLATE_ENDPOINT, controller.SaveAgentPromptTemplate)

	app.Post(CREATE_AGENT_CONTEXT_ENDPOINT, controller.CreateAgentContext)
	app.Get(LIST_AGENT_CONTEXT_KEYS_ENDPOINT, controller.ListAgentContextKeys)
	app.Post(ADD_MESSAGE_TO_AGENT_CONTEXT_ENDPOINT, controller.AddMessageToAgentContext)

	app.Post(SETUP_AGENT_LLM_SETTINGS_ENDPOINT, controller.SetupAgentLLMSettings)
	app.Get(GET_AGENT_LLM_SETTINGS_ENDPOINT, controller.GetAgentLLMSettings)

	app.Post(SETUP_AGENT_DATABASE_SETTINGS_FIRST_STEP_ENDPOINT, controller.SetupAgentDatabaseSettingsFirstStep)
	app.Post(SETUP_AGENT_DATABASE_SETTINGS_SECOND_STEP_ENDPOINT, controller.SetupAgentDatabaseSettingsSecondStep)

	app.Get(LIST_AGENT_MCP_CONNECTIONS_ENDPOINT, controller.ListAgentMCPConnections)
	app.Post(CREATE_AGENT_MCP_CONNECTION_ENDPOINT, controller.CreateAgentMCPConnection)
	app.Delete(REMOVE_AGENT_MCP_CONNECTION_ENDPOINT, controller.RemoveAgentMCPConnection)
}

func (controller *AgentController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Get(GET_AGENT_ENDPOINT, controller.GetAgent)
	app.Get(LIST_AGENTS_ENDPOINT, controller.ListAgents)
}

// GetAgent godoc
// @Summary      Get agent information
// @Description  Retrieves details of a specific agent
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.AgentResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId} [get]
func (controller *AgentController) GetAgent(c *fiber.Ctx) error {
	controller.logger.Info("Get agent request received")
	controller.logger.Debug("Request", zap.Any("params", c.AllParams()))

	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.getAgent.Get(request)
	if err != nil {
		controller.logger.Error("Failed to get agent", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Get agent request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// ListAgents godoc
// @Summary      List agents
// @Description  Returns all agents managed by the orchestrator
// @Tags         agent
// @Produce      json
// @Param        request query requests.ListAgentsRequest true "list agents request"
// @Success      200 {object} responses.ListAgentsResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agents [get]
func (controller *AgentController) ListAgents(c *fiber.Ctx) error {
	controller.logger.Info("List agents request received")
	controller.logger.Debug("Request", zap.Any("queries", c.Queries()))

	request := new(requests.ListAgentsRequest)
	if err := c.QueryParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	resp, err := controller.listAgents.List(request)
	if err != nil {
		controller.logger.Error("Failed to list agents", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("List agents request completed successfully", zap.Int("agents_count", len(resp.Agents)))
	controller.logger.Debug(
		"List agents response",
		zap.Any("page", resp.Page),
		zap.Any("totalElements", resp.TotalElements),
		zap.Any("elementsPerPage", resp.ElementsPerPage),
		zapext.ArrayTruncateEnd("agents", resp.Agents, 20),
	)

	return c.JSON(resp)
}

// AgentCreationFirstStep godoc
// @Summary      Start agent creation
// @Description  Upload agent image and create agent metadata
// @Tags         agent
// @Accept       multipart/form-data
// @Produce      json
// @Param        request formData requests.CreateAgentFirstStepRequest true "agent creation request"
// @Param 	  	 image formData file true "agent image"
// @Success      200 {object} responses.AgentResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/create/first-step [post]
func (controller *AgentController) CreateAgentFirstStep(c *fiber.Ctx) error {
	controller.logger.Info("Agent creation first step request received")
	controller.logger.Debug("Request", zap.Any("body", c.Body()))

	request := new(requests.CreateAgentFirstStepRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	image, err := c.FormFile("image")
	if err != nil {
		controller.logger.Error("Failed to parse image file", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse image file")
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.Image = image
	request.AuthenticatedWallet = authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.createAgentFirstStep.Execute(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Agent creation first step request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// AgentCreationSecondStep godoc
// @Summary      Finish agent creation
// @Description  Completes the agent creation process
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        request body requests.AgentByIdRequest true "agent id"
// @Success      200 {object} responses.CreateAgentSecondStepResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/create/second-step [post]
func (controller *AgentController) CreateAgentSecondStep(c *fiber.Ctx) error {
	controller.logger.Info("Agent creation second step request received")
	controller.logger.Debug("Request",
		zap.Any("body", c.Body()),
		zap.Any("headers", c.GetReqHeaders()),
	)

	request := new(requests.AgentByIdRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.createAgentSecondStep.Execute(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Agent creation second step request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// IsInternalAgent godoc
// @Summary      Check if agent is internal
// @Description  Returns true if the agent is managed internally
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.IsInternalAgentResponse
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/is-internal [get]
func (controller *AgentController) IsInternalAgent(c *fiber.Ctx) error {
	controller.logger.Info("Received request to check if agent is internal")
	controller.logger.Debug("Request",
		zap.Any("queries", c.Queries()),
		zap.Any("params", c.AllParams()),
		zap.Any("headers", c.GetReqHeaders()),
	)

	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.isInternalAgent.IsInternal(request)
	if err != nil {
		controller.logger.Error("Failed to check if agent is internal", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully checked if agent is internal")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// Compute godoc
// @Summary      Compute agent response
// @Description  Runs the agent using the provided input
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        request body requests.AgentComputeRequest true "compute request"
// @Success      200 {object} responses.AgentComputeResponses
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/compute [post]
func (controller *AgentController) Compute(c *fiber.Ctx) error {
	controller.logger.Info("Received request to compute agent")
	controller.logger.Debug("Request",
		zap.Any("body", c.Body()),
		zap.Any("params", c.AllParams()),
		zap.Any("queries", c.Queries()),
		zap.Any("headers", c.GetReqHeaders()),
	)

	request := new(requests.AgentComputeRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.agentCompute.Compute(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to compute agent responses", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully computed agent")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// GetAgentPromptTemplate godoc
// @Summary      Get agent prompt template
// @Description  Retrieves a specific prompt template for an agent
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        key path string true "template key"
// @Success      200 {string} string
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/prompt-template/{key} [get]
func (controller *AgentController) GetAgentPromptTemplate(c *fiber.Ctx) error {
	controller.logger.Info("Get agent prompt template request received")
	controller.logger.Debug("Request",
		zap.Any("queries", c.Queries()),
		zap.Any("path", c.Path()),
	)

	request := new(requests.GetAgentPromptTemplateRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	promptTemplate, err := controller.getAgentPromptTemplate.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get agent prompt template", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Get agent prompt template request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", promptTemplate))

	return c.JSON(promptTemplate)
}

// SaveAgentPromptTemplate godoc
// @Summary      Save agent prompt template
// @Description  Stores a prompt template for a given agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        key path string true "template key"
// @Param        request body requests.SaveAgentPromptTemplateRequest true "template"
// @Success      200 {string} string "OK"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/prompt-template/{key} [post]
func (controller *AgentController) SaveAgentPromptTemplate(c *fiber.Ctx) error {
	controller.logger.Info("Save agent prompt template request received")
	controller.logger.Debug("Request",
		zap.Any("body", c.Body()),
		zap.Any("path", c.Path()),
	)

	request := new(requests.SaveAgentPromptTemplateRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.saveAgentPromptTemplate.Save(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to save agent prompt template", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Save agent prompt template request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}

// GetAgentContextKeys godoc
// @Summary      List agent context keys
// @Description  Returns all context keys for an agent
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        request query requests.ListAgentContextKeysRequest true "context"
// @Success      200 {array} string
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/contexts/keys [get]
func (controller *AgentController) ListAgentContextKeys(c *fiber.Ctx) error {
	controller.logger.Info("List agent contexts request received")
	controller.logger.Debug(
		"Request",
		zap.Any("params", c.AllParams()),
		zap.Any("queries", c.Queries()),
	)

	request := new(requests.ListAgentContextKeysRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}
	if err := c.QueryParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	resp, err := controller.listAgentContextKeys.List(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get agent contexts", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Get agent contexts request completed successfully")
	controller.logger.Debug(
		"List agent contexts response",
		zap.Any("page", resp.Page),
		zap.Any("totalElements", resp.TotalElements),
		zap.Any("elementsPerPage", resp.ElementsPerPage),
		zapext.ArrayTruncateEnd("contextKeys", resp.ContextKeys, 20),
	)

	return c.JSON(resp)
}

// CreateAgentContext godoc
// @Summary      Create agent context
// @Description  Creates a new context for an agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        request body requests.CreateAgentContextRequest true "context"
// @Success      200 {string} string "OK"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/context [post]
func (controller *AgentController) CreateAgentContext(c *fiber.Ctx) error {
	controller.logger.Info("Create agent context request received")
	controller.logger.Debug(
		"Request",
		zap.Any("body", c.Body()),
		zap.Any("path", c.Path()),
	)

	request := new(requests.CreateAgentContextRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.createAgentContext.Create(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to create agent context", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Create agent context request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}

// AddMessageToAgentContext godoc
// @Summary      Add message to agent context
// @Description  Appends a message to a specific agent context
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        contextKey path string true "context key"
// @Param        request body requests.AddMessageToAgentContextRequest true "message"
// @Success      200 {string} string "OK"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/context/{contextKey}/message [post]
func (controller *AgentController) AddMessageToAgentContext(c *fiber.Ctx) error {
	controller.logger.Info("Add message to agent context request received")
	controller.logger.Debug(
		"Request",
		zap.Any("body", c.Body()),
		zap.Any("path", c.Path()),
	)

	request := new(requests.AddMessageToAgentContextRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.addMessageToAgentContext.Add(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to add message to agent context", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Add message to agent context request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}

// SetupAgentLLMSettings godoc
// @Summary      Configure agent LLM settings
// @Description  Sets up the LLM parameters for the agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        request body requests.SetupAgentLLMSettingsRequest true "settings"
// @Success      200 {string} string "OK"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/llm/settings [post]
func (controller *AgentController) SetupAgentLLMSettings(c *fiber.Ctx) error {
	controller.logger.Info("Setup agent LLM settings request received")
	controller.logger.Debug(
		"Request",
		zap.Any("body", c.Body()),
		zap.Any("path", c.Path()),
	)

	request := new(requests.SetupAgentLLMSettingsRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.setupAgentLLMSettings.Setup(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to setup agent LLM settings", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Setup agent LLM settings request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}

// GetAgentLLMSettings godoc
// @Summary      Get agent LLM settings
// @Description  Retrieves the LLM configuration for an agent
// @Tags         agent
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.GetAgentLLMSettingsResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/llm/settings [get]
func (controller *AgentController) GetAgentLLMSettings(c *fiber.Ctx) error {
	controller.logger.Info("Get agent LLM settings request received")
	controller.logger.Debug(
		"Request",
		zap.Any("path", c.Path()),
	)

	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	settings, err := controller.getAgentLLMSettings.Get(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to get agent LLM settings", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Get agent LLM settings request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", settings))

	return c.JSON(settings)
}

// SetupAgentDatabaseSettingsFirstStep godoc
// @Summary      Setup agent database first step
// @Description  Configures database credentials for the agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        request body requests.SetupAgentDatabaseSettingsFirstStepRequest true "settings"
// @Success      200 {object} responses.SetupAgentDatabaseSettingsFirstStepResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/database/settings/first-step [post]
func (controller *AgentController) SetupAgentDatabaseSettingsFirstStep(c *fiber.Ctx) error {
	controller.logger.Info("Setup agent database settings first step request received")
	controller.logger.Debug(
		"Request",
		zap.Any("body", c.Body()),
		zap.Any("path", c.Path()),
	)

	request := new(requests.SetupAgentDatabaseSettingsFirstStepRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	result, err := controller.setupAgentDatabaseSettingsFirstStep.Setup(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to setup agent database settings first step", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Setup agent database settings first step request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", result))

	return c.JSON(result)
}

// SetupAgentDatabaseSettingsSecondStep godoc
// @Summary      Finalize agent database setup
// @Description  Applies database settings to the agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        request body requests.AgentByCollectionTokenIdRequest true "agent"
// @Success      200 {string} string "OK"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/database/settings/second-step [post]
func (controller *AgentController) SetupAgentDatabaseSettingsSecondStep(c *fiber.Ctx) error {
	controller.logger.Info("Setup agent database settings second step request received")
	controller.logger.Debug(
		"Request",
		zap.Any("body", c.Body()),
		zap.Any("path", c.Path()),
	)

	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.setupAgentDatabaseSettingsSecondStep.Setup(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to setup agent database settings second step", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Setup agent database settings second step request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}

// ListAgentMCPConnections godoc
// @Summary      List agent MCP connections
// @Description  Retrieves the MCP connections for a specific agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {array} responses.MCPConnectionResponse "OK"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/mcp/connections [get]
func (controller *AgentController) ListAgentMCPConnections(c *fiber.Ctx) error {
	controller.logger.Info("List agent MCP connections request received")
	controller.logger.Debug("Request", zap.Any("path", c.Path()))

	request := new(requests.ListAgentMCPConnectionsRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse params from request", zap.Error(err))
		return err
	}
	if err := c.QueryParser(request); err != nil {
		controller.logger.Error("Failed to parse query from request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	connections, err := controller.listAgentMCPConnections.List(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to list agent MCP connections", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("List agent MCP connections request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", connections))

	return c.JSON(connections)
}

// RemoveAgentMCPConnection godoc
// @Summary      Remove agent MCP connection
// @Description  Removes a specific MCP connection for an agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        mcpConnectionId path string true "mcp connection id"
// @Success      200 {string} string "No Content"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/mcp/connection/{mcpConnectionId} [delete]
func (controller *AgentController) RemoveAgentMCPConnection(c *fiber.Ctx) error {
	controller.logger.Info("Remove agent MCP connection request received")
	controller.logger.Debug("Request", zap.Any("path", c.Path()))

	request := new(requests.RemoveAgentMCPConnectionRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.removeAgentMCPConnection.Remove(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to remove agent MCP connection", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Remove agent MCP connection request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}

// CreateAgentMCPConnection godoc
// @Summary      Create agent MCP connection
// @Description  Creates a new MCP connection for an agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {string} string "No Content"
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/mcp/connection [post]
func (controller *AgentController) CreateAgentMCPConnection(c *fiber.Ctx) error {
	controller.logger.Info("Create agent MCP connection request received")
	controller.logger.Debug("Request", zap.Any("body", c.Body()), zap.Any("path", c.Path()))

	request := new(requests.CreateAgentMCPConnectionRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request parameters", zap.Error(err))
		return err
	}

	authenticatedWallet, ok := c.Locals("authenticated-wallet").(common.Address)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	request.AuthenticatedWallet = &authenticatedWallet

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	err := controller.createAgentMCPConnection.Create(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to create agent MCP connection", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Create agent MCP connection request completed successfully")

	return c.SendStatus(fiber.StatusOK)
}
