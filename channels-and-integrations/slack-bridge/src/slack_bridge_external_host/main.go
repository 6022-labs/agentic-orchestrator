package main

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/configurations"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load(".env")

	container := configurations.ConfigureDI()
	configurations.ConfigureLogging(container)
	configurations.ConfigureServer(container)
	configurations.ConfigureDatabase(container)

	err := container.Invoke(func(app *fiber.App, logger *zap.Logger, hostSettings *settings.HostSettings) error {
		apiPort := fmt.Sprintf("%d", hostSettings.ApiPort)
		logger.Info("Starting server on port " + apiPort)
		return app.Listen(fmt.Sprintf(":%s", apiPort))
	})

	if err != nil {
		panic(err)
	}
}
