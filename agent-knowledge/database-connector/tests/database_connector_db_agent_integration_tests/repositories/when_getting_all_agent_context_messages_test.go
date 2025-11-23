package repositories_test

import (
	"sort"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAllContextMessagesTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAllContextMessagesTestCases = []*WhenGettingAllContextMessagesTestCase{
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

type WhenGettingAllContextMessagesTestingSuite struct {
	sut *repositories.AgentContextMessageRepository

	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAllContextMessagesBeforeEach(testCase *WhenGettingAllContextMessagesTestCase) *WhenGettingAllContextMessagesTestingSuite {
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

	sut := repositories.NewAgentContextMessageRepository(dbClient)

	agentContextRepository := repositories.NewAgentContextRepository(dbClient)

	return &WhenGettingAllContextMessagesTestingSuite{
		sut: sut,

		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAllContextMessagesAfterEach(suite *WhenGettingAllContextMessagesTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAllContextMessages(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAllContextMessagesTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return empty list", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllContextMessagesBeforeEach(testCase)
					defer WhenGettingAllContextMessagesAfterEach(suite)

					contextKeys, err := suite.sut.GetAll()

					assert.Nil(t, err)
					assert.Empty(t, contextKeys)
				})
			})

			t.Run("Given entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return list of context keys", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllContextMessagesBeforeEach(testCase)
					defer WhenGettingAllContextMessagesAfterEach(suite)

					context := model_builders.NewAgentContextBuilder().
						WithKey("key1").
						Build()

					suite.agentContextRepository.Save(context)

					message1 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(context.ID).
						WithContent("message1").
						Build()

					message2 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(context.ID).
						WithContent("message2").
						Build()

					suite.sut.Save(message1)
					suite.sut.Save(message2)

					messages, err := suite.sut.GetAll()

					sort.SliceStable(messages, func(i, j int) bool {
						return messages[i].Content < messages[j].Content
					})

					assert.Nil(t, err)
					assert.Len(t, messages, 2)
					assert.Equal(t, message1.Content, messages[0].Content)
					assert.Equal(t, message2.Content, messages[1].Content)
					assert.Equal(t, message1.ID, messages[0].ID)
					assert.Equal(t, message2.ID, messages[1].ID)
				})
			})
		})
	}
}
