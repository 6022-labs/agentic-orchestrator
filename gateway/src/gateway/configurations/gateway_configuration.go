package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/signature"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"go.uber.org/dig"
)

func AddGatewayConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(
		metrics.NewGatewayMetrics,
		dig.As(new(metrics_interfaces.GatewayMetricsInterface)),
	); err != nil {
		panic(err)
	}

	// Settings
	err := container.Provide(settings.NewJwtSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewGatewaySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewRateLimitSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewIpfsGatewaysSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewBlockchainConfigSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewSiweSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewWalletSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		signature.NewECDSASigner,
		dig.As(new(interfaces.ECDSASignerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		jwt.NewJwtManager,
		dig.As(new(interfaces.JwtManagerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSignatureVerifier,
		dig.As(new(interfaces.SignatureVerifierInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSignedRequestMessageBuilder,
		dig.As(new(interfaces.SignedRequestMessageBuilderInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSiweMessageBuilder,
		dig.As(new(interfaces.SiweMessageBuilderInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentAccessVerifier,
		dig.As(new(interfaces.AgentAccessVerifierInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewConversationAccessVerifier,
		dig.As(new(interfaces.ConversationAccessVerifierInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		auth.NewNonceStore,
		dig.As(new(interfaces.NonceStoreInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Use cases
	err = container.Provide(
		use_cases.NewGetAuthChallenge,
		dig.As(new(use_cases.GetAuthChallengeInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewAuthenticateSignedRequest,
		dig.As(new(use_cases.AuthenticateSignedRequestInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewVerifyAuth,
		dig.As(new(use_cases.VerifyAuthInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgent,
		dig.As(new(use_cases.GetAgentInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListAgents,
		dig.As(new(use_cases.ListAgentsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListAgentCollections,
		dig.As(new(use_cases.ListAgentCollectionsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewVerifyWalletForBridge,
		dig.As(new(use_cases.VerifyWalletForBridgeInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetConfig,
		dig.As(new(use_cases.GetConfigInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentPromptTemplate,
		dig.As(new(use_cases.GetAgentPromptTemplateInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSaveAgentPromptTemplate,
		dig.As(new(use_cases.SaveAgentPromptTemplateInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentCard,
		dig.As(new(use_cases.GetAgentCardInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentUrl,
		dig.As(new(use_cases.GetAgentUrlInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListConversations,
		dig.As(new(use_cases.ListConversationsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetConversation,
		dig.As(new(use_cases.GetConversationInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetConversationMessages,
		dig.As(new(use_cases.GetConversationMessagesInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListConversationParticipants,
		dig.As(new(use_cases.ListConversationParticipantsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewForceConcludeConversation,
		dig.As(new(use_cases.ForceConcludeConversationInterface)),
	)
	if err != nil {
		panic(err)
	}

}
