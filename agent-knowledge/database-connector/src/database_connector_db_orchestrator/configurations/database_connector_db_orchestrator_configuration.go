package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/factories/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator/factories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"

	db_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator/repositories"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddDatabaseConnectorDbOrchestratorConfiguration(container *dig.Container) {
	// Factories
	err := container.Provide(
		factories.NewOrchestratorDatabaseClientFactory,
		dig.As(new(factories_db_clients.OrchestratorDatabaseClientFactoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Repositories
	err = container.Provide(
		db_repositories.NewDatabaseConnectionRepository,
		dig.As(new(repositories.DatabaseConnectionRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewDatabaseConnectionMigrationRepository,
		dig.As(new(repositories.DatabaseConnectionMigrationRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Settings
	err = container.Provide(
		settings.NewDatabaseSettings,
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		settings.NewMigrationsSettings,
	)
	if err != nil {
		panic(err)
	}

	// Provide the orchestrator database client factory
	err = container.Provide(func(orchestratorDatabaseClientFactory factories_db_clients.OrchestratorDatabaseClientFactoryInterface, logger *zap.Logger) db_clients.OrchestratorDatabaseClientInterface {
		dbClient, err := orchestratorDatabaseClientFactory.Create()
		if err != nil {
			logger.Fatal("Failed to create orchestrator database client", zap.Error(err))
		}

		return dbClient
	})
	if err != nil {
		panic(err)
	}
}
