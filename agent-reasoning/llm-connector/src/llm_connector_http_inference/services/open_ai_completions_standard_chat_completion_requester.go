package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_inference/clients"
	orchestrator_agent_llm_connector_http_inference_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_inference/utils"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type OpenAICompletionsStandardChatCompletionRequesterInterface interface {
	ChatCompletion(
		baseUrl, apiKey string,
		llmConnection *models.LLMConnection,
		messages []orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAICompletionsStandardRequestMessage,
		mcpConnections []*models.MCPConnection,
	) (*string, error)
}

type OpenAICompletionsStandardChatCompletionRequester struct {
	logger                  *zap.Logger
	mcpRequester            interfaces.McpRequesterInterface
	openAICompletionsClient clients.OpenAICompletionsStandardClientInterface
}

func NewOpenAICompletionsStandardChatCompletionRequester(
	logger *zap.Logger,
	mcpRequester interfaces.McpRequesterInterface,
	openAICompletionsClient clients.OpenAICompletionsStandardClientInterface,
) *OpenAICompletionsStandardChatCompletionRequester {
	return &OpenAICompletionsStandardChatCompletionRequester{
		logger:                  logger,
		mcpRequester:            mcpRequester,
		openAICompletionsClient: openAICompletionsClient,
	}
}

func (r *OpenAICompletionsStandardChatCompletionRequester) ChatCompletion(
	baseUrl, apiKey string,
	llmConnection *models.LLMConnection,
	messages []orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAICompletionsStandardRequestMessage,
	mcpConnections []*models.MCPConnection,
) (*string, error) {
	externalRequest := orchestrator_agent_llm_connector_http_inference_models.NewExternalOpenAICompletionsStandardRequest(llmConnection.Model, messages)

	mcpsToolsAsFunc := make(map[*models.MCPConnection][]*interfaces.McpTool)

	for _, mcpConnection := range mcpConnections {
		mcpTools, err := r.mcpRequester.ListTools(mcpConnection)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch MCP tools for connection '%s': %w", mcpConnection.Name, err)
		}

		mcpsToolsAsFunc[mcpConnection] = mcpTools
	}

	for mcpConnection, tools := range mcpsToolsAsFunc {
		externalRequest.AddMcpConnectionAsFunctionCallTools(mcpConnection, tools)
	}

	r.logger.Debug("Request to Completions endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
	)

	response, err := r.openAICompletionsClient.ChatCompletion(baseUrl, apiKey, externalRequest)
	if err != nil {
		return nil, err
	}

	toolCalls := lo.FilterMap(
		response.Choices,
		func(item orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAICompletionsStandardResponseChoice, _ int) ([]orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAICompletionsStandardResponseChoiceMessageToolCall, bool) {
			return item.Message.ToolCalls, len(item.Message.ToolCalls) > 0
		},
	)

	r.logger.Debug("Response from Completions endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
		zap.Int("toolUseCount", len(toolCalls)),
	)

	flattenToolCalls := lo.Flatten(toolCalls)

	// For the moment, tool calls are only related to MCP.
	// Try to find the related MCP of this tool call, call it and re-call the LLM with the tool result.
	if len(flattenToolCalls) > 0 {
		for _, toolCall := range flattenToolCalls {
			messages = append(messages, orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAICompletionsStandardRequestMessage{
				Role: "assistant",
				ToolCalls: []orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAICompletionsStandardRequestMessageToolCall{
					{
						ID:       toolCall.ID,
						Type:     "function", // Can't use toolCall.Type here because some API didn't fill it properly (Mistral...)
						Function: toolCall.Function,
					},
				},
			})

			relatedMcpConnection, mcpTool, err := utils.ExtractMcpInfoFromFunctionName(toolCall.Function.Name, mcpConnections, mcpsToolsAsFunc)
			if err != nil {
				return nil, fmt.Errorf("failed to map function name to MCP tool: %w", err)
			}

			r.logger.Debug("Calling MCP tool from tool call",
				zap.String("agentId", llmConnection.AgentID.String()),
				zap.String("mcpToolName", mcpTool.Name),
				zap.String("mcpConnectionName", relatedMcpConnection.Name),
			)

			parsedArguments, err := toolCall.Function.ParsedArguments()
			if err != nil {
				return nil, fmt.Errorf("failed to parse tool call arguments: %w", err)
			}

			toolResult, err := r.mcpRequester.CallTool(relatedMcpConnection, mcpTool, parsedArguments)
			if err != nil {
				return nil, fmt.Errorf("failed to call tool from MCP: %w", err)
			}

			// Append messages to include tool result and re-call complete
			toolResultMessage, err := orchestrator_agent_llm_connector_http_inference_models.NewExternalOpenAICompletionsStandardRequestMessageFromToolCallOutput(
				toolCall.ID,
				toolResult,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create message from tool call output: %w", err)
			}

			messages = append(messages, *toolResultMessage)
		}

		// Then re-call the complete function with updated messages
		return r.ChatCompletion(
			baseUrl,
			apiKey,
			llmConnection,
			messages,
			mcpConnections,
		)
	}

	if len(response.Choices) == 0 {
		r.logger.Warn(
			"No choices found in response",
			zap.String("agentId", llmConnection.AgentID.String()),
		)

		return nil, nil
	}

	return &response.Choices[0].Message.Content, nil
}
