package configurations

import (
	common_metrics "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics"
	common_metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics/interfaces"
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/mvc/middlewares"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_mvc"
	"go.uber.org/dig"
)

func AddSlackBridgeInternalMvcConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(func() common_metrics_interfaces.ApiRequestMetricsInterface {
		return common_metrics.NewApiRequestMetrics("agentic_orchestrator_slack_bridge_internal_api")
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewApiRequestMetricsMiddleware); err != nil {
		panic(err)
	}

	// Controllers
	err := container.Provide(
		slack_bridge_internal_mvc.NewWalletController,
		dig.Group("controllers"),
		dig.As(new(slack_bridge_internal_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		slack_bridge_internal_mvc.NewBrokerCallbackController,
		dig.Group("controllers"),
		dig.As(new(slack_bridge_internal_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
