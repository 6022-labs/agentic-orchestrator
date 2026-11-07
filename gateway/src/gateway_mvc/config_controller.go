package orchestrator_gateway_mvc

import (
	_ "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const GET_CONFIG_ENDPOINT = "/config"

type ConfigController struct {
	logger      *zap.Logger
	errorWriter *mvc.ErrorWriter
	getConfig   use_cases.GetConfigInterface
}

func NewConfigController(
	logger *zap.Logger,
	errorWriter *mvc.ErrorWriter,
	getConfig use_cases.GetConfigInterface,
) *ConfigController {
	return &ConfigController{logger: logger, errorWriter: errorWriter, getConfig: getConfig}
}

func (controller *ConfigController) InitAuthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

func (controller *ConfigController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Get(GET_CONFIG_ENDPOINT, controller.Get)
}

// Get godoc
// @Summary      Public node config
// @Description  The chain config (keyed by chain id) and IPFS gateways the console boots from
// @Tags         config
// @Produce      json
// @Success      200 {object} responses.ConfigResponse
// @Failure      500 {object} errors.InternalError
// @Router       /config [get]
func (controller *ConfigController) Get(c fiber.Ctx) error {
	response, err := controller.getConfig.Get()
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}
	return c.JSON(response)
}
