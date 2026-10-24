package configurations

import (
	common_metrics "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics"
	common_metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/metrics/interfaces"
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/mvc/middlewares"
	orchestrator_slack_bridge_external_mvc "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc/middlewares"
	"go.uber.org/dig"
)

func AddSlackBridgeExternalMvcConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(func() common_metrics_interfaces.ApiRequestMetricsInterface {
		return common_metrics.NewApiRequestMetrics("agentic_orchestrator_slack_bridge_external_api")
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewApiRequestMetricsMiddleware); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewRequestFailureLogMiddleware); err != nil {
		panic(err)
	}

	// Controllers
	err := container.Provide(
		orchestrator_slack_bridge_external_mvc.NewEventHandlerController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_slack_bridge_external_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_slack_bridge_external_mvc.NewCommandHandlerController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_slack_bridge_external_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Middlewares
	err = container.Provide(
		middlewares.NewSlackSignatureMiddleware,
	)
	if err != nil {
		panic(err)
	}
}
