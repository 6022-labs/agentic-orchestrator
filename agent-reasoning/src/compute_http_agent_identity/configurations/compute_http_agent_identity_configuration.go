package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
	http_agent_identity_clients "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_identity/clients"
	http_agent_identity_services "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_identity/services"
	http_agent_identity_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_identity/settings"
	"go.uber.org/dig"
)

func AddComputeHttpAgentIdentityConfiguration(container *dig.Container) {
	err := container.Provide(http_agent_identity_settings.NewAgentIdentityRegistrySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("agent-identity-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_identity_clients.NewAgentIdentityClient,
		dig.As(new(http_agent_identity_clients.AgentIdentityClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_identity_services.NewAgentRegistryRequester,
		dig.As(new(interfaces.AgentRegistryRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
