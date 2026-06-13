package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/images"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"go.uber.org/dig"
)

func AddRegistryConfiguration(container *dig.Container) {
	// Settings
	if err := container.Provide(settings.NewImageGenerationSettings); err != nil {
		panic(err)
	}

	// Services
	err := container.Provide(
		services.NewWalletCreator,
		dig.As(new(interfaces.WalletCreatorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		images.NewDefaultImageGenerator,
		dig.As(new(interfaces.DefaultImageGeneratorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		images.NewIconImageGenerator,
		dig.As(new(interfaces.IconImageGeneratorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		images.NewGenerateAgentImages,
		dig.As(new(interfaces.GenerateAgentImagesInterface)),
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
		services.NewEIP712SignatureHelper,
		dig.As(new(interfaces.EIP712SignatureHelperInterface)),
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
		use_cases.NewIsInternalAgent,
		dig.As(new(use_cases.IsInternalAgentInterface)),
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
		use_cases.NewGetMintedAgentByEnsDomain,
		dig.As(new(use_cases.GetMintedAgentByEnsDomainInterface)),
	)
	if err != nil {
		panic(err)
	}
}
