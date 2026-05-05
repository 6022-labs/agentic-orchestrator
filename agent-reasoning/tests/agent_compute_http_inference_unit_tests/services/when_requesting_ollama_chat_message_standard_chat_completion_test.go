package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	external_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_http_inference_mocks/clients_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRequestingOllamaChatMessageStandardChatCompletionTestingSuite struct {
	sut *services.OllamaChatMessageStandardChatCompletionRequester

	mockMcpRequester *interfaces_mocks.MockMcpRequesterInterface
	mockOllamaClient *clients_mocks.MockOllamaChatMessageStandardClientInterface
}

func WhenRequestingOllamaChatMessageStandardChatCompletionBeforeEach(t *testing.T) *WhenRequestingOllamaChatMessageStandardChatCompletionTestingSuite {
	mockController := gomock.NewController(t)

	mockMcpRequester := interfaces_mocks.NewMockMcpRequesterInterface(mockController)
	mockOllamaClient := clients_mocks.NewMockOllamaChatMessageStandardClientInterface(mockController)

	sut := services.NewOllamaChatMessageStandardChatCompletionRequester(
		zap.NewNop(),
		mockMcpRequester,
		mockOllamaClient,
	)

	return &WhenRequestingOllamaChatMessageStandardChatCompletionTestingSuite{
		sut:              sut,
		mockMcpRequester: mockMcpRequester,
		mockOllamaClient: mockOllamaClient,
	}
}

func TestWhenRequestingOllamaChatMessageStandardChatCompletion(t *testing.T) {
	t.Parallel()

	const baseUrl = "https://ollama.example.com"

	expectedContent := "hi there"

	llmConnection := model_builders.NewLLMConnectionBuilder().
		WithModel("llama3").
		Build()

	userMessage := external_models.ExternalOllamaChatMessageStandardMessage{Role: "user", Content: "hello"}
	messages := []external_models.ExternalOllamaChatMessageStandardMessage{userMessage}

	t.Run("Given no tool calls are requested", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingOllamaChatMessageStandardChatCompletionTestingSuite) {
			suite.mockOllamaClient.EXPECT().
				ChatMessage(baseUrl, gomock.Any()).
				DoAndReturn(func(gotBaseURL string, request *external_models.ExternalOllamaChatMessageStandardRequest) (*external_models.ExternalOllamaChatMessageStandardResponse, error) {
					assert.Equal(t, baseUrl, gotBaseURL)
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Equal(t, messages, request.Messages)
					assert.Nil(t, request.Tools)

					return &external_models.ExternalOllamaChatMessageStandardResponse{
						Model: llmConnection.Model,
						Message: external_models.ExternalOllamaChatMessageStandardMessage{
							Role:      "assistant",
							Content:   expectedContent,
							ToolCalls: nil,
						},
					}, nil
				})
		}

		t.Run("Should return the message content", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOllamaChatMessageStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, llmConnection, messages, nil)
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

		initSuite := func(suite *WhenRequestingOllamaChatMessageStandardChatCompletionTestingSuite) {
			suite.mockMcpRequester.
				EXPECT().
				ListTools(mcpConnection).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOllamaChatMessageStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, llmConnection, messages, []*models.MCPConnection{mcpConnection})
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.ErrorContains(t, err, fmt.Sprintf("failed to fetch MCP tools for connection '%s'", mcpConnection.Name))
		})
	})

	t.Run("Given the response has empty content", func(t *testing.T) {
		t.Parallel()

		doneReason := "length"

		initSuite := func(suite *WhenRequestingOllamaChatMessageStandardChatCompletionTestingSuite) {
			suite.mockOllamaClient.EXPECT().
				ChatMessage(baseUrl, gomock.Any()).
				Return(&external_models.ExternalOllamaChatMessageStandardResponse{
					Model: llmConnection.Model,
					Message: external_models.ExternalOllamaChatMessageStandardMessage{
						Role:    "assistant",
						Content: "",
					},
					DoneReason: &doneReason,
				}, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOllamaChatMessageStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, llmConnection, messages, nil)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.ErrorContains(t, err, "ollama returned empty content")
			assert.ErrorContains(t, err, doneReason)
		})
	})

	t.Run("Given the LLM requests an MCP tool call", func(t *testing.T) {
		t.Parallel()

		mcpConnection := model_builders.NewMCPConnectionBuilder().
			WithName("inventory").
			WithUrl("https://mcp.example.com").
			Build()

		initSuite := func(suite *WhenRequestingOllamaChatMessageStandardChatCompletionTestingSuite) {
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
			toolArguments := map[string]any{"query": "value"}
			toolResult := map[string]any{"status": "ok"}

			suite.mockMcpRequester.
				EXPECT().
				ListTools(mcpConnection).
				Return([]*interfaces.McpTool{mcpTool}, nil).
				Times(2)

			suite.mockMcpRequester.
				EXPECT().
				CallTool(mcpConnection, mcpTool, toolArguments).
				Return(toolResult, nil)

			// First call returns a tool call request
			suite.mockOllamaClient.EXPECT().ChatMessage(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ string, request *external_models.ExternalOllamaChatMessageStandardRequest) (*external_models.ExternalOllamaChatMessageStandardResponse, error) {
					return &external_models.ExternalOllamaChatMessageStandardResponse{
						Model: llmConnection.Model,
						Message: external_models.ExternalOllamaChatMessageStandardMessage{
							Role:    "assistant",
							Content: "",
							ToolCalls: []external_models.ExternalOllamaChatMessageStandardMessageToolCall{
								{
									Function: &external_models.ExternalOllamaChatMessageStandardMessageToolCallFunction{
										Name:      functionName,
										Arguments: map[string]any{"query": "value"},
									},
								},
							},
						},
					}, nil
				})

			// Second call returns the final message
			suite.mockOllamaClient.EXPECT().ChatMessage(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ string, request *external_models.ExternalOllamaChatMessageStandardRequest) (*external_models.ExternalOllamaChatMessageStandardResponse, error) {
					return &external_models.ExternalOllamaChatMessageStandardResponse{
						Model: llmConnection.Model,
						Message: external_models.ExternalOllamaChatMessageStandardMessage{
							Role:      "assistant",
							Content:   expectedContent,
							ToolCalls: nil,
						},
					}, nil
				})
		}

		t.Run("Should call the MCP tool and return the final completion", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOllamaChatMessageStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, llmConnection, messages, []*models.MCPConnection{mcpConnection})
			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, expectedContent, *result)
			}
		})
	})
}
