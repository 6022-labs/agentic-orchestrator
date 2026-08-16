package repositories_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type WhenListingProcessableConversationsForParticipantTestingSuite struct {
	dbName                   string
	dbClient                 db_clients.DatabaseClientInterface
	sut                      *repositories.ConversationRepository
	bridges                  *repositories.BridgeRepository
	participants             *repositories.ParticipantRepository
	conversationParticipants *repositories.ConversationParticipantRepository
	messages                 *repositories.ConversationMessageRepository
}

func WhenListingProcessableConversationsForParticipantBeforeEach() *WhenListingProcessableConversationsForParticipantTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenListingProcessableConversationsForParticipantTestingSuite{
		dbName:                   dbName,
		dbClient:                 dbClient,
		sut:                      repositories.NewConversationRepository(dbClient),
		bridges:                  repositories.NewBridgeRepository(dbClient),
		participants:             repositories.NewParticipantRepository(dbClient),
		conversationParticipants: repositories.NewConversationParticipantRepository(dbClient),
		messages:                 repositories.NewConversationMessageRepository(dbClient),
	}
}

func WhenListingProcessableConversationsForParticipantAfterEach(
	suite *WhenListingProcessableConversationsForParticipantTestingSuite,
) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func (s *WhenListingProcessableConversationsForParticipantTestingSuite) seedConversation(
	bridgeId uuid.UUID,
) *models.Conversation {
	conv := model_builders.NewConversationBuilder().WithBridgeID(bridgeId).Build()
	_, _ = s.sut.Save(conv)
	return conv
}

func (s *WhenListingProcessableConversationsForParticipantTestingSuite) attach(
	conversationId, participantId uuid.UUID,
) {
	_, _ = s.conversationParticipants.Save(model_builders.NewConversationParticipantBuilder().
		WithConversationID(conversationId).
		WithParticipantID(participantId).
		Build())
}

func (s *WhenListingProcessableConversationsForParticipantTestingSuite) postMessage(
	conversationId, participantId uuid.UUID,
	content string,
	createdAt time.Time,
) {
	_, _ = s.messages.Save(model_builders.NewConversationMessageBuilder().
		WithConversationID(conversationId).
		WithParticipantID(participantId).
		WithContent(content).
		WithCreatedAt(createdAt).
		Build())
}

func TestWhenListingProcessableConversationsForParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given the participant is attached to no conversation", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingProcessableConversationsForParticipantBeforeEach()
			defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

			facilitator := model_builders.NewAgentParticipantBuilder().
				Build()
			_, _ = suite.participants.Save(facilitator)

			result, err := suite.sut.ListProcessableForParticipant(
				facilitator,
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				uint64(30),
			)

			assert.NoError(t, err)
			assert.Len(t, result, 0)
		})
	})

	t.Run("Given the conversation is concluded", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return the conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingProcessableConversationsForParticipantBeforeEach()
			defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, _ = suite.bridges.Save(bridge)

			facilitator := model_builders.NewAgentParticipantBuilder().
				Build()
			_, _ = suite.participants.Save(facilitator)

			concluded := model_builders.NewConversationBuilder().
				WithBridgeID(bridge.ID).
				WithStatus(enums.ConversationStatusConcluded).
				Build()
			_, _ = suite.sut.Save(concluded)
			suite.attach(concluded.ID, facilitator.ID)

			result, err := suite.sut.ListProcessableForParticipant(
				facilitator,
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				uint64(30),
			)

			assert.NoError(t, err)
			assert.Len(t, result, 0)
		})
	})

	t.Run("Given the participant is a facilitator", func(t *testing.T) {
		t.Parallel()

		t.Run("And they have never posted in the conversation", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				facilitator := model_builders.NewAgentParticipantBuilder().
					Build()
				_, _ = suite.participants.Save(facilitator)
				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, facilitator.ID)

				result, err := suite.sut.ListProcessableForParticipant(
					facilitator,
					lo.ToPtr(enums.ParticipantRoleFacilitator),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 1)
			})
		})

		t.Run("And their last post is within the tick window", func(t *testing.T) {
			t.Parallel()

			t.Run("Should not return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				facilitator := model_builders.NewAgentParticipantBuilder().
					Build()
				_, _ = suite.participants.Save(facilitator)
				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, facilitator.ID)
				suite.postMessage(conv.ID, facilitator.ID, "recent", time.Now().Add(-5*time.Second))

				result, err := suite.sut.ListProcessableForParticipant(
					facilitator,
					lo.ToPtr(enums.ParticipantRoleFacilitator),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 0)
			})
		})

		t.Run("And their last post is older than the tick window", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				facilitator := model_builders.NewAgentParticipantBuilder().
					Build()
				_, _ = suite.participants.Save(facilitator)
				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, facilitator.ID)
				suite.postMessage(conv.ID, facilitator.ID, "ancient", time.Now().Add(-10*time.Minute))

				result, err := suite.sut.ListProcessableForParticipant(
					facilitator,
					lo.ToPtr(enums.ParticipantRoleFacilitator),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 1)
			})
		})
	})

	t.Run("Given the participant is a human", func(t *testing.T) {
		t.Parallel()

		t.Run("And no message mentions them nor is authored by them", func(t *testing.T) {
			t.Parallel()

			t.Run("Should not return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				human := model_builders.NewAgentParticipantBuilder().
					Build()
				human.EnsDomain = lo.ToPtr("pierre.eth")
				_, _ = suite.participants.Save(human)

				other := model_builders.NewAgentParticipantBuilder().Build()
				_, _ = suite.participants.Save(other)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, human.ID)
				suite.postMessage(conv.ID, other.ID, "unrelated chatter", time.Now())

				result, err := suite.sut.ListProcessableForParticipant(
					human,
					lo.ToPtr(enums.ParticipantRoleHuman),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 0)
			})
		})

		t.Run("And a message mentions their ENS", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				human := model_builders.NewAgentParticipantBuilder().
					Build()
				human.EnsDomain = lo.ToPtr("pierre.eth")
				_, _ = suite.participants.Save(human)
				other := model_builders.NewAgentParticipantBuilder().Build()
				_, _ = suite.participants.Save(other)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, human.ID)
				suite.postMessage(conv.ID, other.ID, "hey pierre.eth check this", time.Now())

				result, err := suite.sut.ListProcessableForParticipant(
					human,
					lo.ToPtr(enums.ParticipantRoleHuman),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 1)
			})
		})

		t.Run("And a message is authored by them", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				human := model_builders.NewAgentParticipantBuilder().
					Build()
				human.EnsDomain = lo.ToPtr("pierre.eth")
				_, _ = suite.participants.Save(human)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, human.ID)
				suite.postMessage(conv.ID, human.ID, "hello team", time.Now())

				result, err := suite.sut.ListProcessableForParticipant(
					human,
					lo.ToPtr(enums.ParticipantRoleHuman),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 1)
			})
		})

		t.Run("And the conversation has both a mention and their own message", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				human := model_builders.NewAgentParticipantBuilder().
					Build()
				human.EnsDomain = lo.ToPtr("pierre.eth")
				_, _ = suite.participants.Save(human)
				other := model_builders.NewAgentParticipantBuilder().Build()
				_, _ = suite.participants.Save(other)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, human.ID)
				suite.postMessage(conv.ID, human.ID, "first post", time.Now().Add(-2*time.Minute))
				suite.postMessage(conv.ID, other.ID, "pierre.eth follow up", time.Now().Add(-1*time.Minute))

				result, err := suite.sut.ListProcessableForParticipant(
					human,
					lo.ToPtr(enums.ParticipantRoleHuman),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 1)
			})
		})
	})

	t.Run("Given the participant is a clone or expert agent", func(t *testing.T) {
		t.Parallel()

		t.Run("And no message mentions them", func(t *testing.T) {
			t.Parallel()

			t.Run("Should not return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				agent := model_builders.NewAgentParticipantBuilder().Build()
				agent.EnsDomain = lo.ToPtr("helper.eth")
				_, _ = suite.participants.Save(agent)
				other := model_builders.NewAgentParticipantBuilder().Build()
				_, _ = suite.participants.Save(other)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, agent.ID)
				suite.postMessage(conv.ID, other.ID, "unrelated", time.Now())

				result, err := suite.sut.ListProcessableForParticipant(
					agent,
					lo.ToPtr(enums.ParticipantRoleClone),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 0)
			})
		})

		t.Run("And the last mention happened before their last reply", func(t *testing.T) {
			t.Parallel()

			t.Run("Should not return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				agent := model_builders.NewAgentParticipantBuilder().Build()
				agent.EnsDomain = lo.ToPtr("helper.eth")
				_, _ = suite.participants.Save(agent)
				other := model_builders.NewAgentParticipantBuilder().Build()
				_, _ = suite.participants.Save(other)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, agent.ID)
				suite.postMessage(conv.ID, other.ID, "helper.eth ping", time.Now().Add(-2*time.Minute))
				suite.postMessage(conv.ID, agent.ID, "on it", time.Now().Add(-1*time.Minute))

				result, err := suite.sut.ListProcessableForParticipant(
					agent,
					lo.ToPtr(enums.ParticipantRoleClone),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 0)
			})
		})

		t.Run("And a new mention landed after their last reply", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				agent := model_builders.NewAgentParticipantBuilder().Build()
				agent.EnsDomain = lo.ToPtr("helper.eth")
				_, _ = suite.participants.Save(agent)
				other := model_builders.NewAgentParticipantBuilder().Build()
				_, _ = suite.participants.Save(other)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, agent.ID)
				suite.postMessage(conv.ID, other.ID, "helper.eth first ping", time.Now().Add(-3*time.Minute))
				suite.postMessage(conv.ID, agent.ID, "on it", time.Now().Add(-2*time.Minute))
				suite.postMessage(conv.ID, other.ID, "helper.eth follow up", time.Now().Add(-1*time.Minute))

				result, err := suite.sut.ListProcessableForParticipant(
					agent,
					lo.ToPtr(enums.ParticipantRoleClone),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 1)
			})
		})

		t.Run("And the only ENS mentions are in the agent's own messages", func(t *testing.T) {
			t.Parallel()

			t.Run("Should not return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				agent := model_builders.NewAgentParticipantBuilder().Build()
				agent.EnsDomain = lo.ToPtr("helper.eth")
				_, _ = suite.participants.Save(agent)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, agent.ID)
				suite.postMessage(conv.ID, agent.ID, "hi I am helper.eth", time.Now().Add(-2*time.Minute))
				suite.postMessage(conv.ID, agent.ID, "still helper.eth here", time.Now().Add(-1*time.Minute))

				result, err := suite.sut.ListProcessableForParticipant(
					agent,
					lo.ToPtr(enums.ParticipantRoleClone),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 0)
			})
		})

		t.Run("And they have never replied but a message mentions them", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return the conversation", func(t *testing.T) {
				t.Parallel()

				suite := WhenListingProcessableConversationsForParticipantBeforeEach()
				defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

				bridge := model_builders.NewBridgeBuilder().Build()
				_, _ = suite.bridges.Save(bridge)
				agent := model_builders.NewAgentParticipantBuilder().Build()
				agent.EnsDomain = lo.ToPtr("helper.eth")
				_, _ = suite.participants.Save(agent)
				other := model_builders.NewAgentParticipantBuilder().Build()
				_, _ = suite.participants.Save(other)

				conv := suite.seedConversation(bridge.ID)
				suite.attach(conv.ID, agent.ID)
				suite.postMessage(conv.ID, other.ID, "helper.eth are you there?", time.Now())

				result, err := suite.sut.ListProcessableForParticipant(
					agent,
					lo.ToPtr(enums.ParticipantRoleClone),
					uint64(30),
				)

				assert.NoError(t, err)
				assert.Len(t, result, 1)
			})
		})
	})

	t.Run("Given another participant is attached to a different conversation", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return the other's conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingProcessableConversationsForParticipantBeforeEach()
			defer WhenListingProcessableConversationsForParticipantAfterEach(suite)

			bridge := model_builders.NewBridgeBuilder().Build()
			_, _ = suite.bridges.Save(bridge)
			facilitator := model_builders.NewAgentParticipantBuilder().
				Build()
			outsider := model_builders.NewAgentParticipantBuilder().Build()
			_, _ = suite.participants.Save(facilitator)
			_, _ = suite.participants.Save(outsider)

			convA := suite.seedConversation(bridge.ID)
			convB := suite.seedConversation(bridge.ID)
			suite.attach(convA.ID, facilitator.ID)
			suite.attach(convB.ID, outsider.ID)

			result, err := suite.sut.ListProcessableForParticipant(
				facilitator,
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				uint64(30),
			)

			assert.NoError(t, err)
			assert.Len(t, result, 1)
			assert.Equal(t, convA.ID, result[0].ID)
		})
	})
}
