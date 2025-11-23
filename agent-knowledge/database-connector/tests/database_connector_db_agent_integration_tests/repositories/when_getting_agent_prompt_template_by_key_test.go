package repositories_test

import (
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
	"gorm.io/gorm"
)

type WhenGettingAgentPromptTemplateByKeyTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAgentPromptTemplateByKeyTestCases = []*WhenGettingAgentPromptTemplateByKeyTestCase{
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

type WhenGettingAgentPromptTemplateByKeyTestingSuite struct {
	sut *repositories.AgentPromptTemplateRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAgentPromptTemplateByKeyBeforeEach(testCase *WhenGettingAgentPromptTemplateByKeyTestCase) *WhenGettingAgentPromptTemplateByKeyTestingSuite {
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

	return &WhenGettingAgentPromptTemplateByKeyTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAgentPromptTemplateByKeyAfterEach(suite *WhenGettingAgentPromptTemplateByKeyTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAgentPromptTemplateByKey(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAgentPromptTemplateByKeyTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entity exists with the given key", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return not found error", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentPromptTemplateByKeyBeforeEach(testCase)
					defer WhenGettingAgentPromptTemplateByKeyAfterEach(suite)

					_, err := suite.sut.GetByKey("non-existing-key")

					assert.Equal(t, gorm.ErrRecordNotFound, err)
				})
			})

			t.Run("Given entity exists with the given key", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return the entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentPromptTemplateByKeyBeforeEach(testCase)
					defer WhenGettingAgentPromptTemplateByKeyAfterEach(suite)

					entityToCreate := model_builders.NewAgentPromptTemplateBuilder().Build()
					createdEntity, _ := suite.sut.Save(entityToCreate)

					retrievedEntity, err := suite.sut.GetByKey(createdEntity.Key)

					assert.Nil(t, err)
					assert.Equal(t, createdEntity.ID, retrievedEntity.ID)
				})
			})
		})
	}
}
