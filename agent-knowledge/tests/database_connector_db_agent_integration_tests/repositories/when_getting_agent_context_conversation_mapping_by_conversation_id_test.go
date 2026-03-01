package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAgentContextConversationMappingByConversationIdTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAgentContextConversationMappingByConversationIdTestCases = []*WhenGettingAgentContextConversationMappingByConversationIdTestCase{
		{
			name:   "PostgreSQL",
			driver: enums.DatabaseConnectionDriverPostgres,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetPostgresContainerBuilder()
			},
		},
		{
			name:   "MariaDB",
			driver: enums.DatabaseConnectionDriverMariaDB,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMariaDBContainerBuilder()
			},
		},
		{
			name:   "MySQL",
			driver: enums.DatabaseConnectionDriverMySQL,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMySQLContainerBuilder()
			},
		},
	}
)

type WhenGettingAgentContextConversationMappingByConversationIdTestingSuite struct {
	sut                    *repositories.AgentContextConversationMappingRepository
	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAgentContextConversationMappingByConversationIdBeforeEach(testCase *WhenGettingAgentContextConversationMappingByConversationIdTestCase) *WhenGettingAgentContextConversationMappingByConversationIdTestingSuite {
	databaseContainerBuilder := testCase.clientBuilderFn()
	dbName, connectionString := databaseContainerBuilder.Create()

	migrationFolderPath, err := common_tests.FindAgentMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, err := database_connector_db_agent.NewAgentDbClient(
		connectionString,
		zap.NewNop(),
		testCase.driver,
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err != nil {
		panic(err)
	}

	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewAgentContextConversationMappingRepository(dbClient)
	agentContextRepository := repositories.NewAgentContextRepository(dbClient)

	return &WhenGettingAgentContextConversationMappingByConversationIdTestingSuite{
		sut:                    sut,
		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAgentContextConversationMappingByConversationIdAfterEach(suite *WhenGettingAgentContextConversationMappingByConversationIdTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAgentContextConversationMappingByChannel(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAgentContextConversationMappingByConversationIdTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entity does not exist for the channel", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return nil value", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextConversationMappingByConversationIdBeforeEach(testCase)
					defer WhenGettingAgentContextConversationMappingByConversationIdAfterEach(suite)

					conversationId := uuid.New().String()
					agentContextChannelMapping, _ := suite.sut.GetByConversationId(conversationId)

					assert.Nil(t, agentContextChannelMapping)
				})
			})

			t.Run("Given entity exist for the channel", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextConversationMappingByConversationIdBeforeEach(testCase)
					defer WhenGettingAgentContextConversationMappingByConversationIdAfterEach(suite)

					conversationId := uuid.New().String()

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextChannelMapping := model_builders.NewAgentContextConversationMappingBuilder().
						WithConversationId(conversationId).
						WithAgentContextId(agentContext.ID).
						Build()

					savedEntity, _ := suite.sut.Save(agentContextChannelMapping)

					agentContextChannelMapping, err := suite.sut.GetByConversationId(conversationId)

					assert.Nil(t, err)
					assert.Equal(t, agentContextChannelMapping.ID, savedEntity.ID)
				})
			})
		})
	}
}
