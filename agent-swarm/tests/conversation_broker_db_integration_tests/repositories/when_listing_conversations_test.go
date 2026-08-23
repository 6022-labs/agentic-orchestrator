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

type WhenListingConversationsTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.ConversationRepository
	bridges  *repositories.BridgeRepository
}

func WhenListingConversationsBeforeEach(t *testing.T) *WhenListingConversationsTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenListingConversationsTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewConversationRepository(dbClient),
		bridges:  repositories.NewBridgeRepository(dbClient),
	}
}

func WhenListingConversationsAfterEach(suite *WhenListingConversationsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenListingConversations(t *testing.T) {
	t.Parallel()

	t.Run("Given no conversations in the database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			defer WhenListingConversationsAfterEach(suite)

			conversations, err := suite.sut.List(context.Background(), 1, 10)

			assert.NoError(t, err)
			assert.Empty(t, conversations)
		})
	})

	t.Run("Given more conversations than one page", func(t *testing.T) {
		t.Parallel()

		t.Run("Should paginate and order by most recently updated", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			defer WhenListingConversationsAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, err := suite.bridges.Save(context.Background(), bridge)
			assert.NoError(t, err)

			for i := 0; i < 5; i++ {
				entity := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
				_, err = suite.sut.Save(context.Background(), entity)
				assert.NoError(t, err)
			}

			// The last saved conversation carries the freshest updated_at.
			latest := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
			_, err = suite.sut.Save(context.Background(), latest)
			assert.NoError(t, err)

			firstPage, err := suite.sut.List(context.Background(), 1, 4)
			assert.NoError(t, err)
			assert.Len(t, firstPage, 4)
			assert.Equal(t, latest.ID, firstPage[0].ID)

			secondPage, err := suite.sut.List(context.Background(), 2, 4)
			assert.NoError(t, err)
			assert.Len(t, secondPage, 2)
		})
	})
}
