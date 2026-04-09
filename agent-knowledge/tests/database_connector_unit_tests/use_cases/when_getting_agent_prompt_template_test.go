package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/model_builders"
	db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/db_clients_mocks"
	factories_db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/factories_mocks/db_clients_mocks"
	factories_repositories_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/factories_mocks/repositories_mocks"
	repositories_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/repositories_mocks"
	services_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.GetAgentPromptTemplate

	mockAgentDBClientFactory                 *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient                        *db_clients_mocks.MockAgentDatabaseClientInterface
	mockDatabaseConnectionRepository         *repositories_mocks.MockDatabaseConnectionRepositoryInterface
	mockAgentDatabaseCredentialsStorage      *services_mocks.MockAgentDatabaseCredentialsStorageInterface
	mockAgentPromptTemplateRepositoryFactory *factories_repositories_mocks.MockAgentPromptTemplateRepositoryFactoryInterface

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
}

func WhenGettingAgentPromptTemplateBeforeEach(t *testing.T) *WhenGettingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockDatabaseConnectionRepository := repositories_mocks.NewMockDatabaseConnectionRepositoryInterface(mockController)
	mockAgentDatabaseCredentialsStorage := services_mocks.NewMockAgentDatabaseCredentialsStorageInterface(mockController)
	mockAgentPromptTemplateRepositoryFactory := factories_repositories_mocks.NewMockAgentPromptTemplateRepositoryFactoryInterface(mockController)
	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(mockController)

	sut := use_cases.NewGetAgentPromptTemplate(
		mockAgentDBClientFactory,
		mockDatabaseConnectionRepository,
		mockAgentDatabaseCredentialsStorage,
		mockAgentPromptTemplateRepositoryFactory,
	)

	return &WhenGettingAgentPromptTemplateTestingSuite{
		sut:                                      sut,
		mockAgentDBClientFactory:                 mockAgentDBClientFactory,
		mockAgentDbClient:                        mockAgentDbClient,
		mockDatabaseConnectionRepository:         mockDatabaseConnectionRepository,
		mockAgentDatabaseCredentialsStorage:      mockAgentDatabaseCredentialsStorage,
		mockAgentPromptTemplateRepositoryFactory: mockAgentPromptTemplateRepositoryFactory,
		mockAgentPromptTemplateRepository:        mockAgentPromptTemplateRepository,
	}
}

func generateValidGetAgentPromptTemplateRequest() *requests.GetAgentPromptTemplateRequest {
	return &requests.GetAgentPromptTemplateRequest{
		AgentId:     uuid.New(),
		TemplateKey: "facilitator-default",
	}
}

func TestWhenGettingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given invalid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.AgentId = uuid.Nil

			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})

		t.Run("Because templateKey is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.TemplateKey = ""

			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidGetAgentPromptTemplateRequest()
			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given error occurs while getting agent prompt template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentPromptTemplateRepository)

			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidGetAgentPromptTemplateRequest()
			resp, err := suite.sut.Get(request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	t.Run("Given agent prompt template not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentPromptTemplateRepository)

			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should empty response", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidGetAgentPromptTemplateRequest()
			resp, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "", resp.PromptTemplate)
		})
	})

	t.Run("Given agent prompt template found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentPromptTemplateRepository)

			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(model_builders.NewAgentPromptTemplateBuilder().WithContent("test-content").Build(), nil)
		}

		t.Run("Should return prompt template", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidGetAgentPromptTemplateRequest()
			resp, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, &responses.GetAgentPromptTemplateResponse{PromptTemplate: "test-content"}, resp)
		})
	})
}
