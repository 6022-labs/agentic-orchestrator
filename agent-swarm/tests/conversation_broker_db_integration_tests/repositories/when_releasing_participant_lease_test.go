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

type WhenReleasingParticipantLeaseTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ParticipantRepository
}

func WhenReleasingParticipantLeaseBeforeEach(t *testing.T) *WhenReleasingParticipantLeaseTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenReleasingParticipantLeaseTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewParticipantRepository(dbClient),
	}
}

func WhenReleasingParticipantLeaseAfterEach(suite *WhenReleasingParticipantLeaseTestingSuite) {
	_ = suite.dbClient.Close()
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenReleasingParticipantLease(t *testing.T) {
	t.Parallel()

	t.Run("Given the lease it was claimed with", func(t *testing.T) {
		t.Parallel()

		t.Run("Should free the participant and record it as processed", func(t *testing.T) {
			t.Parallel()

			suite := WhenReleasingParticipantLeaseBeforeEach(t)
			defer WhenReleasingParticipantLeaseAfterEach(suite)

			p := model_builders.NewAgentParticipantBuilder().WithLease(uuid.New(), time.Now().Add(time.Minute)).Build()
			_, _ = suite.sut.Save(context.Background(), p)

			assert.NoError(t, suite.sut.ReleaseLease(context.Background(), p))

			got, err := suite.sut.GetById(context.Background(), p.ID)
			assert.NoError(t, err)
			assert.Nil(t, got.LeaseID)
			assert.Nil(t, got.LeaseUntil)
			assert.WithinDuration(t, time.Now(), *got.LastProcessedAt, 10*time.Second)
		})
	})

	t.Run("Given a lease another replica took over", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave the new one alone", func(t *testing.T) {
			t.Parallel()

			suite := WhenReleasingParticipantLeaseBeforeEach(t)
			defer WhenReleasingParticipantLeaseAfterEach(suite)

			stale := model_builders.NewAgentParticipantBuilder().
				WithLease(uuid.New(), time.Now().Add(-time.Second)).
				Build()
			_, _ = suite.sut.Save(context.Background(), stale)
			claimed, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20, time.Minute)
			assert.NoError(t, err)

			assert.NoError(t, suite.sut.ReleaseLease(context.Background(), stale))

			got, err := suite.sut.GetById(context.Background(), stale.ID)
			assert.NoError(t, err)
			assert.Equal(t, *claimed.LeaseID, *got.LeaseID)
			assert.Nil(t, got.LastProcessedAt)
		})
	})
}
