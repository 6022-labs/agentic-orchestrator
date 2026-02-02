package interfaces

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/models"
)

type McpTool struct {
	Name        string
	Description string
	InputScheme McpToolInputScheme
}

type McpToolInputScheme struct {
	Type                 string
	Properties           map[string]ToolInputSchemeProperty
	Required             []string
	AdditionalProperties bool
	Schema               string
}

type ToolInputSchemeProperty struct {
	Type        string
	Items       *ToolInputSchemeProperty
	Description *string
}

type McpRequesterInterface interface {
	ListTools(agentId string, mcpConnection *models.MCPConnection) ([]*McpTool, error)
	CallTool(
		agentId string,
		mcpConnection *models.MCPConnection,
		tool *McpTool,
		toolInput map[string]any,
	) (map[string]any, error)
}
