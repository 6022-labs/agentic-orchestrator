package configurations

import (
	"net/http"

	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_identity_registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_identity_registry/settings"
	"go.uber.org/dig"
)

func AddConversationBrokerHttpAgentIdentityRegistryConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewAgentIdentityRegistrySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("agent-identity-registry-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		clients.NewAgentIdentityRegistryClient,
		dig.As(new(clients.AgentIdentityRegistryClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentRegistryRequester,
		dig.As(new(service_interfaces.AgentRegistryRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
