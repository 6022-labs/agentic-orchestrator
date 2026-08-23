package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/host_configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_host/settings"
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

	shutdownTelemetry := host_configurations.ConfigureTelemetry(container, configurations.AppName)
	defer func() { _ = shutdownTelemetry(context.Background()) }()

	configurations.ConfigureServer(container)
	configurations.ConfigureDatabase(container)
	configurations.ConfigureDaemon(container)

	err = container.Invoke(func(
		app *fiber.App,
		logger *zap.Logger,
		hostSettings *settings.HostSettings,
		daemon interfaces.ConversationProcessorDaemonInterface,
	) error {
		port := fmt.Sprintf(":%d", hostSettings.ApiPort)
		logger.Info("Starting server", zap.String("port", port))

		serverErrCh := make(chan error, 1)
		go func() {
			serverErrCh <- app.Listen(port)
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigCh:
			logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

			daemon.Stop()

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
