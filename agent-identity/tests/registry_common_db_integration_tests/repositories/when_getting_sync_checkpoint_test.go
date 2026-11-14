package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"

	"github.com/stretchr/testify/assert"
)

type WhenGettingSyncCheckpointTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.SyncCheckpointRepository
}

func WhenGettingSyncCheckpointBeforeEach(t *testing.T) *WhenGettingSyncCheckpointTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewSyncCheckpointRepository(dbClient)

	return &WhenGettingSyncCheckpointTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSyncCheckpointAfterEach(suite *WhenGettingSyncCheckpointTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenGettingSyncCheckpoint(t *testing.T) {
	t.Parallel()

	t.Run("Given sync checkpoint exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the checkpoint", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSyncCheckpointBeforeEach(t)
			defer WhenGettingSyncCheckpointAfterEach(suite)

			checkpoint := model_builders.NewSyncCheckpointBuilder().WithBlock(123).Build()
			suite.sut.Save(context.Background(), checkpoint)

			fetchedCheckpoint, err := suite.sut.Get(context.Background())

			assert.Nil(t, err)
			assert.Equal(t, uint64(123), fetchedCheckpoint.Block)
		})
	})

	t.Run("Given sync checkpoint does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSyncCheckpointBeforeEach(t)
			defer WhenGettingSyncCheckpointAfterEach(suite)

			checkpoint, err := suite.sut.Get(context.Background())

			assert.Nil(t, err)
			assert.Nil(t, checkpoint)
		})
	})
}
