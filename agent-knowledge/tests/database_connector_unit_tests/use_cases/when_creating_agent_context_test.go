package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/model_builders"
	db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/db_clients_mocks"
	factories_db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/factories_mocks/db_clients_mocks"
	factories_repositories_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/factories_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenCreatingAgentContextTestingSuite struct {
	sut *use_cases.CreateAgentContext

	mockAgentDBClientFactory          *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient                 *db_clients_mocks.MockAgentDatabaseClientInterface
	mockAgentContextRepositoryFactory *factories_repositories_mocks.MockAgentContextRepositoryFactoryInterface

	mockAgentContextRepository *repositories_mocks.MockAgentContextRepositoryInterface
}

func WhenCreatingAgentContextBeforeEach(t *testing.T) *WhenCreatingAgentContextTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockAgentContextRepositoryFactory := factories_repositories_mocks.NewMockAgentContextRepositoryFactoryInterface(mockController)
	mockAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)

	sut := use_cases.NewCreateAgentContext(
		mockAgentDBClientFactory,
		mockAgentContextRepositoryFactory,
	)

	return &WhenCreatingAgentContextTestingSuite{
		sut: sut,

		mockAgentDBClientFactory: mockAgentDBClientFactory,
		mockAgentDbClient:        mockAgentDbClient,

		mockAgentContextRepositoryFactory: mockAgentContextRepositoryFactory,
		mockAgentContextRepository:        mockAgentContextRepository,
	}
}

func generateValidCreateAgentContextRequest() *requests.CreateAgentContextRequest {
	return &requests.CreateAgentContextRequest{
		AgentId: uuid.New(),
		Key:     "test",
	}
}

func TestWhenCreatingAgentContext(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)

			request := generateValidCreateAgentContextRequest()
			request.AgentId = uuid.Nil

			err := suite.sut.Create(request)
			assert.Error(t, err)
		})

		t.Run("Because contextKey is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)

			request := generateValidCreateAgentContextRequest()
			request.Key = ""

			err := suite.sut.Create(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateValidCreateAgentContextRequest()
			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving agent context", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound)

			suite.mockAgentContextRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateValidCreateAgentContextRequest()
			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound)

			suite.mockAgentContextRepository.EXPECT().
				Save(gomock.Any()).
				Return(model_builders.NewAgentContextBuilder().Build(), nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateValidCreateAgentContextRequest()
			err := suite.sut.Create(request)

			assert.NoError(t, err)
		})
	})
}
