package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/llm_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingLLMConnectionSetupStatusTestingSuite struct {
	sut                         *use_cases.GetLLMConnectionSetupStatus
	llmConnectionRepositoryMock *repositories_mocks.MockLLMConnectionRepositoryInterface
}

func WhenGettingLLMConnectionSetupStatusBeforeEach(t *testing.T) *WhenGettingLLMConnectionSetupStatusTestingSuite {
	mockController := gomock.NewController(t)

	llmConnectionRepositoryMock := repositories_mocks.NewMockLLMConnectionRepositoryInterface(mockController)
	sut := use_cases.NewGetLLMConnectionSetupStatus(llmConnectionRepositoryMock)

	return &WhenGettingLLMConnectionSetupStatusTestingSuite{
		sut:                         sut,
		llmConnectionRepositoryMock: llmConnectionRepositoryMock,
	}
}

func TestWhenGettingLLMConnectionSetupStatus(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusBeforeEach(t)
			request := &requests.AgentIdRequest{AgentId: uuid.Nil}

			_, err := suite.sut.Get(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while querying setup status in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusBeforeEach(t)
			req := &requests.AgentIdRequest{AgentId: uuid.New()}
			suite.llmConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(nil, assert.AnError)

			_, err := suite.sut.Get(req)
			assert.Error(t, err)
		})
	})

	t.Run("Given setup status is not found in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil setup status", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusBeforeEach(t)
			req := &requests.AgentIdRequest{AgentId: uuid.New()}
			suite.llmConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(nil, nil)

			resp, err := suite.sut.Get(req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Nil(t, resp.SetupStatus)
			assert.False(t, resp.IsCompleted)
		})
	})

	t.Run("Given setup status is found in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return setup status", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusBeforeEach(t)
			req := &requests.AgentIdRequest{AgentId: uuid.New()}
			status := enums.LLMConnectionSetupStatusEstablished
			suite.llmConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(&status, nil)

			resp, err := suite.sut.Get(req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, &status, resp.SetupStatus)
			assert.True(t, resp.IsCompleted)
		})
	})
}
