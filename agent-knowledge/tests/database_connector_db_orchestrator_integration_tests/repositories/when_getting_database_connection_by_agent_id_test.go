package repositories_test

import (
	"testing"

	orchestrator_agent_database_connector_db_orchestrator "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_orchestrator"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_orchestrator/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingDatabaseConnectionByAgentIdTestSuite struct {
	sut *repositories.DatabaseConnectionRepository

	dbName                   string
	postgresContainerBuilder *container_builders.PostgresContainerBuilder
}

func WhenGettingDatabaseConnectionByAgentIdBeforeEach() *WhenGettingDatabaseConnectionByAgentIdTestSuite {
	postgresContainerBuilder := container_builders.GetPostgresContainerBuilder()
	dbName, connectionString := postgresContainerBuilder.Create()

	migrationFolderPath, err := common_tests.FindOrchestratorMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, err := orchestrator_agent_database_connector_db_orchestrator.NewOrchestratorDatabaseClient(
		connectionString,
		zap.NewNop(),
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

	sut := repositories.NewDatabaseConnectionRepository(dbClient)

	return &WhenGettingDatabaseConnectionByAgentIdTestSuite{
		sut: sut,

		dbName:                   dbName,
		postgresContainerBuilder: postgresContainerBuilder,
	}
}

func WhenGettingDatabaseConnectionByAgentIdAfterEach(suite *WhenGettingDatabaseConnectionByAgentIdTestSuite) {
	suite.postgresContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingDatabaseConnectionByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingDatabaseConnectionByAgentIdBeforeEach()
			defer WhenGettingDatabaseConnectionByAgentIdAfterEach(suite)

			existingEntity := model_builders.NewDatabaseConnectionBuilder().
				Build()

			createdEntity, _ := suite.sut.Save(existingEntity)

			fetchedEntity, err := suite.sut.GetByAgentId(existingEntity.AgentID)

			assert.Nil(t, err)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingDatabaseConnectionByAgentIdBeforeEach()
			defer WhenGettingDatabaseConnectionByAgentIdAfterEach(suite)

			_, err := suite.sut.GetByAgentId(uuid.New())

			assert.NotNil(t, err)
		})
	})
}
