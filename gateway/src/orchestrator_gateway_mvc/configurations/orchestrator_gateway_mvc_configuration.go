package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_mvc/middlewares"
	"go.uber.org/dig"
)

func AddOrchestratorGatewayMvcConfiguration(container *dig.Container) {
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
		orchestrator_gateway_mvc.NewImportAgentController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_gateway_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_gateway_mvc.NewExportAgentController,
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
}
