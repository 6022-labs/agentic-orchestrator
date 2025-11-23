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

type WhenGettingAllContextsTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAllContextsTestCases = []*WhenGettingAllContextsTestCase{
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

type WhenGettingAllContextsTestingSuite struct {
	sut *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAllContextsBeforeEach(testCase *WhenGettingAllContextsTestCase) *WhenGettingAllContextsTestingSuite {
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

	sut := repositories.NewAgentContextRepository(dbClient)

	return &WhenGettingAllContextsTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAllContextsAfterEach(suite *WhenGettingAllContextsTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAllContexts(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAllContextsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return empty list", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllContextsBeforeEach(testCase)
					defer WhenGettingAllContextsAfterEach(suite)

					contextKeys, err := suite.sut.GetAll()

					assert.Nil(t, err)
					assert.Empty(t, contextKeys)
				})
			})

			t.Run("Given entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return list of context keys", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllContextsBeforeEach(testCase)
					defer WhenGettingAllContextsAfterEach(suite)

					context1 := model_builders.NewAgentContextBuilder().
						WithKey("key1").
						Build()

					context2 := model_builders.NewAgentContextBuilder().
						WithKey("key2").
						Build()

					suite.sut.Save(context1)
					suite.sut.Save(context2)

					contexts, err := suite.sut.GetAll()

					sort.SliceStable(contexts, func(i, j int) bool {
						return contexts[i].Key < contexts[j].Key
					})

					assert.Nil(t, err)
					assert.Len(t, contexts, 2)
					assert.Equal(t, context1.Key, contexts[0].Key)
					assert.Equal(t, context2.Key, contexts[1].Key)
					assert.Equal(t, context1.ID, contexts[0].ID)
					assert.Equal(t, context2.ID, contexts[1].ID)
				})
			})
		})
	}
}
