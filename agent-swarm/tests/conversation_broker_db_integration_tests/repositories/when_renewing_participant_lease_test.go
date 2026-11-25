package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenRenewingParticipantLeaseTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ParticipantRepository
}

func WhenRenewingParticipantLeaseBeforeEach(t *testing.T) *WhenRenewingParticipantLeaseTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenRenewingParticipantLeaseTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewParticipantRepository(dbClient),
	}
}

func WhenRenewingParticipantLeaseAfterEach(suite *WhenRenewingParticipantLeaseTestingSuite) {
	_ = suite.dbClient.Close()
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenRenewingParticipantLease(t *testing.T) {
	t.Parallel()

	t.Run("Given the lease it was claimed with", func(t *testing.T) {
		t.Parallel()

		t.Run("Should extend it", func(t *testing.T) {
			t.Parallel()

			suite := WhenRenewingParticipantLeaseBeforeEach(t)
			defer WhenRenewingParticipantLeaseAfterEach(suite)

			p := model_builders.NewAgentParticipantBuilder().WithLease(uuid.New(), time.Now().Add(time.Second)).Build()
			_, _ = suite.sut.Save(context.Background(), p)

			held, err := suite.sut.RenewLease(context.Background(), p, time.Minute)
			assert.NoError(t, err)
			assert.True(t, held)

			got, err := suite.sut.GetById(context.Background(), p.ID)
			assert.NoError(t, err)
			assert.WithinDuration(t, time.Now().Add(time.Minute), *got.LeaseUntil, 10*time.Second)
		})
	})

	t.Run("Given a lease another replica took over", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report it lost and leave the new one alone", func(t *testing.T) {
			t.Parallel()

			suite := WhenRenewingParticipantLeaseBeforeEach(t)
			defer WhenRenewingParticipantLeaseAfterEach(suite)

			stale := model_builders.NewAgentParticipantBuilder().
				WithLease(uuid.New(), time.Now().Add(-time.Second)).
				Build()
			_, _ = suite.sut.Save(context.Background(), stale)
			claimed, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20, time.Minute)
			assert.NoError(t, err)

			held, err := suite.sut.RenewLease(context.Background(), stale, time.Hour)
			assert.NoError(t, err)
			assert.False(t, held)

			got, err := suite.sut.GetById(context.Background(), stale.ID)
			assert.NoError(t, err)
			assert.Equal(t, *claimed.LeaseID, *got.LeaseID)
			assert.True(t, claimed.LeaseUntil.Equal(*got.LeaseUntil))
		})
	})
}
