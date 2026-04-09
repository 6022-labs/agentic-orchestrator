package services

import (
	"encoding/json"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_mcp/clients"
)

type MCPRequester struct {
	client clients.McpClientInterface
}

func NewMCPRequester(client clients.McpClientInterface) *MCPRequester {
	return &MCPRequester{
		client: client,
	}
}

func (r *MCPRequester) ListTools(mcpConnection *models.MCPConnection) ([]*interfaces.McpTool, error) {
	if mcpConnection == nil {
		return nil, fmt.Errorf("MCP connection is nil")
	}

	result, err := r.client.ListTools(mcpConnection.Url, mcpConnection.HeaderAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to list MCP tools: %w", err)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("MCP list tools returned error: %s", result.Error.Message)
	}

	return result.Result.ToMcpTools(), nil
}

func (r *MCPRequester) CallTool(
	mcpConnection *models.MCPConnection,
	tool *interfaces.McpTool,
	toolInput map[string]any,
) (map[string]any, error) {
	if mcpConnection == nil {
		return nil, fmt.Errorf("MCP connection is nil")
	}
	if tool == nil {
		return nil, fmt.Errorf("MCP tool is nil")
	}

	fixedInput := r.fixToolInputBySchema(tool, toolInput)

	result, err := r.client.CallTool(mcpConnection.Url, tool.Name, fixedInput, mcpConnection.HeaderAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to call MCP tool '%s': %w", tool.Name, err)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("MCP call tool '%s' returned error: %s", tool.Name, result.Error.Message)
	}

	var output map[string]any
	if err := json.Unmarshal(result.Result, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MCP tool '%s' output: %w", tool.Name, err)
	}

	return output, nil
}

func (r *MCPRequester) fixToolInputBySchema(tool *interfaces.McpTool, input map[string]any) map[string]any {
	if input == nil {
		return nil
	}

	if len(tool.InputScheme.Properties) == 0 {
		return input
	}

	// Work on a shallow copy to avoid mutating the caller's map
	fixed := make(map[string]any, len(input))
	for k, v := range input {
		fixed[k] = v
	}

	for name, prop := range tool.InputScheme.Properties {
		val, ok := fixed[name]
		if !ok {
			continue
		}

		switch prop.Type {
		case "string":
			// If schema expects string but value is number/bool/[]byte, coerce to string
			switch typedVal := val.(type) {
			case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, bool:
				fixed[name] = fmt.Sprint(typedVal)
			default:
				// no-op
			}
		default:
			// Additional coercions can be added here if needed (e.g. numeric strings to numbers).
		}
	}

	return fixed
}
