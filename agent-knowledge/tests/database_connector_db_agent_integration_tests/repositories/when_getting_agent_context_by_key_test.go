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
	"gorm.io/gorm"
)

type WhenGettingAgentContextByKeyTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAgentContextByKeyTestCases = []*WhenGettingAgentContextByKeyTestCase{
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

type WhenGettingAgentContextByKeyTestingSuite struct {
	sut *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAgentContextByKeyBeforeEach(testCase *WhenGettingAgentContextByKeyTestCase) *WhenGettingAgentContextByKeyTestingSuite {
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

	return &WhenGettingAgentContextByKeyTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAgentContextByKeyAfterEach(suite *WhenGettingAgentContextByKeyTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAgentContextByKey(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAgentContextByKeyTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entity does not exist for this key", func(t *testing.T) {
				t.Parallel()

				t.Run(("Should return not found error"), func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextByKeyBeforeEach(testCase)
					defer WhenGettingAgentContextByKeyAfterEach(suite)

					agentContextKey := uuid.New().String()
					agentContext, err := suite.sut.GetByKey(agentContextKey)

					assert.Nil(t, agentContext)
					assert.NotNil(t, err)
					assert.Equal(t, gorm.ErrRecordNotFound, err)
				})
			})

			t.Run("Given entity exists for this key", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return the entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextByKeyBeforeEach(testCase)
					defer WhenGettingAgentContextByKeyAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					savedAgentContext, _ := suite.sut.Save(agentContext)

					agentContext, err := suite.sut.GetByKey(savedAgentContext.Key)

					assert.Nil(t, err)
					assert.Equal(t, savedAgentContext.ID, agentContext.ID)
					assert.Equal(t, savedAgentContext.Key, agentContext.Key)
				})
			})
		})
	}
}
