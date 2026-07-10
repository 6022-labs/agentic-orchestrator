package configurations

import (
	"net/http"
	"time"

	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_blockchain/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_blockchain/settings"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

func AddConversationBrokerBlockchainConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewEnsResolverSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(ensResolverSettings *settings.EnsResolverSettings) (clients.EthCallerInterface, error) {
		return ethclient.Dial(ensResolverSettings.RpcUrl)
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(ensResolverSettings *settings.EnsResolverSettings) *http.Client {
		return &http.Client{Timeout: time.Duration(ensResolverSettings.GatewayRequestTimeoutSeconds) * time.Second}
	}, dig.Name("ens-gateway-http-client"))
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
