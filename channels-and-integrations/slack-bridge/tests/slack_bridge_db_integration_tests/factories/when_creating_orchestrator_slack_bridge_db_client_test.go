package factories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/factories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite struct {
	dbName string

	databaseSettings *settings.DatabaseSettings

	sut *factories.DatabaseClientFactory
}

func WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(
	t *testing.T,
) *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	databaseSettings := &settings.DatabaseSettings{
		Url: connectionString,
	}
	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: "",
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	sut := factories.NewDatabaseClientFactory(
		zap.NewNop(),
		dbQueryMetrics,
		migrationsSettings,
		databaseSettings,
	)

	return &WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite{
		sut:              sut,
		dbName:           dbName,
		databaseSettings: databaseSettings,
	}
}

func WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(
	suite *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite,
) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCreatingOrchestratorSlackBridgeDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Should return a new orchestrator db client", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(t)
		defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

		dbClient, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()

		assert.Nil(t, err)
		assert.NotNil(t, dbClient)
	})
}
