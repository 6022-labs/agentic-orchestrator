package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingLLMConnectionSettingsTestingSuite struct {
	sut *use_cases.GetLLMConnectionSettings

	llmConnectionRepositoryMock *repositories_mocks.MockLLMConnectionRepositoryInterface
}

func WhenGettingLLMConnectionSettingsBeforeEach(t *testing.T) *WhenGettingLLMConnectionSettingsTestingSuite {
	mockController := gomock.NewController(t)

	llmConnectionRepositoryMock := repositories_mocks.NewMockLLMConnectionRepositoryInterface(mockController)

	sut := use_cases.NewGetLLMConnectionSettings(llmConnectionRepositoryMock)

	return &WhenGettingLLMConnectionSettingsTestingSuite{
		sut:                         sut,
		llmConnectionRepositoryMock: llmConnectionRepositoryMock,
	}
}

func TestWhenGettingLLMConnectionSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSettingsBeforeEach(t)

			request := &requests.AgentIdRequest{
				AgentId: uuid.Nil,
			}

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while querying llm connection in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSettingsBeforeEach(t)
			request := &requests.AgentIdRequest{
				AgentId: uuid.New(),
			}

			suite.llmConnectionRepositoryMock.EXPECT().GetByAgentId(request.AgentId).Return(nil, assert.AnError)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given llm connection is not found in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSettingsBeforeEach(t)
			request := &requests.AgentIdRequest{
				AgentId: uuid.New(),
			}

			suite.llmConnectionRepositoryMock.EXPECT().GetByAgentId(request.AgentId).Return(nil, nil)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given llm connection is found in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return llm connection metadata", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSettingsBeforeEach(t)
			request := &requests.AgentIdRequest{
				AgentId: uuid.New(),
			}

			existingConnection := model_builders.NewLLMConnectionBuilder().
				WithBaseUrl("https://custom.host").
				Build()
			suite.llmConnectionRepositoryMock.EXPECT().
				GetByAgentId(request.AgentId).
				Return(existingConnection, nil)

			response, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.Equal(t, existingConnection.Model, response.Model)
			assert.Equal(t, existingConnection.BaseUrl, response.BaseUrl)
			assert.Equal(t, existingConnection.Platform, response.Platform)
		})
	})
}
