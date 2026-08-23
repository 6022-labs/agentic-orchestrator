package configurations

import (
	"net/http"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services/interfaces"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
)

func AddSlackBridgeHttpAgentSwarmConversationBrokerConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewConversationBrokerSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		func(conversationBrokerSettings *settings.ConversationBrokerSettings) (clients.ClientWithResponsesInterface, error) {
			httpClient := &http.Client{
				Timeout: time.Duration(conversationBrokerSettings.RequestTimeoutSeconds) * time.Second,
				Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
					http.DefaultTransport,
					common_metrics.NewExternalHttpMetrics(
						"agentic_orchestrator_slack_bridge_http_agent_swarm_conversation_broker",
					),
				),
					otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
				),
			}
			return clients.NewClientWithResponses(
				conversationBrokerSettings.BaseUrl,
				clients.WithHTTPClient(httpClient),
			)
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
