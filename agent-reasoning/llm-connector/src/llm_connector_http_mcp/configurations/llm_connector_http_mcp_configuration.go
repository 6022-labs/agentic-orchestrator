package configurations

import (
	"net/http"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_mcp/clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_mcp/services"
	"go.uber.org/dig"
)

func AddLLMConnectorHttpMcpConfiguration(container *dig.Container) {
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
