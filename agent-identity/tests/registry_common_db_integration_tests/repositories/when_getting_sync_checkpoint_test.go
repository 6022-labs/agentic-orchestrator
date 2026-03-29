package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

type WhenGettingSyncCheckpointTestSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.SyncCheckpointRepository
}

func WhenGettingSyncCheckpointBeforeEach() *WhenGettingSyncCheckpointTestSuite {
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

	return &WhenGettingSyncCheckpointTestSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSyncCheckpointAfterEach(suite *WhenGettingSyncCheckpointTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingSyncCheckpoint(t *testing.T) {
	t.Parallel()

	t.Run("Given sync checkpoint exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the checkpoint", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSyncCheckpointBeforeEach()
			defer WhenGettingSyncCheckpointAfterEach(suite)

			checkpoint := &models.SyncCheckpoint{Block: 123}
			suite.sut.Save(checkpoint)

			fetchedCheckpoint, err := suite.sut.Get()

			assert.Nil(t, err)
			assert.Equal(t, uint64(123), fetchedCheckpoint.Block)
		})
	})

	t.Run("Given sync checkpoint does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return record not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSyncCheckpointBeforeEach()
			defer WhenGettingSyncCheckpointAfterEach(suite)

			_, err := suite.sut.Get()

			assert.NotNil(t, err)
			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		})
	})
}
