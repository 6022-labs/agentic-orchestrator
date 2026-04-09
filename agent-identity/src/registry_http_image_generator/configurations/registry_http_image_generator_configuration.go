package configurations

import (
	"net/http"

	registry_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_image_generator/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_image_generator/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_image_generator/settings"
	"go.uber.org/dig"
)

func AddRegistryHttpImageGeneratorConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewImagesGeneratorSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("images-generator-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		clients.NewImagesGeneratorClient,
		dig.As(new(clients.ImagesGeneratorClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewImagesGeneratorRequester,
		dig.As(new(registry_interfaces.ImagesGeneratorRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
