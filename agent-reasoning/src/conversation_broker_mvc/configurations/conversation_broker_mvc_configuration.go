package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_mvc"
	"go.uber.org/dig"
)

func AddConversationBrokerMvcConfiguration(container *dig.Container) {
	err := container.Provide(
		conversation_broker_mvc.NewBridgesController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		conversation_broker_mvc.NewParticipantsController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		conversation_broker_mvc.NewConversationsController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		conversation_broker_mvc.NewPromptTemplatesController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
