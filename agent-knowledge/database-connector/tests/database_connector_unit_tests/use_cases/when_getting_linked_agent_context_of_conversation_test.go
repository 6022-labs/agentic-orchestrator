package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/responses"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/model_builders"
	db_clients_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/db_clients_mocks"
	factories_db_clients_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/db_clients_mocks"
	factories_repositories_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingLinkedAgentContextOfConversationTestingSuite struct {
	sut *use_cases.GetLinkedAgentContextOfConversation

	mockAgentDBClientFactory                             *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient                                    *db_clients_mocks.MockAgentDatabaseClientInterface
	mockAgentContextRepositoryFactory                    *factories_repositories_mocks.MockAgentContextRepositoryFactoryInterface
	mockAgentContextConversationMappingRepositoryFactory *factories_repositories_mocks.MockAgentContextConversationMappingRepositoryFactoryInterface

	mockAgentContextRepository                    *repositories_mocks.MockAgentContextRepositoryInterface
	mockAgentContextConversationMappingRepository *repositories_mocks.MockAgentContextConversationMappingRepositoryInterface
}

func WhenGettingLinkedAgentContextOfConversationBeforeEach(t *testing.T) *WhenGettingLinkedAgentContextOfConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockAgentContextRepositoryFactory := factories_repositories_mocks.NewMockAgentContextRepositoryFactoryInterface(mockController)
	mockAgentContextConversationMappingRepositoryFactory := factories_repositories_mocks.NewMockAgentContextConversationMappingRepositoryFactoryInterface(mockController)

	mockAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)
	mockAgentContextConversationMappingRepository := repositories_mocks.NewMockAgentContextConversationMappingRepositoryInterface(mockController)

	sut := use_cases.NewGetLinkedAgentContextOfConversation(
		mockAgentDBClientFactory,
		mockAgentContextRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory,
	)

	return &WhenGettingLinkedAgentContextOfConversationTestingSuite{
		sut: sut,

		mockAgentDBClientFactory:                             mockAgentDBClientFactory,
		mockAgentDbClient:                                    mockAgentDbClient,
		mockAgentContextRepositoryFactory:                    mockAgentContextRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory: mockAgentContextConversationMappingRepositoryFactory,

		mockAgentContextRepository:                    mockAgentContextRepository,
		mockAgentContextConversationMappingRepository: mockAgentContextConversationMappingRepository,
	}
}

func generateValidGetLinkedAgentContextOfConversationRequest() *requests.GetLinkedAgentContextOfConversationRequest {
	return &requests.GetLinkedAgentContextOfConversationRequest{
		AgentId:        uuid.New(),
		ConversationId: uuid.New().String(),
	}
}

func TestWhenGettingLinkedAgentContextOfConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			request.AgentId = uuid.Nil

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			request.ConversationId = ""

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingLinkedAgentContextOfConversationTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)
			initSuite(sut)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			_, err := sut.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting conversation mapping", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingLinkedAgentContextOfConversationTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)
			initSuite(sut)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			_, err := sut.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given conversation mapping not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingLinkedAgentContextOfConversationTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return not found status", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)
			initSuite(sut)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			response, err := sut.sut.Get(request)

			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Empty(t, response.Key)
			assert.Equal(t, responses.GetLinkedAgentContextOfConversationResponseStatusNotFound, response.Status)
		})
	})

	t.Run("Given error occurs while getting agent context", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingLinkedAgentContextOfConversationTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)
			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			conversationMapping := model_builders.NewAgentContextConversationMappingBuilder().Build()

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(conversationMapping, nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetById(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)
			initSuite(sut)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			_, err := sut.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent context not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingLinkedAgentContextOfConversationTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)
			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			conversationMapping := model_builders.NewAgentContextConversationMappingBuilder().Build()

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(conversationMapping, nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)
			initSuite(sut)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			_, err := sut.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent context found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingLinkedAgentContextOfConversationTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)
			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			conversationMapping := model_builders.NewAgentContextConversationMappingBuilder().Build()

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(conversationMapping, nil)

			agentContext := model_builders.NewAgentContextBuilder().Build()

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetById(gomock.Any()).
				Return(agentContext, nil)
		}

		t.Run("Should return agent context", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingLinkedAgentContextOfConversationBeforeEach(t)
			initSuite(sut)

			request := generateValidGetLinkedAgentContextOfConversationRequest()
			response, err := sut.sut.Get(request)

			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.NotEmpty(t, response.Key)
			assert.Equal(t, responses.GetLinkedAgentContextOfConversationResponseStatusFound, response.Status)
		})
	})
}
