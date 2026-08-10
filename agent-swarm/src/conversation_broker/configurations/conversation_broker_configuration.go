package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/prompts"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"go.uber.org/dig"
)

func AddConversationBrokerConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewProcessingSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewRegistryAgentResolutionSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewRemoteAgentResolutionSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewX402PaymentSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewWalletSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewX402PaymentPayloadSigner,
		dig.As(new(service_interfaces.X402PaymentPayloadSignerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewX402PaymentChallengeCache,
		dig.As(new(service_interfaces.X402PaymentChallengeCacheInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewRemoteAgentCardResolver,
		dig.As(new(service_interfaces.RemoteAgentCardResolverInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentUrlResolver,
		dig.As(new(service_interfaces.AgentUrlResolverInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewRegistryAgentResolver,
		dig.As(new(service_interfaces.RegistryAgentResolverInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewParticipantRoleResolver,
		dig.As(new(service_interfaces.ParticipantRoleResolverInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		prompts.NewConversationPromptBuilder,
		dig.As(new(service_interfaces.ConversationPromptBuilderInterface)),
	)
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
		services.NewEmojiDecoder,
		dig.As(new(service_interfaces.EmojiDecoderInterface)),
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

	err = container.Provide(use_cases.NewListConversations, dig.As(new(use_cases.ListConversationsInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewGetConversationMessages, dig.As(new(use_cases.GetConversationMessagesInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewListConversationParticipants, dig.As(new(use_cases.ListConversationParticipantsInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewGetAgentPromptTemplate, dig.As(new(use_cases.GetAgentPromptTemplateInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewSaveAgentPromptTemplate, dig.As(new(use_cases.SaveAgentPromptTemplateInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewGetRemoteAgentCard, dig.As(new(use_cases.GetRemoteAgentCardInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewResolveAgentUrl, dig.As(new(use_cases.ResolveAgentUrlInterface)))
	if err != nil {
		panic(err)
	}

	err = container.Provide(use_cases.NewUpdateConversationDescription, dig.As(new(use_cases.UpdateConversationDescriptionInterface)))
	if err != nil {
		panic(err)
	}
}
