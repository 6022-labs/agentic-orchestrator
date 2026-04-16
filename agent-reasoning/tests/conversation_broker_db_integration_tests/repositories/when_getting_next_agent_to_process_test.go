package repositories_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type WhenGettingNextAgentToProcessTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ParticipantRepository
}

func WhenGettingNextAgentToProcessBeforeEach() *WhenGettingNextAgentToProcessTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenGettingNextAgentToProcessTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewParticipantRepository(dbClient),
	}
}

func WhenGettingNextAgentToProcessAfterEach(suite *WhenGettingNextAgentToProcessTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingNextAgentToProcess(t *testing.T) {
	t.Parallel()

	t.Run("Given no participant in database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return record-not-found", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			next, err := suite.sut.GetNextAgentToProcess(20)

			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			assert.Nil(t, next)
		})
	})

	t.Run("Given an agent participant is currently processing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			p := model_builders.NewAgentParticipantBuilder().WithIsProcessing(true).Build()
			_, _ = suite.sut.Save(p)

			next, err := suite.sut.GetNextAgentToProcess(20)

			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			assert.Nil(t, next)
		})
	})

	t.Run("Given an agent participant ran inside the processing interval", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			recent := time.Now().UTC().Add(-5 * time.Second)
			p := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(recent).Build()
			_, _ = suite.sut.Save(p)

			next, err := suite.sut.GetNextAgentToProcess(20)

			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			assert.Nil(t, next)
		})
	})

	t.Run("Given multiple eligible agent participants with non-null last_processed_at", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the one with the oldest last_processed_at", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			older := time.Now().UTC().Add(-40 * time.Second)
			newer := time.Now().UTC().Add(-30 * time.Second)

			pOld := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(older).Build()
			pRecent := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(newer).Build()

			_, _ = suite.sut.Save(pRecent)
			_, _ = suite.sut.Save(pOld)

			next, err := suite.sut.GetNextAgentToProcess(20)

			assert.NoError(t, err)
			assert.Equal(t, pOld.ID, next.ID)
		})
	})

	t.Run("Given agent participants with NULL and non-NULL last_processed_at", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the one with NULL last_processed_at first", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			older := time.Now().UTC().Add(-40 * time.Second)

			pNeverRun := model_builders.NewAgentParticipantBuilder().Build()
			pOld := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(older).Build()

			_, _ = suite.sut.Save(pOld)
			_, _ = suite.sut.Save(pNeverRun)

			next, err := suite.sut.GetNextAgentToProcess(20)

			assert.NoError(t, err)
			assert.Equal(t, pNeverRun.ID, next.ID)
		})
	})

	t.Run("Given only user participants exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return record-not-found", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			human := model_builders.NewUserParticipantBuilder().Build()
			_, _ = suite.sut.Save(human)

			next, err := suite.sut.GetNextAgentToProcess(20)

			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			assert.Nil(t, next)
		})
	})
}
