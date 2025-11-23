package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	orchestrator_http_ipfs_clients "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/clients"
	orchestrator_http_ipfs_services "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/services"
	orchestrator_http_ipfs_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/settings"
	"go.uber.org/dig"
)

func AddOrchestratorGatewayHttpIpfsConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		orchestrator_http_ipfs_settings.NewIPFSSettings,
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		orchestrator_http_ipfs_services.NewIPFSFileStorage,
		dig.As(new(interfaces.IPFSFileStorageInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		orchestrator_http_ipfs_clients.NewIPFSClient,
		dig.As(new(orchestrator_http_ipfs_clients.IPFSClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("ipfs-http-client"))
	if err != nil {
		panic(err)
	}
}
