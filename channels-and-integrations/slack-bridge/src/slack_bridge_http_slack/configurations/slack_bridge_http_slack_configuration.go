package configurations

import (
	"net/http"

	common_metrics "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	external_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	slack_bridge_http_slack_clients "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/clients"
	slack_bridge_http_slack_services "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/services"
	slack_bridge_http_slack_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/settings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
)

func AddSlackBridgeHttpSlackConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		slack_bridge_http_slack_settings.NewIpfsGatewaySettings,
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		slack_bridge_http_slack_clients.NewSlackClient,
		dig.As(new(slack_bridge_http_slack_clients.SlackClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		slack_bridge_http_slack_services.NewSlackUserRequester,
		dig.As(new(external_interfaces.SlackUserRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		slack_bridge_http_slack_services.NewSlackThreadRequester,
		dig.As(new(interfaces.SlackThreadRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		slack_bridge_http_slack_services.NewSlackChannelRequester,
		dig.As(new(external_interfaces.SlackChannelRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return &http.Client{
			Transport: otelhttp.NewTransport(common_metrics.NewHttpMetricsRoundTripper(
				http.DefaultTransport,
				common_metrics.NewExternalHttpMetrics("agentic_orchestrator_slack_bridge_http_slack"),
			),
				otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
			),
		}
	}, dig.Name("slack-http-client"))
	if err != nil {
		panic(err)
	}
}
