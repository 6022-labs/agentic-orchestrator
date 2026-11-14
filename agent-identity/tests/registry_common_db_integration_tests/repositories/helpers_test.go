package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_db_mocks/metrics_mocks/interfaces_mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// newMigratedDatabase creates a fresh migrated database and returns its name for teardown plus its client.
func newMigratedDatabase(t *testing.T) (string, db_clients.DatabaseClientInterface) {
	dbName, connectionString := container_builders.GetPostgresContainerBuilder().Create()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, err := registry_common_db.NewDatabaseClient(
		&settings.DatabaseSettings{Url: connectionString},
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
		dbQueryMetrics,
	)
	if err != nil {
		panic(err)
	}
	if err := dbClient.Migrate(); err != nil {
		panic(err)
	}

	return dbName, dbClient
}
