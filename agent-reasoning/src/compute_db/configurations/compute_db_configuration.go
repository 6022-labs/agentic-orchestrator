package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/factories"
	db_repositories "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddComputeDbConfiguration(container *dig.Container) {
	// Factories
	err := container.Provide(factories.NewDatabaseClientFactory)
	if err != nil {
		panic(err)
	}

	// Repositories
	err = container.Provide(
		db_repositories.NewLLMConnectionRepository,
		dig.As(new(repositories.LLMConnectionRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewMCPConnectionRepository,
		dig.As(new(repositories.MCPConnectionRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Settings
	err = container.Provide(
		settings.NewOrchestratorDatabaseSettings,
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
	err = container.Provide(func(orchestratorDBClientFactory *factories.DatabaseClientFactory, logger *zap.Logger) db_clients.DatabaseClientInterface {
		dbClient, err := orchestratorDBClientFactory.CreateDatabaseClient()
		if err != nil {
			logger.Fatal("Failed to create orchestrator database client", zap.Error(err))
		}

		return dbClient
	})
	if err != nil {
		panic(err)
	}
}
