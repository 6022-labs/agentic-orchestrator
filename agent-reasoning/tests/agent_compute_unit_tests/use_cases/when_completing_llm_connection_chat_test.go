package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCompletingLLMConnectionChatTestingSuite struct {
	sut *use_cases.LLMConnectionChatCompletion

	mockAgentLLMSecretStorage               *services_mocks.MockAgentLLMSecretStorageInterface
	mockChatCompletionRequester             *interfaces_mocks.MockChatCompletionRequesterInterface
	mockLLMConnectionRepository             *repositories_mocks.MockLLMConnectionRepositoryInterface
	mockMCPConnectionRepository             *repositories_mocks.MockMCPConnectionRepositoryInterface
	mockMCPConnectionHeaderAuthValueStorage *services_mocks.MockMCPConnectionHeaderAuthValueStorageInterface
}

func WhenCompletingLLMConnectionChatBeforeEach(t *testing.T) *WhenCompletingLLMConnectionChatTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentLLMSecretStorage := services_mocks.NewMockAgentLLMSecretStorageInterface(mockController)
	mockChatCompletionRequester := interfaces_mocks.NewMockChatCompletionRequesterInterface(mockController)
	mockLLMConnectionRepository := repositories_mocks.NewMockLLMConnectionRepositoryInterface(mockController)
	mockMCPConnectionRepository := repositories_mocks.NewMockMCPConnectionRepositoryInterface(mockController)
	mockMCPConnectionHeaderAuthValueStorage := services_mocks.NewMockMCPConnectionHeaderAuthValueStorageInterface(mockController)

	sut := use_cases.NewLLMConnectionChatCompletion(
		mockAgentLLMSecretStorage,
		mockChatCompletionRequester,
		mockLLMConnectionRepository,
		mockMCPConnectionRepository,
		mockMCPConnectionHeaderAuthValueStorage,
	)

	return &WhenCompletingLLMConnectionChatTestingSuite{
		sut: sut,

		mockAgentLLMSecretStorage:               mockAgentLLMSecretStorage,
		mockChatCompletionRequester:             mockChatCompletionRequester,
		mockLLMConnectionRepository:             mockLLMConnectionRepository,
		mockMCPConnectionRepository:             mockMCPConnectionRepository,
		mockMCPConnectionHeaderAuthValueStorage: mockMCPConnectionHeaderAuthValueStorage,
	}
}

func generateValidLLMConnectionChatCompletionRequest() *requests.LLMConnectionChatCompletionRequest {
	return &requests.LLMConnectionChatCompletionRequest{
		AgentId: uuid.New(),
		Messages: []*models.ChatCompletionMessage{
			{Role: "user", Content: "Hello"},
		},
	}
}

func TestWhenCompletingLLMConnectionChat(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCompletingLLMConnectionChatBeforeEach(t)

			request := generateValidLLMConnectionChatCompletionRequest()
			request.AgentId = uuid.Nil

			result, err := suite.sut.Complete(request)

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "agentId")
		})

		t.Run("Because messages is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCompletingLLMConnectionChatBeforeEach(t)

			request := generateValidLLMConnectionChatCompletionRequest()
			request.Messages = []*models.ChatCompletionMessage{}

			result, err := suite.sut.Complete(request)

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "messages")
		})
	})

	t.Run("Given error occurs while getting llm connection", func(t *testing.T) {
		t.Parallel()

		suite := WhenCompletingLLMConnectionChatBeforeEach(t)

		request := generateValidLLMConnectionChatCompletionRequest()

		suite.mockLLMConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(nil, assert.AnError)

		result, err := suite.sut.Complete(request)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("Given llm connection is not found", func(t *testing.T) {
		t.Parallel()

		suite := WhenCompletingLLMConnectionChatBeforeEach(t)

		request := generateValidLLMConnectionChatCompletionRequest()

		suite.mockLLMConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(nil, nil)

		result, err := suite.sut.Complete(request)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "llm connection not found")
	})

	t.Run("Given error occurs while getting mcp connections", func(t *testing.T) {
		t.Parallel()

		suite := WhenCompletingLLMConnectionChatBeforeEach(t)

		request := generateValidLLMConnectionChatCompletionRequest()
		llmConnection := model_builders.NewLLMConnectionBuilder().
			WithAgentId(request.AgentId).
			Build()

		suite.mockLLMConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(llmConnection, nil)

		suite.mockMCPConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(nil, assert.AnError)

		result, err := suite.sut.Complete(request)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("Given error occurs while getting api key", func(t *testing.T) {
		t.Parallel()

		suite := WhenCompletingLLMConnectionChatBeforeEach(t)

		request := generateValidLLMConnectionChatCompletionRequest()
		llmConnection := model_builders.NewLLMConnectionBuilder().
			WithAgentId(request.AgentId).
			WithPlatform(enums.LLMConnectionPlatformOpenAI).
			Build()

		suite.mockLLMConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(llmConnection, nil)

		suite.mockMCPConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return([]*models.MCPConnection{}, nil)

		suite.mockAgentLLMSecretStorage.EXPECT().
			Get(request.AgentId).
			Return(nil, assert.AnError)

		result, err := suite.sut.Complete(request)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("Given completion is successful with cloud platform", func(t *testing.T) {
		t.Parallel()

		suite := WhenCompletingLLMConnectionChatBeforeEach(t)

		request := generateValidLLMConnectionChatCompletionRequest()
		llmConnection := model_builders.NewLLMConnectionBuilder().
			WithAgentId(request.AgentId).
			WithPlatform(enums.LLMConnectionPlatformOpenAI).
			Build()

		apiKey := "test-api-key"
		expectedCompletion := "Hello! How can I help?"

		suite.mockLLMConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(llmConnection, nil)

		suite.mockMCPConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return([]*models.MCPConnection{}, nil)

		suite.mockAgentLLMSecretStorage.EXPECT().
			Get(request.AgentId).
			Return(&apiKey, nil)

		suite.mockChatCompletionRequester.EXPECT().
			ChatCompletion(llmConnection, gomock.Any(), request.Messages, &apiKey).
			Return(&expectedCompletion, nil)

		result, err := suite.sut.Complete(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCompletion, *result)
	})

	t.Run("Given completion is successful with local platform (ollama)", func(t *testing.T) {
		t.Parallel()

		suite := WhenCompletingLLMConnectionChatBeforeEach(t)

		request := generateValidLLMConnectionChatCompletionRequest()
		llmConnection := model_builders.NewLLMConnectionBuilder().
			WithAgentId(request.AgentId).
			WithPlatform(enums.LLMConnectionPlatformOllama).
			Build()

		expectedCompletion := "Hello from Ollama!"

		suite.mockLLMConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(llmConnection, nil)

		suite.mockMCPConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return([]*models.MCPConnection{}, nil)

		suite.mockChatCompletionRequester.EXPECT().
			ChatCompletion(llmConnection, gomock.Any(), request.Messages, nil).
			Return(&expectedCompletion, nil)

		result, err := suite.sut.Complete(request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCompletion, *result)
	})

	t.Run("Given error occurs while getting header auth for mcp connection", func(t *testing.T) {
		t.Parallel()

		suite := WhenCompletingLLMConnectionChatBeforeEach(t)

		request := generateValidLLMConnectionChatCompletionRequest()
		llmConnection := model_builders.NewLLMConnectionBuilder().
			WithAgentId(request.AgentId).
			WithPlatform(enums.LLMConnectionPlatformOpenAI).
			Build()

		mcpConnection := model_builders.NewMCPConnectionBuilder().
			WithAgentId(request.AgentId).
			Build()

		suite.mockLLMConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return(llmConnection, nil)

		suite.mockMCPConnectionRepository.EXPECT().
			GetByAgentId(request.AgentId).
			Return([]*models.MCPConnection{mcpConnection}, nil)

		suite.mockMCPConnectionHeaderAuthValueStorage.EXPECT().
			Get(mcpConnection).
			Return(nil, assert.AnError)

		result, err := suite.sut.Complete(request)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}
