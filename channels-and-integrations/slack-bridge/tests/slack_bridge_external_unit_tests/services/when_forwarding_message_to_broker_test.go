package services_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	broker_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenForwardingMessageToBrokerTestingSuite struct {
	sut *services.BrokerMessageForwarder

	bridgeId uuid.UUID

	mockBrokerRequester             *broker_requester_mocks.MockConversationBrokerRequesterInterface
	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockSlackMessageRepository      *repositories_mocks.MockSlackMessageRepositoryInterface
}

func WhenForwardingMessageToBrokerBeforeEach(t *testing.T) *WhenForwardingMessageToBrokerTestingSuite {
	mockController := gomock.NewController(t)

	bridgeId := uuid.New()
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)
	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockSlackMessageRepository := repositories_mocks.NewMockSlackMessageRepositoryInterface(mockController)

	sut := services.NewBrokerMessageForwarder(
		zap.NewNop(),
		&broker_settings.ConversationBrokerSettings{BridgeID: bridgeId},
		mockBrokerRequester,
		mockSlackConversationRepository,
		mockSlackMessageRepository,
	)

	return &WhenForwardingMessageToBrokerTestingSuite{
		sut:                             sut,
		bridgeId:                        bridgeId,
		mockBrokerRequester:             mockBrokerRequester,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockSlackMessageRepository:      mockSlackMessageRepository,
	}
}

func TestWhenForwardingMessageToBroker(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack message has already been forwarded", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForwardingMessageToBrokerTestingSuite, channelId, msgTs string) {
			suite.mockSlackMessageRepository.EXPECT().
				GetBySlackIds(channelId, msgTs).
				Return(&models.SlackMessage{}, nil)
		}

		t.Run("Should skip without contacting the broker", func(t *testing.T) {
			t.Parallel()

			suite := WhenForwardingMessageToBrokerBeforeEach(t)
			channel := model_builders.NewSlackChannelBuilder().Build()
			brokerUserId := uuid.New()
			user := model_builders.NewSlackUserBuilder().Build()
			user.BrokerParticipantId = &brokerUserId

			initSuite(suite, channel.Id, "msg.1")

			err := suite.sut.Forward(channel, "thread.1", "msg.1", "hi", user, nil)
			assert.NoError(t, err)
		})
	})

	t.Run("Given an agent author and the conversation already exists", func(t *testing.T) {
		t.Parallel()

		brokerAgentId := uuid.New()
		brokerConvId := uuid.New()
		brokerMsgId := uuid.New()

		initSuite := func(suite *WhenForwardingMessageToBrokerTestingSuite, channelId, threadTs, msgTs string) {
			suite.mockSlackMessageRepository.EXPECT().GetBySlackIds(channelId, msgTs).Return(nil, nil)
			suite.mockSlackConversationRepository.EXPECT().
				GetBySlackIds(channelId, threadTs).
				Return(&models.SlackConversation{
					SlackChannelId:       channelId,
					SlackThreadTs:        threadTs,
					BrokerConversationID: brokerConvId,
				}, nil)
			suite.mockBrokerRequester.EXPECT().
				IngestMessage(brokerConvId, brokerAgentId, "agent reply").
				Return(brokerMsgId, nil)
			suite.mockSlackMessageRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(m *models.SlackMessage) (*models.SlackMessage, error) {
					assert.Equal(t, channelId, m.SlackChannelId)
					assert.Equal(t, msgTs, m.SlackMessageTs)
					assert.Equal(t, brokerMsgId, m.BrokerMessageID)
					return m, nil
				})
		}

		t.Run("Should ingest the message into the existing conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenForwardingMessageToBrokerBeforeEach(t)
			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()
			agent.BrokerParticipantId = &brokerAgentId

			initSuite(suite, channel.Id, "thread.1", "msg.1")

			err := suite.sut.Forward(channel, "thread.1", "msg.1", "agent reply", nil, agent)
			assert.NoError(t, err)
		})
	})

	t.Run("Given a slack user author with a cached broker id and no conversation yet", func(t *testing.T) {
		t.Parallel()

		brokerUserId := uuid.New()
		brokerConvId := uuid.New()
		brokerMsgId := uuid.New()

		initSuite := func(suite *WhenForwardingMessageToBrokerTestingSuite, channel *models.SlackChannel, threadTs, msgTs string) {
			suite.mockSlackMessageRepository.EXPECT().GetBySlackIds(channel.Id, msgTs).Return(nil, nil)

			// No existing conversation → register with channel agents + author.
			suite.mockSlackConversationRepository.EXPECT().GetBySlackIds(channel.Id, threadTs).Return(nil, nil)
			suite.mockBrokerRequester.EXPECT().
				RegisterConversation(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(bridgeId uuid.UUID, description *string, participantIds []uuid.UUID) (uuid.UUID, error) {
					assert.Equal(t, suite.bridgeId, bridgeId)
					assert.Contains(t, participantIds, brokerUserId)
					return brokerConvId, nil
				})
			suite.mockSlackConversationRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(c *models.SlackConversation) (*models.SlackConversation, error) {
					assert.Equal(t, brokerConvId, c.BrokerConversationID)
					return c, nil
				})

			suite.mockBrokerRequester.EXPECT().
				IngestMessage(brokerConvId, brokerUserId, "hello").
				Return(brokerMsgId, nil)
			suite.mockSlackMessageRepository.EXPECT().Save(gomock.Any()).
				DoAndReturn(func(m *models.SlackMessage) (*models.SlackMessage, error) { return m, nil })
		}

		t.Run("Should register the conversation, ingest, and save the message mapping", func(t *testing.T) {
			t.Parallel()

			suite := WhenForwardingMessageToBrokerBeforeEach(t)
			channel := model_builders.NewSlackChannelBuilder().Build()
			user := model_builders.NewSlackUserBuilder().Build()
			user.BrokerParticipantId = &brokerUserId

			initSuite(suite, channel, "thread.1", "msg.1")

			err := suite.sut.Forward(channel, "thread.1", "msg.1", "hello", user, nil)
			assert.NoError(t, err)
		})
	})

	t.Run("Given a slack user author without a cached broker id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForwardingMessageToBrokerTestingSuite, channelId, msgTs string) {
			suite.mockSlackMessageRepository.EXPECT().GetBySlackIds(channelId, msgTs).Return(nil, nil)
		}

		t.Run("Should return an error — SlackUserSynchronizer must populate the id upstream", func(t *testing.T) {
			t.Parallel()

			suite := WhenForwardingMessageToBrokerBeforeEach(t)
			channel := model_builders.NewSlackChannelBuilder().Build()
			user := model_builders.NewSlackUserBuilder().Build()

			initSuite(suite, channel.Id, "msg.1")

			err := suite.sut.Forward(channel, "thread.1", "msg.1", "hello", user, nil)
			assert.Error(t, err)
		})
	})

	t.Run("Given an agent author without a cached broker id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForwardingMessageToBrokerTestingSuite, channelId, msgTs string) {
			suite.mockSlackMessageRepository.EXPECT().GetBySlackIds(channelId, msgTs).Return(nil, nil)
		}

		t.Run("Should return an error — SlackAgentSynchronizer must populate the id upstream", func(t *testing.T) {
			t.Parallel()

			suite := WhenForwardingMessageToBrokerBeforeEach(t)
			channel := model_builders.NewSlackChannelBuilder().Build()
			agent := model_builders.NewSlackAgentBuilder().Build()

			initSuite(suite, channel.Id, "msg.1")

			err := suite.sut.Forward(channel, "thread.1", "msg.1", "hi", nil, agent)
			assert.Error(t, err)
		})
	})

	t.Run("Given the broker IngestMessage call fails", func(t *testing.T) {
		t.Parallel()

		brokerUserId := uuid.New()
		brokerConvId := uuid.New()

		initSuite := func(suite *WhenForwardingMessageToBrokerTestingSuite, channelId, threadTs, msgTs string) {
			suite.mockSlackMessageRepository.EXPECT().GetBySlackIds(channelId, msgTs).Return(nil, nil)
			suite.mockSlackConversationRepository.EXPECT().
				GetBySlackIds(channelId, threadTs).
				Return(&models.SlackConversation{BrokerConversationID: brokerConvId}, nil)
			suite.mockBrokerRequester.EXPECT().IngestMessage(brokerConvId, brokerUserId, "hi").Return(uuid.Nil, errors.New("boom"))
		}

		t.Run("Should propagate the error and not save a message mapping", func(t *testing.T) {
			t.Parallel()

			suite := WhenForwardingMessageToBrokerBeforeEach(t)
			channel := model_builders.NewSlackChannelBuilder().Build()
			user := model_builders.NewSlackUserBuilder().Build()
			user.BrokerParticipantId = &brokerUserId

			initSuite(suite, channel.Id, "thread.1", "msg.1")

			err := suite.sut.Forward(channel, "thread.1", "msg.1", "hi", user, nil)
			assert.Error(t, err)
		})
	})
}
