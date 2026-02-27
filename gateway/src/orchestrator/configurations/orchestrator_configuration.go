package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/auth/signature"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/settings"
	"go.uber.org/dig"
)

func AddOrchestratorConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewBlockchainSettings)
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
}
