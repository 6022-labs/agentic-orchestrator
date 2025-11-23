package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/db_clients_mocks"
	factories_db_clients_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/db_clients_mocks"
	factories_repositories_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.SaveAgentPromptTemplate

	mockAgentDBClientFactory                 *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentPromptTemplateRepositoryFactory *factories_repositories_mocks.MockAgentPromptTemplateRepositoryFactoryInterface

	mockAgentPromptTemplateRepository *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
	mockAgentDbClient                 *db_clients_mocks.MockAgentDatabaseClientInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentPromptTemplateRepositoryFactory := factories_repositories_mocks.NewMockAgentPromptTemplateRepositoryFactoryInterface(mockController)
	mockAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)

	sut := use_cases.NewSaveAgentPromptTemplate(
		mockAgentDBClientFactory,
		mockAgentPromptTemplateRepositoryFactory,
	)

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentDBClientFactory:                 mockAgentDBClientFactory,
		mockAgentPromptTemplateRepositoryFactory: mockAgentPromptTemplateRepositoryFactory,

		mockAgentPromptTemplateRepository: mockAgentPromptTemplateRepository,
		mockAgentDbClient:                 mockAgentDbClient,
	}
}

func generateValidSaveAgentPromptTemplateRequest() *requests.SaveAgentPromptTemplateRequest {
	return &requests.SaveAgentPromptTemplateRequest{
		AgentId: uuid.New(),
		Key:     enums.AgentPromptTemplateKeyFacilitatorDefault,
		Content: "template",
	}
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given invalid request", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agent ID is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.AgentId = uuid.Nil

			err := suite.sut.Save(request)

			assert.Error(t, err)
		})

		t.Run("Because template is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Content = ""

			err := suite.sut.Save(request)

			assert.Error(t, err)
		})

		t.Run("Because key is invalid", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			request := generateValidSaveAgentPromptTemplateRequest()
			request.Key = "invalid"

			err := suite.sut.Save(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidSaveAgentPromptTemplateRequest()
			err := suite.sut.Save(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving agent context", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
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

			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidSaveAgentPromptTemplateRequest()
			err := suite.sut.Save(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent prompt template does not exist for this key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
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

			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any()).
				Return(model_builders.NewAgentPromptTemplateBuilder().Build(), nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidSaveAgentPromptTemplateRequest()
			err := suite.sut.Save(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given agent prompt template already exists for this key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)
			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentPromptTemplateRepository)

			suite.mockAgentPromptTemplateRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(model_builders.NewAgentPromptTemplateBuilder().Build(), nil)

			suite.mockAgentPromptTemplateRepository.EXPECT().
				Save(gomock.Any()).
				Return(model_builders.NewAgentPromptTemplateBuilder().Build(), nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite)

			request := generateValidSaveAgentPromptTemplateRequest()
			err := suite.sut.Save(request)

			assert.NoError(t, err)
		})
	})
}
