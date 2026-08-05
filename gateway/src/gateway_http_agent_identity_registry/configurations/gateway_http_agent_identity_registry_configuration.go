package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	http_registry_clients "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	http_registry_services "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/services"
	http_registry_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/settings"
	"go.uber.org/dig"
)

func AddGatewayHttpAgentIdentityRegistryConfiguration(container *dig.Container) {
	err := container.Provide(http_registry_settings.NewAgentIdentityRegistrySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		func(settings *http_registry_settings.AgentIdentityRegistrySettings) (http_registry_clients.ClientWithResponsesInterface, error) {
			return http_registry_clients.NewClientWithResponses(settings.BaseUrl)
		},
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_registry_services.NewAgentRegistryRequester,
		dig.As(new(interfaces.AgentRegistryRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
