package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const REGISTER_BRIDGE_ENDPOINT = "/bridges"

type BridgesController struct {
	logger         *zap.Logger
	registerBridge use_cases.RegisterBridgeInterface
}

func NewBridgesController(logger *zap.Logger, registerBridge use_cases.RegisterBridgeInterface) *BridgesController {
	return &BridgesController{logger: logger, registerBridge: registerBridge}
}

func (c *BridgesController) InitRoutes(app fiber.Router) {
	app.Post(REGISTER_BRIDGE_ENDPOINT, c.RegisterBridge)
}

// @Summary Register Bridge
// @Description Upserts a bridge: 200 with the existing row when the given `id` exists, 201 with the new row otherwise.
// @Tags Bridges
// @Accept json
// @Produce json
// @Param request body requests.RegisterBridgeRequest true "Register Bridge Request"
// @Success 200 {object} responses.BridgeResponse
// @Success 201 {object} responses.BridgeResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /bridges [post]
func (c *BridgesController) RegisterBridge(ctx *fiber.Ctx) error {
	req := new(requests.RegisterBridgeRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	bridge, existed, err := c.registerBridge.Register(ctx.UserContext(), req)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	status := fiber.StatusCreated
	if existed {
		status = fiber.StatusOK
	}
	return ctx.Status(status).JSON(bridge)
}
