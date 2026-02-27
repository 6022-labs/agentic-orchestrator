package factories

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/db_clients"
	orchestrator_db_orchestrator "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/settings"
	"go.uber.org/zap"
)

type DatabaseClientFactory struct {
	logger             *zap.Logger
	databaseSettings   *settings.DatabaseSettings
	migrationsSettings *settings.MigrationsSettings
}

func NewDatabaseClientFactory(
	logger *zap.Logger,
	databaseSettings *settings.DatabaseSettings,
	migrationsSettings *settings.MigrationsSettings,
) *DatabaseClientFactory {
	return &DatabaseClientFactory{
		logger:             logger,
		databaseSettings:   databaseSettings,
		migrationsSettings: migrationsSettings,
	}
}

func (factory *DatabaseClientFactory) CreateDatabaseClient() (db_clients.DatabaseClientInterface, error) {
	return orchestrator_db_orchestrator.NewDatabaseClient(
		factory.databaseSettings.Url,
		factory.logger,
		factory.migrationsSettings,
	)
}
