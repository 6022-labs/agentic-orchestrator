package models

import (
	"encoding/json"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_inference/utils"
)

type ExternalAnthropicMessagesStandardRequest struct {
	Model     string `json:"model"`
	MaxTokens uint   `json:"max_tokens"`

	Tools      []ExternalAnthropicMessagesStandardRequestTool      `json:"tools,omitempty"`
	System     []ExternalAnthropicMessagesStandardRequestSystem    `json:"system,omitempty"`
	Messages   []ExternalAnthropicMessagesStandardMessage          `json:"messages,omitempty"`
	McpServers []ExternalAnthropicMessagesStandardRequestMCPServer `json:"mcp_servers,omitempty"`
}

type ExternalAnthropicMessagesStandardRequestTool struct {
	Name        string                                                  `json:"name"`
	Description string                                                  `json:"description"`
	InputSchema ExternalAnthropicMessagesStandardRequestToolInputSchema `json:"input_schema"`
}

type ExternalAnthropicMessagesStandardRequestToolInputSchema struct {
	Type       string                            `json:"type"`
	Required   []string                          `json:"required"`
	Properties ExternalToolInputSchemaProperties `json:"properties"`
}

type ExternalAnthropicMessagesStandardRequestSystem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ExternalAnthropicMessagesStandardRequestMCPServer struct {
	Type               string  `json:"type"`
	Url                string  `json:"url"`
	Name               string  `json:"name"`
	AuthorizationToken *string `json:"authorization_token,omitempty"`
}

func NewExternalAnthropicMessagesStandardRequest(
	model string,
	messages []ExternalAnthropicMessagesStandardMessage,
	system []ExternalAnthropicMessagesStandardRequestSystem,
) *ExternalAnthropicMessagesStandardRequest {
	return &ExternalAnthropicMessagesStandardRequest{
		Model:     model,
		MaxTokens: 32000,
		Messages:  messages,
		System:    system,
	}
}

func NewExternalAnthropicMessagesStandardRequestSystemFromChatCompletionMessage(
	message *models.ChatCompletionMessage,
) ExternalAnthropicMessagesStandardRequestSystem {
	return ExternalAnthropicMessagesStandardRequestSystem{
		Type: "text",
		Text: message.Content,
	}
}

func NewExternalAnthropicMessagesStandardRequestSystemsFromChatCompletionMessages(
	messages []*models.ChatCompletionMessage,
) []ExternalAnthropicMessagesStandardRequestSystem {
	inputs := make([]ExternalAnthropicMessagesStandardRequestSystem, 0, len(messages))
	for _, message := range messages {
		if message.Role != "system" {
			continue
		}

		inputs = append(inputs, NewExternalAnthropicMessagesStandardRequestSystemFromChatCompletionMessage(message))
	}

	return inputs
}

func (r *ExternalAnthropicMessagesStandardRequest) AddMcpServerFromMcpConnection(
	mcpConnection *models.MCPConnection,
) {
	var authToken string
	if mcpConnection != nil {
		if v, ok := mcpConnection.HeaderAuth["Authorization"]; ok {
			authToken = v
		}
	}

	authToken, _ = strings.CutPrefix(authToken, "Bearer ")

	if r.McpServers == nil {
		r.McpServers = make([]ExternalAnthropicMessagesStandardRequestMCPServer, 0)
	}

	newMcpServer := ExternalAnthropicMessagesStandardRequestMCPServer{
		Type: "url",
		Url:  mcpConnection.Url,
		Name: mcpConnection.Name,
	}

	if len(strings.TrimSpace(authToken)) != 0 {
		newMcpServer.AuthorizationToken = &authToken
	}

	r.McpServers = append(r.McpServers, newMcpServer)
}

func (r *ExternalAnthropicMessagesStandardRequest) AddMcpConnectionAsFunctionCallTools(
	mcpConnection *models.MCPConnection,
	relatedTools []*interfaces.McpTool,
) {
	if len(relatedTools) == 0 {
		return
	}

	if r.Tools == nil {
		r.Tools = make([]ExternalAnthropicMessagesStandardRequestTool, 0)
	}

	for _, mcpTool := range relatedTools {
		toolName := utils.McpToolToFunctionName(mcpConnection, mcpTool)

		r.Tools = append(r.Tools, ExternalAnthropicMessagesStandardRequestTool{
			Name:        toolName,
			Description: mcpTool.Description,
			InputSchema: ExternalAnthropicMessagesStandardRequestToolInputSchema{
				Type:     mcpTool.InputScheme.Type,
				Required: mcpTool.InputScheme.Required,

				Properties: NewExternalToolInputSchemaPropertiesFromInternalModels(mcpTool.InputScheme.Properties),
			},
		})
	}
}

func NewExternalAnthropicMessagesStandardRequestMessagesFromToolCallOutput(
	toolUseId string,
	output map[string]any,
) (*ExternalAnthropicMessagesStandardMessage, error) {
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	outputStr := string(outputBytes)

	return &ExternalAnthropicMessagesStandardMessage{
		Role: "user",
		Content: []ExternalAnthropicMessagesStandardMessageContent{
			{
				Type:      "tool_result",
				ToolUseId: &toolUseId,
				Content:   &outputStr,
			},
		},
	}, nil
}
