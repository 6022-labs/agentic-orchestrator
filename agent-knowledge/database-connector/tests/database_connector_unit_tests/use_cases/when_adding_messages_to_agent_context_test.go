package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/model_builders"
	db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/db_clients_mocks"
	factories_db_clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/db_clients_mocks"
	factories_repositories_mocks "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/factories_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenAddingMessagesToAgentContextTestingSuite struct {
	sut *use_cases.AddMessagesToAgentContext

	mockAgentDBClientFactory *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockAgentDbClient        *db_clients_mocks.MockAgentDatabaseClientInterface

	mockAgentContextRepositoryFactory                    *factories_repositories_mocks.MockAgentContextRepositoryFactoryInterface
	mockAgentContextMessageRepositoryFactory             *factories_repositories_mocks.MockAgentContextMessageRepositoryFactoryInterface
	mockAgentContextConversationMappingRepositoryFactory *factories_repositories_mocks.MockAgentContextConversationMappingRepositoryFactoryInterface

	mockAgentContextRepository                    *repositories_mocks.MockAgentContextRepositoryInterface
	mockAgentContextMessageRepository             *repositories_mocks.MockAgentContextMessageRepositoryInterface
	mockAgentContextConversationMappingRepository *repositories_mocks.MockAgentContextConversationMappingRepositoryInterface
}

func WhenAddingMessagesToAgentContextBeforeEach(t *testing.T) *WhenAddingMessagesToAgentContextTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)
	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)

	mockAgentContextRepositoryFactory := factories_repositories_mocks.NewMockAgentContextRepositoryFactoryInterface(mockController)
	mockAgentContextMessageRepositoryFactory := factories_repositories_mocks.NewMockAgentContextMessageRepositoryFactoryInterface(mockController)
	mockAgentContextConversationMappingRepositoryFactory := factories_repositories_mocks.NewMockAgentContextConversationMappingRepositoryFactoryInterface(mockController)

	mockAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)
	mockAgentContextMessageRepository := repositories_mocks.NewMockAgentContextMessageRepositoryInterface(mockController)
	mockAgentContextConversationMappingRepository := repositories_mocks.NewMockAgentContextConversationMappingRepositoryInterface(mockController)

	sut := use_cases.NewAddMessagesToAgentContext(
		mockAgentDBClientFactory,
		mockAgentContextRepositoryFactory,
		mockAgentContextMessageRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory,
	)

	return &WhenAddingMessagesToAgentContextTestingSuite{
		sut: sut,

		mockAgentDBClientFactory: mockAgentDBClientFactory,
		mockAgentDbClient:        mockAgentDbClient,

		mockAgentContextRepositoryFactory:                    mockAgentContextRepositoryFactory,
		mockAgentContextMessageRepositoryFactory:             mockAgentContextMessageRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory: mockAgentContextConversationMappingRepositoryFactory,

		mockAgentContextRepository:                    mockAgentContextRepository,
		mockAgentContextMessageRepository:             mockAgentContextMessageRepository,
		mockAgentContextConversationMappingRepository: mockAgentContextConversationMappingRepository,
	}
}

func generateAddMessageToAgentContextRequest() *requests.AddMessagesToAgentContextRequest {
	conversationId := "test-conversation-id"

	return &requests.AddMessagesToAgentContextRequest{
		AgentId:        uuid.New(),
		ConversationId: &conversationId,
		Messages: []requests.Message{
			{
				Content:  "test-content",
				Author:   "test-author",
				Role:     enums.AgentContextMessageRoleUser,
				Platform: enums.AgentContextMessagePlatformSlack,
			},
		},
	}
}

func TestWhenAddingMessageToAgentContext(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			request.AgentId = uuid.Nil

			err := sut.sut.Add(request)
			assert.Error(t, err)
		})

		t.Run("Because messages are empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			request.Messages = nil

			err := sut.sut.Add(request)
			assert.Error(t, err)
		})

		t.Run("Because conversation id and context key are both empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			request.ConversationId = nil
			request.ContextKey = nil

			err := sut.sut.Add(request)
			assert.Error(t, err)
		})

		t.Run("Because both conversation id and context key are set", func(t *testing.T) {
			t.Parallel()

			sut := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			contextKey := "test-context-key"
			request.ContextKey = &contextKey

			err := sut.sut.Add(request)
			assert.Error(t, err)
		})

		t.Run("Because one of the messages contains an empty content", func(t *testing.T) {
			t.Parallel()

			sut := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			request.Messages[0].Content = ""

			err := sut.sut.Add(request)
			assert.Error(t, err)
		})

		t.Run("Because one of the messages contains an empty platform", func(t *testing.T) {
			t.Parallel()

			sut := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			request.Messages[0].Platform = ""

			err := sut.sut.Add(request)
			assert.Error(t, err)
		})

		t.Run("Because one of the messages contains an invalid role", func(t *testing.T) {
			t.Parallel()

			sut := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			request.Messages[0].Role = "invalid-role"

			err := sut.sut.Add(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateAddMessageToAgentContextRequest()
			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given conversation id is set but error occurs while getting conversation mapping", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateAddMessageToAgentContextRequest()
			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given conversation id is set but conversation mapping does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			suite.mockAgentContextConversationMappingRepository.EXPECT().
				GetByConversationId(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateAddMessageToAgentContextRequest()
			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given context key is set but agent context does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			contextKey := "test-context-key"
			suite.mockAgentContextRepository.EXPECT().
				GetByKey(contextKey).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateAddMessageToAgentContextRequest()
			request.ConversationId = nil
			contextKey := "test-context-key"
			request.ContextKey = &contextKey

			err := suite.sut.Add(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting existing messages of agent context", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			contextKey := "test-context-key"
			agentContext := model_builders.NewAgentContextBuilder().Build()
			suite.mockAgentContextRepository.EXPECT().
				GetByKey(contextKey).
				Return(agentContext, nil)

			suite.mockAgentContextMessageRepository.EXPECT().
				GetByAgentContextId(agentContext.ID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateAddMessageToAgentContextRequest()
			request.ConversationId = nil
			contextKey := "test-context-key"
			request.ContextKey = &contextKey

			err := suite.sut.Add(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving message", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			contextKey := "test-context-key"
			agentContext := model_builders.NewAgentContextBuilder().Build()
			suite.mockAgentContextRepository.EXPECT().
				GetByKey(contextKey).
				Return(agentContext, nil)

			suite.mockAgentContextMessageRepository.EXPECT().
				GetByAgentContextId(agentContext.ID).
				Return(nil, nil)

			suite.mockAgentContextMessageRepository.EXPECT().
				SaveMultiple(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)
			initSuite(suite)

			request := generateAddMessageToAgentContextRequest()
			request.ConversationId = nil
			contextKey := "test-context-key"
			request.ContextKey = &contextKey

			err := suite.sut.Add(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given message has already equivalent hash stored", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite, request *requests.AddMessagesToAgentContextRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			contextKey := "test-context-key"
			agentContext := model_builders.NewAgentContextBuilder().Build()
			suite.mockAgentContextRepository.EXPECT().
				GetByKey(contextKey).
				Return(agentContext, nil)

			suite.mockAgentContextMessageRepository.EXPECT().
				GetByAgentContextId(agentContext.ID).
				DoAndReturn(func(agentContextId uuid.UUID) ([]*models.AgentContextMessage, error) {
					messages := make([]*models.AgentContextMessage, len(request.Messages))
					for i, message := range request.Messages {
						messages[i] = message.ToAgentContextMessage(agentContextId)
						err := messages[i].ComputeHash()
						if err != nil {
							panic(err)
						}
					}

					return messages, nil
				})

			suite.mockAgentContextMessageRepository.EXPECT().
				SaveMultiple(gomock.Any()).
				DoAndReturn(func(messages []*models.AgentContextMessage) error {
					assert.Len(t, messages, 0)
					return nil
				})
		}

		t.Run("Should not save message", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)

			request := generateAddMessageToAgentContextRequest()
			request.ConversationId = nil
			contextKey := "test-context-key"
			request.ContextKey = &contextKey

			initSuite(suite, request)

			err := suite.sut.Add(request)
			assert.NoError(t, err)
		})
	})

	t.Run("Given message has not already equivalent hash stored", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessagesToAgentContextTestingSuite, request *requests.AddMessagesToAgentContextRequest) {
			suite.mockAgentDBClientFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentContextRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextRepository)

			suite.mockAgentContextMessageRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextMessageRepository)

			suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
				Create(gomock.Any()).
				Return(suite.mockAgentContextConversationMappingRepository)

			contextKey := "test-context-key"
			agentContext := model_builders.NewAgentContextBuilder().Build()
			suite.mockAgentContextRepository.EXPECT().
				GetByKey(contextKey).
				Return(agentContext, nil)

			suite.mockAgentContextMessageRepository.EXPECT().
				GetByAgentContextId(agentContext.ID).
				Return(nil, nil)

			suite.mockAgentContextMessageRepository.EXPECT().
				SaveMultiple(gomock.Any()).
				DoAndReturn(func(messages []*models.AgentContextMessage) error {
					assert.Len(t, messages, len(request.Messages))
					return nil
				})
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessagesToAgentContextBeforeEach(t)
			request := generateAddMessageToAgentContextRequest()
			request.ConversationId = nil
			contextKey := "test-context-key"
			request.ContextKey = &contextKey

			initSuite(suite, request)

			err := suite.sut.Add(request)
			assert.NoError(t, err)
		})
	})
}
