package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	infra_factories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/factories"
	infra_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	synchronizer_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/repositories"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddRegistryCommonDbConfiguration(container *dig.Container) {
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
		infra_repositories.NewAgentAddressRepository,
		dig.As(new(repositories.AgentAddressRepositoryInterface)),
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

	err = container.Provide(
		infra_repositories.NewSyncCheckpointRepository,
		dig.As(new(synchronizer_repositories.SyncCheckpointRepositoryInterface)),
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
	err = container.Provide(
		func(databaseClientFactory *infra_factories.DatabaseClientFactory, logger *zap.Logger) db_clients.DatabaseClientInterface {
			dbClient, err := databaseClientFactory.CreateDatabaseClient()
			if err != nil {
				logger.Fatal("Failed to create orchestrator database client", zap.Error(err))
			}

			return dbClient
		},
	)
	if err != nil {
		panic(err)
	}
}
