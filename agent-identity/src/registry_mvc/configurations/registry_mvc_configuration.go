package configurations

import (
	common_metrics "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/metrics/interfaces"
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/mvc/middlewares"
	registry_mvc "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_mvc"
	"go.uber.org/dig"
)

func AddRegistryMvcConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(func() metrics_interfaces.ApiRequestMetricsInterface {
		return common_metrics.NewApiRequestMetrics("agentic_orchestrator_registry_api")
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewApiRequestMetricsMiddleware); err != nil {
		panic(err)
	}

	err := container.Provide(
		registry_mvc.NewAgentController,
		dig.Group("controllers"),
		dig.As(new(registry_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		registry_mvc.NewAgentCollectionController,
		dig.Group("controllers"),
		dig.As(new(registry_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
