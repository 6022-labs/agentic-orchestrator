package factories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/factories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite struct {
	dbName string

	databaseSettings *settings.DatabaseSettings

	sut *factories.DatabaseClientFactory
}

func WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach() *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	databaseSettings := &settings.DatabaseSettings{
		Url: connectionString,
	}
	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: "",
	}

	sut := factories.NewDatabaseClientFactory(
		zap.NewNop(),
		migrationsSettings,
		databaseSettings,
	)

	return &WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite{
		sut:              sut,
		dbName:           dbName,
		databaseSettings: databaseSettings,
	}
}

func WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCreatingOrchestratorSlackBridgeDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Should return a new orchestrator db client", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach()
		defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

		dbClient, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()

		assert.Nil(t, err)
		assert.NotNil(t, dbClient)
	})
}
