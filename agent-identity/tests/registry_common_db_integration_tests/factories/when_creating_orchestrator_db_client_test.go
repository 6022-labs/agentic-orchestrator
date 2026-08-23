package factories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/factories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenCreatingOrchestratorDbClientTestingSuite struct {
	dbName string

	migrationsSettings           *settings.MigrationsSettings
	orchestratorDatabaseSettings *settings.DatabaseSettings
	sut                          *factories.DatabaseClientFactory
}

func WhenCreatingOrchestratorDbClientBeforeEach(t *testing.T) *WhenCreatingOrchestratorDbClientTestingSuite {
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
		databaseSettings,
		migrationsSettings,
	)

	return &WhenCreatingOrchestratorDbClientTestingSuite{
		sut:    sut,
		dbName: dbName,

		migrationsSettings:           migrationsSettings,
		orchestratorDatabaseSettings: databaseSettings,
	}
}

func WhenCreatingOrchestratorDbClientAfterEach(suite *WhenCreatingOrchestratorDbClientTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCreatingOrchestratorDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Should return a new orchestrator db client", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingOrchestratorDbClientBeforeEach(t)
		defer WhenCreatingOrchestratorDbClientAfterEach(suite)

		dbClient, err := suite.sut.CreateDatabaseClient()

		assert.Nil(t, err)
		assert.NotNil(t, dbClient)
	})
}
