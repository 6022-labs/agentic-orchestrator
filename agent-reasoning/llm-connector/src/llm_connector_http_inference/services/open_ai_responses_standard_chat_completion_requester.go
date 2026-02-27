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

type OpenAIResponsesStandardChatCompletionRequesterInterface interface {
	ChatCompletion(
		baseUrl, apiKey string,
		hasDirectAccessToMcp bool,
		llmConnection *models.LLMConnection,
		inputs []orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAIResponsesStandardRequestInput,
		mcpConnections []*models.MCPConnection,
		previousResponseId *string,
	) (*string, error)
}

type OpenAIResponsesStandardChatCompletionRequester struct {
	logger                *zap.Logger
	mcpRequester          interfaces.McpRequesterInterface
	openAIResponsesClient clients.OpenAIResponsesStandardClientInterface
}

func NewOpenAIResponsesStandardChatCompletionRequester(
	logger *zap.Logger,
	mcpRequester interfaces.McpRequesterInterface,
	openAIResponsesClient clients.OpenAIResponsesStandardClientInterface,
) *OpenAIResponsesStandardChatCompletionRequester {
	return &OpenAIResponsesStandardChatCompletionRequester{
		logger:                logger,
		mcpRequester:          mcpRequester,
		openAIResponsesClient: openAIResponsesClient,
	}
}

func (r *OpenAIResponsesStandardChatCompletionRequester) ChatCompletion(
	baseUrl, apiKey string,
	hasDirectAccessToMcp bool,
	llmConnection *models.LLMConnection,
	inputs []orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAIResponsesStandardRequestInput,
	mcpConnections []*models.MCPConnection,
	previousResponseId *string,
) (*string, error) {
	externalRequest := orchestrator_agent_llm_connector_http_inference_models.NewExternalOpenAIResponsesStandardRequest(llmConnection.Model, inputs, previousResponseId)

	mcpsToolsAsFunc := make(map[*models.MCPConnection][]*interfaces.McpTool)

	if hasDirectAccessToMcp {
		externalRequest.AddMcpConnectionsAsMcpTools(mcpConnections)
	} else {
		for _, mcpConnection := range mcpConnections {
			mcpTools, err := r.mcpRequester.ListTools(mcpConnection)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch MCP tools for connection '%s': %w", mcpConnection.Name, err)
			}

			mcpsToolsAsFunc[mcpConnection] = mcpTools
		}
	}

	for mcpConnection, tools := range mcpsToolsAsFunc {
		externalRequest.AddMcpConnectionAsFunctionCallTools(mcpConnection, tools)
	}

	r.logger.Debug("Request to Responses endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
	)

	response, err := r.openAIResponsesClient.Responses(baseUrl, apiKey, externalRequest)
	if err != nil {
		return nil, err
	}

	functionCallOutputs := lo.Filter(response.Output, func(output orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAIResponsesStandardResponseOutput, _ int) bool {
		return output.Type == "function_call"
	})

	r.logger.Debug("Response from Responses endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
		zap.Int("toolUseCount", len(functionCallOutputs)),
	)

	// For the moment, function call are only related to MCP.
	// Try to find the related MCP of this function call, call it and re-call the LLM with the tool result.
	if len(functionCallOutputs) > 0 {
		newInputs := make([]orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAIResponsesStandardRequestInput, 0)

		for _, output := range functionCallOutputs {
			relatedMcpConnection, mcpTool, err := utils.ExtractMcpInfoFromFunctionName(*output.Name, mcpConnections, mcpsToolsAsFunc)
			if err != nil {
				return nil, fmt.Errorf("failed to map function name to MCP tool: %w", err)
			}

			r.logger.Debug("Calling MCP tool from function call",
				zap.String("agentId", llmConnection.AgentID.String()),
				zap.String("mcpToolName", mcpTool.Name),
				zap.String("mcpConnectionName", relatedMcpConnection.Name),
			)

			parsedArguments, err := output.ParsedArguments()
			if err != nil {
				return nil, fmt.Errorf("failed to parse function call arguments: %w", err)
			}

			toolResult, err := r.mcpRequester.CallTool(relatedMcpConnection, mcpTool, parsedArguments)
			if err != nil {
				return nil, fmt.Errorf("failed to call tool from MCP: %w", err)
			}

			// Append messages to include tool result and re-call complete
			toolResultMessage, err := orchestrator_agent_llm_connector_http_inference_models.NewExternalOpenAIResponsesStandardRequestFromToolCallOutput(
				*output.CallId,
				toolResult,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create message from tool call output: %w", err)
			}

			newInputs = append(newInputs, *toolResultMessage)
		}

		// Then re-call the complete function with updated messages
		return r.ChatCompletion(
			baseUrl,
			apiKey,
			hasDirectAccessToMcp,
			llmConnection,
			newInputs,
			mcpConnections,
			&response.ID,
		)
	}

	messageOutputs := lo.Filter(response.Output, func(output orchestrator_agent_llm_connector_http_inference_models.ExternalOpenAIResponsesStandardResponseOutput, _ int) bool {
		return output.Type == "message"
	})

	if len(messageOutputs) == 0 {
		r.logger.Warn(
			"No output with type 'message' found in response",
			zap.String("agentId", llmConnection.AgentID.String()),
		)
		return nil, nil
	}

	if len(messageOutputs) > 1 {
		r.logger.Warn(
			"Multiple message outputs found, taking the first one",
			zap.String("agentId", llmConnection.AgentID.String()),
		)
	}

	return &messageOutputs[0].Content[0].Text, nil
}
