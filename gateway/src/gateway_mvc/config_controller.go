package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const GET_CONFIG_ENDPOINT = "/config"

type ConfigController struct {
	logger              *zap.Logger
	getBlockchainConfig use_cases.GetBlockchainConfigInterface
}

func NewConfigController(
	logger *zap.Logger,
	getBlockchainConfig use_cases.GetBlockchainConfigInterface,
) *ConfigController {
	return &ConfigController{logger: logger, getBlockchainConfig: getBlockchainConfig}
}

func (controller *ConfigController) InitAuthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

func (controller *ConfigController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Get(GET_CONFIG_ENDPOINT, controller.Get)
}

// Get godoc
// @Summary      Public chain config
// @Description  The EVM chain config (keyed by chain id) the console boots from
// @Tags         config
// @Produce      json
// @Success      200 {object} responses.BlockchainConfigResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /config [get]
func (controller *ConfigController) Get(c fiber.Ctx) error {
	response, err := controller.getBlockchainConfig.Get()
	if err != nil {
		controller.logger.Error("Failed to get config", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.JSON(response)
}
