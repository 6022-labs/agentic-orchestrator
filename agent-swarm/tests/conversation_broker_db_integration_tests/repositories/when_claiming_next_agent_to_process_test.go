package repositories_test

import (
	"context"

	"sync"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenClaimingNextAgentToProcessTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ParticipantRepository
}

func WhenClaimingNextAgentToProcessBeforeEach(t *testing.T) *WhenClaimingNextAgentToProcessTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenClaimingNextAgentToProcessTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewParticipantRepository(dbClient),
	}
}

func WhenClaimingNextAgentToProcessAfterEach(suite *WhenClaimingNextAgentToProcessTestingSuite) {
	_ = suite.dbClient.Close()
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenClaimingNextAgentToProcess(t *testing.T) {
	t.Parallel()

	t.Run("Given no participant in database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			next, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)

			assert.Nil(t, err)
			assert.Nil(t, next)
		})
	})

	t.Run("Given an agent participant is currently processing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			p := model_builders.NewAgentParticipantBuilder().WithIsProcessing(true).Build()
			_, _ = suite.sut.Save(context.Background(), p)

			next, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)

			assert.Nil(t, err)
			assert.Nil(t, next)
		})
	})

	t.Run("Given an agent participant ran inside the processing interval", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			recent := time.Now().UTC().Add(-5 * time.Second)
			p := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(recent).Build()
			_, _ = suite.sut.Save(context.Background(), p)

			next, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)

			assert.Nil(t, err)
			assert.Nil(t, next)
		})
	})

	t.Run("Given multiple eligible agent participants with non-null last_processed_at", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the one with the oldest last_processed_at", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			older := time.Now().UTC().Add(-40 * time.Second)
			newer := time.Now().UTC().Add(-30 * time.Second)

			pOld := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(older).Build()
			pRecent := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(newer).Build()

			_, _ = suite.sut.Save(context.Background(), pRecent)
			_, _ = suite.sut.Save(context.Background(), pOld)

			next, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)

			assert.NoError(t, err)
			assert.Equal(t, pOld.ID, next.ID)
		})
	})

	t.Run("Given agent participants with NULL and non-NULL last_processed_at", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the one with NULL last_processed_at first", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			older := time.Now().UTC().Add(-40 * time.Second)

			pNeverRun := model_builders.NewAgentParticipantBuilder().Build()
			pOld := model_builders.NewAgentParticipantBuilder().WithLastProcessedAt(older).Build()

			_, _ = suite.sut.Save(context.Background(), pOld)
			_, _ = suite.sut.Save(context.Background(), pNeverRun)

			next, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)

			assert.NoError(t, err)
			assert.Equal(t, pNeverRun.ID, next.ID)
		})
	})

	t.Run("Given only user participants exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			human := model_builders.NewUserParticipantBuilder().Build()
			_, _ = suite.sut.Save(context.Background(), human)

			next, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)

			assert.Nil(t, err)
			assert.Nil(t, next)
		})
	})

	t.Run("Given an eligible agent participant", func(t *testing.T) {
		t.Parallel()

		t.Run("Should mark it processing and not return it again", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			p := model_builders.NewAgentParticipantBuilder().Build()
			_, _ = suite.sut.Save(context.Background(), p)

			claimed, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)
			assert.NoError(t, err)
			assert.Equal(t, p.ID, claimed.ID)
			assert.True(t, claimed.IsProcessing)
			assert.Equal(t, *p.EnsDomain, *claimed.EnsDomain)

			again, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)
			assert.NoError(t, err)
			assert.Nil(t, again)
		})
	})

	t.Run("Given several replicas claiming at once", func(t *testing.T) {
		t.Parallel()

		t.Run("Should give each agent participant to a single one", func(t *testing.T) {
			t.Parallel()

			suite := WhenClaimingNextAgentToProcessBeforeEach(t)
			defer WhenClaimingNextAgentToProcessAfterEach(suite)

			const agents = 3
			for range agents {
				_, _ = suite.sut.Save(context.Background(), model_builders.NewAgentParticipantBuilder().Build())
			}

			var mu sync.Mutex
			claims := map[uuid.UUID]int{}
			wg := &sync.WaitGroup{}
			start := make(chan struct{})
			for range 3 * agents {
				wg.Go(func() {
					<-start
					claimed, err := suite.sut.ClaimNextAgentToProcess(context.Background(), 20)
					assert.NoError(t, err)
					if claimed == nil {
						return
					}
					mu.Lock()
					defer mu.Unlock()
					claims[claimed.ID]++
				})
			}
			close(start)
			wg.Wait()

			for id, count := range claims {
				assert.Equal(t, 1, count, id)
			}
			assert.NotEmpty(t, claims)
		})
	})
}
