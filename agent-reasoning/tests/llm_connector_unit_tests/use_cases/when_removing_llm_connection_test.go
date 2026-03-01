package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRemovingLLMConnectionTestingSuite struct {
	sut *use_cases.RemoveLLMConnection

	llmConnectionRepositoryMock *repositories_mocks.MockLLMConnectionRepositoryInterface
}

func WhenRemovingLLMConnectionBeforeEach(t *testing.T) *WhenRemovingLLMConnectionTestingSuite {
	mockController := gomock.NewController(t)

	llmConnectionRepositoryMock := repositories_mocks.NewMockLLMConnectionRepositoryInterface(mockController)

	sut := use_cases.NewRemoveLLMConnection(llmConnectionRepositoryMock)

	return &WhenRemovingLLMConnectionTestingSuite{
		sut:                         sut,
		llmConnectionRepositoryMock: llmConnectionRepositoryMock,
	}
}

func TestWhenRemovingLLMConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingLLMConnectionBeforeEach(t)

			request := &requests.AgentIdRequest{
				AgentId: uuid.Nil,
			}

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while removing llm connection in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingLLMConnectionBeforeEach(t)

			request := &requests.AgentIdRequest{
				AgentId: uuid.New(),
			}

			suite.llmConnectionRepositoryMock.EXPECT().DeleteByAgentId(request.AgentId).Return(assert.AnError)

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given llm connection is removed successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingLLMConnectionBeforeEach(t)

			request := &requests.AgentIdRequest{
				AgentId: uuid.New(),
			}

			suite.llmConnectionRepositoryMock.EXPECT().DeleteByAgentId(request.AgentId).Return(nil)

			err := suite.sut.Remove(request)

			assert.NoError(t, err)
		})
	})
}
