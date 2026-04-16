package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
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
// @Description Upserts a bridge instance. When the caller supplies `id` and the row already exists the broker returns 200 with the existing row; when the row does not exist (whether or not an id was provided) the broker returns 201 with the new row.
// @Tags Bridges
// @Accept json
// @Produce json
// @Param request body requests.RegisterBridgeRequest true "Register Bridge Request"
// @Success 200 {object} models.Bridge
// @Success 201 {object} models.Bridge
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {string} string "Internal Server Error"
// @Router /bridges [post]
func (c *BridgesController) RegisterBridge(ctx *fiber.Ctx) error {
	req := new(requests.RegisterBridgeRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	bridge, existed, err := c.registerBridge.Register(req)
	if err != nil {
		return mapError(c.logger, ctx, err)
	}
	status := fiber.StatusCreated
	if existed {
		status = fiber.StatusOK
	}
	return ctx.Status(status).JSON(bridge)
}

// mapError centralizes the error-to-HTTP mapping for all controllers.
func mapError(logger *zap.Logger, ctx *fiber.Ctx, err error) error {
	if validationErr, ok := err.(*errors.ValidationError); ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
	}
	if notFoundErr, ok := err.(*errors.NotFoundError); ok {
		return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
	}
	logger.Error("internal error", zap.Error(err))
	return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
}
