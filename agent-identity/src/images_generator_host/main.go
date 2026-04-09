package main

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator_host/settings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Orchestrator images generator Host API
// @version 1.0
// @description This is the API for the orchestrator images generator host.
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

	err = container.Invoke(func(app *fiber.App, logger *zap.Logger, hostSettings *settings.HostSettings) error {
		apiPort := fmt.Sprintf("%d", hostSettings.ApiPort)
		logger.Info("Starting server on port " + apiPort)
		return app.Listen(fmt.Sprintf(":%s", apiPort))
	})

	if err != nil {
		panic(err)
	}
}
