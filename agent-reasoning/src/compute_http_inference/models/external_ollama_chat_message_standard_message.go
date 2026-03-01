package models

import (
	"encoding/json"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
)

type ExternalOllamaChatMessageStandardMessage struct {
	Role      string                                             `json:"role"`
	Content   string                                             `json:"content"`
	ToolName  *string                                            `json:"tool_name,omitempty"`
	ToolCalls []ExternalOllamaChatMessageStandardMessageToolCall `json:"tool_calls,omitempty"`
}

type ExternalOllamaChatMessageStandardMessageToolCall struct {
	Function *ExternalOllamaChatMessageStandardMessageToolCallFunction `json:"function,omitempty"`
}

type ExternalOllamaChatMessageStandardMessageToolCallFunction struct {
	Name        string                                                             `json:"name"`
	Index       *int                                                               `json:"index,omitempty"`
	Arguments   map[string]any                                                     `json:"arguments,omitempty"`
	Parameters  *ExternalOllamaChatMessageStandardMessageToolCallFunctionArguments `json:"parameters,omitempty"`
	Description *string                                                            `json:"description,omitempty"`
}

type ExternalOllamaChatMessageStandardMessageToolCallFunctionArguments struct {
	Type       string                            `json:"type"`
	Required   []string                          `json:"required"`
	Properties ExternalToolInputSchemaProperties `json:"properties"`
}

func NewExternalOllamaChatMessageStandardMessageFromChatCompletionMessage(
	message *models.ChatCompletionMessage,
) ExternalOllamaChatMessageStandardMessage {
	return ExternalOllamaChatMessageStandardMessage{
		Role:    message.Role,
		Content: message.Content,
	}
}

func NewExternalOllamaChatMessageStandardMessagesFromChatCompletionMessages(
	messages []*models.ChatCompletionMessage,
) []ExternalOllamaChatMessageStandardMessage {
	result := make([]ExternalOllamaChatMessageStandardMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, NewExternalOllamaChatMessageStandardMessageFromChatCompletionMessage(message))
	}

	return result
}

func NewExternalOllamaChatMessageStandardMessageFromToolCallOutput(
	toolName string,
	output map[string]any,
) (*ExternalOllamaChatMessageStandardMessage, error) {
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	outputStr := string(outputBytes)

	return &ExternalOllamaChatMessageStandardMessage{
		Role:     "tool",
		ToolName: &toolName,
		Content:  outputStr,
	}, nil
}
