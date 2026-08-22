package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	REGISTER_AGENT_ENDPOINT = "/participants/agents"
	REGISTER_USER_ENDPOINT  = "/participants/users"
)

type ParticipantsController struct {
	logger        *zap.Logger
	registerAgent use_cases.RegisterAgentInterface
	registerUser  use_cases.RegisterUserInterface
}

func NewParticipantsController(
	logger *zap.Logger,
	registerAgent use_cases.RegisterAgentInterface,
	registerUser use_cases.RegisterUserInterface,
) *ParticipantsController {
	return &ParticipantsController{
		logger:        logger,
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
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /participants/agents [post]
func (c *ParticipantsController) RegisterAgent(ctx *fiber.Ctx) error {
	req := new(requests.RegisterAgentRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	p, existed, err := c.registerAgent.Register(ctx.UserContext(), req)
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
// @Failure 500 {string} string "Internal Server Error"
// @Router /participants/users [post]
func (c *ParticipantsController) RegisterUser(ctx *fiber.Ctx) error {
	req := new(requests.RegisterUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	p, err := c.registerUser.Register(ctx.UserContext(), req)
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
	return ctx.Status(fiber.StatusCreated).JSON(p)
}
