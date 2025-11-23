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

type WhenDeletingMCPConnectionTestingSuite struct {
	sut      *repositories.MCPConnectionRepository
	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenDeletingMCPConnectionBeforeEach() *WhenDeletingMCPConnectionTestingSuite {
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

	return &WhenDeletingMCPConnectionTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenDeletingMCPConnectionAfterEach(suite *WhenDeletingMCPConnectionTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenDeletingMCPConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingMCPConnectionBeforeEach()
			defer WhenDeletingMCPConnectionAfterEach(suite)

			nonExistentEntity := model_builders.NewMCPConnectionBuilder().WithID(uuid.New()).Build()
			err := suite.sut.Delete(nonExistentEntity)

			assert.Nil(t, err)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should delete entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingMCPConnectionBeforeEach()
			defer WhenDeletingMCPConnectionAfterEach(suite)

			existingEntity := model_builders.NewMCPConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(existingEntity)

			err := suite.sut.Delete(createdEntity)
			assert.Nil(t, err)

			fetchedEntity, err := suite.sut.GetById(createdEntity.ID)
			assert.Nil(t, fetchedEntity)
			assert.NotNil(t, err)
		})
	})
}
