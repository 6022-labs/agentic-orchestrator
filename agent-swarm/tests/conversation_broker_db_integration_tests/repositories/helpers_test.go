package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	conversation_broker_db "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_db_mocks/metrics_mocks/interfaces_mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// newMigratedDatabase spins up a fresh Postgres DB, runs migrations, and
// returns the dbName for teardown plus the migrated client.
func newMigratedDatabase(t *testing.T) (string, db_clients.DatabaseClientInterface) {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindConversationBrokerMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, err := conversation_broker_db.NewConversationBrokerDbClient(
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
