package utils

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/services/interfaces"
	"github.com/samber/lo"
)

const (
	// We use the underscore as separator
	// MCP connection names cannot contain underscores, so there will be no ambiguity
	McpToolToFunctionNameTemplate  = "%s_%s"
	McpToolToFunctionNameSeparator = "_"
)

func McpToolToFunctionName(mcpConnection *models.MCPConnection, mcpTool *interfaces.McpTool) string {
	return fmt.Sprintf(McpToolToFunctionNameTemplate, mcpConnection.Name, mcpTool.Name)
}

func ExtractMcpInfoFromFunctionName(functionName string, mcpConnections []*models.MCPConnection, mcpsTools map[*models.MCPConnection][]*interfaces.McpTool) (*models.MCPConnection, *interfaces.McpTool, error) {
	parts := strings.SplitN(functionName, McpToolToFunctionNameSeparator, 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return nil, nil, fmt.Errorf("invalid function name format: %q", functionName)
	}
	connName := parts[0]
	toolName := parts[1]

	mcpConnection, found := lo.Find(mcpConnections, func(mcpConnection *models.MCPConnection) bool {
		return mcpConnection.Name == connName
	})
	if !found {
		return nil, nil, fmt.Errorf("no MCP connection found for function name: %q", functionName)
	}

	mcpTool, found := lo.Find(mcpsTools[mcpConnection], func(mcpTool *interfaces.McpTool) bool {
		return mcpTool.Name == toolName
	})
	if !found {
		return nil, nil, fmt.Errorf("no MCP tool found for function name: %q", functionName)
	}

	return mcpConnection, mcpTool, nil
}
