package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingDatabaseConnectionSetupStatusTestingSuite struct {
	sut                              *use_cases.GetDatabaseConnectionSetupStatus
	databaseConnectionRepositoryMock *repositories_mocks.MockDatabaseConnectionRepositoryInterface
}

func WhenGettingDatabaseConnectionSetupStatusBeforeEach(t *testing.T) *WhenGettingDatabaseConnectionSetupStatusTestingSuite {
	mockController := gomock.NewController(t)

	databaseConnectionRepositoryMock := repositories_mocks.NewMockDatabaseConnectionRepositoryInterface(mockController)
	sut := use_cases.NewGetDatabaseConnectionSetupStatus(databaseConnectionRepositoryMock)

	return &WhenGettingDatabaseConnectionSetupStatusTestingSuite{
		sut:                              sut,
		databaseConnectionRepositoryMock: databaseConnectionRepositoryMock,
	}
}

func TestWhenGettingDatabaseConnectionSetupStatus(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingDatabaseConnectionSetupStatusBeforeEach(t)
			request := &requests.AgentIdRequest{AgentId: uuid.Nil}

			_, err := suite.sut.Get(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while querying setup status in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingDatabaseConnectionSetupStatusBeforeEach(t)
			req := &requests.AgentIdRequest{AgentId: uuid.New()}
			suite.databaseConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(nil, assert.AnError)

			_, err := suite.sut.Get(req)
			assert.Error(t, err)
		})
	})

	t.Run("Given setup status is not found in persistence", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil setup status", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingDatabaseConnectionSetupStatusBeforeEach(t)
			req := &requests.AgentIdRequest{AgentId: uuid.New()}
			suite.databaseConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(nil, nil)

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

			suite := WhenGettingDatabaseConnectionSetupStatusBeforeEach(t)
			req := &requests.AgentIdRequest{AgentId: uuid.New()}
			status := enums.DatabaseConnectionSetupStatusEstablished
			suite.databaseConnectionRepositoryMock.EXPECT().GetSetupStatusByAgentId(req.AgentId).Return(&status, nil)

			resp, err := suite.sut.Get(req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, &status, resp.SetupStatus)
			assert.True(t, resp.IsCompleted)
		})
	})
}
