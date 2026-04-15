package configurations

import (
	"net/http"

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

	// Client + requester service registrations land in C5.
}
