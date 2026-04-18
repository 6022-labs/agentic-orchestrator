package configurations

import (
	"net/http"
	"time"

	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/settings"
	"go.uber.org/dig"
)

func AddConversationBrokerHttpChannelsAndIntegrationsBridgeConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewBridgeCallbackSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(s *settings.BridgeCallbackSettings) *http.Client {
		return &http.Client{
			Timeout: time.Duration(s.RequestTimeoutSeconds) * time.Second,
		}
	}, dig.Name("bridge-callback-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		clients.NewBridgeCallbackClient,
		dig.As(new(clients.BridgeCallbackClientInterface)),
	)
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
