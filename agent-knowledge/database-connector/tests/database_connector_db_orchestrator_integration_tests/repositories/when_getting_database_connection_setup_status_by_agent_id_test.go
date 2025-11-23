package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
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

type WhenGettingDatabaseConnectionSetupStatusByAgentIdTestingSuite struct {
	sut *repositories.DatabaseConnectionRepository

	dbName                   string
	postgresContainerBuilder *container_builders.PostgresContainerBuilder
}

func WhenGettingDatabaseConnectionSetupStatusByAgentIdBeforeEach() *WhenGettingDatabaseConnectionSetupStatusByAgentIdTestingSuite {
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

	return &WhenGettingDatabaseConnectionSetupStatusByAgentIdTestingSuite{
		sut: sut,

		dbName:                   dbName,
		postgresContainerBuilder: postgresContainerBuilder,
	}
}

func WhenGettingDatabaseConnectionSetupStatusByAgentIdAfterEach(suite *WhenGettingDatabaseConnectionSetupStatusByAgentIdTestingSuite) {
	suite.postgresContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenGettingDatabaseConnectionSetupStatusByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return setup status", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingDatabaseConnectionSetupStatusByAgentIdBeforeEach()
			defer WhenGettingDatabaseConnectionSetupStatusByAgentIdAfterEach(suite)

			existingEntity := model_builders.NewDatabaseConnectionBuilder().WithSetupStatus(enums.DatabaseConnectionSetupStatusEstablished).Build()
			createdEntity, _ := suite.sut.Save(existingEntity)

			status, err := suite.sut.GetSetupStatusByAgentId(createdEntity.AgentID)

			assert.Nil(t, err)
			assert.NotNil(t, status)
			assert.Equal(t, enums.DatabaseConnectionSetupStatusEstablished, *status)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingDatabaseConnectionSetupStatusByAgentIdBeforeEach()
			defer WhenGettingDatabaseConnectionSetupStatusByAgentIdAfterEach(suite)

			nonExistentAgentID := uuid.New()
			status, err := suite.sut.GetSetupStatusByAgentId(nonExistentAgentID)

			assert.Nil(t, err)
			assert.Nil(t, status)
		})
	})
}
