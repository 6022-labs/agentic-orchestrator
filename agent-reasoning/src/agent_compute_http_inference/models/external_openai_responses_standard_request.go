package models

import (
	"encoding/json"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/utils"
)

type ExternalOpenAIResponsesStandardRequest struct {
	Model              string                                        `json:"model"`
	Input              []ExternalOpenAIResponsesStandardRequestInput `json:"input"`
	Tools              []ExternalOpenAIResponsesStandardRequestTool  `json:"tools"`
	ToolChoice         string                                        `json:"tool_choice,omitempty"`
	PreviousResponseID *string                                       `json:"previous_response_id,omitempty"`
}

type ExternalOpenAIResponsesStandardRequestInput struct {
	Type string `json:"type"`

	// Type `message` fields
	Role    *string `json:"role,omitempty"`
	Content *string `json:"content,omitempty"`

	// Type `function_call_output` fields
	CallId *string `json:"call_id,omitempty"`
	Output *string `json:"output,omitempty"`
}

type ExternalOpenAIResponsesStandardToolRequireApproval string

const (
	ExternalOpenAIResponsesStandardToolRequireApprovalNever  ExternalOpenAIResponsesStandardToolRequireApproval = "never"
	ExternalOpenAIResponsesStandardToolRequireApprovalAlways ExternalOpenAIResponsesStandardToolRequireApproval = "always"
)

func (e ExternalOpenAIResponsesStandardToolRequireApproval) ToPointer() *ExternalOpenAIResponsesStandardToolRequireApproval {
	return &e
}

type ExternalOpenAIResponsesStandardRequestTool struct {
	// Common fields
	Type string `json:"type"`

	// Function type tool fields
	Name        *string                                                       `json:"name,omitempty"`
	Strict      *bool                                                         `json:"strict,omitempty"`
	Parameters  *ExternalOpenAIResponsesStandardRequestToolFunctionParameters `json:"parameters,omitempty"`
	Description *string                                                       `json:"description,omitempty"`

	// Mcp type tool fields
	Headers         *map[string]string                                  `json:"headers,omitempty"`
	ServerUrl       *string                                             `json:"server_url,omitempty"`
	ServerLabel     *string                                             `json:"server_label,omitempty"`
	RequireApproval *ExternalOpenAIResponsesStandardToolRequireApproval `json:"require_approval,omitempty"`
}

type ExternalOpenAIResponsesStandardRequestToolFunctionParameters struct {
	Type                 string                            `json:"type"`
	Properties           ExternalToolInputSchemaProperties `json:"properties"`
	Required             []string                          `json:"required"`
	AdditionalProperties bool                              `json:"additionalProperties"`
}

func NewExternalOpenAIResponsesStandardRequest(
	model string,
	inputs []ExternalOpenAIResponsesStandardRequestInput,
	previousResponseId *string,
) *ExternalOpenAIResponsesStandardRequest {
	return &ExternalOpenAIResponsesStandardRequest{
		Model:              model,
		ToolChoice:         "auto",
		Input:              inputs,
		PreviousResponseID: previousResponseId,
	}
}

func NewExternalOpenAICompletionsStandardRequestInputFromChatCompletionMessage(
	message *models.ChatCompletionMessage,
) ExternalOpenAIResponsesStandardRequestInput {
	return ExternalOpenAIResponsesStandardRequestInput{
		Type:    "message",
		Role:    &message.Role,
		Content: &message.Content,
	}
}

func NewExternalOpenAIResponsesStandardRequestFromToolCallOutput(
	callId string,
	output map[string]any,
) (*ExternalOpenAIResponsesStandardRequestInput, error) {
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	outputStr := string(outputBytes)

	return &ExternalOpenAIResponsesStandardRequestInput{
		Type:   "function_call_output",
		CallId: &callId,
		Output: &outputStr,
	}, nil
}

func NewExternalOpenAICompletionsStandardRequestInputsFromChatCompletionMessages(
	messages []*models.ChatCompletionMessage,
) []ExternalOpenAIResponsesStandardRequestInput {
	inputs := make([]ExternalOpenAIResponsesStandardRequestInput, 0, len(messages))
	for _, message := range messages {
		inputs = append(inputs, NewExternalOpenAICompletionsStandardRequestInputFromChatCompletionMessage(message))
	}

	return inputs
}

func (r *ExternalOpenAIResponsesStandardRequest) AddMcpConnectionsAsMcpTools(
	mcpConnections []*models.MCPConnection,
) {
	if len(mcpConnections) == 0 {
		return
	}

	if r.Tools == nil {
		r.Tools = make([]ExternalOpenAIResponsesStandardRequestTool, 0)
	}

	for _, conn := range mcpConnections {
		r.Tools = append(r.Tools, ExternalOpenAIResponsesStandardRequestTool{
			Type: "mcp",

			ServerLabel: &conn.Name,
			Headers:     &conn.HeaderAuth,

			ServerUrl:       &conn.Url,
			RequireApproval: ExternalOpenAIResponsesStandardToolRequireApprovalNever.ToPointer(),
		})
	}
}

func (r *ExternalOpenAIResponsesStandardRequest) AddMcpConnectionAsFunctionCallTools(
	mcpConnection *models.MCPConnection,
	relatedTools []*interfaces.McpTool,
) {
	if len(relatedTools) == 0 {
		return
	}

	if r.Tools == nil {
		r.Tools = make([]ExternalOpenAIResponsesStandardRequestTool, 0)
	}

	for _, mcpTool := range relatedTools {
		// MCPs are not always as well defined, so we set strict to false by default
		// But also not always all fields are required
		strict := false
		name := utils.McpToolToFunctionName(mcpConnection, mcpTool)

		toolToAdd := ExternalOpenAIResponsesStandardRequestTool{
			Type:        "function",
			Name:        &name,
			Strict:      &strict,
			Description: &mcpTool.Description,
			Parameters: &ExternalOpenAIResponsesStandardRequestToolFunctionParameters{
				Type:                 mcpTool.InputScheme.Type,
				Required:             mcpTool.InputScheme.Required,
				AdditionalProperties: mcpTool.InputScheme.AdditionalProperties,

				Properties: NewExternalToolInputSchemaPropertiesFromInternalModels(mcpTool.InputScheme.Properties),
			},
		}

		if toolToAdd.Parameters.Required == nil {
			toolToAdd.Parameters.Required = make([]string, 0)
		}

		r.Tools = append(r.Tools, toolToAdd)
	}
}
