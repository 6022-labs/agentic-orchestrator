package configurations

import (
	registry_mvc "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_mvc"
	"go.uber.org/dig"
)

func AddRegistryMvcConfiguration(container *dig.Container) {
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
