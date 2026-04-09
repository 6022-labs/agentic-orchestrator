package repositories_test

import (
	"sort"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAllContextConversationMappingsTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAllContextConversationMappingsTestCases = []*WhenGettingAllContextConversationMappingsTestCase{
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

type WhenGettingAllContextConversationMappingsTestingSuite struct {
	sut *repositories.AgentContextConversationMappingRepository

	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAllContextConversationMappingsBeforeEach(testCase *WhenGettingAllContextConversationMappingsTestCase) *WhenGettingAllContextConversationMappingsTestingSuite {
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

	return &WhenGettingAllContextConversationMappingsTestingSuite{
		sut: sut,

		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAllContextConversationMappingsAfterEach(suite *WhenGettingAllContextConversationMappingsTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAllContextConversationMappings(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAllContextConversationMappingsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return empty list", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllContextConversationMappingsBeforeEach(testCase)
					defer WhenGettingAllContextConversationMappingsAfterEach(suite)

					contextKeys, err := suite.sut.GetAll()

					assert.Nil(t, err)
					assert.Empty(t, contextKeys)
				})
			})

			t.Run("Given entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return list of context keys", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllContextConversationMappingsBeforeEach(testCase)
					defer WhenGettingAllContextConversationMappingsAfterEach(suite)

					context := model_builders.NewAgentContextBuilder().
						WithKey("key1").
						Build()

					suite.agentContextRepository.Save(context)

					conversationMapping1 := model_builders.NewAgentContextConversationMappingBuilder().
						WithAgentContextId(context.ID).
						WithConversationId("conversation1").
						Build()

					conversationMapping2 := model_builders.NewAgentContextConversationMappingBuilder().
						WithAgentContextId(context.ID).
						WithConversationId("conversation2").
						Build()

					suite.sut.Save(conversationMapping1)
					suite.sut.Save(conversationMapping2)

					conversationMappings, err := suite.sut.GetAll()

					sort.SliceStable(conversationMappings, func(i, j int) bool {
						return conversationMappings[i].ConversationId < conversationMappings[j].ConversationId
					})

					assert.Nil(t, err)
					assert.Len(t, conversationMappings, 2)
					assert.Equal(t, conversationMapping1.ConversationId, conversationMappings[0].ConversationId)
					assert.Equal(t, conversationMapping2.ConversationId, conversationMappings[1].ConversationId)
				})
			})
		})
	}
}
