package factories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	conversation_broker_db "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/settings"
	"go.uber.org/zap"
)

type DatabaseClientFactory struct {
	logger             *zap.Logger
	dbQueryMetrics     metrics_interfaces.DbQueryMetricsInterface
	migrationsSettings *settings.MigrationsSettings
	databaseSettings   *settings.DatabaseSettings
}

func NewDatabaseClientFactory(
	logger *zap.Logger,
	dbQueryMetrics metrics_interfaces.DbQueryMetricsInterface,
	migrationsSettings *settings.MigrationsSettings,
	databaseSettings *settings.DatabaseSettings,
) *DatabaseClientFactory {
	return &DatabaseClientFactory{
		logger:             logger,
		dbQueryMetrics:     dbQueryMetrics,
		migrationsSettings: migrationsSettings,
		databaseSettings:   databaseSettings,
	}
}

func (factory *DatabaseClientFactory) CreateDatabaseClient() (db_clients.DatabaseClientInterface, error) {
	return conversation_broker_db.NewConversationBrokerDbClient(
		factory.databaseSettings,
		factory.logger,
		factory.migrationsSettings,
		factory.dbQueryMetrics,
	)
}
