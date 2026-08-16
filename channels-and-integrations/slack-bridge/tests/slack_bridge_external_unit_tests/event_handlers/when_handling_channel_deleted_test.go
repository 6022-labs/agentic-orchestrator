package event_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingChannelDeletedTestingSuite struct {
	sut *event_handlers.ChannelDeletedHandler

	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester             *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenHandlingChannelDeletedBeforeEach(t *testing.T) *WhenHandlingChannelDeletedTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := event_handlers.NewChannelDeletedHandler(zap.NewNop(), mockSlackConversationRepository, mockBrokerRequester)

	return &WhenHandlingChannelDeletedTestingSuite{
		sut:                             sut,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:             mockBrokerRequester,
	}
}

func TestWhenHandlingChannelDeleted(t *testing.T) {
	t.Parallel()

	t.Run("Given the channel has linked broker conversations", func(t *testing.T) {
		t.Parallel()

		convId := uuid.New()

		initSuite := func(suite *WhenHandlingChannelDeletedTestingSuite, channelId string) {
			suite.mockSlackConversationRepository.EXPECT().
				ListByChannelId(channelId).
				Return([]*models.SlackConversation{
					{SlackChannelId: channelId, SlackThreadTs: "1.0", BrokerConversationID: convId},
				}, nil)
			suite.mockBrokerRequester.EXPECT().ForceConcludeConversation(convId).Return(nil)
		}

		t.Run("Should force-conclude every broker conversation in the channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelDeletedBeforeEach(t)
			req := &requests.ChannelDeletedRequest{}
			req.Event.Channel = "C123"

			initSuite(suite, req.Event.Channel)

			err := suite.sut.Handle(req)
			assert.NoError(t, err)
		})
	})
}
