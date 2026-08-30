package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	domain_repositories "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/factories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/metrics"
	db_metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/metrics/interfaces"
	db_repositories "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/settings"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddConversationBrokerDbConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(
		func() *metrics.DbQueryMetrics { return metrics.NewDbQueryMetrics(semconv.DBSystemNamePostgreSQL) },
		dig.As(new(db_metrics_interfaces.DbQueryMetricsInterface)),
	); err != nil {
		panic(err)
	}

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

	err = container.Provide(
		func(factory *factories.DatabaseClientFactory, logger *zap.Logger) db_clients.DatabaseClientInterface {
			dbClient, err := factory.CreateDatabaseClient()
			if err != nil {
				logger.Fatal("Failed to create conversation broker database client", zap.Error(err))
			}
			return dbClient
		},
	)
	if err != nil {
		panic(err)
	}

	// Repositories
	err = container.Provide(
		db_repositories.NewBridgeRepository,
		dig.As(new(domain_repositories.BridgeRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewParticipantRepository,
		dig.As(new(domain_repositories.ParticipantRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewConversationRepository,
		dig.As(new(domain_repositories.ConversationRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewConversationParticipantRepository,
		dig.As(new(domain_repositories.ConversationParticipantRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewConversationMemoryBindingRepository,
		dig.As(new(domain_repositories.ConversationMemoryBindingRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewConversationMessageRepository,
		dig.As(new(domain_repositories.ConversationMessageRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_repositories.NewAgentPromptTemplateRepository,
		dig.As(new(domain_repositories.AgentPromptTemplateRepositoryInterface)),
	)
	if err != nil {
		panic(err)
	}
}
