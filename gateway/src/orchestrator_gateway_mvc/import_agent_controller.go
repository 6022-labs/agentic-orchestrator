package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	_ "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	IMPORT_AGENT_FIRST_STEP_ENDPOINT  = "/agent/:agentCollectionAddress/:agentCollectionTokenId/importation/first-step"
	IMPORT_AGENT_SECOND_STEP_ENDPOINT = "/agent/:agentCollectionAddress/:agentCollectionTokenId/importation/second-step"
	IMPORT_AGENT_THIRD_STEP_ENDPOINT  = "/agent/:agentCollectionAddress/:agentCollectionTokenId/importation/third-step"
)

type ImportAgentController struct {
	logger                        *zap.Logger
	importExternalAgentFirstStep  use_cases.ImportExternalAgentFirstStepInterface
	importExternalAgentSecondStep use_cases.ImportExternalAgentSecondStepInterface
	importExternalAgentThirdStep  use_cases.ImportExternalAgentThirdStepInterface
}

func NewImportAgentController(
	logger *zap.Logger,
	importExternalAgentFirstStep use_cases.ImportExternalAgentFirstStepInterface,
	importExternalAgentSecondStep use_cases.ImportExternalAgentSecondStepInterface,
	importExternalAgentThirdStep use_cases.ImportExternalAgentThirdStepInterface,
) *ImportAgentController {
	return &ImportAgentController{
		logger:                        logger,
		importExternalAgentFirstStep:  importExternalAgentFirstStep,
		importExternalAgentSecondStep: importExternalAgentSecondStep,
		importExternalAgentThirdStep:  importExternalAgentThirdStep,
	}
}

func (controller *ImportAgentController) InitAuthenticatedRoutes(app fiber.Router) {
	app.Post(IMPORT_AGENT_FIRST_STEP_ENDPOINT, controller.AgentImportationFirstStep)
	app.Post(IMPORT_AGENT_SECOND_STEP_ENDPOINT, controller.AgentImportationSecondStep)
	app.Post(IMPORT_AGENT_THIRD_STEP_ENDPOINT, controller.AgentImportationThirdStep)
}

func (controller *ImportAgentController) InitUnauthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

// AgentImportationFirstStep godoc
// @Summary      Start agent importation
// @Description  Begins the importation process of an external agent
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.ImportExternalAgentFirstStepResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/import/first-step [post]
func (controller *ImportAgentController) AgentImportationFirstStep(c *fiber.Ctx) error {
	controller.logger.Info("Agent importation first step request received")
	controller.logger.Debug("Request body", zap.Any("body", c.Body()))

	request := new(requests.AgentByCollectionTokenIdRequest)
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

	response, err := controller.importExternalAgentFirstStep.Execute(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Agent importation first step request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// AgentImportationSecondStep godoc
// @Summary      Continue agent importation
// @Description  Performs the second step of agent importation
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Param        request body requests.ImportExternalAgentSecondStepRequest true "importation second step request"
// @Success      200 {object} responses.ImportExternalAgentSecondStepResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/import/second-step [post]
func (controller *ImportAgentController) AgentImportationSecondStep(c *fiber.Ctx) error {
	controller.logger.Info("Agent importation second step request received")
	controller.logger.Debug("Request body", zap.Any("body", c.Body()))

	request := new(requests.ImportExternalAgentSecondStepRequest)
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

	response, err := controller.importExternalAgentSecondStep.Execute(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Agent importation second step request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}

// AgentImportationThirdStep godoc
// @Summary      Finalize agent importation
// @Description  Completes the external agent importation
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        agentCollectionAddress path string true "collection address"
// @Param        agentCollectionTokenId path string true "token id"
// @Success      200 {object} responses.ImportExternalAgentThirdStepResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /agent/{agentCollectionAddress}/{agentCollectionTokenId}/import/third-step [post]
func (controller *ImportAgentController) AgentImportationThirdStep(c *fiber.Ctx) error {
	controller.logger.Info("Agent importation third step request received")
	controller.logger.Debug("Request body", zap.Any("body", c.Body()))

	request := new(requests.AgentByCollectionTokenIdRequest)
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

	response, err := controller.importExternalAgentThirdStep.Execute(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to start agent creation", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Agent importation third step request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}
