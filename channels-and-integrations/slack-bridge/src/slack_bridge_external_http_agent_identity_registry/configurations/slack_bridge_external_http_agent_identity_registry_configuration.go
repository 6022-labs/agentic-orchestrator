package configurations

import (
	"net/http"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics"
	external_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	http_registry_clients "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/clients"
	http_registry_services "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/services"
	http_registry_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/settings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
)

func AddSlackBridgeExternalHttpAgentIdentityRegistryConfiguration(container *dig.Container) {
	err := container.Provide(http_registry_settings.NewAgentIdentityRegistrySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		func(agentIdentityRegistrySettings *http_registry_settings.AgentIdentityRegistrySettings) (http_registry_clients.ClientWithResponsesInterface, error) {
			httpClient := &http.Client{
				Timeout: time.Duration(agentIdentityRegistrySettings.RequestTimeoutSeconds) * time.Second,
				Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
					http.DefaultTransport,
					common_metrics.NewExternalHttpMetrics(
						"agentic_orchestrator_slack_bridge_http_agent_identity_registry",
					),
				),
					otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
				),
			}
			return http_registry_clients.NewClientWithResponses(
				agentIdentityRegistrySettings.BaseUrl,
				http_registry_clients.WithHTTPClient(httpClient),
			)
		},
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_registry_services.NewAgentRegistryRequester,
		dig.As(new(external_interfaces.AgentRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
