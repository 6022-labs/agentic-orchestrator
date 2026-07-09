package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/signature"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"go.uber.org/dig"
)

func AddGatewayConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewJwtSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewRateLimitSettings)
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
		services.NewSignedMessageValidator,
		dig.As(new(interfaces.SignedMessageValidatorInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Auth
	err = container.Provide(
		jwt.NewNonceGenerator,
		dig.As(new(interfaces.NonceGeneratorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		jwt.NewSIWEDetailsExtractor,
		dig.As(new(interfaces.SIWEDetailsExtractorInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Use cases
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
		use_cases.NewInitializeJwtAuthentication,
		dig.As(new(use_cases.InitializeJwtAuthenticationInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewFinishJwtAuthentication,
		dig.As(new(use_cases.FinishJwtAuthenticationInterface)),
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

}
