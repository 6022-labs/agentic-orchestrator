package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/factories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/settings"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddConversationBrokerDbConfiguration(container *dig.Container) {
	err := container.Provide(factories.NewDatabaseClientFactory)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewConversationBrokerDatabaseSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewMigrationsSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(factory *factories.DatabaseClientFactory, logger *zap.Logger) db_clients.DatabaseClientInterface {
		dbClient, err := factory.CreateDatabaseClient()
		if err != nil {
			logger.Fatal("Failed to create conversation broker database client", zap.Error(err))
		}
		return dbClient
	})
	if err != nil {
		panic(err)
	}

	// Repository registrations land in C3.
}
