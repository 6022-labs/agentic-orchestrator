package configurations

import (
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

func AddConversationBrokerBlockchainConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewBlockchainSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewRpcNodeSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewAgentEnsRegistrySettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(rpcNodeSettings *settings.RpcNodeSettings) (clients.EthClientInterface, error) {
		return ethclient.Dial(rpcNodeSettings.HttpUrl)
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewEnsUrlResolver,
		dig.As(new(service_interfaces.EnsUrlResolverInterface)),
	)
	if err != nil {
		panic(err)
	}
}
