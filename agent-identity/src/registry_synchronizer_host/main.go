package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/host_configurations"
	common_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_host/configurations"
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

	k, err := common_configurations.LoadKoanfConfig()
	if err != nil {
		panic(err)
	}

	container := configurations.ConfigureDI(k)
	common_configurations.ConfigureLogging(container, k)

	shutdownTelemetry := host_configurations.ConfigureTelemetry(container, configurations.AppName)
	defer func() { _ = shutdownTelemetry(context.Background()) }()

	configurations.ConfigureServer(container)
	common_configurations.ConfigureDatabase(container)

	err = container.Invoke(func(
		logger *zap.Logger,
		synchronizer *configurations.Synchronizer,
	) error {
		logger.Info("Starting synchronizer application", zap.String("name", configurations.AppName))

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		synchronizerErrCh := make(chan error, 1)
		go func() {
			synchronizerErrCh <- synchronizer.Start(ctx)
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigCh:
			logger.Info("Shutdown signal received", zap.String("signal", sig.String()))
			cancel()

			if err := <-synchronizerErrCh; err != nil {
				logger.Error("Synchronizer stopped with error", zap.Error(err))
			} else {
				logger.Info("Synchronizer stopped")
			}

			logger.Info("Shutdown complete")
			_ = logger.Sync()
			return nil

		case err := <-synchronizerErrCh:
			if err != nil {
				return fmt.Errorf("synchronizer exited unexpectedly: %w", err)
			}
			return nil
		}
	})
	if err != nil {
		panic(err)
	}
}
