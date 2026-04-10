package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	orchestrator_agent_compute_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingLLMConnectionSetupStatusByAgentIdTestingSuite struct {
	sut      *repositories.LLMConnectionRepository
	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenGettingLLMConnectionSetupStatusByAgentIdBeforeEach() *WhenGettingLLMConnectionSetupStatusByAgentIdTestingSuite {
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

	return &WhenGettingLLMConnectionSetupStatusByAgentIdTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingLLMConnectionSetupStatusByAgentIdAfterEach(suite *WhenGettingLLMConnectionSetupStatusByAgentIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingLLMConnectionSetupStatusByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return setup status", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusByAgentIdBeforeEach()
			defer WhenGettingLLMConnectionSetupStatusByAgentIdAfterEach(suite)

			existingEntity := model_builders.NewLLMConnectionBuilder().WithSetupStatus(enums.LLMConnectionSetupStatusEstablished).Build()
			createdEntity, _ := suite.sut.Save(existingEntity)

			status, err := suite.sut.GetSetupStatusByAgentId(createdEntity.AgentID)

			assert.Nil(t, err)
			assert.NotNil(t, status)
			assert.Equal(t, enums.LLMConnectionSetupStatusEstablished, *status)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingLLMConnectionSetupStatusByAgentIdBeforeEach()
			defer WhenGettingLLMConnectionSetupStatusByAgentIdAfterEach(suite)

			nonExistentAgentID := uuid.New()
			status, err := suite.sut.GetSetupStatusByAgentId(nonExistentAgentID)

			assert.Nil(t, err)
			assert.Nil(t, status)
		})
	})
}
