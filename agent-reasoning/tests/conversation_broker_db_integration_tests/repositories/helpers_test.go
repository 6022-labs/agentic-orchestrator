package repositories_test

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	conversation_broker_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"go.uber.org/zap"
)

// newMigratedDatabase spins up a fresh Postgres DB, runs migrations, and
// returns the dbName for teardown plus the migrated client.
func newMigratedDatabase() (string, db_clients.DatabaseClientInterface) {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindConversationBrokerMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, err := conversation_broker_db.NewConversationBrokerDbClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err != nil {
		panic(err)
	}
	if err := dbClient.Migrate(); err != nil {
		panic(err)
	}

	return dbName, dbClient
}
