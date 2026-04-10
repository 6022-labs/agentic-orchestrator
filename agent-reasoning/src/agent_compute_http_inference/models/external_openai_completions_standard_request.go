package models

import (
	"encoding/json"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/utils"
)

type ExternalOpenAICompletionsStandardRequest struct {
	Model    string                                            `json:"model"`
	Messages []ExternalOpenAICompletionsStandardRequestMessage `json:"messages"`
	Tools    []ExternalOpenAICompletionsStandardRequestTool    `json:"tools,omitempty"`
}

type ExternalOpenAICompletionsStandardRequestMessage struct {
	// Common fields
	Role    string  `json:"role"`
	Content *string `json:"content,omitempty"`

	// `Tool call` fields
	ToolCalls []ExternalOpenAICompletionsStandardRequestMessageToolCall `json:"tool_calls,omitempty"`

	// `Tool call output` fields
	ToolCallId *string `json:"tool_call_id,omitempty"`
}

type ExternalOpenAICompletionsStandardRequestMessageToolCall struct {
	ID       string                                            `json:"id"`
	Type     string                                            `json:"type"` // typically "function"
	Function ExternalOpenAICompletionsStandardToolCallFunction `json:"function"`
}

type ExternalOpenAICompletionsStandardRequestTool struct {
	Type     string                                               `json:"type"`
	Function ExternalOpenAICompletionsStandardRequestToolFunction `json:"function"`
}

type ExternalOpenAICompletionsStandardRequestToolFunction struct {
	Name        string                                                         `json:"name"`
	Strict      *bool                                                          `json:"strict,omitempty"`
	Description string                                                         `json:"description"`
	Parameters  ExternalOpenAICompletionsStandardRequestToolFunctionParameters `json:"parameters"`
}

type ExternalOpenAICompletionsStandardRequestToolFunctionParameters struct {
	Type                 string                            `json:"type"`
	Properties           ExternalToolInputSchemaProperties `json:"properties"`
	Required             []string                          `json:"required"`
	AdditionalProperties bool                              `json:"additionalProperties,omitempty"`
}

func NewExternalOpenAICompletionsStandardRequest(
	model string,
	messages []ExternalOpenAICompletionsStandardRequestMessage,
) *ExternalOpenAICompletionsStandardRequest {
	return &ExternalOpenAICompletionsStandardRequest{
		Model:    model,
		Messages: messages,
	}
}

func NewExternalOpenAICompletionsStandardRequestMessageFromChatCompletionMessage(
	message *models.ChatCompletionMessage,
) ExternalOpenAICompletionsStandardRequestMessage {
	return ExternalOpenAICompletionsStandardRequestMessage{
		Role:    message.Role,
		Content: &message.Content,
	}
}

func NewExternalOpenAICompletionsStandardRequestMessageFromToolCallOutput(
	callId string,
	output map[string]any,
) (*ExternalOpenAICompletionsStandardRequestMessage, error) {
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	outputStr := string(outputBytes)

	return &ExternalOpenAICompletionsStandardRequestMessage{
		Role:       "tool",
		ToolCallId: &callId,
		Content:    &outputStr,
	}, nil
}

func NewExternalOpenAICompletionsStandardRequestMessagesFromChatCompletionMessages(
	messages []*models.ChatCompletionMessage,
) []ExternalOpenAICompletionsStandardRequestMessage {
	result := make([]ExternalOpenAICompletionsStandardRequestMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, NewExternalOpenAICompletionsStandardRequestMessageFromChatCompletionMessage(message))
	}

	return result
}

func (r *ExternalOpenAICompletionsStandardRequest) AddMcpConnectionAsFunctionCallTools(
	mcpConnection *models.MCPConnection,
	relatedTools []*interfaces.McpTool,
) {
	if len(relatedTools) == 0 {
		return
	}

	if r.Tools == nil {
		r.Tools = make([]ExternalOpenAICompletionsStandardRequestTool, 0)
	}

	for _, mcpTool := range relatedTools {
		// MCPs are not always as well defined, so we set strict to false by default
		// But also not always all fields are required
		strict := false
		name := utils.McpToolToFunctionName(mcpConnection, mcpTool)

		r.Tools = append(r.Tools, ExternalOpenAICompletionsStandardRequestTool{
			Type: "function",
			Function: ExternalOpenAICompletionsStandardRequestToolFunction{
				Name:        name,
				Strict:      &strict,
				Description: mcpTool.Description,
				Parameters: ExternalOpenAICompletionsStandardRequestToolFunctionParameters{
					Type:                 mcpTool.InputScheme.Type,
					Required:             mcpTool.InputScheme.Required,
					AdditionalProperties: mcpTool.InputScheme.AdditionalProperties,

					Properties: NewExternalToolInputSchemaPropertiesFromInternalModels(mcpTool.InputScheme.Properties),
				},
			},
		})
	}
}
