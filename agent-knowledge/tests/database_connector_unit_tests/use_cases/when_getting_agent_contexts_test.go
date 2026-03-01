package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/db_clients_mocks"
	factories_db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/factories_mocks/db_clients_mocks"
	factories_repositories_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/factories_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentContextsTestingSuite struct {
	sut *use_cases.GetAgentContexts

	mockAgentDBClientFactory          *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient                 *db_clients_mocks.MockAgentDatabaseClientInterface
	mockAgentContextRepositoryFactory *factories_repositories_mocks.MockAgentContextRepositoryFactoryInterface

	mockAgentContextRepository *repositories_mocks.MockAgentContextRepositoryInterface
}

func WhenGettingAgentContextsBeforeEach(t *testing.T) *WhenGettingAgentContextsTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockAgentContextRepositoryFactory := factories_repositories_mocks.NewMockAgentContextRepositoryFactoryInterface(mockController)
	mockAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)

	sut := use_cases.NewGetAgentContexts(
		mockAgentDBClientFactory,
		mockAgentContextRepositoryFactory,
	)

	return &WhenGettingAgentContextsTestingSuite{
		sut: sut,

		mockAgentDBClientFactory:          mockAgentDBClientFactory,
		mockAgentDbClient:                 mockAgentDbClient,
		mockAgentContextRepository:        mockAgentContextRepository,
		mockAgentContextRepositoryFactory: mockAgentContextRepositoryFactory,
	}
}

func generateValidGetAgentContextsRequest() *requests.GetAgentContextsRequest {
	return &requests.GetAgentContextsRequest{
		AgentId: uuid.New(),
	}
}

func TestWhenGettingAgentContexts(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextsBeforeEach(t)
			request := generateValidGetAgentContextsRequest()
			request.AgentId = uuid.Nil

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must be provided")
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextsTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().Create(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextsBeforeEach(t)
			initSuite(suite)
			request := &requests.GetAgentContextsRequest{AgentId: uuid.New()}

			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given error occurs while searching agent contexts", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextsTestingSuite, request *requests.GetAgentContextsRequest) {
			suite.mockAgentDBClientFactory.EXPECT().Create(gomock.Any()).Return(suite.mockAgentDbClient, nil)
			suite.mockAgentDbClient.EXPECT().Close().Return(nil)
			suite.mockAgentContextRepositoryFactory.EXPECT().Create(gomock.Any()).Return(suite.mockAgentContextRepository)
			suite.mockAgentContextRepository.EXPECT().Search(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextsBeforeEach(t)
			request := &requests.GetAgentContextsRequest{AgentId: uuid.New()}
			initSuite(suite, request)

			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given agent contexts are found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextsTestingSuite, request *requests.GetAgentContextsRequest) {
			suite.mockAgentDBClientFactory.EXPECT().Create(gomock.Any()).Return(suite.mockAgentDbClient, nil)
			suite.mockAgentDbClient.EXPECT().Close().Return(nil)
			suite.mockAgentContextRepositoryFactory.EXPECT().Create(gomock.Any()).Return(suite.mockAgentContextRepository)
			suite.mockAgentContextRepository.EXPECT().Search(gomock.Any()).Return([]*models.AgentContext{
				&models.AgentContext{Key: "key1"},
				&models.AgentContext{Key: "key2"},
			}, nil)
		}

		t.Run("Should return agent contexts response", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextsBeforeEach(t)
			request := &requests.GetAgentContextsRequest{AgentId: uuid.New()}
			initSuite(suite, request)

			resp, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Len(t, resp, 2)
		})
	})
}
