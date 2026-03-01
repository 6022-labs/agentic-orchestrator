package models

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"

type ExternalAnthropicMessagesStandardMessage struct {
	// Common fields
	Role    string                                            `json:"role"`
	Content []ExternalAnthropicMessagesStandardMessageContent `json:"content,omitempty"`
}

type ExternalAnthropicMessagesStandardMessageContent struct {
	// Common fields
	Type string `json:"type"`

	// `text` type fields
	Text *string `json:"text,omitempty"`

	// `thinking` type fields
	Signature *string `json:"signature,omitempty"`
	Thinking  *string `json:"thinking,omitempty"`

	// `redacted_thinking` type fields
	Data *string `json:"data,omitempty"`

	// `tool_use`,
	// `mcp_tool_use`,
	// `server_tool_use` fields
	Id    *string        `json:"id,omitempty"`
	Name  *string        `json:"name,omitempty"`
	Input map[string]any `json:"input,omitempty"`

	// `tool_result`,
	// `mcp_tool_result`,
	// `web_fetch_tool_result`,
	// `web_search_tool_result`,
	// `code_execution_tool_result`,
	// `bash_code_execution_tool_result`,
	// `text_editor_code_execution_tool_result` type fields
	Content   any     `json:"content,omitempty"`
	ToolUseId *string `json:"tool_use_id,omitempty"`

	// `mcp_tool_use` fields
	ServerName *string `json:"server_name,omitempty"`

	// `container_upload` fields
	FileId *string `json:"file_id,omitempty"`
}

func NewExternalAnthropicMessagesStandardMessageFromChatCompletionMessage(
	message *models.ChatCompletionMessage,
) ExternalAnthropicMessagesStandardMessage {
	return ExternalAnthropicMessagesStandardMessage{
		Role: message.Role,
		Content: []ExternalAnthropicMessagesStandardMessageContent{
			{
				Type: "text",
				Text: &message.Content,
			},
		},
	}
}

func NewExternalAnthropicMessagesStandardMessagesFromChatCompletionMessages(
	messages []*models.ChatCompletionMessage,
) []ExternalAnthropicMessagesStandardMessage {
	inputs := make([]ExternalAnthropicMessagesStandardMessage, 0, len(messages))
	for _, message := range messages {
		if message.Role != "user" && message.Role != "assistant" {
			continue
		}

		inputs = append(inputs, NewExternalAnthropicMessagesStandardMessageFromChatCompletionMessage(message))
	}

	return inputs
}
