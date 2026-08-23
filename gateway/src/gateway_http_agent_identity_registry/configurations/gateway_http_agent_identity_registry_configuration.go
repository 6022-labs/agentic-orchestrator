package configurations

import (
	"net/http"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/gateway/src/common/metrics"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	http_registry_clients "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	http_registry_services "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/services"
	http_registry_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/settings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
)

func AddGatewayHttpAgentIdentityRegistryConfiguration(container *dig.Container) {
	err := container.Provide(http_registry_settings.NewAgentIdentityRegistrySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		func(settings *http_registry_settings.AgentIdentityRegistrySettings) (http_registry_clients.ClientWithResponsesInterface, error) {
			httpClient := &http.Client{Timeout: time.Duration(settings.RequestTimeoutSeconds) * time.Second,
				Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
					http.DefaultTransport,
					common_metrics.NewExternalHttpMetrics("agentic_orchestrator_gateway_http_agent_identity_registry"),
				),
					otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
				),
			}
			return http_registry_clients.NewClientWithResponses(
				settings.BaseUrl,
				http_registry_clients.WithHTTPClient(httpClient),
			)
		},
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		http_registry_services.NewAgentRegistryRequester,
		dig.As(new(interfaces.AgentRegistryRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
