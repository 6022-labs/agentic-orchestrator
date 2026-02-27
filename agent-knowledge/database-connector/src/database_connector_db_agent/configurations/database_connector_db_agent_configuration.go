package configurations

import (
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/factories/repositories"
	db_agent_factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/factories/db_clients"
	db_agent_factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/factories/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/settings"

	"go.uber.org/dig"
)

func AddDatabaseConnectorDbAgentConfiguration(container *dig.Container) {
	// DB client factories
	err := container.Provide(
		db_agent_factories_db_clients.NewAgentDatabaseClientFactory,
		dig.As(new(factories_db_clients.AgentDatabaseClientFactoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Settings
	err = container.Provide(
		settings.NewMigrationsSettings,
	)
	if err != nil {
		panic(err)
	}

	// Repository factories
	err = container.Provide(
		db_agent_factories_repositories.NewAgentContextRepositoryFactory,
		dig.As(new(factories_repositories.AgentContextRepositoryFactoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_agent_factories_repositories.NewAgentPromptTemplateRepositoryFactory,
		dig.As(new(factories_repositories.AgentPromptTemplateRepositoryFactoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_agent_factories_repositories.NewAgentContextMessageRepositoryFactory,
		dig.As(new(factories_repositories.AgentContextMessageRepositoryFactoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		db_agent_factories_repositories.NewAgentContextConversationMappingRepositoryFactory,
		dig.As(new(factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface)),
	)
	if err != nil {
		panic(err)
	}
}
