package services

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/clients"
	orchestrator_agent_compute_http_inference_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/utils"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type AnthropicMessagesStandardChatCompletionRequesterInterface interface {
	ChatCompletion(
		baseUrl, apiKey string,
		hasDirectAccessToMcp bool,
		llmConnection *models.LLMConnection,
		messages []orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardMessage,
		systems []orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardRequestSystem,
		mcpConnections []*models.MCPConnection,
	) (*string, error)
}

type AnthropicMessagesStandardChatCompletionRequester struct {
	logger                  *zap.Logger
	mcpRequester            interfaces.McpRequesterInterface
	anthropicMessagesClient clients.AnthropicMessagesStandardClientInterface
}

func NewAnthropicMessagesStandardChatCompletionRequester(
	logger *zap.Logger,
	mcpRequester interfaces.McpRequesterInterface,
	anthropicMessagesClient clients.AnthropicMessagesStandardClientInterface,
) *AnthropicMessagesStandardChatCompletionRequester {
	return &AnthropicMessagesStandardChatCompletionRequester{
		logger:                  logger,
		mcpRequester:            mcpRequester,
		anthropicMessagesClient: anthropicMessagesClient,
	}
}

func (r *AnthropicMessagesStandardChatCompletionRequester) ChatCompletion(
	baseUrl, apiKey string,
	hasDirectAccessToMcp bool,
	llmConnection *models.LLMConnection,
	messages []orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardMessage,
	systems []orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardRequestSystem,
	mcpConnections []*models.MCPConnection,
) (*string, error) {
	externalRequest := orchestrator_agent_compute_http_inference_models.NewExternalAnthropicMessagesStandardRequest(
		llmConnection.Model,
		messages,
		systems,
	)

	alreadyAddedMcps := make(map[string]bool)

	if hasDirectAccessToMcp {
		// Even if the platform has direct access to MCP
		// Anthropic only support MCP servers with OAuth "Authorization" header auth
		for _, mcpConnection := range mcpConnections {
			auth, ok := mcpConnection.HeaderAuth["Authorization"]
			if len(mcpConnection.HeaderAuth) == 0 || (ok && len(mcpConnection.HeaderAuth) == 1 && strings.HasPrefix(auth, "Bearer ")) {
				externalRequest.AddMcpServerFromMcpConnection(mcpConnection)
				alreadyAddedMcps[mcpConnection.Name] = true
			}
		}
	}

	mcpsToolsAsFunc := make(map[*models.MCPConnection][]*interfaces.McpTool)

	// Add remaining MCP connections as function call tools
	for _, mcpConnection := range mcpConnections {
		if _, alreadyAdded := alreadyAddedMcps[mcpConnection.Name]; alreadyAdded {
			continue
		}

		mcpTools, err := r.mcpRequester.ListTools(mcpConnection)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch MCP tools for connection '%s': %w", mcpConnection.Name, err)
		}

		mcpsToolsAsFunc[mcpConnection] = mcpTools
	}

	for mcpConnection, tools := range mcpsToolsAsFunc {
		externalRequest.AddMcpConnectionAsFunctionCallTools(mcpConnection, tools)
	}

	r.logger.Debug("Request to Anthropic Messages endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
	)

	response, err := r.anthropicMessagesClient.Messages(baseUrl, apiKey, externalRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call Anthropic Messages API: %w", err)
	}

	toolUseContents := lo.Filter(response.Content, func(content orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardMessageContent, _ int) bool {
		return content.Type == "tool_use"
	})

	r.logger.Debug("Response from Anthropic Messages endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
		zap.Int("contentItems", len(response.Content)),
		zap.Int("toolUseCount", len(toolUseContents)),
		zap.Strings("contentTypes", lo.Map(response.Content, func(content orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardMessageContent, _ int) string {
			return content.Type
		})),
	)

	if len(toolUseContents) > 0 {
		messages = append(messages, response.ExternalAnthropicMessagesStandardMessage)

		for _, toolCall := range toolUseContents {
			relatedMcpConnection, mcpTool, err := utils.ExtractMcpInfoFromFunctionName(*toolCall.Name, mcpConnections, mcpsToolsAsFunc)
			if err != nil {
				return nil, fmt.Errorf("failed to map function name to MCP tool: %w", err)
			}

			r.logger.Debug("Calling MCP tool from tool use",
				zap.String("agentId", llmConnection.AgentID.String()),
				zap.String("mcpToolName", mcpTool.Name),
				zap.String("mcpConnectionName", relatedMcpConnection.Name),
			)

			toolResult, err := r.mcpRequester.CallTool(relatedMcpConnection, mcpTool, toolCall.Input)
			if err != nil {
				return nil, fmt.Errorf("failed to call tool from MCP: %w", err)
			}

			// Append messages to include tool result and re-call complete
			toolResultMessage, err := orchestrator_agent_compute_http_inference_models.NewExternalAnthropicMessagesStandardRequestMessagesFromToolCallOutput(
				*toolCall.Id,
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
			hasDirectAccessToMcp,
			llmConnection,
			messages,
			systems,
			mcpConnections,
		)
	}

	textContents := lo.Filter(response.Content, func(content orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardMessageContent, _ int) bool {
		return content.Type == "text"
	})

	if len(textContents) == 0 {
		r.logger.Warn(
			"No content with type 'text' found in response",
			zap.String("agentId", llmConnection.AgentID.String()),
			zap.Stringp("stopReason", response.StopReason),
			zap.Int("contentItems", len(response.Content)),
			zap.Strings("contentTypes", lo.Map(response.Content, func(content orchestrator_agent_compute_http_inference_models.ExternalAnthropicMessagesStandardMessageContent, _ int) string {
				return content.Type
			})),
		)

		return nil, nil
	}

	return textContents[len(textContents)-1].Text, nil
}
