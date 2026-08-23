package configurations

import (
	"net/http"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/metrics"
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/settings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
)

func AddConversationBrokerHttpRemoteAgentConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(
		metrics.NewRemoteAgentMetrics,
		dig.As(new(metrics_interfaces.RemoteAgentMetricsInterface)),
	); err != nil {
		panic(err)
	}

	err := container.Provide(settings.NewRemoteAgentSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(remoteAgentSettings *settings.RemoteAgentSettings) *http.Client {
		return &http.Client{
			Timeout: time.Duration(remoteAgentSettings.RequestTimeoutSeconds) * time.Second,
			Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
				http.DefaultTransport,
				common_metrics.NewExternalHttpMetrics("agentic_orchestrator_conversation_broker_http_remote_agent"),
			),
				otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
			),
		}
	}, dig.Name("agent-node-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(remoteAgentSettings *settings.RemoteAgentSettings) *http.Client {
		return &http.Client{
			Timeout: time.Duration(remoteAgentSettings.WellKnownRequestTimeoutSeconds) * time.Second,
			Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
				http.DefaultTransport,
				common_metrics.NewExternalHttpMetrics(
					"agentic_orchestrator_conversation_broker_http_remote_agent_well_known",
				),
			),
				otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
			),
		}
	}, dig.Name("agent-node-well-known-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		clients.NewRemoteAgentClient,
		dig.As(new(clients.RemoteAgentClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewRemoteAgentCardRequester,
		dig.As(new(service_interfaces.RemoteAgentCardRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewRemoteAgentCompletionsRequester,
		dig.As(new(service_interfaces.RemoteAgentCompletionsRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
