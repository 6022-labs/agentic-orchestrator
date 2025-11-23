package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
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

type WhenCreatingMultipleAgentPromptTemplatesTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenCreatingMultipleAgentPromptTemplatesTestCases = []*WhenCreatingMultipleAgentPromptTemplatesTestCase{
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

type WhenCreatingMultipleAgentPromptTemplatesTestingSuite struct {
	sut *repositories.AgentPromptTemplateRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenCreatingMultipleAgentPromptTemplatesBeforeEach(testCase *WhenCreatingMultipleAgentPromptTemplatesTestCase) *WhenCreatingMultipleAgentPromptTemplatesTestingSuite {
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

	return &WhenCreatingMultipleAgentPromptTemplatesTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenCreatingMultipleAgentPromptTemplatesAfterEach(suite *WhenCreatingMultipleAgentPromptTemplatesTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenCreatingMultipleAgentPromptTemplates(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenCreatingMultipleAgentPromptTemplatesTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should not return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingMultipleAgentPromptTemplatesBeforeEach(testCase)
					defer WhenCreatingMultipleAgentPromptTemplatesAfterEach(suite)

					contextsToCreate := []*models.AgentPromptTemplate{
						model_builders.NewAgentPromptTemplateBuilder().
							WithKey("key1").
							Build(),
						model_builders.NewAgentPromptTemplateBuilder().
							WithKey("key2").
							Build(),
						model_builders.NewAgentPromptTemplateBuilder().
							WithKey("key3").
							Build(),
					}

					_, err := suite.sut.CreateMultiple(contextsToCreate)
					assert.Nil(t, err)
				})
			})

			t.Run("Given entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingMultipleAgentPromptTemplatesBeforeEach(testCase)
					defer WhenCreatingMultipleAgentPromptTemplatesAfterEach(suite)

					contextsToCreate := []*models.AgentPromptTemplate{
						model_builders.NewAgentPromptTemplateBuilder().
							WithKey("key1").
							Build(),
						model_builders.NewAgentPromptTemplateBuilder().
							WithKey("key2").
							Build(),
						model_builders.NewAgentPromptTemplateBuilder().
							WithKey("key3").
							Build(),
					}

					_, err := suite.sut.CreateMultiple(contextsToCreate)
					if err != nil {
						t.Fatal(err)
					}

					_, err = suite.sut.CreateMultiple(contextsToCreate)

					assert.NotNil(t, err)
				})
			})
		})
	}
}
