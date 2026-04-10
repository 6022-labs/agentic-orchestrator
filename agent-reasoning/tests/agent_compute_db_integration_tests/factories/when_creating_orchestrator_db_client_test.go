package factories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db/factories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCreatingOrchestratorDbClientTestSuite struct {
	dbName string

	migrationsSettings           *settings.MigrationsSettings
	orchestratorDatabaseSettings *settings.DatabaseSettings
	sut                          *factories.DatabaseClientFactory
}

func WhenCreatingOrchestratorDbClientBeforeEach() *WhenCreatingOrchestratorDbClientTestSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	orchestratorDatabaseSettings := &settings.DatabaseSettings{
		Url: connectionString,
	}
	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: "",
	}

	sut := factories.NewDatabaseClientFactory(
		zap.NewNop(),
		migrationsSettings,
		orchestratorDatabaseSettings,
	)

	return &WhenCreatingOrchestratorDbClientTestSuite{
		sut:    sut,
		dbName: dbName,

		migrationsSettings:           migrationsSettings,
		orchestratorDatabaseSettings: orchestratorDatabaseSettings,
	}
}

func WhenCreatingOrchestratorDbClientAfterEach(suite *WhenCreatingOrchestratorDbClientTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCreatingOrchestratorDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Should return a new orchestrator db client", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingOrchestratorDbClientBeforeEach()
		defer WhenCreatingOrchestratorDbClientAfterEach(suite)

		dbClient, err := suite.sut.CreateDatabaseClient()

		assert.Nil(t, err)
		assert.NotNil(t, dbClient)
	})
}
