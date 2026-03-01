package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/clients"
	orchestrator_agent_llm_connector_http_inference_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/utils"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type AmazonBedrockConverseStandardChatCompletionRequesterInterface interface {
	ChatCompletion(
		baseUrl, apiKey string,
		llmConnection *models.LLMConnection,
		messages []orchestrator_agent_llm_connector_http_inference_models.ExternalAmazonBedrockConverseStandardMessage,
		messagesSystem []orchestrator_agent_llm_connector_http_inference_models.ExternalAmazonBedrockConverseStandardRequestSystem,
		mcpConnections []*models.MCPConnection,
	) (*string, error)
}

type AmazonBedrockConverseStandardChatCompletionRequester struct {
	logger                      *zap.Logger
	mcpRequester                interfaces.McpRequesterInterface
	amazonBedrockConverseClient clients.AmazonBedrockConverseStandardClientInterface
}

func NewAmazonBedrockConverseStandardChatCompletionRequester(
	logger *zap.Logger,
	mcpRequester interfaces.McpRequesterInterface,
	amazonBedrockConverseClient clients.AmazonBedrockConverseStandardClientInterface,
) *AmazonBedrockConverseStandardChatCompletionRequester {
	return &AmazonBedrockConverseStandardChatCompletionRequester{
		logger:                      logger,
		mcpRequester:                mcpRequester,
		amazonBedrockConverseClient: amazonBedrockConverseClient,
	}
}

func (r *AmazonBedrockConverseStandardChatCompletionRequester) ChatCompletion(
	baseUrl, apiKey string,
	llmConnection *models.LLMConnection,
	messages []orchestrator_agent_llm_connector_http_inference_models.ExternalAmazonBedrockConverseStandardMessage,
	messagesSystem []orchestrator_agent_llm_connector_http_inference_models.ExternalAmazonBedrockConverseStandardRequestSystem,
	mcpConnections []*models.MCPConnection,
) (*string, error) {
	externalRequest := orchestrator_agent_llm_connector_http_inference_models.NewExternalAmazonBedrockConverseStandardRequest(llmConnection.Model, messages, messagesSystem)

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

	r.logger.Debug("Request to Amazon Bedrock Converse Messages endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
	)

	response, err := r.amazonBedrockConverseClient.Messages(baseUrl, apiKey, externalRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call Amazon Bedrock Converse API: %w", err)
	}

	toolUseContents := lo.Filter(response.Output.Message.Content, func(content orchestrator_agent_llm_connector_http_inference_models.ExternalAmazonBedrockConverseStandardMessageContent, _ int) bool {
		return content.ToolUse != nil
	})

	r.logger.Debug("Response from Amazon Bedrock Converse Messages endpoint",
		zap.String("agentId", llmConnection.AgentID.String()),
		zap.Int("toolUseCount", len(toolUseContents)),
	)

	if len(toolUseContents) > 0 {
		messages = append(messages, response.Output.Message)

		toolUseResults := make(map[string]map[string]any)

		for _, toolCall := range toolUseContents {
			relatedMcpConnection, mcpTool, err := utils.ExtractMcpInfoFromFunctionName(toolCall.ToolUse.Name, mcpConnections, mcpsToolsAsFunc)
			if err != nil {
				return nil, fmt.Errorf("failed to map function name to MCP tool: %w", err)
			}

			r.logger.Debug("Calling MCP tool from tool use",
				zap.String("agentId", llmConnection.AgentID.String()),
				zap.String("mcpToolName", mcpTool.Name),
				zap.String("mcpConnectionName", relatedMcpConnection.Name),
			)

			toolResult, err := r.mcpRequester.CallTool(relatedMcpConnection, mcpTool, toolCall.ToolUse.Input)
			if err != nil {
				return nil, fmt.Errorf("failed to call tool from MCP: %w", err)
			}

			toolUseResults[toolCall.ToolUse.ToolUseId] = toolResult
		}

		// Append messages to include tool result and re-call complete
		messages = append(messages, orchestrator_agent_llm_connector_http_inference_models.NewExternalAmazonBedrockConverseStandardMessageFromToolCallOutputs(
			toolUseResults,
		))

		// Then re-call the complete function with updated messages
		return r.ChatCompletion(
			baseUrl,
			apiKey,
			llmConnection,
			messages,
			messagesSystem,
			mcpConnections,
		)
	}

	textContents := lo.Filter(response.Output.Message.Content, func(content orchestrator_agent_llm_connector_http_inference_models.ExternalAmazonBedrockConverseStandardMessageContent, _ int) bool {
		return content.Text != nil
	})

	if len(textContents) == 0 {
		r.logger.Warn(
			"No content with type 'text' found in response",
			zap.String("agentId", llmConnection.AgentID.String()),
		)

		return nil, nil
	}

	return textContents[len(textContents)-1].Text, nil
}
