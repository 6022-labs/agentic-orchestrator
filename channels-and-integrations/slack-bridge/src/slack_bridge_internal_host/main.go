package main

import (
	"fmt"

	orchestrator_slack_bridge_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/configurations"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Slack Bridge Internal API
// @version 1.0
// @description This is the API documentation for the Slack Bridge Internal Host.
// @BasePath /
func main() {
	godotenv.Load(".env")

	config, err := orchestrator_slack_bridge_configurations.LoadKoanfConfig()
	if err != nil {
		panic(err)
	}

	container := configurations.ConfigureDI(config)
	configurations.ConfigureLogging(container)
	configurations.ConfigureServer(container)
	configurations.ConfigureDatabase(container)

	err = container.Invoke(func(
		app *fiber.App,
		logger *zap.Logger,
		hostSettings *settings.HostSettings,
		bridgeRegistrar interfaces.BridgeRegistrarInterface,
	) error {
		if _, err := bridgeRegistrar.Ensure(); err != nil {
			return fmt.Errorf("failed to register bridge with conversation broker: %w", err)
		}

		port := fmt.Sprintf(":%d", hostSettings.ApiPort)
		logger.Info("Starting server", zap.String("port", port))
		return app.Listen(port)
	})
	if err != nil {
		panic(err)
	}
}
