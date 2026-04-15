package repositories_test

import (
	"testing"
	"time"

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

	t.Run("Given a minted participant", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be retrievable by id and by collection key", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingParticipantBeforeEach()
			defer WhenSavingParticipantAfterEach(suite)

			entity := model_builders.NewParticipantBuilder().Build()
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

	t.Run("Given a non-minted participant", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be retrievable by id", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingParticipantBeforeEach()
			defer WhenSavingParticipantAfterEach(suite)

			entity := model_builders.NewNonMintedParticipantBuilder().Build()
			_, err := suite.sut.Save(entity)
			assert.NoError(t, err)

			got, err := suite.sut.GetById(entity.ID)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, got.ID)
			assert.NotNil(t, got.DisplayName)
		})
	})

	t.Run("Given multiple minted participants eligible for processing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the one with oldest last_processed_at (nulls first)", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingParticipantBeforeEach()
			defer WhenSavingParticipantAfterEach(suite)

			older := time.Now().Add(-1 * time.Hour)
			newer := time.Now().Add(-10 * time.Minute)

			pNeverRun := model_builders.NewParticipantBuilder().Build()
			pOld := model_builders.NewParticipantBuilder().WithLastProcessedAt(older).Build()
			pRecent := model_builders.NewParticipantBuilder().WithLastProcessedAt(newer).Build()

			_, _ = suite.sut.Save(pOld)
			_, _ = suite.sut.Save(pRecent)
			_, _ = suite.sut.Save(pNeverRun)

			next, err := suite.sut.GetNextMintedParticipantToProcess(60)
			assert.NoError(t, err)
			assert.Equal(t, pNeverRun.ID, next.ID)
		})

		t.Run("Should ignore non-minted humans", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingParticipantBeforeEach()
			defer WhenSavingParticipantAfterEach(suite)

			human := model_builders.NewNonMintedParticipantBuilder().Build()
			_, _ = suite.sut.Save(human)

			_, err := suite.sut.GetNextMintedParticipantToProcess(60)
			assert.Error(t, err) // gorm.ErrRecordNotFound
		})
	})
}
