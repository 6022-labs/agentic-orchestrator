package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/db_clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	infra_factories "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/factories"
	infra_repositories "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddOrchestratorDbConfiguration(container *dig.Container) {
	// Factories
	err := container.Provide(infra_factories.NewDatabaseClientFactory)
	if err != nil {
		panic(err)
	}

	// Repositories
	err = container.Provide(
		infra_repositories.NewAgentRepository,
		dig.As(new(repositories.AgentRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		infra_repositories.NewAgentCollectionRepository,
		dig.As(new(repositories.AgentCollectionRepositoryInterface)),
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
	err = container.Provide(func(databaseClientFactory *infra_factories.DatabaseClientFactory, logger *zap.Logger) db_clients.DatabaseClientInterface {
		dbClient, err := databaseClientFactory.CreateDatabaseClient()
		if err != nil {
			logger.Fatal("Failed to create orchestrator database client", zap.Error(err))
		}

		return dbClient
	})
	if err != nil {
		panic(err)
	}
}
