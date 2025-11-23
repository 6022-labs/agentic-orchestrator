package repositories

import (
	"testing"

	orchestrator_agent_database_connector_db_orchestrator "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenDeletingDatabaseConnectionByAgentIdTestingSuite struct {
	sut *repositories.DatabaseConnectionRepository

	dbName                   string
	postgresContainerBuilder *container_builders.PostgresContainerBuilder
}

func WhenDeletingDatabaseConnectionByAgentIdBeforeEach() *WhenDeletingDatabaseConnectionByAgentIdTestingSuite {
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

	return &WhenDeletingDatabaseConnectionByAgentIdTestingSuite{
		sut: sut,

		dbName:                   dbName,
		postgresContainerBuilder: postgresContainerBuilder,
	}
}

func WhenDeletingDatabaseConnectionByAgentIdAfterEach(suite *WhenDeletingDatabaseConnectionByAgentIdTestingSuite) {
	suite.postgresContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenDeletingDatabaseConnectionByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingDatabaseConnectionByAgentIdBeforeEach()
			defer WhenDeletingDatabaseConnectionByAgentIdAfterEach(suite)

			err := suite.sut.DeleteByAgentId(uuid.New())
			assert.NoError(t, err)
		})
	})

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should delete entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingDatabaseConnectionByAgentIdBeforeEach()
			defer WhenDeletingDatabaseConnectionByAgentIdAfterEach(suite)

			existingEntity := model_builders.NewDatabaseConnectionBuilder().
				Build()

			_, _ = suite.sut.Save(existingEntity)

			err := suite.sut.DeleteByAgentId(existingEntity.AgentID)
			assert.NoError(t, err)

			_, err = suite.sut.GetByAgentId(existingEntity.AgentID)
			assert.Error(t, err)
		})
	})
}
