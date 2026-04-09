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

type WhenSavingAgentContextMessageTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenSavingAgentContextMessageTestCases = []*WhenSavingAgentContextMessageTestCase{
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

type WhenSavingAgentContextMessageTestingSuite struct {
	sut *repositories.AgentContextMessageRepository

	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenSavingAgentContextMessageBeforeEach(testCase *WhenSavingAgentContextMessageTestCase) *WhenSavingAgentContextMessageTestingSuite {
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

	return &WhenSavingAgentContextMessageTestingSuite{
		sut: sut,

		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenSavingAgentContextMessageAfterEach(suite *WhenSavingAgentContextMessageTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenSavingAgentContextMessage(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenSavingAgentContextMessageTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entity does not exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return created entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextMessageBeforeEach(testCase)
					defer WhenSavingAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					entityToCreate := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					createdEntity, err := suite.sut.Save(entityToCreate)

					assert.Nil(t, err)
					assert.Equal(t, entityToCreate.ID, createdEntity.ID)
					assert.Equal(t, entityToCreate.Role, createdEntity.Role)
					assert.Equal(t, entityToCreate.Hash, createdEntity.Hash)
					assert.Equal(t, entityToCreate.Author, createdEntity.Author)
					assert.Equal(t, entityToCreate.AgentContextId, createdEntity.AgentContextId)
				})

				t.Run("Should store entity in database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextMessageBeforeEach(testCase)
					defer WhenSavingAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					entityToCreate := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					createdEntity, _ := suite.sut.Save(entityToCreate)

					suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

					assert.NotEqual(t, uuid.Nil, createdEntity.ID)

					assert.Equal(t, entityToCreate.ID, createdEntity.ID)
					assert.Equal(t, entityToCreate.ID, createdEntity.ID)
					assert.Equal(t, entityToCreate.Role, createdEntity.Role)
					assert.Equal(t, entityToCreate.Hash, createdEntity.Hash)
					assert.Equal(t, entityToCreate.Author, createdEntity.Author)
					assert.Equal(t, entityToCreate.AgentContextId, createdEntity.AgentContextId)
				})
			})

			t.Run("Given entity exists", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return updated entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextMessageBeforeEach(testCase)
					defer WhenSavingAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					entityToCreate := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()
					createdEntity, _ := suite.sut.Save(entityToCreate)

					entityToUpdate := model_builders.NewAgentContextMessageBuilder().
						WithID(createdEntity.ID).
						WithAgentContextId(agentContext.ID).
						Build()
					updatedEntity, err := suite.sut.Save(entityToUpdate)

					assert.Nil(t, err)
					assert.Equal(t, entityToUpdate.ID, updatedEntity.ID)
					assert.Equal(t, entityToCreate.Role, updatedEntity.Role)
					assert.Equal(t, entityToCreate.Hash, updatedEntity.Hash)
					assert.Equal(t, entityToCreate.Author, updatedEntity.Author)
					assert.Equal(t, entityToCreate.AgentContextId, updatedEntity.AgentContextId)
				})

				t.Run("Should update entity in database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingAgentContextMessageBeforeEach(testCase)
					defer WhenSavingAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					entityToCreate := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()
					createdEntity, _ := suite.sut.Save(entityToCreate)

					entityToUpdate := model_builders.NewAgentContextMessageBuilder().
						WithID(createdEntity.ID).
						WithAgentContextId(agentContext.ID).
						Build()
					suite.sut.Save(entityToUpdate)

					suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

					assert.Equal(t, entityToCreate.ID, entityToUpdate.ID)
					assert.Equal(t, entityToUpdate.ID, entityToUpdate.ID)
					assert.Equal(t, entityToCreate.Role, entityToUpdate.Role)
					assert.Equal(t, entityToCreate.Hash, entityToUpdate.Hash)
					assert.Equal(t, entityToCreate.Author, entityToUpdate.Author)
					assert.Equal(t, entityToCreate.AgentContextId, entityToUpdate.AgentContextId)
				})
			})
		})
	}
}
