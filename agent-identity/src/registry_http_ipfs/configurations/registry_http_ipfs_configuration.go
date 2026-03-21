package configurations

import (
	"net/http"

	registry_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/settings"
	"go.uber.org/dig"
)

func AddRegistryHttpIPFSConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewIPFSSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("ipfs-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		clients.NewIPFSClient,
		dig.As(new(clients.IPFSClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewIPFSFileStorage,
		dig.As(new(registry_interfaces.IPFSFileStorageInterface)),
	)
	if err != nil {
		panic(err)
	}
}
