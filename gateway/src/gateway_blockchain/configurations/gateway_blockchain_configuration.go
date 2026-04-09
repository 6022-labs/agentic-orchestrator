package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	orchestrator_blockchain_services "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_blockchain/services"
	orchestrator_blockchain_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_blockchain/settings"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

func AddGatewayBlockchainConfiguration(container *dig.Container) {
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
