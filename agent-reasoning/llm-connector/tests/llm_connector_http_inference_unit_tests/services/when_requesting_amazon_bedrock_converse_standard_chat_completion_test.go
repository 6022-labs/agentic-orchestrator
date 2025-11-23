package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"
	external_models "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/llm_connector_http_inference_mocks/clients_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/llm_connector_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRequestingAmazonBedrockConverseStandardChatCompletionTestingSuite struct {
	sut *services.AmazonBedrockConverseStandardChatCompletionRequester

	mockMcpRequester                *interfaces_mocks.MockMcpRequesterInterface
	mockAmazonBedrockConverseClient *clients_mocks.MockAmazonBedrockConverseStandardClientInterface
}

func WhenRequestingAmazonBedrockConverseStandardChatCompletionBeforeEach(t *testing.T) *WhenRequestingAmazonBedrockConverseStandardChatCompletionTestingSuite {
	mockController := gomock.NewController(t)

	mockMcpRequester := interfaces_mocks.NewMockMcpRequesterInterface(mockController)
	mockAmazonBedrockConverseClient := clients_mocks.NewMockAmazonBedrockConverseStandardClientInterface(mockController)

	sut := services.NewAmazonBedrockConverseStandardChatCompletionRequester(
		zap.NewNop(),
		mockMcpRequester,
		mockAmazonBedrockConverseClient,
	)

	return &WhenRequestingAmazonBedrockConverseStandardChatCompletionTestingSuite{
		sut:                             sut,
		mockMcpRequester:                mockMcpRequester,
		mockAmazonBedrockConverseClient: mockAmazonBedrockConverseClient,
	}
}

func TestWhenRequestingAmazonBedrockConverseStandardChatCompletion(t *testing.T) {
	t.Parallel()

	const (
		baseUrl = "https://bedrock.example.com"
		apiKey  = "test-key"
	)

	expectedContent := "hi there"

	llmConnection := model_builders.NewLLMConnectionBuilder().
		WithModel("amazon.titan-text-g1").
		Build()

	userText := "hello"
	messages := []external_models.ExternalAmazonBedrockConverseStandardMessage{
		{
			Role: "user",
			Content: []external_models.ExternalAmazonBedrockConverseStandardMessageContent{
				{Text: &userText},
			},
		},
	}

	t.Run("Given no tool calls are requested", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingAmazonBedrockConverseStandardChatCompletionTestingSuite) {
			suite.mockAmazonBedrockConverseClient.
				EXPECT().
				Messages(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalAmazonBedrockConverseStandardRequest) (*external_models.ExternalAmazonBedrockConverseStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.ModelId)
					if assert.Len(t, request.Messages, 1) {
						assert.Equal(t, "user", request.Messages[0].Role)
						if assert.Len(t, request.Messages[0].Content, 1) {
							assert.NotNil(t, request.Messages[0].Content[0].Text)
							assert.Equal(t, userText, *request.Messages[0].Content[0].Text)
						}
					}

					mid := "intermediate"
					last := expectedContent
					return &external_models.ExternalAmazonBedrockConverseStandardResponse{
						Output: external_models.ExternalAmazonBedrockConverseStandardResponseOutput{
							Message: external_models.ExternalAmazonBedrockConverseStandardMessage{
								Role: "assistant",
								Content: []external_models.ExternalAmazonBedrockConverseStandardMessageContent{
									{Text: &mid},
									{Text: &last},
								},
							},
						},
					}, nil
				})
		}

		t.Run("Should return the last text content", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAmazonBedrockConverseStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, llmConnection, messages, nil, nil)
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

		initSuite := func(suite *WhenRequestingAmazonBedrockConverseStandardChatCompletionTestingSuite) {
			suite.mockMcpRequester.
				EXPECT().
				ListTools(mcpConnection).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAmazonBedrockConverseStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, llmConnection, messages, nil, []*models.MCPConnection{mcpConnection})
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

		initSuite := func(suite *WhenRequestingAmazonBedrockConverseStandardChatCompletionTestingSuite) {
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

			firstCall := suite.mockAmazonBedrockConverseClient.
				EXPECT().
				Messages(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalAmazonBedrockConverseStandardRequest) (*external_models.ExternalAmazonBedrockConverseStandardResponse, error) {
					if assert.NotNil(t, request.ToolConfig) {
						assert.NotNil(t, request.ToolConfig.ToolChoice.Auto)
						if assert.Len(t, request.ToolConfig.Tools, 1) {
							assert.Equal(t, functionName, request.ToolConfig.Tools[0].ToolSpec.Name)
						}
					}

					return &external_models.ExternalAmazonBedrockConverseStandardResponse{
						Output: external_models.ExternalAmazonBedrockConverseStandardResponseOutput{
							Message: external_models.ExternalAmazonBedrockConverseStandardMessage{
								Role: "assistant",
								Content: []external_models.ExternalAmazonBedrockConverseStandardMessageContent{
									{
										ToolUse: &external_models.ExternalAmazonBedrockConverseStandardMessageContentToolUse{
											Name:      functionName,
											Input:     toolInput,
											ToolUseId: "call_1",
										},
									},
								},
							},
						},
					}, nil
				})

			secondCall := suite.mockAmazonBedrockConverseClient.
				EXPECT().
				Messages(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalAmazonBedrockConverseStandardRequest) (*external_models.ExternalAmazonBedrockConverseStandardResponse, error) {
					assert.GreaterOrEqual(t, len(request.Messages), 2)

					last := expectedContent
					return &external_models.ExternalAmazonBedrockConverseStandardResponse{
						Output: external_models.ExternalAmazonBedrockConverseStandardResponseOutput{
							Message: external_models.ExternalAmazonBedrockConverseStandardMessage{
								Role: "assistant",
								Content: []external_models.ExternalAmazonBedrockConverseStandardMessageContent{
									{Text: &last},
								},
							},
						},
					}, nil
				})

			gomock.InOrder(firstCall, secondCall)
		}

		t.Run("Should call the MCP tool and return the final completion", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingAmazonBedrockConverseStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, llmConnection, messages, nil, []*models.MCPConnection{mcpConnection})
			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, expectedContent, *result)
			}
		})
	})
}
