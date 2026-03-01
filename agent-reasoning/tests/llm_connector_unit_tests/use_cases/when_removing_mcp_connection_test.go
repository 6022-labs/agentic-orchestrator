package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenRemovingMCPConnectionTestingSuite struct {
	sut *use_cases.RemoveMCPConnection

	mcpConnectionRepositoryMock             *repositories_mocks.MockMCPConnectionRepositoryInterface
	mcpConnectionHeaderAuthValueStorageMock *services_mocks.MockMCPConnectionHeaderAuthValueStorageInterface
}

func WhenRemovingMCPConnectionBeforeEach(t *testing.T) *WhenRemovingMCPConnectionTestingSuite {
	mockController := gomock.NewController(t)

	mcpConnectionRepositoryMock := repositories_mocks.NewMockMCPConnectionRepositoryInterface(mockController)
	mcpConnectionHeaderAuthValueStorageMock := services_mocks.NewMockMCPConnectionHeaderAuthValueStorageInterface(mockController)

	sut := use_cases.NewRemoveMCPConnection(mcpConnectionRepositoryMock, mcpConnectionHeaderAuthValueStorageMock)

	return &WhenRemovingMCPConnectionTestingSuite{
		sut: sut,

		mcpConnectionRepositoryMock:             mcpConnectionRepositoryMock,
		mcpConnectionHeaderAuthValueStorageMock: mcpConnectionHeaderAuthValueStorageMock,
	}
}

func generateValidRemoveMCPConnectionRequest() *requests.IdRequest {
	return &requests.IdRequest{
		ID: uuid.New(),
	}
}

func TestWhenRemovingMCPConnection(t *testing.T) {
	t.Parallel()
	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingMCPConnectionBeforeEach(t)
			invalidRequest := generateValidRemoveMCPConnectionRequest()
			invalidRequest.ID = uuid.Nil

			err := suite.sut.Remove(invalidRequest)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while finding MCPConnection in repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingMCPConnectionTestingSuite, request *requests.IdRequest) {
			suite.mcpConnectionRepositoryMock.EXPECT().GetById(request.ID).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingMCPConnectionBeforeEach(t)
			validRequest := generateValidRemoveMCPConnectionRequest()
			initSuite(suite, validRequest)

			err := suite.sut.Remove(validRequest)

			assert.Error(t, err)
		})
	})

	t.Run("Given MCPConnection not found in repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingMCPConnectionTestingSuite, request *requests.IdRequest) {
			suite.mcpConnectionRepositoryMock.EXPECT().GetById(request.ID).Return(nil, gorm.ErrRecordNotFound)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingMCPConnectionBeforeEach(t)
			validRequest := generateValidRemoveMCPConnectionRequest()
			initSuite(suite, validRequest)

			err := suite.sut.Remove(validRequest)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while deleting secret for MCPConnection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingMCPConnectionTestingSuite, request *requests.IdRequest) {
			mcpConnection := model_builders.NewMCPConnectionBuilder().
				WithID(request.ID).
				Build()

			suite.mcpConnectionRepositoryMock.EXPECT().GetById(mcpConnection.ID).Return(mcpConnection, nil)
			suite.mcpConnectionHeaderAuthValueStorageMock.EXPECT().Delete(mcpConnection).Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingMCPConnectionBeforeEach(t)
			validRequest := generateValidRemoveMCPConnectionRequest()
			initSuite(suite, validRequest)

			err := suite.sut.Remove(validRequest)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while deleting MCPConnection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingMCPConnectionTestingSuite, request *requests.IdRequest) {
			mcpConnection := model_builders.NewMCPConnectionBuilder().
				WithID(request.ID).
				Build()

			suite.mcpConnectionRepositoryMock.EXPECT().GetById(mcpConnection.ID).Return(mcpConnection, nil)
			suite.mcpConnectionHeaderAuthValueStorageMock.EXPECT().Delete(mcpConnection).Return(nil)
			suite.mcpConnectionRepositoryMock.EXPECT().Delete(mcpConnection).Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingMCPConnectionBeforeEach(t)
			validRequest := generateValidRemoveMCPConnectionRequest()
			initSuite(suite, validRequest)

			err := suite.sut.Remove(validRequest)

			assert.Error(t, err)
		})
	})

	t.Run("Given MCPConnection is removed successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("With auth headers", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenRemovingMCPConnectionTestingSuite, request *requests.IdRequest) {
				mcpConnection := model_builders.NewMCPConnectionBuilder().
					WithID(request.ID).
					Build()

				suite.mcpConnectionRepositoryMock.EXPECT().GetById(mcpConnection.ID).Return(mcpConnection, nil)
				suite.mcpConnectionHeaderAuthValueStorageMock.EXPECT().Delete(mcpConnection).Return(nil)
				suite.mcpConnectionRepositoryMock.EXPECT().Delete(mcpConnection).Return(nil)
			}

			t.Run("Should return nil", func(t *testing.T) {
				t.Parallel()

				suite := WhenRemovingMCPConnectionBeforeEach(t)
				validRequest := generateValidRemoveMCPConnectionRequest()
				initSuite(suite, validRequest)

				err := suite.sut.Remove(validRequest)

				assert.NoError(t, err)
			})
		})

		t.Run("Without auth headers", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenRemovingMCPConnectionTestingSuite, request *requests.IdRequest) {
				mcpConnection := model_builders.NewMCPConnectionBuilder().
					WithID(request.ID).
					WithHeaderAuthKeys(nil).
					Build()

				suite.mcpConnectionRepositoryMock.EXPECT().GetById(mcpConnection.ID).Return(mcpConnection, nil)
				suite.mcpConnectionRepositoryMock.EXPECT().Delete(mcpConnection).Return(nil)
			}

			t.Run("Should return nil", func(t *testing.T) {
				t.Parallel()

				suite := WhenRemovingMCPConnectionBeforeEach(t)
				validRequest := generateValidRemoveMCPConnectionRequest()
				initSuite(suite, validRequest)

				err := suite.sut.Remove(validRequest)

				assert.NoError(t, err)
			})
		})
	})
}
