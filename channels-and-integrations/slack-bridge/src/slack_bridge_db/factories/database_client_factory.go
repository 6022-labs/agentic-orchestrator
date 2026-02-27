package factories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"go.uber.org/zap"
)

type DatabaseClientFactory struct {
	logger                                  *zap.Logger
	migrationsSettings                      *settings.MigrationsSettings
	orchestratorSlackBridgeDatabaseSettings *settings.DatabaseSettings
}

func NewDatabaseClientFactory(
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
	orchestratorSlackBridgeDatabaseSettings *settings.DatabaseSettings,
) *DatabaseClientFactory {
	return &DatabaseClientFactory{
		logger:                                  logger,
		migrationsSettings:                      migrationsSettings,
		orchestratorSlackBridgeDatabaseSettings: orchestratorSlackBridgeDatabaseSettings,
	}
}

func (factory *DatabaseClientFactory) CreateOrchestratorSlackBridgeDBClient() (*slack_bridge_db.DatabaseClient, error) {
	return slack_bridge_db.NewDatabaseClient(
		factory.orchestratorSlackBridgeDatabaseSettings.Url,
		factory.logger,
		factory.migrationsSettings,
	)
}
