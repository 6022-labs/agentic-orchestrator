package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"
	external_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector_http_inference/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/tests/llm_connector_http_inference_mocks/clients_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/tests/llm_connector_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRequestingAnthropicMessagesStandardChatCompletionTestingSuite struct {
	sut *services.AnthropicMessagesStandardChatCompletionRequester

	mockMcpRequester            *interfaces_mocks.MockMcpRequesterInterface
	mockAnthropicMessagesClient *clients_mocks.MockAnthropicMessagesStandardClientInterface
}

func WhenRequestingAnthropicMessagesStandardChatCompletionBeforeEach(t *testing.T) *WhenRequestingAnthropicMessagesStandardChatCompletionTestingSuite {
	mockController := gomock.NewController(t)

	mockMcpRequester := interfaces_mocks.NewMockMcpRequesterInterface(mockController)
	mockAnthropicMessagesClient := clients_mocks.NewMockAnthropicMessagesStandardClientInterface(mockController)

	sut := services.NewAnthropicMessagesStandardChatCompletionRequester(
		zap.NewNop(),
		mockMcpRequester,
		mockAnthropicMessagesClient,
	)

	return &WhenRequestingAnthropicMessagesStandardChatCompletionTestingSuite{
		sut:                         sut,
		mockMcpRequester:            mockMcpRequester,
		mockAnthropicMessagesClient: mockAnthropicMessagesClient,
	}
}

func TestWhenRequestingAnthropicMessagesStandardChatCompletion(t *testing.T) {
	t.Parallel()

	const (
		baseUrl = "https://anthropic.example.com"
		apiKey  = "test-key"
	)

	expectedContent := "hi there"

	llmConnection := model_builders.NewLLMConnectionBuilder().
		WithModel("claude-3-5-sonnet").
		Build()

	userContent := "hello"
	messages := []external_models.ExternalAnthropicMessagesStandardMessage{
		{
			Role: "user",
			Content: []external_models.ExternalAnthropicMessagesStandardMessageContent{
				{Type: "text", Text: &userContent},
			},
		},
	}

	t.Run("Given no tool calls are requested", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingAnthropicMessagesStandardChatCompletionTestingSuite) {
			suite.mockAnthropicMessagesClient.
				EXPECT().
				Messages(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalAnthropicMessagesStandardRequest) (*external_models.ExternalAnthropicMessagesStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Equal(t, messages, request.Messages)
					assert.Nil(t, request.Tools)
					assert.Nil(t, request.McpServers)

					return &external_models.ExternalAnthropicMessagesStandardResponse{
						ExternalAnthropicMessagesStandardMessage: external_models.ExternalAnthropicMessagesStandardMessage{
							Role: "assistant",
							Content: []external_models.ExternalAnthropicMessagesStandardMessageContent{
								{Type: "text", Text: &expectedContent},
							},
						},
					}, nil
				})
		}

		t.Run("Should return the last text content", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAnthropicMessagesStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, false, llmConnection, messages, nil, nil)

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, expectedContent, *result)
			}
		})
	})

	t.Run("Given retrieving MCP tools fails", func(t *testing.T) {
		t.Parallel()

		mcpConnection := model_builders.NewMCPConnectionBuilder().
			WithName("inventory").
			WithUrl("https://mcp.example.com").
			Build()
		mcpConnection.HeaderAuth = map[string]string{"authorization": "Bearer token"}

		initSuite := func(suite *WhenRequestingAnthropicMessagesStandardChatCompletionTestingSuite) {
			suite.mockMcpRequester.
				EXPECT().
				ListTools(mcpConnection).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAnthropicMessagesStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, false, llmConnection, messages, nil, []*models.MCPConnection{mcpConnection})
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.ErrorContains(t, err, fmt.Sprintf("failed to fetch MCP tools for connection '%s'", mcpConnection.Name))
		})
	})

	t.Run("Given the LLM requests an MCP tool call", func(t *testing.T) {
		t.Parallel()

		mcpConnection := model_builders.NewMCPConnectionBuilder().
			WithName("inventory").
			WithUrl("https://mcp.example.com").
			Build()
		mcpConnection.HeaderAuth = map[string]string{"authorization": "Bearer token"}

		initSuite := func(suite *WhenRequestingAnthropicMessagesStandardChatCompletionTestingSuite) {
			mcpTool := &interfaces.McpTool{
				Name:        "fetch_items",
				Description: "Fetch items",
				InputScheme: interfaces.McpToolInputScheme{
					Type: "object",
					Properties: map[string]interfaces.ToolInputSchemeProperty{
						"query": {Type: "string"},
					},
					Required:             []string{"query"},
					AdditionalProperties: false,
				},
			}
			functionName := fmt.Sprintf("%s_%s", mcpConnection.Name, mcpTool.Name)
			toolInput := map[string]any{"query": "value"}
			toolResult := map[string]any{"status": "ok"}

			suite.mockMcpRequester.
				EXPECT().
				ListTools(mcpConnection).
				Return([]*interfaces.McpTool{
					mcpTool,
				}, nil).
				Times(2)

			suite.mockMcpRequester.
				EXPECT().
				CallTool(mcpConnection, mcpTool, toolInput).
				Return(toolResult, nil)

			firstCall := suite.mockAnthropicMessagesClient.
				EXPECT().
				Messages(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalAnthropicMessagesStandardRequest) (*external_models.ExternalAnthropicMessagesStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Equal(t, messages, request.Messages)
					if assert.Len(t, request.Tools, 1) {
						assert.Equal(t, functionName, request.Tools[0].Name)
					}
					assert.Nil(t, request.McpServers)

					callId := "call_1"

					return &external_models.ExternalAnthropicMessagesStandardResponse{
						ExternalAnthropicMessagesStandardMessage: external_models.ExternalAnthropicMessagesStandardMessage{
							Role: "assistant",
							Content: []external_models.ExternalAnthropicMessagesStandardMessageContent{
								{
									Type:  "tool_use",
									Id:    &callId,
									Name:  &functionName,
									Input: toolInput,
								},
							},
						},
					}, nil
				})

			secondCall := suite.mockAnthropicMessagesClient.
				EXPECT().
				Messages(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalAnthropicMessagesStandardRequest) (*external_models.ExternalAnthropicMessagesStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Len(t, request.Messages, len(messages)+2)

					assistantMessage := request.Messages[len(messages)]
					assert.Equal(t, "assistant", assistantMessage.Role)
					if assert.Len(t, assistantMessage.Content, 1) {
						c := assistantMessage.Content[0]
						assert.Equal(t, "tool_use", c.Type)
						if assert.NotNil(t, c.Name) {
							assert.Equal(t, functionName, *c.Name)
						}
						assert.NotNil(t, c.Input)
					}

					toolOutputMessage := request.Messages[len(messages)+1]
					assert.Equal(t, "user", toolOutputMessage.Role)
					if assert.Len(t, toolOutputMessage.Content, 1) {
						c := toolOutputMessage.Content[0]
						assert.Equal(t, "tool_result", c.Type)
						if assert.NotNil(t, c.ToolUseId) {
							assert.Equal(t, "call_1", *c.ToolUseId)
						}

						sp, ok := c.Content.(*string)
						if assert.True(t, ok) && assert.NotNil(t, sp) {
							assert.Contains(t, *sp, "\"status\"")
							assert.Contains(t, *sp, "\"ok\"")
						}
					}

					return &external_models.ExternalAnthropicMessagesStandardResponse{
						ExternalAnthropicMessagesStandardMessage: external_models.ExternalAnthropicMessagesStandardMessage{
							Role: "assistant",
							Content: []external_models.ExternalAnthropicMessagesStandardMessageContent{
								{Type: "text", Text: &expectedContent},
							},
						},
					}, nil
				})

			gomock.InOrder(firstCall, secondCall)
		}

		t.Run("Should call the MCP tool and return the final completion", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAnthropicMessagesStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, false, llmConnection, messages, nil, []*models.MCPConnection{mcpConnection})
			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, expectedContent, *result)
			}
		})
	})
}
