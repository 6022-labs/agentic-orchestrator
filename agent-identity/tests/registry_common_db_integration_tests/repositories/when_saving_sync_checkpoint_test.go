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

type WhenSavingSyncCheckpointTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.SyncCheckpointRepository
}

func WhenSavingSyncCheckpointBeforeEach(t *testing.T) *WhenSavingSyncCheckpointTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewSyncCheckpointRepository(dbClient)

	return &WhenSavingSyncCheckpointTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenSavingSyncCheckpointAfterEach(suite *WhenSavingSyncCheckpointTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenSavingSyncCheckpoint(t *testing.T) {
	t.Parallel()

	t.Run("Given sync checkpoint does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should save and return the checkpoint", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingSyncCheckpointBeforeEach(t)
			defer WhenSavingSyncCheckpointAfterEach(suite)

			checkpoint := model_builders.NewSyncCheckpointBuilder().WithBlock(42).Build()
			savedCheckpoint, err := suite.sut.Save(context.Background(), checkpoint)

			assert.Nil(t, err)
			assert.NotNil(t, savedCheckpoint)
			assert.Equal(t, uint64(42), savedCheckpoint.Block)
		})
	})
}
