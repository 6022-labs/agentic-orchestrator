package configurations

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/transports"

	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
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
		services.NewBlockchainRequester,
		dig.As(new(interfaces.BlockchainRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentCollectionRequester,
		dig.As(new(interfaces.AgentCollectionBlockchainRequesterInterface)),
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

	// Metrics
	err = container.Provide(
		metrics.NewBlockchainRpcMetrics,
		dig.As(new(metrics_interfaces.BlockchainRpcMetricsInterface)),
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

	// Only the http client is time-bounded: the ws one carries long-lived subscriptions.
	err = container.Provide(func(
		rpcNodeSettings *settings.RpcNodeSettings,
		blockchainSettings *common_settings.BlockchainSettings,
		blockchainRpcMetrics metrics_interfaces.BlockchainRpcMetricsInterface,
	) (*ethclient.Client, error) {
		chainId := strconv.FormatInt(blockchainSettings.ChainId, 10)
		rpcRoundTripper, err := transports.NewRpcRoundTripper(
			http.DefaultTransport,
			rpcNodeSettings.HttpUrl,
			blockchainRpcMetrics,
			chainId,
		)
		if err != nil {
			return nil, err
		}
		httpClient := &http.Client{
			Timeout: time.Duration(rpcNodeSettings.RequestTimeoutSeconds) * time.Second,
			Transport: otelhttp.NewTransport(
				rpcRoundTripper,
				otelhttp.WithMeterProvider(metricnoop.NewMeterProvider()),
			),
		}

		rpcClient, err := rpc.DialOptions(
			context.Background(),
			rpcRoundTripper.Origin(),
			rpc.WithHTTPClient(httpClient),
		)
		if err != nil {
			return nil, err
		}

		return ethclient.NewClient(rpcClient), nil
	}, dig.Name("http"))
	if err != nil {
		panic(err)
	}
}
