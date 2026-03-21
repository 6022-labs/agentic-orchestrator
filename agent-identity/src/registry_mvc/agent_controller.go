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
