package models

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/utils"
)

type ExternalOllamaChatMessageStandardRequest struct {
	Model    string                                         `json:"model"`
	Messages []ExternalOllamaChatMessageStandardMessage     `json:"messages"`
	Tools    []ExternalOllamaChatMessageStandardRequestTool `json:"tools,omitempty"`
	Think    *bool                                          `json:"think,omitempty"`
	Stream   *bool                                          `json:"stream,omitempty"`
}

type ExternalOllamaChatMessageStandardRequestTool struct {
	Type     string                                                   `json:"type"`
	Function ExternalOllamaChatMessageStandardMessageToolCallFunction `json:"function"`
}

func NewExternalOllamaChatMessageStandardRequest(model string, messages []ExternalOllamaChatMessageStandardMessage) *ExternalOllamaChatMessageStandardRequest {
	stream := false

	return &ExternalOllamaChatMessageStandardRequest{
		Model:    model,
		Stream:   &stream,
		Messages: messages,
	}
}

func (r *ExternalOllamaChatMessageStandardRequest) AddMcpConnectionAsFunctionCallTools(mcpConnection *models.MCPConnection, tools []*interfaces.McpTool) {
	if len(tools) == 0 {
		return
	}

	if r.Tools == nil {
		r.Tools = make([]ExternalOllamaChatMessageStandardRequestTool, 0)
	}

	for _, tool := range tools {
		name := utils.McpToolToFunctionName(mcpConnection, tool)

		r.Tools = append(r.Tools, ExternalOllamaChatMessageStandardRequestTool{
			Type: "function",
			Function: ExternalOllamaChatMessageStandardMessageToolCallFunction{
				Name:        name,
				Description: &tool.Description,
				Parameters: &ExternalOllamaChatMessageStandardMessageToolCallFunctionArguments{
					Type:       tool.InputScheme.Type,
					Required:   tool.InputScheme.Required,
					Properties: NewExternalToolInputSchemaPropertiesFromInternalModels(tool.InputScheme.Properties),
				},
			},
		})
	}
}
