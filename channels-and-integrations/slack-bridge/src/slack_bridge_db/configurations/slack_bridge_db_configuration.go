package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/factories"
	db_repositories "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddSlackBridgeDbConfiguration(container *dig.Container) {
	// Factories
	err := container.Provide(factories.NewDatabaseClientFactory)
	if err != nil {
		panic(err)
	}

	// Repositories
	err = container.Provide(
		db_repositories.NewSlackUserRepository,
		dig.As(new(repositories.SlackUserRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewSlackChannelRepository,
		dig.As(new(repositories.SlackChannelRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewSlackMessageRepository,
		dig.As(new(repositories.SlackMessageRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewSlackThreadRepository,
		dig.As(new(repositories.SlackThreadRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewSlackAgentRepository,
		dig.As(new(repositories.SlackAgentRepositoryInterface)),
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
	err = container.Provide(func(orchestratorSlackBridgeDBClientFactory *factories.DatabaseClientFactory, logger *zap.Logger) *slack_bridge_db.DatabaseClient {
		dbClient, err := orchestratorSlackBridgeDBClientFactory.CreateOrchestratorSlackBridgeDBClient()
		if err != nil {
			logger.Fatal("Failed to create orchestrator database client", zap.Error(err))
		}

		return dbClient
	})
	if err != nil {
		panic(err)
	}
}
