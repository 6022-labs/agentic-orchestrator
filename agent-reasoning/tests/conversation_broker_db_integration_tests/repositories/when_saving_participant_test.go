package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenSavingParticipantTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ParticipantRepository
}

func WhenSavingParticipantBeforeEach() *WhenSavingParticipantTestingSuite {
	dbName, dbClient := newMigratedDatabase()
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

			suite := WhenSavingParticipantBeforeEach()
			defer WhenSavingParticipantAfterEach(suite)

			entity := model_builders.NewAgentParticipantBuilder().Build()
			_, err := suite.sut.Save(entity)
			assert.NoError(t, err)

			byId, err := suite.sut.GetById(entity.ID)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, byId.ID)

			byKey, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(
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

			suite := WhenSavingParticipantBeforeEach()
			defer WhenSavingParticipantAfterEach(suite)

			entity := model_builders.NewUserParticipantBuilder().Build()
			_, err := suite.sut.Save(entity)
			assert.NoError(t, err)

			got, err := suite.sut.GetById(entity.ID)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, got.ID)
			assert.NotNil(t, got.DisplayName)
		})
	})

}
