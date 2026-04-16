package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenSavingConversationParticipantTestingSuite struct {
	dbName       string
	dbClient     db_clients.DatabaseClientInterface
	sut          *repositories.ConversationParticipantRepository
	bridges      *repositories.BridgeRepository
	conversations *repositories.ConversationRepository
	participants *repositories.ParticipantRepository
}

func WhenSavingConversationParticipantBeforeEach() *WhenSavingConversationParticipantTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenSavingConversationParticipantTestingSuite{
		dbName:        dbName,
		dbClient:      dbClient,
		sut:           repositories.NewConversationParticipantRepository(dbClient),
		bridges:       repositories.NewBridgeRepository(dbClient),
		conversations: repositories.NewConversationRepository(dbClient),
		participants:  repositories.NewParticipantRepository(dbClient),
	}
}

func WhenSavingConversationParticipantAfterEach(suite *WhenSavingConversationParticipantTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingConversationParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given a conversation and a participant exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should attach, list and detach the participant", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingConversationParticipantBeforeEach()
			defer WhenSavingConversationParticipantAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, _ = suite.bridges.Save(bridge)

			conv := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
			_, _ = suite.conversations.Save(conv)

			part := model_builders.NewAgentParticipantBuilder().Build()
			_, _ = suite.participants.Save(part)

			link := model_builders.NewConversationParticipantBuilder().
				WithConversationID(conv.ID).
				WithParticipantID(part.ID).
				Build()

			_, err := suite.sut.Save(link)
			assert.NoError(t, err)

			list, err := suite.sut.ListByConversationId(conv.ID)
			assert.NoError(t, err)
			assert.Len(t, list, 1)
			assert.Equal(t, part.ID, list[0].ParticipantID)

			err = suite.sut.DeleteByConversationIdAndParticipantId(conv.ID, part.ID)
			assert.NoError(t, err)

			list, err = suite.sut.ListByConversationId(conv.ID)
			assert.NoError(t, err)
			assert.Len(t, list, 0)
		})
	})
}
