package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	orchestrator_http_images_generator_clients "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_images_generator/clients"
	orchestrator_http_images_generator_services "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_images_generator/services"
	orchestrator_http_images_generator_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_images_generator/settings"
	"go.uber.org/dig"
)

func AddGatewayHttpAgentIdentityImagesGeneratorConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		orchestrator_http_images_generator_settings.NewImagesGeneratorSettings,
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		orchestrator_http_images_generator_clients.NewImagesGeneratorClient,
		dig.As(new(orchestrator_http_images_generator_clients.ImagesGeneratorClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		orchestrator_http_images_generator_services.NewImagesGeneratorRequester,
		dig.As(new(interfaces.ImagesGeneratorRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("images-generator-http-client"))
	if err != nil {
		panic(err)
	}
}
