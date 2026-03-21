package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_platform_bridge/clients"
	orchestrator_http_platform_bridge_services "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_platform_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_platform_bridge/settings"
	"go.uber.org/dig"
)

func AddGatewayHttpPlatformBridgeConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		settings.NewPlatformBridgesSettings,
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		orchestrator_http_platform_bridge_services.NewPlatformBridgeRequester,
		dig.As(new(interfaces.PlatformBridgeRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		clients.NewPlatformBridgeClient,
		dig.As(new(clients.PlatformBridgeClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("platform-bridge-http-client"))
	if err != nil {
		panic(err)
	}
}
