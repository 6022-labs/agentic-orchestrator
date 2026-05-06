package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/repositories_mocks"
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

		req := &requests.AgentIdRequest{AgentId: uuid.New()}

		initSuite := func(suite *WhenGettingLLMConnectionSetupStatusTestingSuite) {
			suite.llmConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(req)
			assert.Error(t, err)
		})
	})

	t.Run("Given setup status is not found in persistence", func(t *testing.T) {
		t.Parallel()

		req := &requests.AgentIdRequest{AgentId: uuid.New()}

		initSuite := func(suite *WhenGettingLLMConnectionSetupStatusTestingSuite) {
			suite.llmConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(nil, nil)
		}

		t.Run("Should return nil setup status", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusBeforeEach(t)
			initSuite(suite)

			resp, err := suite.sut.Get(req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Nil(t, resp.SetupStatus)
			assert.False(t, resp.IsCompleted)
		})
	})

	t.Run("Given setup status is found in persistence", func(t *testing.T) {
		t.Parallel()

		req := &requests.AgentIdRequest{AgentId: uuid.New()}
		status := enums.LLMConnectionSetupStatusEstablished

		initSuite := func(suite *WhenGettingLLMConnectionSetupStatusTestingSuite) {
			suite.llmConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(&status, nil)
		}

		t.Run("Should return setup status", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusBeforeEach(t)
			initSuite(suite)

			resp, err := suite.sut.Get(req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, &status, resp.SetupStatus)
			assert.True(t, resp.IsCompleted)
		})
	})
}
