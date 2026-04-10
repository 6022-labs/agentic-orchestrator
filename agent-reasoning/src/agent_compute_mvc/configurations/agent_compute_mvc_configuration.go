package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_mvc"
	"go.uber.org/dig"
)

func AddAgentComputeMvcConfiguration(container *dig.Container) {
	err := container.Provide(
		agent_compute_mvc.NewChatCompletionController,
		dig.Group("controllers"),
		dig.As(new(agent_compute_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		agent_compute_mvc.NewLLMConnectionController,
		dig.Group("controllers"),
		dig.As(new(agent_compute_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		agent_compute_mvc.NewMCPConnectionController,
		dig.Group("controllers"),
		dig.As(new(agent_compute_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		agent_compute_mvc.NewAgentComputeController,
		dig.Group("controllers"),
		dig.As(new(agent_compute_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
