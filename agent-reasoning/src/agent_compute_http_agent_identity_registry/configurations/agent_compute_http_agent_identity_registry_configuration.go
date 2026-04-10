package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	http_registry_clients "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_identity_registry/clients"
	http_registry_services "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_identity_registry/services"
	http_registry_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_identity_registry/settings"
	"go.uber.org/dig"
)

func AddAgentComputeHttpAgentIdentityRegistryConfiguration(container *dig.Container) {
	err := container.Provide(http_registry_settings.NewAgentIdentityRegistrySettings)
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
		http_registry_clients.NewAgentIdentityRegistryClient,
		dig.As(new(http_registry_clients.AgentIdentityRegistryClientInterface)),
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
