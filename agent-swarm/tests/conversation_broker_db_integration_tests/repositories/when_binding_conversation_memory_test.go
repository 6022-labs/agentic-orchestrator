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

type WhenBindingConversationMemoryTestingSuite struct {
	dbName        string
	dbClient      db_clients.DatabaseClientInterface
	sut           *repositories.ConversationMemoryBindingRepository
	bridges       *repositories.BridgeRepository
	conversations *repositories.ConversationRepository
	participants  *repositories.ParticipantRepository
}

func WhenBindingConversationMemoryBeforeEach(t *testing.T) *WhenBindingConversationMemoryTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)
	return &WhenBindingConversationMemoryTestingSuite{
		dbName:        dbName,
		dbClient:      dbClient,
		sut:           repositories.NewConversationMemoryBindingRepository(dbClient),
		bridges:       repositories.NewBridgeRepository(dbClient),
		conversations: repositories.NewConversationRepository(dbClient),
		participants:  repositories.NewParticipantRepository(dbClient),
	}
}

func WhenBindingConversationMemoryAfterEach(suite *WhenBindingConversationMemoryTestingSuite) {
	// The pool outlives the dropped database otherwise.
	_ = suite.dbClient.Close()
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

// seed creates the rows the binding's foreign keys require.
func (suite *WhenBindingConversationMemoryTestingSuite) seed(
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

func TestWhenBindingConversationMemory(t *testing.T) {
	t.Parallel()

	t.Run("Given the pair is unbound", func(t *testing.T) {
		t.Parallel()

		t.Run("Should persist the session", func(t *testing.T) {
			t.Parallel()

			suite := WhenBindingConversationMemoryBeforeEach(t)
			defer WhenBindingConversationMemoryAfterEach(suite)

			conversation, participant, _ := suite.seed(t)
			sessionId := uuid.New()

			err := suite.sut.Bind(context.Background(), &models.ConversationMemoryBinding{
				ConversationID:  conversation.ID,
				ParticipantID:   participant.ID,
				MemorySessionID: sessionId,
			})

			assert.NoError(t, err)

			binding, err := suite.sut.GetByConversationAndParticipant(
				context.Background(),
				conversation.ID,
				participant.ID,
			)
			assert.NoError(t, err)
			assert.Equal(t, sessionId, binding.MemorySessionID)
			assert.False(t, binding.BoundAt.IsZero())
		})
	})

	t.Run("Given the pair is already bound", func(t *testing.T) {
		t.Parallel()

		t.Run("Should move the binding to the new session", func(t *testing.T) {
			t.Parallel()

			suite := WhenBindingConversationMemoryBeforeEach(t)
			defer WhenBindingConversationMemoryAfterEach(suite)

			conversation, participant, _ := suite.seed(t)
			first := uuid.New()
			second := uuid.New()

			err := suite.sut.Bind(context.Background(), &models.ConversationMemoryBinding{
				ConversationID:  conversation.ID,
				ParticipantID:   participant.ID,
				MemorySessionID: first,
			})
			assert.NoError(t, err)

			err = suite.sut.Bind(context.Background(), &models.ConversationMemoryBinding{
				ConversationID:  conversation.ID,
				ParticipantID:   participant.ID,
				MemorySessionID: second,
			})

			assert.NoError(t, err)

			binding, err := suite.sut.GetByConversationAndParticipant(
				context.Background(),
				conversation.ID,
				participant.ID,
			)
			assert.NoError(t, err)
			// A derived key may move.
			assert.Equal(t, second, binding.MemorySessionID)
		})
	})

	t.Run("Given two participants in the same conversation", func(t *testing.T) {
		t.Parallel()

		t.Run("Should keep each participant's session separate", func(t *testing.T) {
			t.Parallel()

			suite := WhenBindingConversationMemoryBeforeEach(t)
			defer WhenBindingConversationMemoryAfterEach(suite)

			conversation, first, second := suite.seed(t)
			firstSession := uuid.New()
			secondSession := uuid.New()

			err := suite.sut.Bind(context.Background(), &models.ConversationMemoryBinding{
				ConversationID:  conversation.ID,
				ParticipantID:   first.ID,
				MemorySessionID: firstSession,
			})
			assert.NoError(t, err)

			err = suite.sut.Bind(context.Background(), &models.ConversationMemoryBinding{
				ConversationID:  conversation.ID,
				ParticipantID:   second.ID,
				MemorySessionID: secondSession,
			})
			assert.NoError(t, err)

			firstBinding, err := suite.sut.GetByConversationAndParticipant(
				context.Background(),
				conversation.ID,
				first.ID,
			)
			assert.NoError(t, err)
			assert.Equal(t, firstSession, firstBinding.MemorySessionID)

			secondBinding, err := suite.sut.GetByConversationAndParticipant(
				context.Background(),
				conversation.ID,
				second.ID,
			)
			assert.NoError(t, err)
			assert.Equal(t, secondSession, secondBinding.MemorySessionID)
		})
	})
}
