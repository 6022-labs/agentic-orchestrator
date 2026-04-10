package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/db_clients"
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

type WhenCheckingMCPConnectionExistenceByAgentIdAndUrlTestingSuite struct {
	sut      *repositories.MCPConnectionRepository
	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenCheckingMCPConnectionExistenceByAgentIdAndUrlBeforeEach() *WhenCheckingMCPConnectionExistenceByAgentIdAndUrlTestingSuite {
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

	sut := repositories.NewMCPConnectionRepository(dbClient)

	return &WhenCheckingMCPConnectionExistenceByAgentIdAndUrlTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenCheckingMCPConnectionExistenceByAgentIdAndUrlAfterEach(suite *WhenCheckingMCPConnectionExistenceByAgentIdAndUrlTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCheckingMCPConnectionExistenceByAgentIdAndUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists for this agent and this url", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return true", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndUrlBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndUrlAfterEach(suite)

			agentID := uuid.New()
			url := "https://example.com/mcp"
			entity := model_builders.NewMCPConnectionBuilder().WithAgentId(agentID).WithUrl(url).Build()
			suite.sut.Save(entity)

			exists, err := suite.sut.ExistsByAgentIdAndUrl(agentID, url)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.True(t, *exists)
		})
	})

	t.Run("Given entity does not exist for this agent and this url", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndUrlBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndUrlAfterEach(suite)

			agentID := uuid.New()
			url := "https://nonexistent.com/mcp"

			exists, err := suite.sut.ExistsByAgentIdAndUrl(agentID, url)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.False(t, *exists)
		})
	})

	t.Run("Given entity exists for this agent but not this url", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndUrlBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndUrlAfterEach(suite)

			agentID := uuid.New()
			url := "https://otherexample.com/mcp"
			otherUrl := "https://example.com/mcp"
			entity := model_builders.NewMCPConnectionBuilder().WithAgentId(agentID).WithUrl(otherUrl).Build()
			suite.sut.Save(entity)

			exists, err := suite.sut.ExistsByAgentIdAndUrl(agentID, url)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.False(t, *exists)
		})
	})

	t.Run("Given entity exists for this url but not this agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndUrlBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndUrlAfterEach(suite)

			agentID := uuid.New()
			otherAgentID := uuid.New()
			url := "https://example.com/mcp"
			entity := model_builders.NewMCPConnectionBuilder().WithAgentId(otherAgentID).WithUrl(url).Build()
			suite.sut.Save(entity)

			exists, err := suite.sut.ExistsByAgentIdAndUrl(agentID, url)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.False(t, *exists)
		})
	})
}
