package registry_mvc

import (
	registry_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	LIST_AGENTS_ENDPOINT             = "/agents"
	GET_AGENT_ENDPOINT               = "/agent/:agentCollectionAddress/:agentCollectionTokenId"
	GET_AGENT_BY_ENS_DOMAIN_ENDPOINT = "/agent/ens/:ensDomain"
	AGENT_IS_INTERNAL_ENDPOINT       = "/agent/:agentCollectionAddress/:agentCollectionTokenId/is-internal"
	CREATE_AGENT_FIRST_STEP_ENDPOINT = "/agent/create/first-step"
	CREATE_AGENT_SECOND_STEP_ENDPOINT = "/agent/create/second-step"
)

type AgentController struct {
	logger                    *zap.Logger
	getAgent                  use_cases.GetAgentInterface
	getMintedAgentByEnsDomain use_cases.GetMintedAgentByEnsDomainInterface
	listAgents                use_cases.ListAgentsInterface
	isInternalAgent           use_cases.IsInternalAgentInterface
	createAgentFirstStep      use_cases.CreateAgentFirstStepInterface
	createAgentSecondStep     use_cases.CreateAgentSecondStepInterface
}

func NewAgentController(
	logger *zap.Logger,
	getAgent use_cases.GetAgentInterface,
	getMintedAgentByEnsDomain use_cases.GetMintedAgentByEnsDomainInterface,
	listAgents use_cases.ListAgentsInterface,
	isInternalAgent use_cases.IsInternalAgentInterface,
	createAgentFirstStep use_cases.CreateAgentFirstStepInterface,
	createAgentSecondStep use_cases.CreateAgentSecondStepInterface,
) *AgentController {
	return &AgentController{
		logger:                    logger,
		getAgent:                  getAgent,
		getMintedAgentByEnsDomain: getMintedAgentByEnsDomain,
		listAgents:                listAgents,
		isInternalAgent:           isInternalAgent,
		createAgentFirstStep:      createAgentFirstStep,
		createAgentSecondStep:     createAgentSecondStep,
	}
}

func (controller *AgentController) RegisterRoutes(app fiber.Router) {
	app.Get(GET_AGENT_ENDPOINT, controller.GetAgent)
	app.Get(GET_AGENT_BY_ENS_DOMAIN_ENDPOINT, controller.GetMintedAgentByEnsDomain)
	app.Get(LIST_AGENTS_ENDPOINT, controller.ListAgents)
	app.Get(AGENT_IS_INTERNAL_ENDPOINT, controller.IsInternalAgent)
	app.Post(CREATE_AGENT_FIRST_STEP_ENDPOINT, controller.CreateAgentFirstStep)
	app.Post(CREATE_AGENT_SECOND_STEP_ENDPOINT, controller.CreateAgentSecondStep)
}

// GetAgent godoc
// @Summary      Get agent
// @Description  Retrieves details of a specific agent by collection address and token ID
// @Tags         Agent
// @Produce      json
// @Param        agentCollectionAddress path string true "Agent collection address"
// @Param        agentCollectionTokenId path int true "Agent collection token ID"
// @Success      200 {object} responses.AgentResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId} [get]
func (controller *AgentController) GetAgent(c *fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	response, err := controller.getAgent.Get(request)
	if err != nil {
		controller.logger.Error("Failed to get agent", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// ListAgents godoc
// @Summary      List agents
// @Description  Retrieves a paginated list of agents with optional filtering
// @Tags         Agent
// @Produce      json
// @Param        page query int false "Page number"
// @Param        elementsPerPage query int false "Elements per page"
// @Param        sortDirection query string false "Sort direction (asc, desc)"
// @Param        sortBy query string false "Sort by field"
// @Param        isInternal query bool false "Filter by internal status"
// @Success      200 {object} responses.ListAgentsResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agents [get]
func (controller *AgentController) ListAgents(c *fiber.Ctx) error {
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

	return c.JSON(resp)
}

// IsInternalAgent godoc
// @Summary      Check if agent is internal
// @Description  Checks whether a specific agent is an internal agent
// @Tags         Agent
// @Produce      json
// @Param        agentCollectionAddress path string true "Agent collection address"
// @Param        agentCollectionTokenId path int true "Agent collection token ID"
// @Success      200 {object} responses.IsInternalAgentResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/is-internal [get]
func (controller *AgentController) IsInternalAgent(c *fiber.Ctx) error {
	request := new(requests.AgentByCollectionTokenIdRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	response, err := controller.isInternalAgent.IsInternal(request)
	if err != nil {
		controller.logger.Error("Failed to check if agent is internal", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// CreateAgentFirstStep godoc
// @Summary      Start agent creation
// @Description  Uploads agent image and creates agent metadata as the first step of agent creation
// @Tags         Agent
// @Accept       multipart/form-data
// @Produce      json
// @Param        authenticatedWallet formData string true "Authenticated wallet address"
// @Param        collectionAddress formData string true "Collection address"
// @Param        name formData string true "Agent name"
// @Param        role formData string true "Agent role"
// @Param        mission formData string true "Agent mission"
// @Param        walletAddress formData string false "Agent wallet address"
// @Param        cloneOfAgentId formData int false "ID of agent to clone"
// @Param        image formData file true "Agent image"
// @Success      200 {object} responses.AgentResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/create/first-step [post]
func (controller *AgentController) CreateAgentFirstStep(c *fiber.Ctx) error {
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

	request.Image = image
	if common.IsHexAddress(request.AuthenticatedWalletRaw) {
		request.AuthenticatedWallet = common.HexToAddress(request.AuthenticatedWalletRaw)
	}

	response, err := controller.createAgentFirstStep.Execute(request)
	if err != nil {
		if validationErr, ok := err.(*registry_errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// CreateAgentSecondStep godoc
// @Summary      Complete agent creation
// @Description  Completes agent creation by generating the mint signature
// @Tags         Agent
// @Accept       json
// @Produce      json
// @Param        request body requests.AgentByIdRequest true "Agent ID and authenticated wallet"
// @Success      200 {object} responses.CreateAgentSecondStepResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/create/second-step [post]
func (controller *AgentController) CreateAgentSecondStep(c *fiber.Ctx) error {
	request := new(requests.AgentByIdRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	if common.IsHexAddress(request.AuthenticatedWalletRaw) {
		request.AuthenticatedWallet = common.HexToAddress(request.AuthenticatedWalletRaw)
	}

	response, err := controller.createAgentSecondStep.Execute(request)
	if err != nil {
		if validationErr, ok := err.(*registry_errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*registry_errors.NotFoundError); ok {
			return c.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}

		controller.logger.Error("Failed to continue agent creation", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// GetMintedAgentByEnsDomain godoc
// @Summary      Get agent by ENS domain
// @Description  Retrieves a minted agent by its ENS domain name
// @Tags         Agent
// @Produce      json
// @Param        ensDomain path string true "ENS domain name"
// @Success      200 {object} responses.AgentResponse
// @Failure      404 {string} string "Agent not found"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/ens/{ensDomain} [get]
func (controller *AgentController) GetMintedAgentByEnsDomain(c *fiber.Ctx) error {
	request := new(requests.GetAgentByEnsDomainRequest)
	if err := c.ParamsParser(request); err != nil {
		controller.logger.Error("Failed to parse request params", zap.Error(err))
		return err
	}

	response, err := controller.getMintedAgentByEnsDomain.Get(request)
	if err != nil {
		controller.logger.Error("Failed to get agent by ENS domain", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	if response == nil {
		return c.Status(fiber.StatusNotFound).SendString("Agent not found")
	}

	return c.JSON(response)
}
