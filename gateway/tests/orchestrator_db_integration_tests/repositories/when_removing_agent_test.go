package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenRemovingAgentTestingSuite struct {
	dbName string

	sut *repositories.AgentRepository
}

func WhenRemovingAgentBeforeEach() *WhenRemovingAgentTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := orchestrator_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenRemovingAgentTestingSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenRemovingAgentAfterEach(suite *WhenRemovingAgentTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenRemovingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentBeforeEach()
			defer WhenRemovingAgentAfterEach(suite)

			nonExistentAgent := model_builders.NewAgentBuilder().Build()

			err := suite.sut.Remove(nonExistentAgent)

			assert.Nil(t, err)
		})
	})

	t.Run("Given agent exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentBeforeEach()
			defer WhenRemovingAgentAfterEach(suite)

			agent := model_builders.NewAgentBuilder().Build()
			_, err := suite.sut.Save(agent)
			if err != nil {
				t.Fatalf("Error saving agent: %v", err)
			}

			err = suite.sut.Remove(agent)

			assert.Nil(t, err)

			fetchedAgent, err := suite.sut.GetById(agent.ID)

			assert.NotNil(t, err)
			assert.Nil(t, fetchedAgent)
		})
	})
}
