package event_handlers_test

import (
	"context"

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

type WhenHandlingChannelLeftTestingSuite struct {
	sut *event_handlers.ChannelLeftHandler

	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester             *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenHandlingChannelLeftBeforeEach(t *testing.T) *WhenHandlingChannelLeftTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := event_handlers.NewChannelLeftHandler(zap.NewNop(), mockSlackConversationRepository, mockBrokerRequester)

	return &WhenHandlingChannelLeftTestingSuite{
		sut:                             sut,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:             mockBrokerRequester,
	}
}

func TestWhenHandlingChannelLeft(t *testing.T) {
	t.Parallel()

	t.Run("Given the channel has a single linked broker conversation", func(t *testing.T) {
		t.Parallel()

		convId := uuid.New()

		initSuite := func(suite *WhenHandlingChannelLeftTestingSuite, channelId string) {
			suite.mockSlackConversationRepository.EXPECT().
				ListByChannelId(gomock.Any(), channelId).
				Return([]*models.SlackConversation{
					{SlackChannelId: channelId, SlackThreadTs: "1.0", BrokerConversationID: convId},
				}, nil)
			suite.mockBrokerRequester.EXPECT().ForceConcludeConversation(gomock.Any(), convId).Return(nil)
		}

		t.Run("Should force-conclude the broker conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelLeftBeforeEach(t)
			req := &requests.ChannelLeftRequest{}
			req.Event.Channel = "C123"

			initSuite(suite, req.Event.Channel)

			err := suite.sut.Handle(context.Background(), req)
			assert.NoError(t, err)
		})
	})

	t.Run("Given the channel has no linked broker conversations", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelLeftTestingSuite, channelId string) {
			suite.mockSlackConversationRepository.EXPECT().ListByChannelId(gomock.Any(), channelId).Return(nil, nil)
		}

		t.Run("Should not call the broker", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelLeftBeforeEach(t)
			req := &requests.ChannelLeftRequest{}
			req.Event.Channel = "C123"

			initSuite(suite, req.Event.Channel)

			err := suite.sut.Handle(context.Background(), req)
			assert.NoError(t, err)
		})
	})
}
