package configurations

import (
	"net/http"
	"strings"
	"time"

	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/settings"
	"go.uber.org/dig"
)

func AddConversationBrokerHttpChannelsAndIntegrationsBridgeConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewBridgeCallbackSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(bridgeCallbackSettings *settings.BridgeCallbackSettings) clients.ClientFactory {
		httpClient := &http.Client{
			Timeout: time.Duration(bridgeCallbackSettings.RequestTimeoutSeconds) * time.Second,
		}
		return func(baseUrl string) (clients.ClientWithResponsesInterface, error) {
			return clients.NewClientWithResponses(strings.TrimRight(baseUrl, "/"), clients.WithHTTPClient(httpClient))
		}
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewBridgeCallbackRequester,
		dig.As(new(service_interfaces.BridgeCallbackRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
