package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ConfigureDatabase(container *dig.Container) {
	err := container.Invoke(func(orchestratorDbClient *slack_bridge_db.DatabaseClient, logger *zap.Logger) {
		err := orchestratorDbClient.Migrate()
		if err != nil {
			logger.Fatal("Failed to migrate database", zap.Error(err))
		}
	})
	if err != nil {
		panic(err)
	}
}
