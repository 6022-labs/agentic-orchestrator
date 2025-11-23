package orchestrator_sc_sync_mvc

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	UPDATE_AGENT_OWNER_ENDPOINT         = "/agent/update-agent-owner"
	UPDATE_AGENT_ATTRIBUTE_ENDPOINT     = "/agent/update-agent-attribute"
	FINISH_AGENT_CREATION_ENDPOINT      = "/agent/finish-creation"
	FINISH_AGENT_IMPORTATION_ENDPOINT   = "/agent/finish-importation"
	REMOVE_AGENT_FROM_PROPOSAL_ENDPOINT = "/agent/remove-from-proposal"
	PAIR_PROPOSAL_TO_AGENT_ENDPOINT     = "/agent/pair-proposal-to-agent"
	UPSERT_AGENT_IMAGE_ENDPOINT         = "/agent/upsert-image"
)

type UseCaseController struct {
	logger                  *zap.Logger
	updateAgentOwner        use_cases.UpdateAgentOwnerInterface
	upsertAgentImage        use_cases.UpsertAgentImageInterface
	pairProposalToAgent     use_cases.PairProposalToAgentInterface
	finishAgentCreation     use_cases.FinishAgentCreationInterface
	updateAgentAttribute    use_cases.UpdateAgentAttributeInterface
	finishAgentImportation  use_cases.FinishAgentImportationInterface
	removeAgentFromProposal use_cases.RemoveAgentFromProposalInterface
}

func NewUseCaseController(
	logger *zap.Logger,
	updateAgentOwner use_cases.UpdateAgentOwnerInterface,
	upsertAgentImage use_cases.UpsertAgentImageInterface,
	pairProposalToAgent use_cases.PairProposalToAgentInterface,
	finishAgentCreation use_cases.FinishAgentCreationInterface,
	updateAgentAttribute use_cases.UpdateAgentAttributeInterface,
	finishAgentImportation use_cases.FinishAgentImportationInterface,
	removeAgentFromProposal use_cases.RemoveAgentFromProposalInterface,
) *UseCaseController {
	return &UseCaseController{
		logger:                  logger,
		updateAgentOwner:        updateAgentOwner,
		upsertAgentImage:        upsertAgentImage,
		pairProposalToAgent:     pairProposalToAgent,
		finishAgentCreation:     finishAgentCreation,
		updateAgentAttribute:    updateAgentAttribute,
		finishAgentImportation:  finishAgentImportation,
		removeAgentFromProposal: removeAgentFromProposal,
	}
}

func (controller *UseCaseController) InitRoutes(app fiber.Router) {
	app.Post(FINISH_AGENT_CREATION_ENDPOINT, controller.FinishAgentCreation)
	app.Post(FINISH_AGENT_IMPORTATION_ENDPOINT, controller.FinishAgentImportation)
	app.Post(UPDATE_AGENT_OWNER_ENDPOINT, controller.UpdateAgentOwner)
	app.Post(UPDATE_AGENT_ATTRIBUTE_ENDPOINT, controller.UpdateAgentAttribute)
	app.Post(PAIR_PROPOSAL_TO_AGENT_ENDPOINT, controller.PairProposalToAgent)
	app.Post(REMOVE_AGENT_FROM_PROPOSAL_ENDPOINT, controller.RemoveAgentFromProposal)
	app.Post(UPSERT_AGENT_IMAGE_ENDPOINT, controller.UpsertAgentImage)
}

// UpdateAgentOwner godoc
// @Summary      Update agent owner
// @Description  Updates the ownership of an agent
// @Tags         sc-sync
// @Accept       json
// @Produce      json
// @Param        request body requests.UpdateAgentOwnerRequest true "owner info"
// @Success      200 {string} string "OK"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/update-agent-owner [post]
func (controller *UseCaseController) UpdateAgentOwner(c *fiber.Ctx) error {
	request := new(requests.UpdateAgentOwnerRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.updateAgentOwner.Update(request)
	if err != nil {
		controller.logger.Error("Error updating agent owner", zap.Error(err))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpsertAgentImage godoc
// @Summary      Upsert agent image
// @Description  Uploads or updates the agent image
// @Tags         sc-sync
// @Accept       json
// @Produce      json
// @Param        request body requests.UpsertAgentImageRequest true "image"
// @Success      200 {string} string "OK"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/upsert-image [post]
func (controller *UseCaseController) UpsertAgentImage(c *fiber.Ctx) error {
	request := new(requests.UpsertAgentImageRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.upsertAgentImage.Upsert(request)
	if err != nil {
		controller.logger.Error("Error upserting agent image", zap.Error(err))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateAgentMission godoc
// @Summary      Update agent attribute
// @Description  Updates the attribute details of an agent
// @Tags         sc-sync
// @Accept       json
// @Produce      json
// @Param        request body requests.UpdateAgentAttributeRequest true "attribute info"
// @Success      200 {string} string "OK"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/update-agent-attribute [post]
func (controller *UseCaseController) UpdateAgentAttribute(c *fiber.Ctx) error {
	request := new(requests.UpdateAgentAttributeRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.updateAgentAttribute.Update(request)
	if err != nil {
		controller.logger.Error("Error updating agent attribute", zap.Error(err))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// PairProposalToAgent godoc
// @Summary      Pair proposal to agent
// @Description  Links a proposal with an agent
// @Tags         sc-sync
// @Accept       json
// @Produce      json
// @Param        request body requests.PairProposalToAgentRequest true "pairing"
// @Success      200 {string} string "OK"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/pair-proposal-to-agent [post]
func (controller *UseCaseController) PairProposalToAgent(c *fiber.Ctx) error {
	request := new(requests.PairProposalToAgentRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.pairProposalToAgent.Pair(request)
	if err != nil {
		controller.logger.Error("Error pairing proposal to agent", zap.Error(err))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// FinishAgentCreation godoc
// @Summary      Finish agent creation
// @Description  Marks the agent creation process as complete
// @Tags         sc-sync
// @Accept       json
// @Produce      json
// @Param        request body requests.FinishAgentCreationRequest true "info"
// @Success      200 {string} string "OK"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/finish-creation [post]
func (controller *UseCaseController) FinishAgentCreation(c *fiber.Ctx) error {
	request := new(requests.FinishAgentCreationRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.finishAgentCreation.Finish(request)
	if err != nil {
		controller.logger.Error("Error finishing agent creation", zap.Error(err))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// FinishAgentImportation godoc
// @Summary      Finish agent importation
// @Description  Marks the external agent importation as complete
// @Tags         sc-sync
// @Accept       json
// @Produce      json
// @Param        request body requests.FinishAgentImportationRequest true "info"
// @Success      200 {string} string "OK"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/finish-importation [post]
func (controller *UseCaseController) FinishAgentImportation(c *fiber.Ctx) error {
	request := new(requests.FinishAgentImportationRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.finishAgentImportation.Finish(request)
	if err != nil {
		controller.logger.Error("Error finishing agent importation", zap.Error(err))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// RemoveAgentFromProposal godoc
// @Summary      Remove agent from proposal
// @Description  Detaches an agent from a proposal
// @Tags         sc-sync
// @Accept       json
// @Produce      json
// @Param        request body requests.RemoveAgentFromProposalRequest true "info"
// @Success      200 {string} string "OK"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /agent/remove-from-proposal [post]
func (controller *UseCaseController) RemoveAgentFromProposal(c *fiber.Ctx) error {
	request := new(requests.RemoveAgentFromProposalRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	err := controller.removeAgentFromProposal.Remove(request)
	if err != nil {
		controller.logger.Error("Error removing agent from proposal", zap.Error(err))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
