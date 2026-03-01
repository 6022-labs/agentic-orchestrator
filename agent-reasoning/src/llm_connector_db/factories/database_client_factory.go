package factories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/db_clients"
	orchestrator_agent_llm_connector_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db/settings"
	"go.uber.org/zap"
)

type DatabaseClientFactory struct {
	logger                       *zap.Logger
	migrationsSettings           *settings.MigrationsSettings
	orchestratorDatabaseSettings *settings.DatabaseSettings
}

func NewDatabaseClientFactory(
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
	orchestratorDatabaseSettings *settings.DatabaseSettings,
) *DatabaseClientFactory {
	return &DatabaseClientFactory{
		logger:                       logger,
		migrationsSettings:           migrationsSettings,
		orchestratorDatabaseSettings: orchestratorDatabaseSettings,
	}
}

func (factory *DatabaseClientFactory) CreateDatabaseClient() (db_clients.DatabaseClientInterface, error) {
	return orchestrator_agent_llm_connector_db.NewOrchestratorDbClient(
		factory.orchestratorDatabaseSettings.Url,
		factory.logger,
		factory.migrationsSettings,
	)
}
