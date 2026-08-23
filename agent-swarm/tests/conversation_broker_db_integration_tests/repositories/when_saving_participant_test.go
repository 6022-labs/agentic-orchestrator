package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenSavingParticipantTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ParticipantRepository
}

func WhenSavingParticipantBeforeEach(t *testing.T) *WhenSavingParticipantTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenSavingParticipantTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewParticipantRepository(dbClient),
	}
}

func WhenSavingParticipantAfterEach(suite *WhenSavingParticipantTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given an agent participant", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be retrievable by id and by collection key", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingParticipantBeforeEach(t)
			defer WhenSavingParticipantAfterEach(suite)

			entity := model_builders.NewAgentParticipantBuilder().Build()
			_, err := suite.sut.Save(context.Background(), entity)
			assert.NoError(t, err)

			byId, err := suite.sut.GetById(context.Background(), entity.ID)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, byId.ID)

			byKey, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(context.Background(),
				*entity.CollectionAddress, *entity.CollectionTokenId,
			)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, byKey.ID)
		})
	})

	t.Run("Given a user participant", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be retrievable by id", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingParticipantBeforeEach(t)
			defer WhenSavingParticipantAfterEach(suite)

			entity := model_builders.NewUserParticipantBuilder().Build()
			_, err := suite.sut.Save(context.Background(), entity)
			assert.NoError(t, err)

			got, err := suite.sut.GetById(context.Background(), entity.ID)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, got.ID)
			assert.NotNil(t, got.DisplayName)
		})
	})

}
