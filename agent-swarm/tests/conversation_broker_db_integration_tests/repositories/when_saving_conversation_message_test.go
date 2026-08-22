package repositories_test

import (
	"context"

	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenSavingConversationMessageTestingSuite struct {
	dbName        string
	dbClient      db_clients.DatabaseClientInterface
	sut           *repositories.ConversationMessageRepository
	bridges       *repositories.BridgeRepository
	conversations *repositories.ConversationRepository
	participants  *repositories.ParticipantRepository
}

func WhenSavingConversationMessageBeforeEach() *WhenSavingConversationMessageTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenSavingConversationMessageTestingSuite{
		dbName:        dbName,
		dbClient:      dbClient,
		sut:           repositories.NewConversationMessageRepository(dbClient),
		bridges:       repositories.NewBridgeRepository(dbClient),
		conversations: repositories.NewConversationRepository(dbClient),
		participants:  repositories.NewParticipantRepository(dbClient),
	}
}

func WhenSavingConversationMessageAfterEach(suite *WhenSavingConversationMessageTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingConversationMessage(t *testing.T) {
	t.Parallel()

	t.Run("Given messages in the same conversation", func(t *testing.T) {
		t.Parallel()

		t.Run("Should list them ordered by created_at ASC", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingConversationMessageBeforeEach()
			defer WhenSavingConversationMessageAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, _ = suite.bridges.Save(context.Background(), bridge)
			conv := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
			_, _ = suite.conversations.Save(context.Background(), conv)
			part := model_builders.NewAgentParticipantBuilder().Build()
			_, _ = suite.participants.Save(context.Background(), part)

			now := time.Now()
			older := model_builders.NewConversationMessageBuilder().
				WithConversationID(conv.ID).
				WithParticipantID(part.ID).
				WithContent("first").
				WithCreatedAt(now.Add(-2 * time.Minute)).
				Build()
			newer := model_builders.NewConversationMessageBuilder().
				WithConversationID(conv.ID).
				WithParticipantID(part.ID).
				WithContent("second").
				WithCreatedAt(now).
				Build()

			// Save newer first to make sure ordering is not insertion-order.
			_, err := suite.sut.Save(context.Background(), newer)
			assert.NoError(t, err)
			_, err = suite.sut.Save(context.Background(), older)
			assert.NoError(t, err)

			list, err := suite.sut.ListByConversationId(context.Background(), conv.ID)
			assert.NoError(t, err)
			assert.Len(t, list, 2)
			assert.Equal(t, "first", list[0].Content)
			assert.Equal(t, "second", list[1].Content)
		})
	})
}
