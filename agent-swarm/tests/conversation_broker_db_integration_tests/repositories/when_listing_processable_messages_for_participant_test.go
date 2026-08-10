package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type WhenListingProcessableMessagesForParticipantTestingSuite struct {
	dbName        string
	dbClient      db_clients.DatabaseClientInterface
	sut           *repositories.ConversationMessageRepository
	bridges       *repositories.BridgeRepository
	conversations *repositories.ConversationRepository
	participants  *repositories.ParticipantRepository
}

func WhenListingProcessableMessagesForParticipantBeforeEach() *WhenListingProcessableMessagesForParticipantTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenListingProcessableMessagesForParticipantTestingSuite{
		dbName:        dbName,
		dbClient:      dbClient,
		sut:           repositories.NewConversationMessageRepository(dbClient),
		bridges:       repositories.NewBridgeRepository(dbClient),
		conversations: repositories.NewConversationRepository(dbClient),
		participants:  repositories.NewParticipantRepository(dbClient),
	}
}

func WhenListingProcessableMessagesForParticipantAfterEach(suite *WhenListingProcessableMessagesForParticipantTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

// seed creates a conversation populated with one human author + one
// AI participant whose ENS domain is `agent.eth`, and three messages:
//   - "hi everyone"          (mentions nobody)
//   - "hey @agent.eth ping"  (mentions the agent)
//   - "follow up"            (no mention) — authored by the human
func (s *WhenListingProcessableMessagesForParticipantTestingSuite) seed(agentEnsDomain string) (*models.Conversation, *models.Participant, *models.Participant, []*models.ConversationMessage) {
	bridge := model_builders.NewBridgeBuilder().Build()
	_, _ = s.bridges.Save(bridge)

	conv := model_builders.NewConversationBuilder().WithBridgeID(bridge.ID).Build()
	_, _ = s.conversations.Save(conv)

	human := model_builders.NewUserParticipantBuilder().Build()
	_, _ = s.participants.Save(human)

	agent := model_builders.NewAgentParticipantBuilder().Build()
	agent.EnsDomain = lo.ToPtr(agentEnsDomain)
	_, _ = s.participants.Save(agent)

	plain := model_builders.NewConversationMessageBuilder().
		WithConversationID(conv.ID).
		WithParticipantID(human.ID).
		WithContent("hi everyone").
		Build()
	mention := model_builders.NewConversationMessageBuilder().
		WithConversationID(conv.ID).
		WithParticipantID(human.ID).
		WithContent("hey " + agentEnsDomain + " ping").
		Build()
	followUp := model_builders.NewConversationMessageBuilder().
		WithConversationID(conv.ID).
		WithParticipantID(human.ID).
		WithContent("follow up").
		Build()
	_, _ = s.sut.Save(plain)
	_, _ = s.sut.Save(mention)
	_, _ = s.sut.Save(followUp)

	return conv, agent, human, []*models.ConversationMessage{plain, mention, followUp}
}

func TestWhenListingProcessableMessagesForParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given the participant is a facilitator", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return every message in the conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingProcessableMessagesForParticipantBeforeEach()
			defer WhenListingProcessableMessagesForParticipantAfterEach(suite)

			conv, agent, _, allMessages := suite.seed("agent.eth")

			result, err := suite.sut.ListProcessableForParticipantAndConversation(agent, lo.ToPtr(enums.ParticipantRoleFacilitator), conv.ID)

			assert.NoError(t, err)
			assert.Len(t, result, len(allMessages))
		})
	})

	t.Run("Given the participant is a clone or expert", func(t *testing.T) {
		t.Parallel()

		t.Run("Should only return messages mentioning their ENS domain", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingProcessableMessagesForParticipantBeforeEach()
			defer WhenListingProcessableMessagesForParticipantAfterEach(suite)

			conv, agent, _, _ := suite.seed("agent.eth")

			result, err := suite.sut.ListProcessableForParticipantAndConversation(agent, lo.ToPtr(enums.ParticipantRoleClone), conv.ID)

			assert.NoError(t, err)
			assert.Len(t, result, 1)
			assert.Contains(t, result[0].Content, "agent.eth")
		})
	})

	t.Run("Given the participant is human", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return their own messages OR messages mentioning their ENS", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingProcessableMessagesForParticipantBeforeEach()
			defer WhenListingProcessableMessagesForParticipantAfterEach(suite)

			conv, _, human, _ := suite.seed("someone.eth")
			human.EnsDomain = lo.ToPtr("human.eth")

			result, err := suite.sut.ListProcessableForParticipantAndConversation(human, lo.ToPtr(enums.ParticipantRoleHuman), conv.ID)

			assert.NoError(t, err)
			// Human authored all 3 messages in seed; mention substring is "someone.eth", not human's ens.
			assert.Len(t, result, 3)
		})
	})
}
