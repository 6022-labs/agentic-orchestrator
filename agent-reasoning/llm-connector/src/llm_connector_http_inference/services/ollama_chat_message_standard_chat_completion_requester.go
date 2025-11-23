package services

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/clients"
	orchestrator_agent_llm_connector_http_inference_models "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/utils"
	"go.uber.org/zap"
)

type OllamaChatMessageStandardChatCompletionRequesterInterface interface {
	ChatCompletion(
		baseUrl string,
		llmConnection *models.LLMConnection,
		messages []orchestrator_agent_llm_connector_http_inference_models.ExternalOllamaChatMessageStandardMessage,
		mcpConnections []*models.MCPConnection,
	) (*string, error)
}

type OllamaChatMessageStandardChatCompletionRequester struct {
	logger                  *zap.Logger
	mcpRequester            interfaces.McpRequesterInterface
	ollamaChatMessageClient clients.OllamaChatMessageStandardClientInterface
}

func NewOllamaChatMessageStandardChatCompletionRequester(
	logger *zap.Logger,
	mcpRequester interfaces.McpRequesterInterface,
	ollamaChatMessageClient clients.OllamaChatMessageStandardClientInterface,
) *OllamaChatMessageStandardChatCompletionRequester {
	return &OllamaChatMessageStandardChatCompletionRequester{
		logger:                  logger,
		mcpRequester:            mcpRequester,
		ollamaChatMessageClient: ollamaChatMessageClient,
	}
}

func (r *OllamaChatMessageStandardChatCompletionRequester) ChatCompletion(
	baseUrl string,
	llmConnection *models.LLMConnection,
	messages []orchestrator_agent_llm_connector_http_inference_models.ExternalOllamaChatMessageStandardMessage,
	mcpConnections []*models.MCPConnection,
) (*string, error) {
	externalRequest := orchestrator_agent_llm_connector_http_inference_models.NewExternalOllamaChatMessageStandardRequest(llmConnection.Model, messages)

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

	r.logger.Debug("Request to Chat Message endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
	)

	response, err := r.ollamaChatMessageClient.ChatMessage(baseUrl, externalRequest)
	if err != nil {
		return nil, err
	}

	toolCalls := response.Message.ToolCalls

	r.logger.Debug("Response from Chat Message endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
		zap.Int("toolUseCount", len(toolCalls)),
	)

	// For the moment, tool calls are only related to MCP.
	// Try to find the related MCP of this tool call, call it and re-call the LLM with the tool result.
	if len(toolCalls) > 0 {
		messages = append(messages, response.Message)

		for _, toolCall := range toolCalls {
			relatedMcpConnection, mcpTool, err := utils.ExtractMcpInfoFromFunctionName(toolCall.Function.Name, mcpConnections, mcpsToolsAsFunc)
			if err != nil {
				return nil, fmt.Errorf("failed to map function name to MCP tool: %w", err)
			}

			r.logger.Debug("Calling MCP tool from tool call",
				zap.String("agentId", llmConnection.AgentID.String()),
				zap.String("mcpToolName", mcpTool.Name),
				zap.String("mcpConnectionName", relatedMcpConnection.Name),
			)

			toolResult, err := r.mcpRequester.CallTool(relatedMcpConnection, mcpTool, toolCall.Function.Arguments)
			if err != nil {
				return nil, fmt.Errorf("failed to call tool from MCP: %w", err)
			}

			// Append messages to include tool result and re-call complete
			toolResultMessage, err := orchestrator_agent_llm_connector_http_inference_models.NewExternalOllamaChatMessageStandardMessageFromToolCallOutput(
				mcpTool.Name,
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
			llmConnection,
			messages,
			mcpConnections,
		)
	}

	if len(response.Message.Content) == 0 {
		r.logger.Warn(
			"Empty content in response message",
			zap.String("agentId", llmConnection.AgentID.String()),
		)

		return nil, nil
	}

	return &response.Message.Content, nil
}
