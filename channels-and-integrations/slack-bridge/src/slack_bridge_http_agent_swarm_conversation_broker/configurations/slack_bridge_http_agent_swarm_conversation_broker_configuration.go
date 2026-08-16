package configurations

import (
	"net/http"
	"time"

	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services/interfaces"
	"go.uber.org/dig"
)

func AddSlackBridgeHttpAgentSwarmConversationBrokerConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewConversationBrokerSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		func(conversationBrokerSettings *settings.ConversationBrokerSettings) (clients.ClientWithResponsesInterface, error) {
			httpClient := &http.Client{Timeout: time.Duration(conversationBrokerSettings.RequestTimeoutSeconds) * time.Second}
			return clients.NewClientWithResponses(conversationBrokerSettings.BaseUrl, clients.WithHTTPClient(httpClient))
		},
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewBridgeRegistrar,
		dig.As(new(interfaces.BridgeRegistrarInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewConversationBrokerRequester,
		dig.As(new(slack_bridge_interfaces.ConversationBrokerRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
