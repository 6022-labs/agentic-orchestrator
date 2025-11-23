package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	external_interfaces "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/clients"
	http_orchestrator_services "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/settings"
	"go.uber.org/dig"
)

func AddSlackBridgeHttpOrchestratorConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewOrchestratorSettings)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		clients.NewOrchestratorGatewayClient,
		dig.As(new(clients.OrchestratorGatewayClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		http_orchestrator_services.NewAgentRequester,
		dig.As(new(external_interfaces.AgentRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_orchestrator_services.NewAgentCompletionsRequester,
		dig.As(new(interfaces.AgentCompletionsRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("orchestrator-main-api-http-client"))
	if err != nil {
		panic(err)
	}
}
