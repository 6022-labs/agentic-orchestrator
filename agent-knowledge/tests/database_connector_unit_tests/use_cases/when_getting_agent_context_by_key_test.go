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

type WhenGettingAgentContextByKeyTestingSuite struct {
	sut *use_cases.GetAgentContextByKey

	mockAgentDBClientFactory          *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient                 *db_clients_mocks.MockAgentDatabaseClientInterface
	mockAgentContextRepositoryFactory *factories_repositories_mocks.MockAgentContextRepositoryFactoryInterface

	mockAgentContextRepository *repositories_mocks.MockAgentContextRepositoryInterface
}

func WhenGettingAgentContextByKeyBeforeEach(t *testing.T) *WhenGettingAgentContextByKeyTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockAgentContextRepositoryFactory := factories_repositories_mocks.NewMockAgentContextRepositoryFactoryInterface(mockController)
	mockAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)

	sut := use_cases.NewGetAgentContextByKey(
		mockAgentDBClientFactory,
		mockAgentContextRepositoryFactory,
	)

	return &WhenGettingAgentContextByKeyTestingSuite{
		sut: sut,

		mockAgentDBClientFactory:          mockAgentDBClientFactory,
		mockAgentDbClient:                 mockAgentDbClient,
		mockAgentContextRepositoryFactory: mockAgentContextRepositoryFactory,

		mockAgentContextRepository: mockAgentContextRepository,
	}
}

func generateValidGetAgentContextByKeyRequest() *requests.GetAgentContextByKeyRequest {
	return &requests.GetAgentContextByKeyRequest{
		AgentId: uuid.New(),
		Key:     uuid.New().String(),
	}
}

func TestWhenGettingAgentContextByKey(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextByKeyBeforeEach(t)

			request := generateValidGetAgentContextByKeyRequest()
			request.AgentId = uuid.Nil

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})

		t.Run("Because key is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextByKeyBeforeEach(t)

			request := generateValidGetAgentContextByKeyRequest()
			request.Key = ""

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextByKeyTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextByKeyBeforeEach(t)
			initSuite(suite)

			request := generateValidGetAgentContextByKeyRequest()
			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given error occurs while getting agent context by key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextByKeyTestingSuite, request *requests.GetAgentContextByKeyRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(request.Key).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextByKeyBeforeEach(t)
			request := generateValidGetAgentContextByKeyRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given agent context not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextByKeyTestingSuite, request *requests.GetAgentContextByKeyRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(request.Key).
				Return(nil, gorm.ErrRecordNotFound)
		}

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextByKeyBeforeEach(t)
			request := generateValidGetAgentContextByKeyRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given agent context found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextByKeyTestingSuite, request *requests.GetAgentContextByKeyRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			agentContext := model_builders.NewAgentContextBuilder().
				WithKey(request.Key).
				Build()

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(request.Key).
				Return(agentContext, nil)
		}

		t.Run("Should return agent context response", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextByKeyBeforeEach(t)
			request := generateValidGetAgentContextByKeyRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, request.Key, resp.Key)
		})
	})
}
