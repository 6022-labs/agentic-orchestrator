package configurations

import (
	common_metrics "github.com/6022-labs/agentic-orchestrator/gateway/src/common/metrics"
	"net/http"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/clients"
	orchestrator_http_platform_bridge_services "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/settings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
)

func AddGatewayHttpChannelsAndIntegrationsBridgeConfiguration(container *dig.Container) {
	err := container.Provide(
		settings.NewPlatformBridgesSettings,
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		func(platformBridgesSettings *settings.PlatformBridgesSettings) (map[enums.BridgePlatform]clients.ClientWithResponsesInterface, error) {
			httpClient := &http.Client{
				Timeout: time.Duration(platformBridgesSettings.SlackBridge.RequestTimeoutSeconds) * time.Second,
				Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
					http.DefaultTransport,
					common_metrics.NewExternalHttpMetrics(
						"agentic_orchestrator_gateway_http_channels_and_integrations_bridge",
					),
				),
					otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
				),
			}
			slackClient, err := clients.NewClientWithResponses(
				platformBridgesSettings.SlackBridge.BaseUrl,
				clients.WithHTTPClient(httpClient),
			)
			if err != nil {
				return nil, err
			}
			return map[enums.BridgePlatform]clients.ClientWithResponsesInterface{
				enums.BridgePlatformSlack: slackClient,
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_http_platform_bridge_services.NewPlatformBridgeRequester,
		dig.As(new(interfaces.PlatformBridgeRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
