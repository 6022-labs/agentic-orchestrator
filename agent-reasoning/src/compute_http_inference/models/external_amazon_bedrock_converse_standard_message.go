package models

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
)

// This type is used by both request and response
type ExternalAmazonBedrockConverseStandardMessage struct {
	Role    string                                                `json:"role"`
	Content []ExternalAmazonBedrockConverseStandardMessageContent `json:"content"`
}

type ExternalAmazonBedrockConverseStandardMessageContent struct {
	Text             *string                                                              `json:"text,omitempty"`
	ReasoningContent *ExternalAmazonBedrockConverseStandardMessageContentReasoningContent `json:"reasoningContent,omitempty"`
	ToolResult       *ExternalAmazonBedrockConverseStandardMessageContentToolResult       `json:"toolResult,omitempty"`
	ToolUse          *ExternalAmazonBedrockConverseStandardMessageContentToolUse          `json:"toolUse,omitempty"`
}

type ExternalAmazonBedrockConverseStandardMessageContentReasoningContent struct {
	ReasoningText *ExternalAmazonBedrockConverseStandardMessageContentReasoningContentReasoningText `json:"reasoningText,omitempty"`
}

type ExternalAmazonBedrockConverseStandardMessageContentReasoningContentReasoningText struct {
	Text      string  `json:"text"`
	Signature *string `json:"signature,omitempty"`
}

type ExternalAmazonBedrockConverseStandardMessageContentToolUse struct {
	Name      string         `json:"name"`
	Input     map[string]any `json:"input"`
	ToolUseId string         `json:"toolUseId"`
}

type ExternalAmazonBedrockConverseStandardMessageContentToolResult struct {
	Content   []ExternalAmazonBedrockConverseStandardMessageContentToolResultContent `json:"content"`
	ToolUseId string                                                                 `json:"toolUseId"`
}

type ExternalAmazonBedrockConverseStandardMessageContentToolResultContent struct {
	Json any     `json:"json,omitempty"`
	Text *string `json:"text,omitempty"`
}

func NewExternalAmazonBedrockConverseStandardMessageFromChatCompletionMessage(
	message *models.ChatCompletionMessage,
) ExternalAmazonBedrockConverseStandardMessage {
	return ExternalAmazonBedrockConverseStandardMessage{
		Role: message.Role,
		Content: []ExternalAmazonBedrockConverseStandardMessageContent{
			{
				Text: &message.Content,
			},
		},
	}
}

func NewExternalAmazonBedrockConverseStandardMessageFromToolCallOutputs(
	outputs map[string]map[string]any,
) ExternalAmazonBedrockConverseStandardMessage {
	contents := make([]ExternalAmazonBedrockConverseStandardMessageContent, 0, len(outputs))

	for toolUseId, toolResult := range outputs {
		contents = append(contents, ExternalAmazonBedrockConverseStandardMessageContent{
			ToolResult: &ExternalAmazonBedrockConverseStandardMessageContentToolResult{
				ToolUseId: toolUseId,
				Content: []ExternalAmazonBedrockConverseStandardMessageContentToolResultContent{
					{
						Json: toolResult,
					},
				},
			},
		})
	}

	return ExternalAmazonBedrockConverseStandardMessage{
		Role:    "user",
		Content: contents,
	}
}
