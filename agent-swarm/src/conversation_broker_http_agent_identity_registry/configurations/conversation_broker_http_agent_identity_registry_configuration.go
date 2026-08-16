package configurations

import (
	"net/http"
	"time"

	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/settings"
	"go.uber.org/dig"
)

func AddConversationBrokerHttpAgentIdentityRegistryConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewAgentIdentityRegistrySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		func(agentIdentityRegistrySettings *settings.AgentIdentityRegistrySettings) (clients.ClientWithResponsesInterface, error) {
			httpClient := &http.Client{
				Timeout: time.Duration(agentIdentityRegistrySettings.RequestTimeoutSeconds) * time.Second,
			}
			return clients.NewClientWithResponses(
				agentIdentityRegistrySettings.BaseUrl,
				clients.WithHTTPClient(httpClient),
			)
		},
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
