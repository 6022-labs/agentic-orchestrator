package main

import (
	common_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/configurations"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_host/configurations"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const appName = "6022-AgentIdentity-RegistrySynchronizer"

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
	configurations.ConfigureServer(container)
	common_configurations.ConfigureDatabase(container)

	err = container.Invoke(func(
		logger *zap.Logger,
		synchronizer *configurations.Synchronizer,
	) error {
		logger.Info("Starting synchronizer application", zap.String("name", appName))

		return synchronizer.Start()
	})
	if err != nil {
		panic(err)
	}
}
