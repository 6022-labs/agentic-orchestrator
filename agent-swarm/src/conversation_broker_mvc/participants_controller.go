package conversation_broker_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v3"
)

const (
	REGISTER_AGENT_ENDPOINT = "/participants/agents"
	REGISTER_USER_ENDPOINT  = "/participants/users"
)

type ParticipantsController struct {
	errorWriter   *mvc.ErrorWriter
	registerAgent use_cases.RegisterAgentInterface
	registerUser  use_cases.RegisterUserInterface
}

func NewParticipantsController(
	errorWriter *mvc.ErrorWriter,
	registerAgent use_cases.RegisterAgentInterface,
	registerUser use_cases.RegisterUserInterface,
) *ParticipantsController {
	return &ParticipantsController{
		errorWriter:   errorWriter,
		registerAgent: registerAgent,
		registerUser:  registerUser,
	}
}

func (c *ParticipantsController) InitRoutes(app fiber.Router) {
	app.Post(REGISTER_AGENT_ENDPOINT, c.RegisterAgent)
	app.Post(REGISTER_USER_ENDPOINT, c.RegisterUser)
}

// @Summary Register Agent Participant
// @Description Upserts an agent participant by NFT identity: 200 if known, 201 if new; role + ENS from the registry.
// @Tags Participants
// @Accept json
// @Produce json
// @Param request body requests.RegisterAgentRequest true "Register Agent Request"
// @Success 200 {object} responses.ParticipantResponse
// @Success 201 {object} responses.ParticipantResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 400 {object} errors.BadRequestError "invalid_body"
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {object} errors.InternalError "internal_error"
// @Router /participants/agents [post]
func (c *ParticipantsController) RegisterAgent(ctx fiber.Ctx) error {
	req := new(requests.RegisterAgentRequest)
	if err := ctx.Bind().Body(req); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_body", err.Error()))
	}
	p, existed, err := c.registerAgent.Register(ctx.Context(), req)
	if err != nil {
		return c.errorWriter.Write(ctx, err)
	}
	status := fiber.StatusCreated
	if existed {
		status = fiber.StatusOK
	}
	return ctx.Status(status).JSON(p)
}

// @Summary Register User Participant
// @Description Creates a user participant by display name only. No dedup — bridges must cache the returned id.
// @Tags Participants
// @Accept json
// @Produce json
// @Param request body requests.RegisterUserRequest true "Register User Request"
// @Success 201 {object} responses.ParticipantResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 400 {object} errors.BadRequestError "invalid_body"
// @Failure 500 {object} errors.InternalError "internal_error"
// @Router /participants/users [post]
func (c *ParticipantsController) RegisterUser(ctx fiber.Ctx) error {
	req := new(requests.RegisterUserRequest)
	if err := ctx.Bind().Body(req); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_body", err.Error()))
	}
	p, err := c.registerUser.Register(ctx.Context(), req)
	if err != nil {
		return c.errorWriter.Write(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(p)
}
