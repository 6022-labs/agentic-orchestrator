package services_test

import (
	"encoding/json"
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

type WhenRequestingOpenAICompletionsStandardChatCompletionTestingSuite struct {
	sut *services.OpenAICompletionsStandardChatCompletionRequester

	mockMcpRequester            *interfaces_mocks.MockMcpRequesterInterface
	mockOpenAICompletionsClient *clients_mocks.MockOpenAICompletionsStandardClientInterface
}

func WhenRequestingOpenAICompletionsStandardChatCompletionBeforeEach(t *testing.T) *WhenRequestingOpenAICompletionsStandardChatCompletionTestingSuite {
	mockController := gomock.NewController(t)

	mockMcpRequester := interfaces_mocks.NewMockMcpRequesterInterface(mockController)
	mockOpenAICompletionsClient := clients_mocks.NewMockOpenAICompletionsStandardClientInterface(mockController)

	sut := services.NewOpenAICompletionsStandardChatCompletionRequester(
		zap.NewNop(),
		mockMcpRequester,
		mockOpenAICompletionsClient,
	)

	return &WhenRequestingOpenAICompletionsStandardChatCompletionTestingSuite{
		sut:                         sut,
		mockMcpRequester:            mockMcpRequester,
		mockOpenAICompletionsClient: mockOpenAICompletionsClient,
	}
}

func TestWhenRequestingOpenAICompletionsStandardChatCompletion(t *testing.T) {
	t.Parallel()

	const (
		baseUrl = "https://openai.example.com"
		apiKey  = "test-key"

		expectedContent = "hi there"
	)

	llmConnection := model_builders.NewLLMConnectionBuilder().
		WithModel("gpt-4").
		Build()

	userMessage := "hello"
	messages := []external_models.ExternalOpenAICompletionsStandardRequestMessage{
		{
			Role:    "user",
			Content: &userMessage,
		},
	}

	t.Run("Given no tool calls are requested", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingOpenAICompletionsStandardChatCompletionTestingSuite) {

			suite.mockOpenAICompletionsClient.
				EXPECT().
				ChatCompletion(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalOpenAICompletionsStandardRequest) (*external_models.ExternalOpenAICompletionsStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Equal(t, messages, request.Messages)
					assert.Nil(t, request.Tools)

					response := &external_models.ExternalOpenAICompletionsStandardResponse{
						Choices: []external_models.ExternalOpenAICompletionsStandardResponseChoice{
							{
								Message: external_models.ExternalOpenAICompletionsStandardResponseChoiceMessage{
									Role:    "assistant",
									Content: expectedContent,
								},
							},
						},
					}

					return response, nil
				})

		}

		t.Run("Should return the first choice content", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOpenAICompletionsStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, llmConnection, messages, nil)
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

		initSuite := func(suite *WhenRequestingOpenAICompletionsStandardChatCompletionTestingSuite) {
			suite.mockMcpRequester.
				EXPECT().
				ListTools(mcpConnection).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOpenAICompletionsStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, llmConnection, messages, []*models.MCPConnection{mcpConnection})
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

		initSuite := func(suite *WhenRequestingOpenAICompletionsStandardChatCompletionTestingSuite) {
			mcpTool := &interfaces.McpTool{
				Name:        "fetch_items",
				Description: "Fetch items",
				InputScheme: interfaces.McpToolInputScheme{
					Type:                 "object",
					Properties:           map[string]interfaces.ToolInputSchemeProperty{},
					Required:             []string{},
					AdditionalProperties: false,
				},
			}
			functionName := fmt.Sprintf("%s_%s", mcpConnection.Name, mcpTool.Name)
			toolArguments := map[string]any{"query": "value"}
			toolArgumentsJson, _ := json.Marshal(toolArguments)

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
				CallTool(mcpConnection, mcpTool, toolArguments).
				Return(toolResult, nil)

			firstCall := suite.mockOpenAICompletionsClient.
				EXPECT().
				ChatCompletion(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalOpenAICompletionsStandardRequest) (*external_models.ExternalOpenAICompletionsStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Equal(t, messages, request.Messages)
					if assert.Len(t, request.Tools, 1) {
						assert.Equal(t, functionName, request.Tools[0].Function.Name)
					}

					response := &external_models.ExternalOpenAICompletionsStandardResponse{
						Choices: []external_models.ExternalOpenAICompletionsStandardResponseChoice{
							{
								Message: external_models.ExternalOpenAICompletionsStandardResponseChoiceMessage{
									Role: "assistant",
									ToolCalls: []external_models.ExternalOpenAICompletionsStandardResponseChoiceMessageToolCall{
										{
											ID: "call_1",
											Function: external_models.ExternalOpenAICompletionsStandardToolCallFunction{
												Name:      functionName,
												Arguments: string(toolArgumentsJson),
											},
										},
									},
								},
							},
						},
					}

					return response, nil
				})

			secondCall := suite.mockOpenAICompletionsClient.
				EXPECT().
				ChatCompletion(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalOpenAICompletionsStandardRequest) (*external_models.ExternalOpenAICompletionsStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Len(t, request.Messages, len(messages)+2)

					assistantToolMessage := request.Messages[len(messages)]
					assert.Equal(t, "assistant", assistantToolMessage.Role)
					if assert.Len(t, assistantToolMessage.ToolCalls, 1) {
						toolCall := assistantToolMessage.ToolCalls[0]
						assert.Equal(t, "call_1", toolCall.ID)
						assert.Equal(t, functionName, toolCall.Function.Name)
						assert.Equal(t, string(toolArgumentsJson), toolCall.Function.Arguments)
					}

					toolOutputMessage := request.Messages[len(messages)+1]
					assert.Equal(t, "tool", toolOutputMessage.Role)
					if assert.NotNil(t, toolOutputMessage.ToolCallId) {
						assert.Equal(t, "call_1", *toolOutputMessage.ToolCallId)
					}
					if assert.NotNil(t, toolOutputMessage.Content) {
						assert.Equal(t, "{\"status\":\"ok\"}", *toolOutputMessage.Content)
					}

					response := &external_models.ExternalOpenAICompletionsStandardResponse{
						Choices: []external_models.ExternalOpenAICompletionsStandardResponseChoice{
							{
								Message: external_models.ExternalOpenAICompletionsStandardResponseChoiceMessage{
									Role:    "assistant",
									Content: expectedContent,
								},
							},
						},
					}
					return response, nil
				})

			gomock.InOrder(firstCall, secondCall)

		}

		t.Run("Should call the MCP tool and return the final completion", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOpenAICompletionsStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, llmConnection, messages, []*models.MCPConnection{mcpConnection})
			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, expectedContent, *result)
			}
		})
	})
}
