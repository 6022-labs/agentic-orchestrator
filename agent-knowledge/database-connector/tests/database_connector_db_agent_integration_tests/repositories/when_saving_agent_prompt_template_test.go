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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenSavingAgentPromptTemplateTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenSavingAgentPromptTemplateTestCases = []*WhenSavingAgentPromptTemplateTestCase{
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

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *repositories.AgentPromptTemplateRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(testCase *WhenSavingAgentPromptTemplateTestCase) *WhenSavingAgentPromptTemplateTestingSuite {
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

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenSavingAgentPromptTemplateAfterEach(suite *WhenSavingAgentPromptTemplateTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenSavingAgentPromptTemplateTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entity does not exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return created entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentPromptTemplateBeforeEach(testCase)
					defer WhenSavingAgentPromptTemplateAfterEach(suite)

					entityToCreate := model_builders.NewAgentPromptTemplateBuilder().Build()
					createdEntity, err := suite.sut.Save(entityToCreate)

					assert.Nil(t, err)
					assert.Equal(t, entityToCreate.ID, createdEntity.ID)
					assert.Equal(t, entityToCreate.Key, createdEntity.Key)
				})

				t.Run("Should store entity in database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentPromptTemplateBeforeEach(testCase)
					defer WhenSavingAgentPromptTemplateAfterEach(suite)

					entityToCreate := model_builders.NewAgentPromptTemplateBuilder().Build()
					createdEntity, _ := suite.sut.Save(entityToCreate)

					suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

					assert.NotEqual(t, uuid.Nil, createdEntity.ID)

					assert.Equal(t, entityToCreate.ID, createdEntity.ID)
					assert.Equal(t, entityToCreate.Key, createdEntity.Key)
				})
			})

			t.Run("Given entity exists", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return updated entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentPromptTemplateBeforeEach(testCase)
					defer WhenSavingAgentPromptTemplateAfterEach(suite)

					entityToCreate := model_builders.NewAgentPromptTemplateBuilder().Build()
					createdEntity, _ := suite.sut.Save(entityToCreate)

					entityToUpdate := model_builders.NewAgentPromptTemplateBuilder().WithID(createdEntity.ID).Build()
					updatedEntity, err := suite.sut.Save(entityToUpdate)

					assert.Nil(t, err)
					assert.Equal(t, entityToUpdate.ID, updatedEntity.ID)
					assert.Equal(t, entityToUpdate.Key, updatedEntity.Key)
				})

				t.Run("Should update entity in database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentPromptTemplateBeforeEach(testCase)
					defer WhenSavingAgentPromptTemplateAfterEach(suite)

					entityToCreate := model_builders.NewAgentPromptTemplateBuilder().Build()
					createdEntity, _ := suite.sut.Save(entityToCreate)

					entityToUpdate := model_builders.NewAgentPromptTemplateBuilder().WithID(createdEntity.ID).Build()
					suite.sut.Save(entityToUpdate)

					suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

					assert.Equal(t, entityToCreate.ID, entityToUpdate.ID)
					assert.Equal(t, entityToCreate.Key, entityToUpdate.Key)
				})
			})
		})
	}
}
