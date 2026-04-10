package models

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/utils"
	"github.com/samber/lo"
)

type ExternalAmazonBedrockConverseStandardRequest struct {
	ModelId         string                                                       `json:"modelId"`
	Messages        []ExternalAmazonBedrockConverseStandardMessage               `json:"messages"`
	System          []ExternalAmazonBedrockConverseStandardRequestSystem         `json:"system"`
	ToolConfig      *ExternalAmazonBedrockConverseStandardRequestToolConfig      `json:"toolConfig,omitempty"`
	InferenceConfig *ExternalAmazonBedrockConverseStandardRequestInferenceConfig `json:"inferenceConfig,omitempty"`
}

type ExternalAmazonBedrockConverseStandardRequestSystem struct {
	Text string `json:"text"`
}

type ExternalAmazonBedrockConverseStandardRequestToolConfig struct {
	ToolChoice ExternalAmazonBedrockConverseStandardRequestToolConfigToolChoice `json:"toolChoice"`
	Tools      []ExternalAmazonBedrockConverseStandardRequestToolConfigTool     `json:"tools"`
}

type ExternalAmazonBedrockConverseStandardRequestToolConfigToolChoice struct {
	Auto *ExternalAmazonBedrockConverseStandardRequestToolConfigToolChoiceAuto `json:"auto,omitempty"`
}

type ExternalAmazonBedrockConverseStandardRequestToolConfigToolChoiceAuto struct {
}

type ExternalAmazonBedrockConverseStandardRequestToolConfigTool struct {
	ToolSpec ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpec `json:"toolSpec"`
}

type ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpec struct {
	Name        string                                                                        `json:"name"`
	Description string                                                                        `json:"description"`
	InputSchema ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpecInputSchema `json:"inputSchema"`
}

type ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpecInputSchema struct {
	Json ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpecInputSchemaJson `json:"json"`
}

type ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpecInputSchemaJson struct {
	Type                 string                            `json:"type"`
	Properties           ExternalToolInputSchemaProperties `json:"properties"`
	Required             []string                          `json:"required"`
	AdditionalProperties bool                              `json:"additionalProperties"`
}

type ExternalAmazonBedrockConverseStandardRequestInferenceConfig struct {
	MaxTokens int `json:"maxTokens"`
}

func NewExternalAmazonBedrockConverseStandardRequest(
	model string,
	messages []ExternalAmazonBedrockConverseStandardMessage,
	system []ExternalAmazonBedrockConverseStandardRequestSystem,
) *ExternalAmazonBedrockConverseStandardRequest {
	return &ExternalAmazonBedrockConverseStandardRequest{
		ModelId:  model,
		System:   system,
		Messages: messages,
		InferenceConfig: &ExternalAmazonBedrockConverseStandardRequestInferenceConfig{
			MaxTokens: 32000,
		},
	}
}

func NewExternalAmazonBedrockConverseStandardRequestSystemFromChatCompletionMessage(
	message *models.ChatCompletionMessage,
) ExternalAmazonBedrockConverseStandardRequestSystem {
	return ExternalAmazonBedrockConverseStandardRequestSystem{
		Text: message.Content,
	}
}

func NewExternalAmazonBedrockConverseStandardRequestMessagesAndSystemFromChatCompletionMessages(
	messages []*models.ChatCompletionMessage,
) (
	[]ExternalAmazonBedrockConverseStandardMessage,
	[]ExternalAmazonBedrockConverseStandardRequestSystem,
) {
	systemMessages := lo.Filter(messages, func(item *models.ChatCompletionMessage, _ int) bool {
		return strings.EqualFold(item.Role, "system")
	})

	otherMessages, _ := lo.Difference(messages, systemMessages)

	return lo.Map(otherMessages, func(item *models.ChatCompletionMessage, _ int) ExternalAmazonBedrockConverseStandardMessage {
			return NewExternalAmazonBedrockConverseStandardMessageFromChatCompletionMessage(item)
		}), lo.Map(systemMessages, func(item *models.ChatCompletionMessage, _ int) ExternalAmazonBedrockConverseStandardRequestSystem {
			return NewExternalAmazonBedrockConverseStandardRequestSystemFromChatCompletionMessage(item)
		})
}

func (r *ExternalAmazonBedrockConverseStandardRequest) AddDefaultToolConfigIfNil() {
	if r.ToolConfig == nil {
		r.ToolConfig = &ExternalAmazonBedrockConverseStandardRequestToolConfig{
			Tools: make([]ExternalAmazonBedrockConverseStandardRequestToolConfigTool, 0),
			ToolChoice: ExternalAmazonBedrockConverseStandardRequestToolConfigToolChoice{
				Auto: &ExternalAmazonBedrockConverseStandardRequestToolConfigToolChoiceAuto{},
			},
		}
	}
}

func (r *ExternalAmazonBedrockConverseStandardRequest) AddMcpConnectionAsFunctionCallTools(
	mcpConnection *models.MCPConnection,
	relatedTools []*interfaces.McpTool,
) {
	if len(relatedTools) == 0 {
		return
	}

	r.AddDefaultToolConfigIfNil()

	for _, mcpTool := range relatedTools {
		toolName := utils.McpToolToFunctionName(mcpConnection, mcpTool)

		functionTool := ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpec{
			Name:        toolName,
			Description: mcpTool.Description,
			InputSchema: ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpecInputSchema{
				Json: ExternalAmazonBedrockConverseStandardRequestToolConfigToolToolSpecInputSchemaJson{
					Type:                 mcpTool.InputScheme.Type,
					Required:             mcpTool.InputScheme.Required,
					AdditionalProperties: mcpTool.InputScheme.AdditionalProperties,

					Properties: NewExternalToolInputSchemaPropertiesFromInternalModels(mcpTool.InputScheme.Properties),
				},
			},
		}

		if functionTool.InputSchema.Json.Required == nil {
			functionTool.InputSchema.Json.Required = []string{}
		}

		r.ToolConfig.Tools = append(r.ToolConfig.Tools, ExternalAmazonBedrockConverseStandardRequestToolConfigTool{
			ToolSpec: functionTool,
		})
	}
}
