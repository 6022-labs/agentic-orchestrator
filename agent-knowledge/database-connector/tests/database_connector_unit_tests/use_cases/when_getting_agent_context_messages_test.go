package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/model_builders"
	db_clients_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/db_clients_mocks"
	factories_db_clients_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/db_clients_mocks"
	factories_repositories_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentContextMessagesTestingSuite struct {
	sut *use_cases.GetAgentContextMessages

	mockAgentDBClientFactory                 *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient                        *db_clients_mocks.MockAgentDatabaseClientInterface
	mockAgentContextMessageRepositoryFactory *factories_repositories_mocks.MockAgentContextMessageRepositoryFactoryInterface

	mockAgentContextMessageRepository *repositories_mocks.MockAgentContextMessageRepositoryInterface
}

func WhenGettingAgentContextMessagesBeforeEach(t *testing.T) *WhenGettingAgentContextMessagesTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockAgentContextMessageRepositoryFactory := factories_repositories_mocks.NewMockAgentContextMessageRepositoryFactoryInterface(mockController)

	mockAgentContextMessageRepository := repositories_mocks.NewMockAgentContextMessageRepositoryInterface(mockController)

	sut := use_cases.NewGetAgentContextMessages(
		mockAgentDBClientFactory,
		mockAgentContextMessageRepositoryFactory,
	)

	return &WhenGettingAgentContextMessagesTestingSuite{
		sut: sut,

		mockAgentDBClientFactory:                 mockAgentDBClientFactory,
		mockAgentDbClient:                        mockAgentDbClient,
		mockAgentContextMessageRepositoryFactory: mockAgentContextMessageRepositoryFactory,

		mockAgentContextMessageRepository: mockAgentContextMessageRepository,
	}
}

func generateValidGetAgentContextMessagesRequest() *requests.GetAgentContextMessagesRequest {
	return &requests.GetAgentContextMessagesRequest{
		AgentId: uuid.New(),
		Key:     uuid.New().String(),
	}
}

func TestWhenGettingAgentContextMessages(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextMessagesBeforeEach(t)

			request := generateValidGetAgentContextMessagesRequest()
			request.AgentId = uuid.Nil

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})

		t.Run("Because key is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentContextMessagesBeforeEach(t)

			request := generateValidGetAgentContextMessagesRequest()
			request.Key = ""

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextMessagesTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingAgentContextMessagesBeforeEach(t)
			initSuite(sut)

			request := generateValidGetAgentContextMessagesRequest()
			_, err := sut.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting messages", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextMessagesTestingSuite, request *requests.GetAgentContextMessagesRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextMessageRepository.EXPECT().
				GetByAgentContextKey(request.Key).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingAgentContextMessagesBeforeEach(t)
			request := generateValidGetAgentContextMessagesRequest()

			initSuite(sut, request)

			_, err := sut.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given messages not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextMessagesTestingSuite, request *requests.GetAgentContextMessagesRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextMessageRepository.EXPECT().
				GetByAgentContextKey(request.Key).
				Return(nil, nil)
		}

		t.Run("Should return empty response", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingAgentContextMessagesBeforeEach(t)
			request := generateValidGetAgentContextMessagesRequest()

			initSuite(sut, request)

			response, err := sut.sut.Get(request)

			assert.NoError(t, err)
			assert.Empty(t, response)
		})
	})

	t.Run("Given messages found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentContextMessagesTestingSuite, request *requests.GetAgentContextMessagesRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			messages := []*models.AgentContextMessage{
				model_builders.NewAgentContextMessageBuilder().Build(),
				model_builders.NewAgentContextMessageBuilder().Build(),
			}

			suite.mockAgentContextMessageRepository.EXPECT().
				GetByAgentContextKey(request.Key).
				Return(messages, nil)
		}

		t.Run("Should return messages", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingAgentContextMessagesBeforeEach(t)
			request := generateValidGetAgentContextMessagesRequest()

			initSuite(sut, request)

			response, err := sut.sut.Get(request)

			assert.NoError(t, err)
			assert.NotEmpty(t, response)
		})
	})
}
