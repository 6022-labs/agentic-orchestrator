package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenGettingAgentContextByIdTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenGettingAgentContextByIdTestCases = []*WhenGettingAgentContextByIdTestCase{
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

type WhenGettingAgentContextByIdTestingSuite struct {
	sut *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenGettingAgentContextByIdBeforeEach(testCase *WhenGettingAgentContextByIdTestCase) *WhenGettingAgentContextByIdTestingSuite {
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

	return &WhenGettingAgentContextByIdTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenGettingAgentContextByIdAfterEach(suite *WhenGettingAgentContextByIdTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingAgentContextById(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenGettingAgentContextByIdTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given entity does not exist for this id", func(t *testing.T) {
				t.Parallel()

				t.Run(("Should return not found error"), func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextByIdBeforeEach(testCase)
					defer WhenGettingAgentContextByIdAfterEach(suite)

					agentContextId := uuid.New()
					agentContext, err := suite.sut.GetById(agentContextId)

					assert.Nil(t, agentContext)
					assert.NotNil(t, err)
					assert.Equal(t, gorm.ErrRecordNotFound, err)
				})
			})

			t.Run("Given entity exists for this id", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return the entity", func(t *testing.T) {
					t.Parallel()

					suite := WhenGettingAgentContextByIdBeforeEach(testCase)
					defer WhenGettingAgentContextByIdAfterEach(suite)

					agentContext := model_builders.NewAgentContextBuilder().Build()
					savedAgentContext, _ := suite.sut.Save(agentContext)

					agentContext, err := suite.sut.GetById(savedAgentContext.ID)

					assert.Nil(t, err)
					assert.Equal(t, savedAgentContext.ID, agentContext.ID)
					assert.Equal(t, savedAgentContext.Key, agentContext.Key)
				})
			})
		})
	}
}
