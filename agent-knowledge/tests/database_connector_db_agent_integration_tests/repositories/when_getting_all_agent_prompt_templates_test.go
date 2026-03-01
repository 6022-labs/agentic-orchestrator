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

type WhenGettingAllPromptTemplatesTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAllPromptTemplatesTestCases = []*WhenGettingAllPromptTemplatesTestCase{
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

type WhenGettingAllPromptTemplatesTestingSuite struct {
	sut *repositories.AgentPromptTemplateRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAllPromptTemplatesBeforeEach(testCase *WhenGettingAllPromptTemplatesTestCase) *WhenGettingAllPromptTemplatesTestingSuite {
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

	sut := repositories.NewAgentPromptTemplateRepository(dbClient)

	return &WhenGettingAllPromptTemplatesTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAllPromptTemplatesAfterEach(suite *WhenGettingAllPromptTemplatesTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAllPromptTemplates(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAllPromptTemplatesTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return empty list", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllPromptTemplatesBeforeEach(testCase)
					defer WhenGettingAllPromptTemplatesAfterEach(suite)

					contextKeys, err := suite.sut.GetAll()

					assert.Nil(t, err)
					assert.Empty(t, contextKeys)
				})
			})

			t.Run("Given entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return list of context keys", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAllPromptTemplatesBeforeEach(testCase)
					defer WhenGettingAllPromptTemplatesAfterEach(suite)

					template1 := model_builders.NewAgentPromptTemplateBuilder().
						WithKey("template1").
						Build()

					template2 := model_builders.NewAgentPromptTemplateBuilder().
						WithKey("template2").
						Build()

					suite.sut.Save(template1)
					suite.sut.Save(template2)

					templates, err := suite.sut.GetAll()

					sort.SliceStable(templates, func(i, j int) bool {
						return templates[i].Key < templates[j].Key
					})

					assert.Nil(t, err)
					assert.Len(t, templates, 2)
					assert.Equal(t, template1.Key, templates[0].Key)
					assert.Equal(t, template2.Key, templates[1].Key)
					assert.Equal(t, template1.ID, templates[0].ID)
					assert.Equal(t, template2.ID, templates[1].ID)
				})
			})
		})
	}
}
