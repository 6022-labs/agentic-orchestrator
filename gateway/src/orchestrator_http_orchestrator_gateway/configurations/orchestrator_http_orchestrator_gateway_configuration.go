package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	http_orchestrator_clients "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_http_orchestrator_gateway/clients"
	http_orchestrator_services "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_http_orchestrator_gateway/services"
	"go.uber.org/dig"
)

func AddOrchestratorHttpOrchestratorGatewayConfiguration(container *dig.Container) {
	// Clients
	err := container.Provide(
		http_orchestrator_clients.NewOrchestratorMainApiClient,
		dig.As(new(http_orchestrator_clients.OrchestratorMainApiClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		http_orchestrator_services.NewExternalAgentInformationFetcher,
		dig.As(new(interfaces.ExternalAgentInformationFetcherInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("orchestrator-main-api-http-client"))
	if err != nil {
		panic(err)
	}
}
