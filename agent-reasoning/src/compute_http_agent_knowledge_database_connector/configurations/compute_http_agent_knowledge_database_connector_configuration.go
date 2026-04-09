package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
	http_agent_database_connector_clients "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge_database_connector/clients"
	http_agent_database_connector_services "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge_database_connector/services"
	http_agent_database_connector_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge_database_connector/settings"
	"go.uber.org/dig"
)

func AddComputeHttpAgentKnowledgeDatabaseConnectorConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		http_agent_database_connector_settings.NewAgentKnowledgeDatabaseConnectorSettings,
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		http_agent_database_connector_clients.NewAgentDatabaseConnectorClient,
		dig.As(new(http_agent_database_connector_clients.AgentDatabaseConnectorClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		http_agent_database_connector_services.NewAgentContextRequester,
		dig.As(new(interfaces.AgentContextRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_database_connector_services.NewAgentContextMessageRequester,
		dig.As(new(interfaces.AgentContextMessageRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_database_connector_services.NewAgentContextConversationMappingRequester,
		dig.As(new(interfaces.AgentContextConversationMappingRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_database_connector_services.NewAgentPromptTemplateRequester,
		dig.As(new(interfaces.AgentPromptTemplateRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("database-connector-http-client"))
	if err != nil {
		panic(err)
	}
}
