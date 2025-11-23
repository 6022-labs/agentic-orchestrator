package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/db_clients"
	orchestrator_agent_llm_connector_db "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_db/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingMCPConnectionsByAgentIdTestingSuite struct {
	sut      *repositories.MCPConnectionRepository
	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenGettingMCPConnectionsByAgentIdBeforeEach() *WhenGettingMCPConnectionsByAgentIdTestingSuite {
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

	sut := repositories.NewMCPConnectionRepository(dbClient)

	return &WhenGettingMCPConnectionsByAgentIdTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingMCPConnectionsByAgentIdAfterEach(suite *WhenGettingMCPConnectionsByAgentIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingMCPConnectionsByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given entities exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return all entities for agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMCPConnectionsByAgentIdBeforeEach()
			defer WhenGettingMCPConnectionsByAgentIdAfterEach(suite)

			agentID := uuid.New()
			entity1 := model_builders.NewMCPConnectionBuilder().WithAgentId(agentID).Build()
			entity2 := model_builders.NewMCPConnectionBuilder().WithAgentId(agentID).Build()
			suite.sut.Save(entity1)
			suite.sut.Save(entity2)

			entities, err := suite.sut.GetByAgentId(agentID)

			assert.Nil(t, err)
			assert.Len(t, entities, 2)
			ids := []uuid.UUID{entities[0].ID, entities[1].ID}
			assert.Contains(t, ids, entity1.ID)
			assert.Contains(t, ids, entity2.ID)
		})
	})

	t.Run("Given no entities exist for agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty slice", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMCPConnectionsByAgentIdBeforeEach()
			defer WhenGettingMCPConnectionsByAgentIdAfterEach(suite)

			nonExistentAgentID := uuid.New()
			entities, err := suite.sut.GetByAgentId(nonExistentAgentID)

			assert.Nil(t, err)
			assert.Empty(t, entities)
		})
	})
}
