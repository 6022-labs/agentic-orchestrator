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

type WhenCountingConversationsTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ConversationRepository
	bridges  *repositories.BridgeRepository
}

func WhenCountingConversationsBeforeEach(t *testing.T) *WhenCountingConversationsTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenCountingConversationsTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewConversationRepository(dbClient),
		bridges:  repositories.NewBridgeRepository(dbClient),
	}
}

func WhenCountingConversationsAfterEach(suite *WhenCountingConversationsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCountingConversations(t *testing.T) {
	t.Parallel()

	t.Run("Given no conversations in the database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return zero", func(t *testing.T) {
			t.Parallel()

			suite := WhenCountingConversationsBeforeEach(t)
			defer WhenCountingConversationsAfterEach(suite)

			count, err := suite.sut.Count(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, uint(0), count)
		})
	})

	t.Run("Given conversations in the database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the total count", func(t *testing.T) {
			t.Parallel()

			suite := WhenCountingConversationsBeforeEach(t)
			defer WhenCountingConversationsAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, err := suite.bridges.Save(context.Background(), bridge)
			assert.NoError(t, err)

			for i := 0; i < 3; i++ {
				entity := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
				_, err = suite.sut.Save(context.Background(), entity)
				assert.NoError(t, err)
			}

			count, err := suite.sut.Count(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, uint(3), count)
		})
	})
}
