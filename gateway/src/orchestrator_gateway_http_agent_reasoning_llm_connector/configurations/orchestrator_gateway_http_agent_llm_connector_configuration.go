package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	http_agent_llm_connector_clients "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_reasoning_llm_connector/clients"
	http_agent_llm_connector_services "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_reasoning_llm_connector/services"
	http_agent_llm_connector_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_reasoning_llm_connector/settings"
	"go.uber.org/dig"
)

func AddOrchestratorGatewayHttpAgentReasoningLlmConnectorConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		http_agent_llm_connector_settings.NewAgentLLMConnectorSettings,
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		http_agent_llm_connector_clients.NewAgentLLMConnectorClient,
		dig.As(new(http_agent_llm_connector_clients.AgentLLMConnectorClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		http_agent_llm_connector_services.NewLLMConnectionRequester,
		dig.As(new(interfaces.LLMConnectionRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_llm_connector_services.NewMCPConnectionRequester,
		dig.As(new(interfaces.MCPConnectionRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_agent_llm_connector_services.NewLLMConnectionCompletionRequester,
		dig.As(new(interfaces.LLMConnectionChatCompletionRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("llm-connector-http-client"))
	if err != nil {
		panic(err)
	}
}
