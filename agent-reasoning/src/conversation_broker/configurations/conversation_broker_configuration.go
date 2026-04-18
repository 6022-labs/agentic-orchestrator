package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services"
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"go.uber.org/dig"
)

func AddConversationBrokerConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewProcessingSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewConversationConclude,
		dig.As(new(service_interfaces.ConversationConcludeInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewConversationProcessor,
		dig.As(new(service_interfaces.ConversationProcessorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewConversationProcessorDaemon,
		dig.As(new(service_interfaces.ConversationProcessorDaemonInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewRegisterBridge, dig.As(new(use_cases.RegisterBridgeInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewRegisterAgent, dig.As(new(use_cases.RegisterAgentInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewRegisterUser, dig.As(new(use_cases.RegisterUserInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewRegisterConversation, dig.As(new(use_cases.RegisterConversationInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewAttachParticipant, dig.As(new(use_cases.AttachParticipantInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewDetachParticipant, dig.As(new(use_cases.DetachParticipantInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewIngestMessage, dig.As(new(use_cases.IngestMessageInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewForceConclude, dig.As(new(use_cases.ForceConcludeInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewGetConversation, dig.As(new(use_cases.GetConversationInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewUpdateConversationDescription, dig.As(new(use_cases.UpdateConversationDescriptionInterface)))
	if err != nil {
		panic(err)
	}
}
