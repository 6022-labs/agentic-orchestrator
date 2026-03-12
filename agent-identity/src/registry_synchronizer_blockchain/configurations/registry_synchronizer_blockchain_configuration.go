package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

func AddRegistrySynchronizerBlockchainConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewRpcNodeSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewAgentCollectionsManagerSettings)
	if err != nil {
		panic(err)
	}

	// Services (requesters)
	err = container.Provide(
		services.NewAgentCollectionRequester,
		dig.As(new(interfaces.AgentCollectionBlockchainRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewOrchestratorDomainResolver,
		dig.As(new(interfaces.OrchestratorDomainResolverInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentCollectionsManagerRequester,
		dig.As(new(interfaces.AgentCollectionsManagerBlockchainRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services (subscription providers)
	err = container.Provide(
		services.NewAgentCollectionEventSubscriptionProvider,
		dig.As(new(interfaces.AgentCollectionEventSubscriptionProviderInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentCollectionsManagerEventSubscriptionProvider,
		dig.As(new(interfaces.AgentCollectionsManagerEventSubscriptionProviderInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Ethereum clients
	err = container.Provide(func(rpcNodeSettings *settings.RpcNodeSettings) (*ethclient.Client, error) {
		return ethclient.Dial(rpcNodeSettings.WsUrl)
	}, dig.Name("ws"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(rpcNodeSettings *settings.RpcNodeSettings) (*ethclient.Client, error) {
		return ethclient.Dial(rpcNodeSettings.HttpUrl)
	}, dig.Name("http"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(rpcNodeSettings *settings.RpcNodeSettings) (*ethclient.Client, error) {
		return ethclient.Dial(rpcNodeSettings.ENSHttpUrl)
	}, dig.Name("ens"))
	if err != nil {
		panic(err)
	}
}
