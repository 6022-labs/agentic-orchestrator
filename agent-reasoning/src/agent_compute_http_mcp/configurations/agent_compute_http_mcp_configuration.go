package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_mcp/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_mcp/services"
	"go.uber.org/dig"
)

func AddAgentComputeHttpMcpConfiguration(container *dig.Container) {
	// Services
	err := container.Provide(
		services.NewMCPRequester,
		dig.As(new(interfaces.McpRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Clients
	err = container.Provide(
		clients.NewMcpClient,
		dig.As(new(clients.McpClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Http clients
	err = container.Provide(func() *http.Client {
		return http.DefaultClient
	}, dig.Name("mcp-http-client"))
	if err != nil {
		panic(err)
	}
}
