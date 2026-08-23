package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/host_configurations"
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

	shutdownTelemetry := host_configurations.ConfigureTelemetry(container, configurations.AppName)
	defer func() { _ = shutdownTelemetry(context.Background()) }()

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

		serverErrCh := make(chan error, 1)
		go func() {
			serverErrCh <- app.Listen(fmt.Sprintf(":%s", apiPort))
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigCh:
			logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

			if err := app.Shutdown(); err != nil {
				logger.Warn("Fiber shutdown returned an error", zap.Error(err))
			}

			if err := <-serverErrCh; err != nil {
				logger.Error("Server stopped with error", zap.Error(err))
			} else {
				logger.Info("Server stopped")
			}

			logger.Info("Shutdown complete")
			_ = logger.Sync()
			return nil

		case err := <-serverErrCh:
			if err != nil {
				return fmt.Errorf("server exited unexpectedly: %w", err)
			}
			return nil
		}
	})

	if err != nil {
		panic(err)
	}
}
