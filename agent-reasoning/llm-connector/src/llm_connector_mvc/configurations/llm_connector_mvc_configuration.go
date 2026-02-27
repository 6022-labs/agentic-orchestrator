package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_mvc"
	"go.uber.org/dig"
)

func AddLLMConnectorMvcConfiguration(container *dig.Container) {
	err := container.Provide(
		llm_connector_mvc.NewChatCompletionController,
		dig.Group("controllers"),
		dig.As(new(llm_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		llm_connector_mvc.NewLLMConnectionController,
		dig.Group("controllers"),
		dig.As(new(llm_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		llm_connector_mvc.NewMCPConnectionController,
		dig.Group("controllers"),
		dig.As(new(llm_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
