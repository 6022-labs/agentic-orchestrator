package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/db_clients"
	orchestrator_agent_llm_connector_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

type WhenGettingLLMConnectionByAgentIdTestingSuite struct {
	sut *repositories.LLMConnectionRepository

	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenGettingLLMConnectionByAgentIdBeforeEach() *WhenGettingLLMConnectionByAgentIdTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, _ := orchestrator_agent_llm_connector_db.NewOrchestratorDbClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	dbClient.Migrate()

	sut := repositories.NewLLMConnectionRepository(dbClient)

	return &WhenGettingLLMConnectionByAgentIdTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingLLMConnectionByAgentIdAfterEach(suite *WhenGettingLLMConnectionByAgentIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingLLMConnectionByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionByAgentIdBeforeEach()
			defer WhenGettingLLMConnectionByAgentIdAfterEach(suite)

			existingEntity := model_builders.NewLLMConnectionBuilder().
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

			suite := WhenGettingLLMConnectionByAgentIdBeforeEach()
			defer WhenGettingLLMConnectionByAgentIdAfterEach(suite)

			nonExistentAgentID := model_builders.NewLLMConnectionBuilder().Build().ID
			_, err := suite.sut.GetByAgentId(nonExistentAgentID)

			assert.NotNil(t, err)
		})
	})
}
