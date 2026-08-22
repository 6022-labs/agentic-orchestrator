package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

type WhenSavingSyncCheckpointTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.SyncCheckpointRepository
}

func WhenSavingSyncCheckpointBeforeEach() *WhenSavingSyncCheckpointTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := registry_common_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := repositories.NewSyncCheckpointRepository(dbClient)

	return &WhenSavingSyncCheckpointTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenSavingSyncCheckpointAfterEach(suite *WhenSavingSyncCheckpointTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingSyncCheckpoint(t *testing.T) {
	t.Parallel()

	t.Run("Given sync checkpoint does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should save and return the checkpoint", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingSyncCheckpointBeforeEach()
			defer WhenSavingSyncCheckpointAfterEach(suite)

			checkpoint := model_builders.NewSyncCheckpointBuilder().WithBlock(42).Build()
			savedCheckpoint, err := suite.sut.Save(context.Background(), checkpoint)

			assert.Nil(t, err)
			assert.NotNil(t, savedCheckpoint)
			assert.Equal(t, uint64(42), savedCheckpoint.Block)
		})
	})
}
