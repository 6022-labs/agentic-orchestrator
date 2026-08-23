package factories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	orchestrator_db_orchestrator "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"go.uber.org/zap"
)

type DatabaseClientFactory struct {
	logger             *zap.Logger
	dbQueryMetrics     metrics_interfaces.DbQueryMetricsInterface
	databaseSettings   *settings.DatabaseSettings
	migrationsSettings *settings.MigrationsSettings
}

func NewDatabaseClientFactory(
	logger *zap.Logger,
	dbQueryMetrics metrics_interfaces.DbQueryMetricsInterface,
	databaseSettings *settings.DatabaseSettings,
	migrationsSettings *settings.MigrationsSettings,
) *DatabaseClientFactory {
	return &DatabaseClientFactory{
		logger:             logger,
		dbQueryMetrics:     dbQueryMetrics,
		databaseSettings:   databaseSettings,
		migrationsSettings: migrationsSettings,
	}
}

func (factory *DatabaseClientFactory) CreateDatabaseClient() (db_clients.DatabaseClientInterface, error) {
	return orchestrator_db_orchestrator.NewDatabaseClient(
		factory.databaseSettings.Url,
		factory.logger,
		factory.migrationsSettings,
		factory.dbQueryMetrics,
	)
}
