package use_cases_test

import (
	"testing"

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

type WhenLinkingConversationToAgentContextTestingSuite struct {
	sut *use_cases.LinkConversationToAgentContext

	mockAgentDBClientFactory                             *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient                                    *db_clients_mocks.MockAgentDatabaseClientInterface
	mockAgentContextRepositoryFactory                    *factories_repositories_mocks.MockAgentContextRepositoryFactoryInterface
	mockAgentContextConversationMappingRepositoryFactory *factories_repositories_mocks.MockAgentContextConversationMappingRepositoryFactoryInterface

	mockAgentContextRepository                    *repositories_mocks.MockAgentContextRepositoryInterface
	mockAgentContextConversationMappingRepository *repositories_mocks.MockAgentContextConversationMappingRepositoryInterface
}

func WhenLinkingConversationToAgentContextBeforeEach(t *testing.T) *WhenLinkingConversationToAgentContextTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockAgentContextRepositoryFactory := factories_repositories_mocks.NewMockAgentContextRepositoryFactoryInterface(mockController)
	mockAgentContextConversationMappingRepositoryFactory := factories_repositories_mocks.NewMockAgentContextConversationMappingRepositoryFactoryInterface(mockController)

	mockAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)
	mockAgentContextConversationMappingRepository := repositories_mocks.NewMockAgentContextConversationMappingRepositoryInterface(mockController)

	sut := use_cases.NewLinkConversationToAgentContext(
		mockAgentDBClientFactory,
		mockAgentContextRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory,
	)

	return &WhenLinkingConversationToAgentContextTestingSuite{
		sut: sut,

		mockAgentDBClientFactory:                             mockAgentDBClientFactory,
		mockAgentDbClient:                                    mockAgentDbClient,
		mockAgentContextRepositoryFactory:                    mockAgentContextRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory: mockAgentContextConversationMappingRepositoryFactory,

		mockAgentContextRepository:                    mockAgentContextRepository,
		mockAgentContextConversationMappingRepository: mockAgentContextConversationMappingRepository,
	}
}

func generateValidLinkConversationToAgentContextRequest() *requests.LinkConversationToAgentContextRequest {
	return &requests.LinkConversationToAgentContextRequest{
		AgentId:        uuid.New(),
		Key:            uuid.New().String(),
		ConversationId: uuid.New().String(),
	}
}

func TestWhenLinkingConversationToAgentContext(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenLinkingConversationToAgentContextBeforeEach(t)

			request := generateValidLinkConversationToAgentContextRequest()
			request.AgentId = uuid.Nil

			err := suite.sut.Link(request)

			assert.Error(t, err)
		})

		t.Run("Because key is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenLinkingConversationToAgentContextBeforeEach(t)

			request := generateValidLinkConversationToAgentContextRequest()
			request.Key = ""

			err := suite.sut.Link(request)

			assert.Error(t, err)
		})

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenLinkingConversationToAgentContextBeforeEach(t)

			request := generateValidLinkConversationToAgentContextRequest()
			request.ConversationId = ""

			err := suite.sut.Link(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenLinkingConversationToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenLinkingConversationToAgentContextBeforeEach(t)
			initSuite(sut)

			request := generateValidLinkConversationToAgentContextRequest()
			err := sut.sut.Link(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting conversation mapping", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenLinkingConversationToAgentContextTestingSuite) {
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

			sut := WhenLinkingConversationToAgentContextBeforeEach(t)
			initSuite(sut)

			request := generateValidLinkConversationToAgentContextRequest()
			err := sut.sut.Link(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given conversation mapping already exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenLinkingConversationToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			convoMapping := model_builders.NewAgentContextConversationMappingBuilder().Build()

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(convoMapping, nil)
		}

		t.Run("Should return no error", func(t *testing.T) {
			t.Parallel()

			sut := WhenLinkingConversationToAgentContextBeforeEach(t)
			initSuite(sut)

			request := generateValidLinkConversationToAgentContextRequest()
			err := sut.sut.Link(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given error occurs while getting agent context by key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenLinkingConversationToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenLinkingConversationToAgentContextBeforeEach(t)
			initSuite(sut)

			request := generateValidLinkConversationToAgentContextRequest()
			err := sut.sut.Link(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent context not found by key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenLinkingConversationToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenLinkingConversationToAgentContextBeforeEach(t)
			initSuite(sut)

			request := generateValidLinkConversationToAgentContextRequest()
			err := sut.sut.Link(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving conversation mapping", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenLinkingConversationToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(model_builders.NewAgentContextBuilder().Build(), nil)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenLinkingConversationToAgentContextBeforeEach(t)
			initSuite(sut)

			request := generateValidLinkConversationToAgentContextRequest()
			err := sut.sut.Link(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given conversation mapping is saved successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenLinkingConversationToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextRepository.EXPECT().
				GetByKey(gomock.Any()).
				Return(model_builders.NewAgentContextBuilder().Build(), nil)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return no error", func(t *testing.T) {
			t.Parallel()

			sut := WhenLinkingConversationToAgentContextBeforeEach(t)
			initSuite(sut)

			request := generateValidLinkConversationToAgentContextRequest()
			err := sut.sut.Link(request)

			assert.NoError(t, err)
		})
	})
}
