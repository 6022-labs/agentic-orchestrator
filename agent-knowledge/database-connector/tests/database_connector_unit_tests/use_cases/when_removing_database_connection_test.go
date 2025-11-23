package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRemovingDatabaseConnectionTestingSuite struct {
	sut *use_cases.RemoveDatabaseConnection

	databaseConnectionRepositoryMock   *repositories_mocks.MockDatabaseConnectionRepositoryInterface
	agentDatabaseCredentialStorageMock *services_mocks.MockAgentDatabaseCredentialsStorageInterface
}

func WhenRemovingDatabaseConnectionBeforeEach(t *testing.T) *WhenRemovingDatabaseConnectionTestingSuite {
	mockController := gomock.NewController(t)

	databaseConnectionRepositoryMock := repositories_mocks.NewMockDatabaseConnectionRepositoryInterface(mockController)
	agentDatabaseCredentialStorageMock := services_mocks.NewMockAgentDatabaseCredentialsStorageInterface(mockController)

	sut := use_cases.NewRemoveDatabaseConnection(
		databaseConnectionRepositoryMock,
		agentDatabaseCredentialStorageMock,
	)

	return &WhenRemovingDatabaseConnectionTestingSuite{
		sut: sut,

		databaseConnectionRepositoryMock:   databaseConnectionRepositoryMock,
		agentDatabaseCredentialStorageMock: agentDatabaseCredentialStorageMock,
	}
}

func TestWhenRemovingDatabaseConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurred while deleting credentials", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingDatabaseConnectionTestingSuite) {
			suite.agentDatabaseCredentialStorageMock.EXPECT().
				Delete(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingDatabaseConnectionBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Remove(&requests.AgentIdRequest{
				AgentId: uuid.New(),
			})

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurred while removing database connection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingDatabaseConnectionTestingSuite) {
			suite.agentDatabaseCredentialStorageMock.EXPECT().
				Delete(gomock.Any()).
				Return(nil)

			suite.databaseConnectionRepositoryMock.EXPECT().
				DeleteByAgentId(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingDatabaseConnectionBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Remove(&requests.AgentIdRequest{
				AgentId: uuid.New(),
			})

			assert.Error(t, err)
		})
	})

	t.Run("Given no error occurred", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingDatabaseConnectionTestingSuite) {
			suite.agentDatabaseCredentialStorageMock.EXPECT().
				Delete(gomock.Any()).
				Return(nil)

			suite.databaseConnectionRepositoryMock.EXPECT().
				DeleteByAgentId(gomock.Any()).
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingDatabaseConnectionBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Remove(&requests.AgentIdRequest{
				AgentId: uuid.New(),
			})

			assert.NoError(t, err)
		})
	})
}
