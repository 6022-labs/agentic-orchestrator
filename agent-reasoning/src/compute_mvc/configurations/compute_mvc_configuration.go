package configurations

import (
	compute_mvc "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_mvc"
	"go.uber.org/dig"
)

func AddComputeMvcConfiguration(container *dig.Container) {
	err := container.Provide(
		compute_mvc.NewChatCompletionController,
		dig.Group("controllers"),
		dig.As(new(compute_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		compute_mvc.NewLLMConnectionController,
		dig.Group("controllers"),
		dig.As(new(compute_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		compute_mvc.NewMCPConnectionController,
		dig.Group("controllers"),
		dig.As(new(compute_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
