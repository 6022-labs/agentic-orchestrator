package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/settings"
	"go.uber.org/dig"
)

func AddConversationBrokerHttpAgentReasoningAgentComputeConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewAgentReasoningAgentComputeSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("agent-compute-http-client"))
	if err != nil {
		panic(err)
	}

	// Client + requester service registrations land in C5.
}
