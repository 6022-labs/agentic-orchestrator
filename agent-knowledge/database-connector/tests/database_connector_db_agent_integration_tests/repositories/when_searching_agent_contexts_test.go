package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent"
	db_agent_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenSearchingAgentContextsTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenSearchingAgentContextsTestCases = []*WhenSearchingAgentContextsTestCase{
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

type WhenSearchingAgentContextsTestingSuite struct {
	sut *db_agent_repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenSearchingAgentContextsBeforeEach(testCase *WhenSearchingAgentContextsTestCase) *WhenSearchingAgentContextsTestingSuite {
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

	sut := db_agent_repositories.NewAgentContextRepository(dbClient)

	return &WhenSearchingAgentContextsTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenSearchingAgentContextsAfterEach(suite *WhenSearchingAgentContextsTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenSearchingAgentContexts(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenSearchingAgentContextsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return empty list", func(t *testing.T) {
					t.Parallel()

					suite := WhenSearchingAgentContextsBeforeEach(testCase)
					defer WhenSearchingAgentContextsAfterEach(suite)

					params := repositories.AgentContextSearchParams{}
					results, err := suite.sut.Search(params)

					assert.Nil(t, err)
					assert.Empty(t, results)
				})
			})

			t.Run("Given entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return all entities when no search params", func(t *testing.T) {
					t.Parallel()

					suite := WhenSearchingAgentContextsBeforeEach(testCase)
					defer WhenSearchingAgentContextsAfterEach(suite)

					context1 := model_builders.NewAgentContextBuilder().WithKey("alpha").Build()
					context2 := model_builders.NewAgentContextBuilder().WithKey("beta").Build()
					context3 := model_builders.NewAgentContextBuilder().WithKey("gamma").Build()

					suite.sut.Save(context1)
					suite.sut.Save(context2)
					suite.sut.Save(context3)

					params := repositories.AgentContextSearchParams{}
					results, err := suite.sut.Search(params)

					assert.Nil(t, err)
					assert.Len(t, results, 3)
				})

				t.Run("Should filter by search param", func(t *testing.T) {
					t.Parallel()

					suite := WhenSearchingAgentContextsBeforeEach(testCase)
					defer WhenSearchingAgentContextsAfterEach(suite)

					context1 := model_builders.NewAgentContextBuilder().WithKey("search_me").Build()
					context2 := model_builders.NewAgentContextBuilder().WithKey("not_me").Build()
					suite.sut.Save(context1)
					suite.sut.Save(context2)

					search := "%search%"
					params := repositories.AgentContextSearchParams{
						Search: &search,
					}
					results, err := suite.sut.Search(params)

					assert.Nil(t, err)
					assert.Len(t, results, 1)
					assert.Equal(t, context1.Key, results[0].Key)
				})

				t.Run("Should sort results", func(t *testing.T) {
					t.Parallel()

					suite := WhenSearchingAgentContextsBeforeEach(testCase)
					defer WhenSearchingAgentContextsAfterEach(suite)

					contextA := model_builders.NewAgentContextBuilder().WithKey("A").Build()
					contextB := model_builders.NewAgentContextBuilder().WithKey("B").Build()
					suite.sut.Save(contextB)
					suite.sut.Save(contextA)

					sortBy := "key"
					sortDirection := "asc"
					params := repositories.AgentContextSearchParams{
						SortBy:        &sortBy,
						SortDirection: &sortDirection,
					}
					results, err := suite.sut.Search(params)

					assert.Nil(t, err)
					assert.GreaterOrEqual(t, len(results), 2)
					assert.Equal(t, "A", results[0].Key)
					assert.Equal(t, "B", results[1].Key)
				})
			})
		})
	}
}
