package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	http_compute_clients "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_http_agent_reasoning_agent_compute/clients"
	http_compute_services "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_http_agent_reasoning_agent_compute/services"
	http_compute_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_http_agent_reasoning_agent_compute/settings"
	"go.uber.org/dig"
)

func AddSlackBridgeInternalHttpAgentReasoningAgentComputeConfiguration(container *dig.Container) {
	err := container.Provide(http_compute_settings.NewAgentReasoningAgentComputeSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_compute_clients.NewAgentReasoningComputeClient,
		dig.As(new(http_compute_clients.AgentReasoningComputeClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_compute_services.NewAgentCompletionsRequester,
		dig.As(new(interfaces.AgentCompletionsRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("agent-reasoning-compute-http-client"))
	if err != nil {
		panic(err)
	}
}
