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

type WhenSavingMultipleAgentContextMessageTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenSavingMultipleAgentContextMessageTestCases = []*WhenSavingMultipleAgentContextMessageTestCase{
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

type WhenSavingMultipleAgentContextMessageTestingSuite struct {
	sut                    *repositories.AgentContextMessageRepository
	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenSavingMultipleAgentContextMessageBeforeEach(testCase *WhenSavingMultipleAgentContextMessageTestCase) *WhenSavingMultipleAgentContextMessageTestingSuite {
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

	return &WhenSavingMultipleAgentContextMessageTestingSuite{
		sut:                    sut,
		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenSavingMultipleAgentContextMessageAfterEach(suite *WhenSavingMultipleAgentContextMessageTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenSavingMultipleAgentContextMessage(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenSavingMultipleAgentContextMessageTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entities does not exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should not return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingMultipleAgentContextMessageBeforeEach(testCase)
					defer WhenSavingMultipleAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextMessage1 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					agentContextMessage2 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					err := suite.sut.SaveMultiple([]*models.AgentContextMessage{agentContextMessage1, agentContextMessage2})

					assert.Nil(t, err)
				})

				t.Run("Should store entity in database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingMultipleAgentContextMessageBeforeEach(testCase)
					defer WhenSavingMultipleAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextMessage1 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					agentContextMessage2 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					suite.sut.SaveMultiple([]*models.AgentContextMessage{agentContextMessage1, agentContextMessage2})

					var createdEntities []*models.AgentContextMessage
					suite.dbClient.GetDB().Find(&createdEntities)

					assert.Len(t, createdEntities, 2)

					for _, entity := range createdEntities {
						containId := entity.ID == agentContextMessage1.ID || entity.ID == agentContextMessage2.ID
						assert.True(t, containId)
					}
				})
			})

			t.Run("Given entities already exists", func(t *testing.T) {
				t.Parallel()

				t.Run("Should not return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingMultipleAgentContextMessageBeforeEach(testCase)
					defer WhenSavingMultipleAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextMessage1 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					agentContextMessage2 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					suite.sut.SaveMultiple([]*models.AgentContextMessage{agentContextMessage1, agentContextMessage2})

					agentContextMessage1.Content = "Updated message 1"
					agentContextMessage2.Content = "Updated message 2"

					err := suite.sut.SaveMultiple([]*models.AgentContextMessage{agentContextMessage1, agentContextMessage2})

					assert.Nil(t, err)
				})

				t.Run("Should update entities in database", func(t *testing.T) {
					t.Parallel()

					suite := WhenSavingMultipleAgentContextMessageBeforeEach(testCase)
					defer WhenSavingMultipleAgentContextMessageAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextMessage1 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					agentContextMessage2 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					suite.sut.SaveMultiple([]*models.AgentContextMessage{agentContextMessage1, agentContextMessage2})

					agentContextMessage1.Content = "Updated message 1"
					agentContextMessage2.Content = "Updated message 2"

					suite.sut.SaveMultiple([]*models.AgentContextMessage{agentContextMessage1, agentContextMessage2})

					var createdEntities []*models.AgentContextMessage
					suite.dbClient.GetDB().Find(&createdEntities)

					for _, entity := range createdEntities {
						containId := entity.ID == agentContextMessage1.ID || entity.ID == agentContextMessage2.ID
						assert.True(t, containId)

						containContent := entity.Content == agentContextMessage1.Content || entity.Content == agentContextMessage2.Content
						assert.True(t, containContent)
					}
				})
			})
		})
	}
}
