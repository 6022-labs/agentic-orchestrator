package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ConfigureDatabase(container *dig.Container) {
	err := container.Invoke(func(orchestratorDatabaseClient db_clients.DatabaseClientInterface, logger *zap.Logger) {
		err := orchestratorDatabaseClient.Migrate()
		if err != nil {
			logger.Fatal("Failed to migrate database", zap.Error(err))
		}
	})
	if err != nil {
		panic(err)
	}
}
