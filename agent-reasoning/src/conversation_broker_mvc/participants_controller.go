package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	REGISTER_MINTED_PARTICIPANT_ENDPOINT     = "/participants/minted"
	REGISTER_NON_MINTED_PARTICIPANT_ENDPOINT = "/participants/non-minted"
)

type ParticipantsController struct {
	logger                       *zap.Logger
	registerMintedParticipant    use_cases.RegisterMintedParticipantInterface
	registerNonMintedParticipant use_cases.RegisterNonMintedParticipantInterface
}

func NewParticipantsController(
	logger *zap.Logger,
	registerMintedParticipant use_cases.RegisterMintedParticipantInterface,
	registerNonMintedParticipant use_cases.RegisterNonMintedParticipantInterface,
) *ParticipantsController {
	return &ParticipantsController{
		logger:                       logger,
		registerMintedParticipant:    registerMintedParticipant,
		registerNonMintedParticipant: registerNonMintedParticipant,
	}
}

func (c *ParticipantsController) InitRoutes(app fiber.Router) {
	app.Post(REGISTER_MINTED_PARTICIPANT_ENDPOINT, c.RegisterMinted)
	app.Post(REGISTER_NON_MINTED_PARTICIPANT_ENDPOINT, c.RegisterNonMinted)
}

// @Summary Register Minted Participant
// @Description Registers a minted participant by NFT identity. Broker fetches role + ENS from the agent-identity registry.
// @Tags Participants
// @Accept json
// @Produce json
// @Param request body requests.RegisterMintedParticipantRequest true "Register Minted Participant Request"
// @Success 201 {object} models.Participant
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /participants/minted [post]
func (c *ParticipantsController) RegisterMinted(ctx *fiber.Ctx) error {
	req := new(requests.RegisterMintedParticipantRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	p, err := c.registerMintedParticipant.Register(req)
	if err != nil {
		return mapError(c.logger, ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(p)
}

// @Summary Register Non-Minted Participant
// @Description Registers a human participant identified only by display name (not minted on-chain).
// @Tags Participants
// @Accept json
// @Produce json
// @Param request body requests.RegisterNonMintedParticipantRequest true "Register Non-Minted Participant Request"
// @Success 201 {object} models.Participant
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /participants/non-minted [post]
func (c *ParticipantsController) RegisterNonMinted(ctx *fiber.Ctx) error {
	req := new(requests.RegisterNonMintedParticipantRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	p, err := c.registerNonMintedParticipant.Register(req)
	if err != nil {
		return mapError(c.logger, ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(p)
}
