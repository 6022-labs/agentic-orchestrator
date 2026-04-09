package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
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

type WhenGettingAgentContextMessageByAgentContextKeyTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAgentContextMessageByAgentContextKeyTestCases = []*WhenGettingAgentContextMessageByAgentContextKeyTestCase{
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

type WhenGettingAgentContextMessageByAgentContextKeyTestingSuite struct {
	sut                    *repositories.AgentContextMessageRepository
	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAgentContextMessageByAgentContextKeyBeforeEach(testCase *WhenGettingAgentContextMessageByAgentContextKeyTestCase) *WhenGettingAgentContextMessageByAgentContextKeyTestingSuite {
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

	return &WhenGettingAgentContextMessageByAgentContextKeyTestingSuite{
		sut:                    sut,
		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAgentContextMessageByAgentContextKeyAfterEach(suite *WhenGettingAgentContextMessageByAgentContextKeyTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAgentContextMessageByAgentContextKey(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAgentContextMessageByAgentContextKeyTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entities does not exist for the agent context key", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return empty list", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextMessageByAgentContextKeyBeforeEach(testCase)
					defer WhenGettingAgentContextMessageByAgentContextKeyAfterEach(suite)

					agentContextKey := uuid.New().String()
					agentContextMessages, err := suite.sut.GetByAgentContextKey(agentContextKey)

					assert.Nil(t, err)
					assert.Empty(t, agentContextMessages)
				})
			})

			t.Run("Given entities exist for the agent context key", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return list of agent context messages", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextMessageByAgentContextKeyBeforeEach(testCase)
					defer WhenGettingAgentContextMessageByAgentContextKeyAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					_, _ = suite.agentContextRepository.Save(agentContext)

					agentContextMessage1 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()
					agentContextMessage2 := model_builders.NewAgentContextMessageBuilder().
						WithAgentContextId(agentContext.ID).
						Build()

					suite.sut.SaveMultiple([]*models.AgentContextMessage{agentContextMessage1, agentContextMessage2})

					agentContextMessages, err := suite.sut.GetByAgentContextKey(agentContext.Key)

					assert.Nil(t, err)
					assert.Len(t, agentContextMessages, 2)
				})
			})
		})
	}
}
