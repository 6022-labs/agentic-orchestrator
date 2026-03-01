package main

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_host/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Orchestrator SC Sync Host
// @version 1.0
// @description This is the API documentation for the Orchestrator SC Sync.
// @BasePath /
func main() {
	godotenv.Load(".env")

	config, err := configurations.LoadKoanfConfig()
	if err != nil {
		panic(err)
	}

	container := configurations.ConfigureDI(config)
	configurations.ConfigureLogging(container)
	configurations.ConfigureDatabase(container)
	configurations.ConfigureServer(container)

	container.Invoke(func(logger *zap.Logger, hostSettings *settings.HostSettings) {
		logger.Info("Starting Orchestrator SC Sync Host")

		if !hostSettings.UseApi && !hostSettings.UseListeners {
			logger.Info("Both API and listeners are disabled. Waiting for 10 seconds before exiting to avoid too much requests to RPC nodes due to HTTP synchronization.")
			time.Sleep(10 * time.Second)
		}
	})

	err = container.Invoke(func(logger *zap.Logger, synchronizeAllAgentCollections use_cases.SynchronizeAllAgentCollectionsInterface) error {
		logger.Info("Synchronizing all agent collections...")
		return synchronizeAllAgentCollections.Synchronize()
	})
	if err != nil {
		panic(err)
	}

	err = container.Invoke(func(logger *zap.Logger, synchronizeAllAgents use_cases.SynchronizeAllAgentsInterface) error {
		logger.Info("Synchronizing all agents...")
		return synchronizeAllAgents.Synchronize()
	})
	if err != nil {
		panic(err)
	}

	err = container.Invoke(func(logger *zap.Logger, hostSettings *settings.HostSettings) error {
		if !hostSettings.UseApi && !hostSettings.UseListeners {
			return nil
		}

		if !hostSettings.UseApi && hostSettings.UseListeners {
			logger.Info("Listeners are enabled, but API is disabled. Run long-running tasks to keep the application running.")
			select {}
		}

		return container.Invoke(func(app *fiber.App) error {
			apiPort := fmt.Sprintf("%d", hostSettings.ApiPort)
			logger.Info("Starting server on port " + apiPort)
			return app.Listen(fmt.Sprintf(":%s", apiPort))
		})
	})

	if err != nil {
		panic(err)
	}
}
