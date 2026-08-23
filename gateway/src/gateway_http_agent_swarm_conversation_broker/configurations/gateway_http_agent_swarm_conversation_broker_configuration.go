package configurations

import (
	"net/http"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/gateway/src/common/metrics"
	gateway_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/settings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
)

// newConversationBrokerClient builds the generated broker client with the configured timeout.
func newConversationBrokerClient(
	conversationBrokerSettings *settings.ConversationBrokerSettings,
) (clients.ClientWithResponsesInterface, error) {
	httpClient := &http.Client{Timeout: time.Duration(conversationBrokerSettings.RequestTimeoutSeconds) * time.Second,
		Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
			http.DefaultTransport,
			common_metrics.NewExternalHttpMetrics("agentic_orchestrator_gateway_http_agent_swarm_conversation_broker"),
		),
			otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
		),
	}
	return clients.NewClientWithResponses(conversationBrokerSettings.BaseUrl, clients.WithHTTPClient(httpClient))
}

func AddGatewayHttpAgentSwarmConversationBrokerConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewConversationBrokerSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(newConversationBrokerClient)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewConversationBrokerRequester,
		dig.As(new(gateway_interfaces.ConversationBrokerRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
