package services_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/services/interfaces"
	external_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_http_inference_mocks/clients_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRequestingOpenAIResponsesStandardChatCompletionTestingSuite struct {
	sut *services.OpenAIResponsesStandardChatCompletionRequester

	mockMcpRequester          *interfaces_mocks.MockMcpRequesterInterface
	mockOpenAIResponsesClient *clients_mocks.MockOpenAIResponsesStandardClientInterface
}

func WhenRequestingOpenAIResponsesStandardChatCompletionBeforeEach(t *testing.T) *WhenRequestingOpenAIResponsesStandardChatCompletionTestingSuite {
	mockController := gomock.NewController(t)

	mockMcpRequester := interfaces_mocks.NewMockMcpRequesterInterface(mockController)
	mockOpenAIResponsesClient := clients_mocks.NewMockOpenAIResponsesStandardClientInterface(mockController)

	sut := services.NewOpenAIResponsesStandardChatCompletionRequester(
		zap.NewNop(),
		mockMcpRequester,
		mockOpenAIResponsesClient,
	)

	return &WhenRequestingOpenAIResponsesStandardChatCompletionTestingSuite{
		sut:                       sut,
		mockMcpRequester:          mockMcpRequester,
		mockOpenAIResponsesClient: mockOpenAIResponsesClient,
	}
}

func TestWhenRequestingOpenAIResponsesStandardChatCompletion(t *testing.T) {
	t.Parallel()

	const (
		baseUrl = "https://openai.example.com"
		apiKey  = "test-key"

		expectedContent = "hi there"
	)

	llmConnection := model_builders.NewLLMConnectionBuilder().
		WithModel("gpt-4").
		Build()

	userRole := "user"
	userContent := "hello"
	inputs := []external_models.ExternalOpenAIResponsesStandardRequestInput{
		{Type: "message", Role: &userRole, Content: &userContent},
	}

	t.Run("Given no tool calls are requested", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingOpenAIResponsesStandardChatCompletionTestingSuite) {
			suite.mockOpenAIResponsesClient.
				EXPECT().
				Responses(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalOpenAIResponsesStandardRequest) (*external_models.ExternalOpenAIResponsesStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Equal(t, inputs, request.Input)
					assert.Nil(t, request.Tools)
					assert.Nil(t, request.PreviousResponseID)

					assistantRole := "assistant"
					return &external_models.ExternalOpenAIResponsesStandardResponse{
						Output: []external_models.ExternalOpenAIResponsesStandardResponseOutput{
							{
								Type: "message",
								Role: &assistantRole,
								Content: []external_models.ExternalOpenAIResponsesStandardResponseOutputContent{
									{Type: "output_text", Text: expectedContent},
								},
							},
						},
					}, nil
				})
		}

		t.Run("Should return the first message output content", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOpenAIResponsesStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, false, llmConnection, inputs, nil, nil)
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

		initSuite := func(suite *WhenRequestingOpenAIResponsesStandardChatCompletionTestingSuite) {
			suite.mockMcpRequester.
				EXPECT().
				ListTools(mcpConnection).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOpenAIResponsesStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, false, llmConnection, inputs, []*models.MCPConnection{mcpConnection}, nil)
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

		initSuite := func(suite *WhenRequestingOpenAIResponsesStandardChatCompletionTestingSuite) {
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
			toolArgumentsJson, _ := json.Marshal(toolArguments)
			toolArgumentsString := string(toolArgumentsJson)
			toolResult := map[string]any{"status": "ok"}
			callID := "call_1"
			firstResponseID := "resp_123"

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

			firstCall := suite.mockOpenAIResponsesClient.
				EXPECT().
				Responses(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalOpenAIResponsesStandardRequest) (*external_models.ExternalOpenAIResponsesStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					assert.Equal(t, inputs, request.Input)
					assert.Nil(t, request.PreviousResponseID)
					if assert.Len(t, request.Tools, 1) {
						tool := request.Tools[0]
						assert.Equal(t, "function", tool.Type)
						if assert.NotNil(t, tool.Name) {
							assert.Equal(t, functionName, *tool.Name)
						}
					}

					return &external_models.ExternalOpenAIResponsesStandardResponse{
						ID: firstResponseID,
						Output: []external_models.ExternalOpenAIResponsesStandardResponseOutput{
							{
								Type:      "function_call",
								Name:      &functionName,
								CallId:    &callID,
								Arguments: &toolArgumentsString,
							},
						},
					}, nil
				})

			secondCall := suite.mockOpenAIResponsesClient.
				EXPECT().
				Responses(baseUrl, apiKey, gomock.Any()).
				DoAndReturn(func(_ string, _ string, request *external_models.ExternalOpenAIResponsesStandardRequest) (*external_models.ExternalOpenAIResponsesStandardResponse, error) {
					assert.Equal(t, llmConnection.Model, request.Model)
					if assert.NotNil(t, request.PreviousResponseID) {
						assert.Equal(t, firstResponseID, *request.PreviousResponseID)
					}
					if assert.Len(t, request.Input, 1) {
						toolOutputInput := request.Input[0]
						assert.Equal(t, "function_call_output", toolOutputInput.Type)
						if assert.NotNil(t, toolOutputInput.CallId) {
							assert.Equal(t, callID, *toolOutputInput.CallId)
						}
						if assert.NotNil(t, toolOutputInput.Output) {
							assert.Equal(t, "{\"status\":\"ok\"}", *toolOutputInput.Output)
						}
					}
					if assert.Len(t, request.Tools, 1) {
						tool := request.Tools[0]
						assert.Equal(t, "function", tool.Type)
						if assert.NotNil(t, tool.Name) {
							assert.Equal(t, functionName, *tool.Name)
						}
					}

					assistantRole := "assistant"
					return &external_models.ExternalOpenAIResponsesStandardResponse{
						Output: []external_models.ExternalOpenAIResponsesStandardResponseOutput{
							{
								Type: "message",
								Role: &assistantRole,
								Content: []external_models.ExternalOpenAIResponsesStandardResponseOutputContent{
									{Type: "output_text", Text: expectedContent},
								},
							},
						},
					}, nil
				})

			gomock.InOrder(firstCall, secondCall)
		}

		t.Run("Should call the MCP tool and return the final completion", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingOpenAIResponsesStandardChatCompletionBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ChatCompletion(baseUrl, apiKey, false, llmConnection, inputs, []*models.MCPConnection{mcpConnection}, nil)
			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, expectedContent, *result)
			}
		})
	})
}
