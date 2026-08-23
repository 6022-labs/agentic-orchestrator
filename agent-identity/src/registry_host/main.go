package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/host_configurations"
	common_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/settings"
	"github.com/gofiber/fiber/v3"
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

	shutdownTelemetry := host_configurations.ConfigureTelemetry(container, configurations.AppName)
	defer func() { _ = shutdownTelemetry(context.Background()) }()

	configurations.ConfigureServer(container)
	common_configurations.ConfigureDatabase(container)

	err = container.Invoke(func(app *fiber.App, logger *zap.Logger, hostSettings *settings.HostSettings) error {
		apiPort := fmt.Sprintf("%d", hostSettings.ApiPort)
		logger.Info("Starting server on port " + apiPort)

		serverErrCh := make(chan error, 1)
		go func() {
			serverErrCh <- app.Listen(fmt.Sprintf(":%s", apiPort), fiber.ListenConfig{ListenerNetwork: "tcp"})
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
