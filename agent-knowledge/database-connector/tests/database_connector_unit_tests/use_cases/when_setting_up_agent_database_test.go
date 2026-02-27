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
	"go.uber.org/zap"
)

type WhenSettingUpAgentDatabaseTestingSuite struct {
	sut *use_cases.SetupAgentDatabase

	mockAgentDBClientFactory         *factories_db_clients_mocks.MockAgentDatabaseClientFactoryInterface
	mockDatabaseConnectionRepository *repositories_mocks.MockDatabaseConnectionRepositoryInterface

	mockAgentContextRepositoryFactory                    *factories_repositories_mocks.MockAgentContextRepositoryFactoryInterface
	mockAgentPromptTemplateRepositoryFactory             *factories_repositories_mocks.MockAgentPromptTemplateRepositoryFactoryInterface
	mockAgentContextMessageRepositoryFactory             *factories_repositories_mocks.MockAgentContextMessageRepositoryFactoryInterface
	mockAgentContextConversationMappingRepositoryFactory *factories_repositories_mocks.MockAgentContextConversationMappingRepositoryFactoryInterface

	mockCopyToAgentContextRepository                    *repositories_mocks.MockAgentContextRepositoryInterface
	mockCopyToAgentPromptTemplateRepository             *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
	mockCopyToAgentContextMessageRepository             *repositories_mocks.MockAgentContextMessageRepositoryInterface
	mockCopyToAgentContextConversationMappingRepository *repositories_mocks.MockAgentContextConversationMappingRepositoryInterface

	mockCopyFromAgentContextRepository                    *repositories_mocks.MockAgentContextRepositoryInterface
	mockCopyFromAgentPromptTemplateRepository             *repositories_mocks.MockAgentPromptTemplateRepositoryInterface
	mockCopyFromAgentContextMessageRepository             *repositories_mocks.MockAgentContextMessageRepositoryInterface
	mockCopyFromAgentContextConversationMappingRepository *repositories_mocks.MockAgentContextConversationMappingRepositoryInterface

	mockAgentDbClient *db_clients_mocks.MockAgentDatabaseClientInterface

	mockCopyToAgentDbClient   *db_clients_mocks.MockAgentDatabaseClientInterface
	mockCopyFromAgentDbClient *db_clients_mocks.MockAgentDatabaseClientInterface
}

func WhenSettingUpAgentDatabaseBeforeEach(t *testing.T) *WhenSettingUpAgentDatabaseTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentDBClientFactory := factories_db_clients_mocks.NewMockAgentDatabaseClientFactoryInterface(mockController)

	mockAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockDatabaseConnectionRepository := repositories_mocks.NewMockDatabaseConnectionRepositoryInterface(mockController)
	mockAgentContextRepositoryFactory := factories_repositories_mocks.NewMockAgentContextRepositoryFactoryInterface(mockController)
	mockAgentPromptTemplateRepositoryFactory := factories_repositories_mocks.NewMockAgentPromptTemplateRepositoryFactoryInterface(mockController)
	mockAgentContextMessageRepositoryFactory := factories_repositories_mocks.NewMockAgentContextMessageRepositoryFactoryInterface(mockController)
	mockAgentContextConversationMappingRepositoryFactory := factories_repositories_mocks.NewMockAgentContextConversationMappingRepositoryFactoryInterface(mockController)

	mockCopyToAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)
	mockCopyFromAgentDbClient := db_clients_mocks.NewMockAgentDatabaseClientInterface(mockController)

	mockCopyToAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)
	mockCopyToAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(mockController)
	mockCopyToAgentContextMessageRepository := repositories_mocks.NewMockAgentContextMessageRepositoryInterface(mockController)
	mockCopyToAgentContextConversationMappingRepository := repositories_mocks.NewMockAgentContextConversationMappingRepositoryInterface(mockController)

	mockCopyFromAgentContextRepository := repositories_mocks.NewMockAgentContextRepositoryInterface(mockController)
	mockCopyFromAgentPromptTemplateRepository := repositories_mocks.NewMockAgentPromptTemplateRepositoryInterface(mockController)
	mockCopyFromAgentContextMessageRepository := repositories_mocks.NewMockAgentContextMessageRepositoryInterface(mockController)
	mockCopyFromAgentContextConversationMappingRepository := repositories_mocks.NewMockAgentContextConversationMappingRepositoryInterface(mockController)

	sut := use_cases.NewSetupAgentDatabase(
		zap.NewNop(),
		mockAgentDBClientFactory,
		mockDatabaseConnectionRepository,
		mockAgentContextRepositoryFactory,
		mockAgentPromptTemplateRepositoryFactory,
		mockAgentContextMessageRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory,
	)

	return &WhenSettingUpAgentDatabaseTestingSuite{
		sut: sut,

		mockAgentDBClientFactory:         mockAgentDBClientFactory,
		mockDatabaseConnectionRepository: mockDatabaseConnectionRepository,

		mockAgentContextRepositoryFactory:                    mockAgentContextRepositoryFactory,
		mockAgentPromptTemplateRepositoryFactory:             mockAgentPromptTemplateRepositoryFactory,
		mockAgentContextMessageRepositoryFactory:             mockAgentContextMessageRepositoryFactory,
		mockAgentContextConversationMappingRepositoryFactory: mockAgentContextConversationMappingRepositoryFactory,

		mockCopyToAgentContextRepository:                    mockCopyToAgentContextRepository,
		mockCopyToAgentPromptTemplateRepository:             mockCopyToAgentPromptTemplateRepository,
		mockCopyToAgentContextMessageRepository:             mockCopyToAgentContextMessageRepository,
		mockCopyToAgentContextConversationMappingRepository: mockCopyToAgentContextConversationMappingRepository,

		mockCopyFromAgentContextRepository:                    mockCopyFromAgentContextRepository,
		mockCopyFromAgentPromptTemplateRepository:             mockCopyFromAgentPromptTemplateRepository,
		mockCopyFromAgentContextMessageRepository:             mockCopyFromAgentContextMessageRepository,
		mockCopyFromAgentContextConversationMappingRepository: mockCopyFromAgentContextConversationMappingRepository,

		mockAgentDbClient:         mockAgentDbClient,
		mockCopyToAgentDbClient:   mockCopyToAgentDbClient,
		mockCopyFromAgentDbClient: mockCopyFromAgentDbClient,
	}
}

func generateValidSetupAgentDatabaseRequest() *requests.SetupAgentDatabaseRequest {
	return &requests.SetupAgentDatabaseRequest{
		AgentId: uuid.New(),
	}
}

func generateValidSetupAgentDatabaseRequestWithCopyFromAgentId() *requests.SetupAgentDatabaseRequest {
	copyFromAgentId := uuid.New()
	return &requests.SetupAgentDatabaseRequest{
		AgentId:         uuid.New(),
		CopyFromAgentId: &copyFromAgentId,
	}
}

func TestWhenSettingUpAgentDatabase(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.Nil,
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting database connection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given database connection is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given database connection is established", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			databaseConnection := model_builders.NewDatabaseConnectionBuilder().
				WithSetupStatus(enums.DatabaseConnectionSetupStatusEstablished).
				Build()

			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(databaseConnection, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given database connection is establishing", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			databaseConnection := model_builders.NewDatabaseConnectionBuilder().
				WithSetupStatus(enums.DatabaseConnectionSetupStatusEstablishing).
				Build()

			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(databaseConnection, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while updating setup status to establishing", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			databaseConnection := model_builders.NewDatabaseConnectionBuilder().
				WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
				Build()

			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(databaseConnection, nil)

			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while generating db client", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			databaseConnection := model_builders.NewDatabaseConnectionBuilder().
				WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
				Build()
			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(databaseConnection, nil)

			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
					return dc, nil
				})

			// Deferred call to update final status to Configured
			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
					return dc, nil
				})

			suite.mockAgentDBClientFactory.EXPECT().
				CreateFromDatabaseConnection(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while running migrations", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			databaseConnection := model_builders.NewDatabaseConnectionBuilder().
				WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
				Build()
			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(databaseConnection, nil)

			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
					return dc, nil
				})

			// Deferred call to update final status to Configured
			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
					return dc, nil
				})

			suite.mockAgentDBClientFactory.EXPECT().
				CreateFromDatabaseConnection(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().
				Close().
				Return(nil)

			suite.mockAgentDbClient.EXPECT().
				Migrate().
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given copy from another agent is requested", func(t *testing.T) {
		t.Parallel()

		t.Run("And error occurs while creating from database client", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithAgentId(request.AgentId).
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(request.AgentId).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while creating to database client", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while querying from database messages", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while querying from database prompt templates", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextMessage{}, nil)

				suite.mockCopyFromAgentPromptTemplateRepository.EXPECT().
					GetAll().
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while querying from database context conversation mappings", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextMessage{}, nil)

				suite.mockCopyFromAgentPromptTemplateRepository.EXPECT().
					GetAll().
					Return([]*models.AgentPromptTemplate{}, nil)

				suite.mockCopyFromAgentContextConversationMappingRepository.EXPECT().
					GetAll().
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while querying from database contexts", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextMessage{}, nil)

				suite.mockCopyFromAgentPromptTemplateRepository.EXPECT().
					GetAll().
					Return([]*models.AgentPromptTemplate{}, nil)

				suite.mockCopyFromAgentContextConversationMappingRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextConversationMapping{}, nil)

				suite.mockCopyFromAgentContextRepository.EXPECT().
					GetAll().
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while saving to database contexts", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextMessage{}, nil)

				suite.mockCopyFromAgentPromptTemplateRepository.EXPECT().
					GetAll().
					Return([]*models.AgentPromptTemplate{}, nil)

				suite.mockCopyFromAgentContextConversationMappingRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextConversationMapping{}, nil)

				suite.mockCopyFromAgentContextRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContext{}, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					BeginTransaction().
					Return(nil, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					RollbackTransaction(gomock.Any()).
					Return(nil)

				suite.mockCopyToAgentContextRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while saving to database prompt templates", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextMessage{}, nil)

				suite.mockCopyFromAgentPromptTemplateRepository.EXPECT().
					GetAll().
					Return([]*models.AgentPromptTemplate{}, nil)

				suite.mockCopyFromAgentContextConversationMappingRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextConversationMapping{}, nil)

				suite.mockCopyFromAgentContextRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContext{}, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					BeginTransaction().
					Return(nil, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					RollbackTransaction(gomock.Any()).
					Return(nil)

				suite.mockCopyToAgentContextRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				suite.mockCopyToAgentPromptTemplateRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while saving to database context conversation mappings", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextMessage{}, nil)

				suite.mockCopyFromAgentPromptTemplateRepository.EXPECT().
					GetAll().
					Return([]*models.AgentPromptTemplate{}, nil)

				suite.mockCopyFromAgentContextConversationMappingRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextConversationMapping{}, nil)

				suite.mockCopyFromAgentContextRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContext{}, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					BeginTransaction().
					Return(nil, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					RollbackTransaction(gomock.Any()).
					Return(nil)

				suite.mockCopyToAgentContextRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				suite.mockCopyToAgentPromptTemplateRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				suite.mockCopyToAgentContextConversationMappingRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while saving to database context messages", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite, request *requests.SetupAgentDatabaseRequest) {
				databaseConnection := model_builders.NewDatabaseConnectionBuilder().
					WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
					Build()
				suite.mockDatabaseConnectionRepository.EXPECT().
					GetByAgentId(gomock.Any()).
					Return(databaseConnection, nil)

				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
						return dc, nil
					})

				// Deferred call to update final status to Configured
				suite.mockDatabaseConnectionRepository.EXPECT().
					Save(gomock.Any()).
					Times(1).
					DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
						assert.Equal(t, enums.DatabaseConnectionSetupStatusConfigured, dc.SetupStatus)
						return dc, nil
					})

				suite.mockAgentDBClientFactory.EXPECT().
					CreateFromDatabaseConnection(gomock.Any()).
					Return(suite.mockAgentDbClient, nil)

				suite.mockAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDbClient.EXPECT().
					Migrate().
					Return(nil)

				// Copy related expectations
				suite.mockAgentDBClientFactory.EXPECT().
					Create(*request.CopyFromAgentId).
					Times(1).
					Return(suite.mockCopyFromAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyFromAgentDbClient).
					Return(suite.mockCopyFromAgentContextConversationMappingRepository)

				suite.mockCopyFromAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockAgentDBClientFactory.EXPECT().
					Create(request.AgentId).
					Times(1).
					Return(suite.mockCopyToAgentDbClient, nil)

				suite.mockAgentContextRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextRepository)

				suite.mockAgentPromptTemplateRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentPromptTemplateRepository)

				suite.mockAgentContextMessageRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextMessageRepository)

				suite.mockAgentContextConversationMappingRepositoryFactory.EXPECT().
					Create(suite.mockCopyToAgentDbClient).
					Return(suite.mockCopyToAgentContextConversationMappingRepository)

				suite.mockCopyToAgentDbClient.EXPECT().
					Close().
					Return(nil)

				suite.mockCopyFromAgentContextMessageRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextMessage{}, nil)

				suite.mockCopyFromAgentPromptTemplateRepository.EXPECT().
					GetAll().
					Return([]*models.AgentPromptTemplate{}, nil)

				suite.mockCopyFromAgentContextConversationMappingRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContextConversationMapping{}, nil)

				suite.mockCopyFromAgentContextRepository.EXPECT().
					GetAll().
					Return([]*models.AgentContext{}, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					BeginTransaction().
					Return(nil, nil)

				suite.mockCopyToAgentDbClient.EXPECT().
					RollbackTransaction(gomock.Any()).
					Return(nil)

				suite.mockCopyToAgentContextRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				suite.mockCopyToAgentPromptTemplateRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				suite.mockCopyToAgentContextConversationMappingRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				suite.mockCopyToAgentContextMessageRepository.EXPECT().
					CreateMultipleUsingTransaction(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSettingUpAgentDatabaseBeforeEach(t)
				request := generateValidSetupAgentDatabaseRequestWithCopyFromAgentId()

				initSuite(sut, request)

				err := sut.sut.Setup(request)

				assert.Error(t, err)
			})
		})
	})

	t.Run("Given error occurs while saving database connection final status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			databaseConnection := model_builders.NewDatabaseConnectionBuilder().
				WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
				Build()
			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(databaseConnection, nil)

			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
					return dc, nil
				})

			// Deferred call to update final status to Configured
			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablished, dc.SetupStatus)
					return dc, assert.AnError
				})

			suite.mockAgentDBClientFactory.EXPECT().
				CreateFromDatabaseConnection(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentDbClient.EXPECT().
				Migrate().
				Return(nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseTestingSuite) {
			databaseConnection := model_builders.NewDatabaseConnectionBuilder().
				WithSetupStatus(enums.DatabaseConnectionSetupStatusConfigured).
				Build()
			suite.mockDatabaseConnectionRepository.EXPECT().
				GetByAgentId(gomock.Any()).
				Return(databaseConnection, nil)

			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablishing, dc.SetupStatus)
					return dc, nil
				})

			// Deferred call to update final status to Configured
			suite.mockDatabaseConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Times(1).
				DoAndReturn(func(dc *models.DatabaseConnection) (*models.DatabaseConnection, error) {
					assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablished, dc.SetupStatus)
					return dc, nil
				})

			suite.mockAgentDBClientFactory.EXPECT().
				CreateFromDatabaseConnection(gomock.Any()).
				Return(suite.mockAgentDbClient, nil)

			suite.mockAgentDbClient.EXPECT().Close().Return(nil)

			suite.mockAgentDbClient.EXPECT().
				Migrate().
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentDatabaseBeforeEach(t)
			initSuite(sut)

			request := &requests.SetupAgentDatabaseRequest{
				AgentId: uuid.New(),
			}

			err := sut.sut.Setup(request)

			assert.NoError(t, err)
		})
	})
}
