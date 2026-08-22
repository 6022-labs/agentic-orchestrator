package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenListingParticipantsByConversationTestingSuite struct {
	dbName                   string
	dbClient                 db_clients.DatabaseClientInterface
	sut                      *repositories.ParticipantRepository
	bridges                  *repositories.BridgeRepository
	conversations            *repositories.ConversationRepository
	conversationParticipants *repositories.ConversationParticipantRepository
}

func WhenListingParticipantsByConversationBeforeEach() *WhenListingParticipantsByConversationTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenListingParticipantsByConversationTestingSuite{
		dbName:                   dbName,
		dbClient:                 dbClient,
		sut:                      repositories.NewParticipantRepository(dbClient),
		bridges:                  repositories.NewBridgeRepository(dbClient),
		conversations:            repositories.NewConversationRepository(dbClient),
		conversationParticipants: repositories.NewConversationParticipantRepository(dbClient),
	}
}

func WhenListingParticipantsByConversationAfterEach(suite *WhenListingParticipantsByConversationTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenListingParticipantsByConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given a conversation has multiple participants attached", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return all of them and exclude participants from other conversations", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingParticipantsByConversationBeforeEach()
			defer WhenListingParticipantsByConversationAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, _ = suite.bridges.Save(context.Background(), bridge)

			conv := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
			otherConv := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
			_, _ = suite.conversations.Save(context.Background(), conv)
			_, _ = suite.conversations.Save(context.Background(), otherConv)

			agent := model_builders.NewAgentParticipantBuilder().Build()
			human := model_builders.NewUserParticipantBuilder().Build()
			outsider := model_builders.NewAgentParticipantBuilder().Build()
			_, _ = suite.sut.Save(context.Background(), agent)
			_, _ = suite.sut.Save(context.Background(), human)
			_, _ = suite.sut.Save(context.Background(), outsider)

			_, _ = suite.conversationParticipants.Save(
				context.Background(),
				model_builders.NewConversationParticipantBuilder().
					WithConversationID(conv.ID).WithParticipantID(agent.ID).Build(),
			)
			_, _ = suite.conversationParticipants.Save(
				context.Background(),
				model_builders.NewConversationParticipantBuilder().
					WithConversationID(conv.ID).WithParticipantID(human.ID).Build(),
			)
			_, _ = suite.conversationParticipants.Save(
				context.Background(),
				model_builders.NewConversationParticipantBuilder().
					WithConversationID(otherConv.ID).WithParticipantID(outsider.ID).Build(),
			)

			result, err := suite.sut.ListByConversationId(context.Background(), conv.ID)

			assert.NoError(t, err)
			assert.Len(t, result, 2)

			ids := []uuid.UUID{result[0].ID, result[1].ID}
			assert.Contains(t, ids, agent.ID)
			assert.Contains(t, ids, human.ID)
			assert.NotContains(t, ids, outsider.ID)
		})
	})

	t.Run("Given a conversation has no participants", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingParticipantsByConversationBeforeEach()
			defer WhenListingParticipantsByConversationAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, _ = suite.bridges.Save(context.Background(), bridge)
			conv := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
			_, _ = suite.conversations.Save(context.Background(), conv)

			result, err := suite.sut.ListByConversationId(context.Background(), conv.ID)

			assert.NoError(t, err)
			assert.Len(t, result, 0)
		})
	})
}
