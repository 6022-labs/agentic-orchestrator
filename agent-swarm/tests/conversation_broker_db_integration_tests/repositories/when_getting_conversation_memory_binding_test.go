package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenGettingConversationMemoryBindingTestingSuite struct {
	dbName        string
	dbClient      db_clients.DatabaseClientInterface
	sut           *repositories.ConversationMemoryBindingRepository
	bridges       *repositories.BridgeRepository
	conversations *repositories.ConversationRepository
	participants  *repositories.ParticipantRepository
}

func WhenGettingConversationMemoryBindingBeforeEach(t *testing.T) *WhenGettingConversationMemoryBindingTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenGettingConversationMemoryBindingTestingSuite{
		dbName:        dbName,
		dbClient:      dbClient,
		sut:           repositories.NewConversationMemoryBindingRepository(dbClient),
		bridges:       repositories.NewBridgeRepository(dbClient),
		conversations: repositories.NewConversationRepository(dbClient),
		participants:  repositories.NewParticipantRepository(dbClient),
	}
}

func WhenGettingConversationMemoryBindingAfterEach(suite *WhenGettingConversationMemoryBindingTestingSuite) {
	// The pool outlives the dropped database otherwise.
	_ = suite.dbClient.Close()
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

// seed creates the rows the binding's foreign keys require.
func (suite *WhenGettingConversationMemoryBindingTestingSuite) seed(
	t *testing.T,
) (*models.Conversation, *models.Participant, *models.Participant) {
	bridge := model_builders.NewBridgeBuilder().Build()
	_, err := suite.bridges.Save(context.Background(), bridge)
	assert.NoError(t, err)

	conversation := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
	_, err = suite.conversations.Save(context.Background(), conversation)
	assert.NoError(t, err)

	first := model_builders.NewAgentParticipantBuilder().Build()
	_, err = suite.participants.Save(context.Background(), first)
	assert.NoError(t, err)

	second := model_builders.NewAgentParticipantBuilder().Build()
	_, err = suite.participants.Save(context.Background(), second)
	assert.NoError(t, err)

	return conversation, first, second
}

func TestWhenGettingConversationMemoryBinding(t *testing.T) {
	t.Parallel()

	t.Run("Given the pair was never bound", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report no binding", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMemoryBindingBeforeEach(t)
			defer WhenGettingConversationMemoryBindingAfterEach(suite)

			conversation, participant, _ := suite.seed(t)

			binding, err := suite.sut.GetByConversationAndParticipant(
				context.Background(),
				conversation.ID,
				participant.ID,
			)

			assert.NoError(t, err)
			assert.Nil(t, binding)
		})
	})

	t.Run("Given another participant of the conversation is bound", func(t *testing.T) {
		t.Parallel()

		t.Run("Should still report no binding for this one", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMemoryBindingBeforeEach(t)
			defer WhenGettingConversationMemoryBindingAfterEach(suite)

			conversation, bound, unbound := suite.seed(t)
			err := suite.sut.Bind(context.Background(), &models.ConversationMemoryBinding{
				ConversationID:  conversation.ID,
				ParticipantID:   bound.ID,
				MemorySessionID: uuid.New(),
			})
			assert.NoError(t, err)

			binding, err := suite.sut.GetByConversationAndParticipant(
				context.Background(),
				conversation.ID,
				unbound.ID,
			)

			assert.NoError(t, err)
			assert.Nil(t, binding)
		})
	})
}
