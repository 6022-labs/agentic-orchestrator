package configurations

import (
	orchestrator_images_generator_mvc "github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator_mvc"
	"go.uber.org/dig"
)

func AddImagesGeneratorMvcConfiguration(container *dig.Container) {
	// Controllers
	err := container.Provide(
		orchestrator_images_generator_mvc.NewImagesController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_images_generator_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
