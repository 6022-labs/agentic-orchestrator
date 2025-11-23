package configurations

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/db_clients"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ConfigureDatabase(container *dig.Container) {
	err := container.Invoke(func(databaseClient db_clients.DatabaseClientInterface, logger *zap.Logger) {
		err := databaseClient.Migrate()
		if err != nil {
			logger.Fatal("Failed to migrate database", zap.Error(err))
		}
	})
	if err != nil {
		panic(err)
	}
}
