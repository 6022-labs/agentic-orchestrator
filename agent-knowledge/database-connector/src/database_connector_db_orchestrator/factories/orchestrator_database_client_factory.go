package factories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	orchestrator_agent_database_connector_db_orchestrator "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"go.uber.org/zap"
)

type OrchestratorDatabaseClientFactory struct {
	logger                       *zap.Logger
	migrationsSettings           *settings.MigrationsSettings
	orchestratorDatabaseSettings *settings.DatabaseSettings
}

func NewOrchestratorDatabaseClientFactory(
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
	orchestratorDatabaseSettings *settings.DatabaseSettings,
) *OrchestratorDatabaseClientFactory {
	return &OrchestratorDatabaseClientFactory{
		logger:                       logger,
		migrationsSettings:           migrationsSettings,
		orchestratorDatabaseSettings: orchestratorDatabaseSettings,
	}
}

func (factory *OrchestratorDatabaseClientFactory) Create() (db_clients.OrchestratorDatabaseClientInterface, error) {
	return orchestrator_agent_database_connector_db_orchestrator.NewOrchestratorDatabaseClient(
		factory.orchestratorDatabaseSettings.Url,
		factory.logger,
		factory.migrationsSettings,
	)
}
