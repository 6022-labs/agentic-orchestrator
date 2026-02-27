package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	orchestrator_blockchain_services "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_blockchain/services"
	orchestrator_blockchain_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_blockchain/settings"
	orchestrator_blockchain_subscribers "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_blockchain/subscribers"
	sc_sync_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/services/interfaces"
	orchestrator_sc_sync_subscribers "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/subscribers"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

func AddOrchestratorBlockchainConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(orchestrator_blockchain_settings.NewRpcNodeSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_blockchain_settings.NewToken6022Settings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_blockchain_settings.NewAgentCollectionsManagerSettings)
	if err != nil {
		panic(err)
	}

	// Subscribers
	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionsManagerCollectionCreatedSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionsManagerCollectionCreatedSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionAgentAttributeUpdatedSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionAgentAttributeUpdatedSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionAgentSignerUpdatedSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionAgentSignerUpdatedSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionMintedSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionMintedSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionTransferSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionTransferSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionAgentImageUpdatedSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionAgentImageUpdatedSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionMintProposalCreatedSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionMintProposalCreatedSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_subscribers.NewAgentCollectionMintProposalRemovedSubscriber,
		dig.As(new(orchestrator_sc_sync_subscribers.AgentCollectionMintProposalRemovedSubscriberInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		orchestrator_blockchain_services.NewToken6022Requester,
		dig.As(new(interfaces.Token6022RequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_services.NewOrchestratorDomainResolver,
		dig.As(new(interfaces.OrchestratorDomainResolverInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_services.NewAgentCollectionRequester,
		dig.As(new(interfaces.AgentCollectionRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_blockchain_services.NewAgentCollectionsManagerRequester,
		dig.As(new(sc_sync_interfaces.AgentCollectionsManagerRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Ethereum clients
	err = container.Provide(func(rpcNodeSettings *orchestrator_blockchain_settings.RpcNodeSettings) (*ethclient.Client, error) {
		return ethclient.Dial(rpcNodeSettings.WsUrl)
	}, dig.Name("ws"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(rpcNodeSettings *orchestrator_blockchain_settings.RpcNodeSettings) (*ethclient.Client, error) {
		return ethclient.Dial(rpcNodeSettings.HttpUrl)
	}, dig.Name("http"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(rpcNodeSettings *orchestrator_blockchain_settings.RpcNodeSettings) (*ethclient.Client, error) {
		return ethclient.Dial(rpcNodeSettings.ENSHttpUrl)
	}, dig.Name("ens"))
	if err != nil {
		panic(err)
	}
}
