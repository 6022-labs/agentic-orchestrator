package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	http_agent_compute_clients "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_agent_compute/clients"
	http_agent_compute_services "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_agent_compute/services"
	http_agent_compute_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_agent_compute/settings"
	"go.uber.org/dig"
)

func AddGatewayHttpAgentReasoningAgentComputeConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		http_agent_compute_settings.NewAgentReasoningAgentComputeSettings,
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		http_agent_compute_clients.NewAgentReasoningComputeClient,
		dig.As(new(http_agent_compute_clients.AgentReasoningComputeClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		http_agent_compute_services.NewLLMConnectionRequester,
		dig.As(new(interfaces.LLMConnectionRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_compute_services.NewMCPConnectionRequester,
		dig.As(new(interfaces.MCPConnectionRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("compute-http-client"))
	if err != nil {
		panic(err)
	}
}
