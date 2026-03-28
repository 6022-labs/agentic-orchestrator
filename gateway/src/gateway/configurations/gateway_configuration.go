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

	err = container.Provide(settings.NewBlockchainSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewOrchestratorSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewAgentCollectionSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewDecentralizedDomainNameProviderSettings)
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
		services.NewWalletCreator,
		dig.As(new(interfaces.WalletCreatorInterface)),
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
		services.NewMintSigner,
		dig.As(new(interfaces.MintSignerInterface)),
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

	err = container.Provide(
		services.NewAgentPrivateKeyStorage,
		dig.As(new(interfaces.AgentPrivateKeyStorageInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewOrchestratorPrivateKeyFetcher,
		dig.As(new(interfaces.OrchestratorPrivateKeyFetcherInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewDecentralizedDomainNameBuilder,
		dig.As(new(interfaces.DecentralizedDomainNameBuilderInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewEIP712SignatureHelper,
		dig.As(new(interfaces.EIP712SignatureHelperInterface)),
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
		use_cases.NewCreateAgentFirstStep,
		dig.As(new(use_cases.CreateAgentFirstStepInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewCreateAgentSecondStep,
		dig.As(new(use_cases.CreateAgentSecondStepInterface)),
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
		use_cases.NewSaveAgentPromptTemplate,
		dig.As(new(use_cases.SaveAgentPromptTemplateInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewAddMessageToAgentContext,
		dig.As(new(use_cases.AddMessageToAgentContextInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewCreateAgentContext,
		dig.As(new(use_cases.CreateAgentContextInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListAgentContextKeys,
		dig.As(new(use_cases.ListAgentContextKeysInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSetupAgentLLMSettings,
		dig.As(new(use_cases.SetupAgentLLMSettingsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSetupAgentDatabaseSettingsFirstStep,
		dig.As(new(use_cases.SetupAgentDatabaseSettingsFirstStepInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSetupAgentDatabaseSettingsSecondStep,
		dig.As(new(use_cases.SetupAgentDatabaseSettingsSecondStepInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentLLMSettings,
		dig.As(new(use_cases.GetAgentLLMSettingsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewExportAgent,
		dig.As(new(use_cases.ExportAgentInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewIsInternalAgent,
		dig.As(new(use_cases.IsInternalAgentInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewCreateAgentMCPConnection,
		dig.As(new(use_cases.CreateAgentMCPConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewRemoveAgentMCPConnection,
		dig.As(new(use_cases.RemoveAgentMCPConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListAgentMCPConnections,
		dig.As(new(use_cases.ListAgentMCPConnectionsInterface)),
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
