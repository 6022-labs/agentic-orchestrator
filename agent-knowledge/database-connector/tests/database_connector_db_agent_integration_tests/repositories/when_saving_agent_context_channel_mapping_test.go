package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenSavingAgentContextChannelMappingTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenSavingAgentContextChannelMappingTestCases = []*WhenSavingAgentContextChannelMappingTestCase{
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

type WhenSavingAgentContextChannelMappingTestingSuite struct {
	sut                    *repositories.AgentContextConversationMappingRepository
	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenSavingAgentContextChannelMappingBeforeEach(testCase *WhenSavingAgentContextChannelMappingTestCase) *WhenSavingAgentContextChannelMappingTestingSuite {
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

	return &WhenSavingAgentContextChannelMappingTestingSuite{
		sut:                    sut,
		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenSavingAgentContextChannelMappingAfterEach(suite *WhenSavingAgentContextChannelMappingTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenSavingAgentContextChannelMapping(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenSavingAgentContextChannelMappingTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entity does not exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return created entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextChannelMappingBeforeEach(testCase)
					defer WhenSavingAgentContextChannelMappingAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextChannelMapping := model_builders.NewAgentContextConversationMappingBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					createdEntity, err := suite.sut.Save(agentContextChannelMapping)

					assert.Nil(t, err)
					assert.Equal(t, agentContextChannelMapping, createdEntity)
				})

				t.Run("Should save entity in the database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextChannelMappingBeforeEach(testCase)
					defer WhenSavingAgentContextChannelMappingAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextChannelMapping := model_builders.NewAgentContextConversationMappingBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					createdEntity, _ := suite.sut.Save(agentContextChannelMapping)

					suite.dbClient.GetDB().First(&agentContextChannelMapping, createdEntity.ID)

					assert.Equal(t, agentContextChannelMapping.ID, createdEntity.ID)
					assert.Equal(t, agentContextChannelMapping.AgentContextId, createdEntity.AgentContextId)
				})
			})

			t.Run("Given entity already exists", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return updated entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextChannelMappingBeforeEach(testCase)
					defer WhenSavingAgentContextChannelMappingAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextChannelMapping := model_builders.NewAgentContextConversationMappingBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					createdEntity, _ := suite.sut.Save(agentContextChannelMapping)

					createdEntity.ConversationId = "new_conversation_id"
					updatedEntity, err := suite.sut.Save(createdEntity)

					assert.Nil(t, err)
					assert.Equal(t, createdEntity, updatedEntity)
				})

				t.Run("Should update entity in the database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextChannelMappingBeforeEach(testCase)
					defer WhenSavingAgentContextChannelMappingAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextChannelMapping := model_builders.NewAgentContextConversationMappingBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					createdEntity, _ := suite.sut.Save(agentContextChannelMapping)

					createdEntity.ConversationId = "new_conversation_id"
					updatedEntity, err := suite.sut.Save(createdEntity)

					assert.Nil(t, err)
					assert.Equal(t, createdEntity.AgentContextId, updatedEntity.AgentContextId)
				})
			})
		})
	}
}
