package factories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/factories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCreatingOrchestratorDbClientTestSuite struct {
	sut *factories.OrchestratorDatabaseClientFactory

	dbName                   string
	postgresContainerBuilder *container_builders.PostgresContainerBuilder

	migrationsSettings           *settings.MigrationsSettings
	orchestratorDatabaseSettings *settings.DatabaseSettings
}

func WhenCreatingOrchestratorDbClientBeforeEach() *WhenCreatingOrchestratorDbClientTestSuite {
	postgresContainerBuilder := container_builders.GetPostgresContainerBuilder()
	dbName, connectionString := postgresContainerBuilder.Create()

	orchestratorDatabaseSettings := &settings.DatabaseSettings{
		Url: connectionString,
	}
	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: "",
	}

	sut := factories.NewOrchestratorDatabaseClientFactory(
		zap.NewNop(),
		migrationsSettings,
		orchestratorDatabaseSettings,
	)

	return &WhenCreatingOrchestratorDbClientTestSuite{
		sut: sut,

		dbName:                   dbName,
		postgresContainerBuilder: postgresContainerBuilder,

		migrationsSettings:           migrationsSettings,
		orchestratorDatabaseSettings: orchestratorDatabaseSettings,
	}
}

func WhenCreatingOrchestratorDbClientAfterEach(suite *WhenCreatingOrchestratorDbClientTestSuite) {
	suite.postgresContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenCreatingOrchestratorDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Should return a new orchestrator db client", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingOrchestratorDbClientBeforeEach()
		defer WhenCreatingOrchestratorDbClientAfterEach(suite)

		dbClient, err := suite.sut.Create()

		assert.Nil(t, err)
		assert.NotNil(t, dbClient)
	})
}
