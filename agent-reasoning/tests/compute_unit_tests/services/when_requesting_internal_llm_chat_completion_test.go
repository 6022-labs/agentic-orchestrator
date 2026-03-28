package services_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/compute_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRequestingInternalLLMChatCompletionTestingSuite struct {
	sut                     *services.InternalLLMConnectionChatCompletionRequester
	llmChatCompletionUseCase *services_mocks.MockLLMChatCompletionUseCase
}

func WhenRequestingInternalLLMChatCompletionBeforeEach(t *testing.T) *WhenRequestingInternalLLMChatCompletionTestingSuite {
	mockController := gomock.NewController(t)

	llmChatCompletionUseCaseMock := services_mocks.NewMockLLMChatCompletionUseCase(mockController)

	sut := services.NewInternalLLMConnectionChatCompletionRequester(llmChatCompletionUseCaseMock)

	return &WhenRequestingInternalLLMChatCompletionTestingSuite{
		sut:                     sut,
		llmChatCompletionUseCase: llmChatCompletionUseCaseMock,
	}
}

func TestWhenRequestingInternalLLMChatCompletion(t *testing.T) {
	t.Parallel()

	t.Run("Given use case returns error", func(t *testing.T) {
		t.Parallel()

		suite := WhenRequestingInternalLLMChatCompletionBeforeEach(t)

		agentId := uuid.New()
		messages := []*models.AgentContextMessage{
			{
				Content:   "Hello",
				Role:      enums.AgentContextMessageRoleUser,
				CreatedAt: time.Now(),
			},
		}

		suite.llmChatCompletionUseCase.EXPECT().
			Complete(gomock.Any()).
			Return(nil, assert.AnError)

		result, err := suite.sut.ChatCompletion(agentId, messages)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("Given messages are sorted by CreatedAt and completion is returned", func(t *testing.T) {
		t.Parallel()

		suite := WhenRequestingInternalLLMChatCompletionBeforeEach(t)

		agentId := uuid.New()
		now := time.Now()
		messages := []*models.AgentContextMessage{
			{
				Content:   "Second message",
				Role:      enums.AgentContextMessageRoleAssistant,
				CreatedAt: now.Add(1 * time.Hour),
			},
			{
				Content:   "First message",
				Role:      enums.AgentContextMessageRoleUser,
				CreatedAt: now,
			},
		}

		expectedCompletion := "completion response"

		suite.llmChatCompletionUseCase.EXPECT().
			Complete(gomock.Any()).
			DoAndReturn(func(request *requests.LLMConnectionChatCompletionRequest) (*string, error) {
				assert.Equal(t, agentId, request.AgentId)
				assert.Len(t, request.Messages, 2)

				// Messages should be sorted by CreatedAt (oldest first)
				assert.Equal(t, "First message", request.Messages[0].Content)
				assert.Equal(t, string(enums.AgentContextMessageRoleUser), request.Messages[0].Role)
				assert.Equal(t, "Second message", request.Messages[1].Content)
				assert.Equal(t, string(enums.AgentContextMessageRoleAssistant), request.Messages[1].Role)

				return &expectedCompletion, nil
			})

		result, err := suite.sut.ChatCompletion(agentId, messages)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCompletion, *result)
	})

	t.Run("Given messages converts roles to ChatCompletionMessage format", func(t *testing.T) {
		t.Parallel()

		suite := WhenRequestingInternalLLMChatCompletionBeforeEach(t)

		agentId := uuid.New()
		messages := []*models.AgentContextMessage{
			{
				Content:   "System prompt",
				Role:      enums.AgentContextMessageRoleSystem,
				CreatedAt: time.Now(),
			},
		}

		expectedCompletion := "response"

		suite.llmChatCompletionUseCase.EXPECT().
			Complete(gomock.Any()).
			DoAndReturn(func(request *requests.LLMConnectionChatCompletionRequest) (*string, error) {
				assert.Len(t, request.Messages, 1)
				assert.Equal(t, "system", request.Messages[0].Role)
				assert.Equal(t, "System prompt", request.Messages[0].Content)

				return &expectedCompletion, nil
			})

		result, err := suite.sut.ChatCompletion(agentId, messages)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCompletion, *result)
	})
}
