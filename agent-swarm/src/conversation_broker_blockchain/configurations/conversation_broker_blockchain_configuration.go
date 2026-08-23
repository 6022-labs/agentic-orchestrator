package configurations

import (
	"context"
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/dig"
	"net/http"
	"strconv"
)

func AddConversationBrokerBlockchainConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(
		metrics.NewBlockchainRpcMetrics,
		dig.As(new(metrics_interfaces.BlockchainRpcMetricsInterface)),
	); err != nil {
		panic(err)
	}

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

	err = container.Provide(func(
		rpcNodeSettings *settings.RpcNodeSettings,
		blockchainSettings *settings.BlockchainSettings,
		blockchainRpcMetrics metrics_interfaces.BlockchainRpcMetricsInterface,
	) (clients.EthClientInterface, error) {
		chainId := strconv.FormatInt(blockchainSettings.ChainId, 10)
		rpcClient, err := rpc.DialOptions(
			context.Background(),
			rpcNodeSettings.HttpUrl,
			rpc.WithHTTPClient(&http.Client{
				Transport: otelhttp.NewTransport(
					metrics.NewRpcRoundTripper(http.DefaultTransport, blockchainRpcMetrics, chainId),
					otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
				),
			}),
		)
		if err != nil {
			return nil, err
		}

		return ethclient.NewClient(rpcClient), nil
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
