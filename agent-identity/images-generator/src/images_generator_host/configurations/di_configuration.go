package configurations

import (
	images_generator_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/configurations"
	images_generator_host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator_host/settings"
	images_generator_mvc_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator_mvc/configurations"

	"github.com/spf13/viper"
	"go.uber.org/dig"
)

func ConfigureDI() *dig.Container {
	container := dig.New()

	container.Provide(images_generator_host_settings.NewHostSettings)
	container.Provide(viper.GetViper)

	images_generator_configurations.AddImagesGeneratorConfiguration(container)
	images_generator_mvc_configurations.AddImagesGeneratorMvcConfiguration(container)

	return container
}
