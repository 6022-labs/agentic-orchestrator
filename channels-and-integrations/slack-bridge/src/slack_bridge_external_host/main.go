package main

import (
	"context"

	"fmt"

	orchestrator_slack_bridge_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/configurations"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Slack Bridge External API
// @version 1.0
// @description This is the API documentation for the Slack Bridge External Host.
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
		if _, err := bridgeRegistrar.Ensure(context.Background()); err != nil {
			return fmt.Errorf("failed to register bridge with conversation broker: %w", err)
		}

		apiPort := fmt.Sprintf("%d", hostSettings.ApiPort)
		logger.Info("Starting server on port " + apiPort)
		return app.Listen(fmt.Sprintf(":%s", apiPort))
	})

	if err != nil {
		panic(err)
	}
}
