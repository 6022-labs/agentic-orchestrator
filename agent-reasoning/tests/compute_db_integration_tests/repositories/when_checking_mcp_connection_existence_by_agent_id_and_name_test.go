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

type WhenCheckingMCPConnectionExistenceByAgentIdAndNameTestingSuite struct {
	sut      *repositories.MCPConnectionRepository
	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenCheckingMCPConnectionExistenceByAgentIdAndNameBeforeEach() *WhenCheckingMCPConnectionExistenceByAgentIdAndNameTestingSuite {
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

	return &WhenCheckingMCPConnectionExistenceByAgentIdAndNameTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenCheckingMCPConnectionExistenceByAgentIdAndNameAfterEach(suite *WhenCheckingMCPConnectionExistenceByAgentIdAndNameTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCheckingMCPConnectionExistenceByAgentIdAndName(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists for agent and name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return true", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndNameBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndNameAfterEach(suite)

			agentID := uuid.New()
			name := "TestConnection"
			entity := model_builders.NewMCPConnectionBuilder().WithAgentId(agentID).WithName(name).Build()
			suite.sut.Save(entity)

			exists, err := suite.sut.ExistsByAgentIdAndName(agentID, name)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.True(t, *exists)
		})
	})

	t.Run("Given entity does not exist for agent and name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndNameBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndNameAfterEach(suite)

			agentID := uuid.New()
			name := "NonExistentConnection"

			exists, err := suite.sut.ExistsByAgentIdAndName(agentID, name)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.False(t, *exists)
		})
	})

	t.Run("Given entity exists for agent but not name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndNameBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndNameAfterEach(suite)

			agentID := uuid.New()
			name := "TestConnection"
			otherName := "OtherConnection"
			entity := model_builders.NewMCPConnectionBuilder().WithAgentId(agentID).WithName(otherName).Build()
			suite.sut.Save(entity)

			exists, err := suite.sut.ExistsByAgentIdAndName(agentID, name)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.False(t, *exists)
		})
	})

	t.Run("Given entity exists for name but not agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingMCPConnectionExistenceByAgentIdAndNameBeforeEach()
			defer WhenCheckingMCPConnectionExistenceByAgentIdAndNameAfterEach(suite)

			agentID := uuid.New()
			otherAgentID := uuid.New()
			name := "TestConnection"
			entity := model_builders.NewMCPConnectionBuilder().WithAgentId(otherAgentID).WithName(name).Build()
			suite.sut.Save(entity)

			exists, err := suite.sut.ExistsByAgentIdAndName(agentID, name)

			assert.Nil(t, err)
			assert.NotNil(t, exists)
			assert.False(t, *exists)
		})
	})
}
