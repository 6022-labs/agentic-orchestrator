package main

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_host/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Orchestrator Conversation Broker Host API
// @version 1.0
// @description Bridge-agnostic conversation broker for the agent reasoning domain.
// @BasePath /
func main() {
	godotenv.Load(".env")

	config, err := configurations.LoadKoanfConfig()
	if err != nil {
		panic(err)
	}

	container := configurations.ConfigureDI(config)
	configurations.ConfigureLogging(container)
	configurations.ConfigureServer(container)
	configurations.ConfigureDatabase(container)

	err = container.Invoke(func(app *fiber.App, logger *zap.Logger, hostSettings *settings.HostSettings) error {
		apiPort := fmt.Sprintf("%d", hostSettings.ApiPort)
		logger.Info("Starting server on port " + apiPort)
		return app.Listen(fmt.Sprintf(":%s", apiPort))
	})

	if err != nil {
		panic(err)
	}
}
