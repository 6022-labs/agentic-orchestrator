package conversation_broker_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v3"
)

const REGISTER_BRIDGE_ENDPOINT = "/bridges"

type BridgesController struct {
	errorWriter    *mvc.ErrorWriter
	registerBridge use_cases.RegisterBridgeInterface
}

func NewBridgesController(
	errorWriter *mvc.ErrorWriter,
	registerBridge use_cases.RegisterBridgeInterface,
) *BridgesController {
	return &BridgesController{errorWriter: errorWriter, registerBridge: registerBridge}
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
// @Failure 400 {object} errors.BadRequestError "invalid_body"
// @Failure 500 {object} errors.InternalError "internal_error"
// @Router /bridges [post]
func (c *BridgesController) RegisterBridge(ctx fiber.Ctx) error {
	req := new(requests.RegisterBridgeRequest)
	if err := ctx.Bind().Body(req); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_body", err.Error()))
	}
	bridge, existed, err := c.registerBridge.Register(ctx.Context(), req)
	if err != nil {
		return c.errorWriter.Write(ctx, err)
	}
	status := fiber.StatusCreated
	if existed {
		status = fiber.StatusOK
	}
	return ctx.Status(status).JSON(bridge)
}
