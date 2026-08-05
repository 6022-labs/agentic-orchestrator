package factories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	conversation_broker_db "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/settings"
	"go.uber.org/zap"
)

type DatabaseClientFactory struct {
	logger             *zap.Logger
	migrationsSettings *settings.MigrationsSettings
	databaseSettings   *settings.DatabaseSettings
}

func NewDatabaseClientFactory(
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
	databaseSettings *settings.DatabaseSettings,
) *DatabaseClientFactory {
	return &DatabaseClientFactory{
		logger:             logger,
		migrationsSettings: migrationsSettings,
		databaseSettings:   databaseSettings,
	}
}

func (factory *DatabaseClientFactory) CreateDatabaseClient() (db_clients.DatabaseClientInterface, error) {
	return conversation_broker_db.NewConversationBrokerDbClient(
		factory.databaseSettings.Url,
		factory.logger,
		factory.migrationsSettings,
	)
}
