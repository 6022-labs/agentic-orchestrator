package main

import (
	"fmt"

	common_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Agent Identity Registry API
// @version 1.0
// @description This is the API documentation for the Agent Identity Registry.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	godotenv.Load(".env")

	config, err := common_configurations.LoadKoanfConfig()
	if err != nil {
		panic(err)
	}

	container := configurations.ConfigureDI(config)
	common_configurations.ConfigureLogging(container, config)
	configurations.ConfigureServer(container)
	common_configurations.ConfigureDatabase(container)

	err = container.Invoke(func(app *fiber.App, logger *zap.Logger, hostSettings *settings.HostSettings) error {
		apiPort := fmt.Sprintf("%d", hostSettings.ApiPort)
		logger.Info("Starting server on port " + apiPort)
		return app.Listen(fmt.Sprintf(":%s", apiPort))
	})

	if err != nil {
		panic(err)
	}
}
