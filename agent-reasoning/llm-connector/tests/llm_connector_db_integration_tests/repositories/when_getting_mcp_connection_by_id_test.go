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

type WhenGettingMCPConnectionByIdTestingSuite struct {
	sut *repositories.MCPConnectionRepository

	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenGettingMCPConnectionByIdBeforeEach() *WhenGettingMCPConnectionByIdTestingSuite {
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

	return &WhenGettingMCPConnectionByIdTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingMCPConnectionByIdAfterEach(suite *WhenGettingMCPConnectionByIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingMCPConnectionById(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMCPConnectionByIdBeforeEach()
			defer WhenGettingMCPConnectionByIdAfterEach(suite)

			entity := model_builders.NewMCPConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entity)

			fetchedEntity, err := suite.sut.GetById(createdEntity.ID)

			assert.Nil(t, err)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
			assert.Equal(t, createdEntity.AgentID, fetchedEntity.AgentID)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMCPConnectionByIdBeforeEach()
			defer WhenGettingMCPConnectionByIdAfterEach(suite)

			nonExistentID := uuid.New()
			entity, err := suite.sut.GetById(nonExistentID)

			assert.Nil(t, entity)
			assert.NotNil(t, err)
		})
	})
}
