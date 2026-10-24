package configurations

import (
	common_metrics "github.com/6022-labs/agentic-orchestrator/gateway/src/common/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/common/metrics/interfaces"
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc/middlewares"
	orchestrator_gateway_mvc "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_mvc/middlewares"
	"go.uber.org/dig"
)

func AddGatewayMvcConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(func() metrics_interfaces.ApiRequestMetricsInterface {
		return common_metrics.NewApiRequestMetrics("agentic_orchestrator_gateway_api")
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewApiRequestMetricsMiddleware); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewRequestFailureLogMiddleware); err != nil {
		panic(err)
	}

	// Middlewares
	err := container.Provide(
		middlewares.NewJWTAuthenticationInjectionMiddleware,
		dig.Group("auth-injection-middlewares"),
		dig.As(new(middlewares.AuthInjectionMiddlewareInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		middlewares.NewSignatureAuthenticationInjectionMiddleware,
		dig.Group("auth-injection-middlewares"),
		dig.As(new(middlewares.AuthInjectionMiddlewareInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		middlewares.NewAuthenticationMandatoryMiddleware,
	)
	if err != nil {
		panic(err)
	}

	// Controllers
	err = container.Provide(
		orchestrator_gateway_mvc.NewAuthController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_gateway_mvc.NewPingController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_gateway_mvc.NewAgentController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_gateway_mvc.NewAgentCollectionController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_gateway_mvc.NewWalletController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_gateway_mvc.NewConversationController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_gateway_mvc.NewConfigController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
