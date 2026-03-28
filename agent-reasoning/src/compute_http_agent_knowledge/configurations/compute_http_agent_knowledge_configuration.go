package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
	http_agent_knowledge_clients "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/clients"
	http_agent_knowledge_services "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/services"
	http_agent_knowledge_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/settings"
	"go.uber.org/dig"
)

func AddComputeHttpAgentKnowledgeConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		http_agent_knowledge_settings.NewAgentKnowledgeSettings,
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		http_agent_knowledge_clients.NewAgentKnowledgeClient,
		dig.As(new(http_agent_knowledge_clients.AgentKnowledgeClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		http_agent_knowledge_services.NewAgentContextRequester,
		dig.As(new(interfaces.AgentContextRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_knowledge_services.NewAgentContextMessageRequester,
		dig.As(new(interfaces.AgentContextMessageRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_knowledge_services.NewAgentContextConversationMappingRequester,
		dig.As(new(interfaces.AgentContextConversationMappingRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_knowledge_services.NewAgentPromptTemplateRequester,
		dig.As(new(interfaces.AgentPromptTemplateRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("agent-knowledge-http-client"))
	if err != nil {
		panic(err)
	}
}
