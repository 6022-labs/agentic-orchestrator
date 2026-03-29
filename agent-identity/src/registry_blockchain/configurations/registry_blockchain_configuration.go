package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_blockchain/settings"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

func AddRegistryBlockchainConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewRpcNodeSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewToken6022Settings)
	if err != nil {
		panic(err)
	}

	// Services (requesters)
	err = container.Provide(
		services.NewToken6022Requester,
		dig.As(new(interfaces.Token6022RequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Ethereum clients
	err = container.Provide(func(rpcNodeSettings *settings.RpcNodeSettings) (*ethclient.Client, error) {
		return ethclient.Dial(rpcNodeSettings.HttpUrl)
	}, dig.Name("http"))
	if err != nil {
		panic(err)
	}
}
