package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/db_clients"
	orchestrator_agent_compute_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenDeletingLLMConnectionByAgentIdTestingSuite struct {
	sut *repositories.LLMConnectionRepository

	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenDeletingLLMConnectionByAgentIdBeforeEach() *WhenDeletingLLMConnectionByAgentIdTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, _ := orchestrator_agent_compute_db.NewOrchestratorDbClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	dbClient.Migrate()

	sut := repositories.NewLLMConnectionRepository(dbClient)

	return &WhenDeletingLLMConnectionByAgentIdTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenDeletingLLMConnectionByAgentIdAfterEach(suite *WhenDeletingLLMConnectionByAgentIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenDeletingLLMConnectionByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingLLMConnectionByAgentIdBeforeEach()
			defer WhenDeletingLLMConnectionByAgentIdAfterEach(suite)

			err := suite.sut.DeleteByAgentId(uuid.New())

			assert.Nil(t, err)
		})
	})

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should delete entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingLLMConnectionByAgentIdBeforeEach()
			defer WhenDeletingLLMConnectionByAgentIdAfterEach(suite)

			existingEntity := model_builders.NewLLMConnectionBuilder().
				Build()

			createdEntity, _ := suite.sut.Save(existingEntity)

			err := suite.sut.DeleteByAgentId(existingEntity.AgentID)

			assert.Nil(t, err)

			fetchedEntity, err := suite.sut.GetByAgentId(createdEntity.AgentID)

			assert.NotNil(t, err)
			assert.Nil(t, fetchedEntity)
		})
	})
}
