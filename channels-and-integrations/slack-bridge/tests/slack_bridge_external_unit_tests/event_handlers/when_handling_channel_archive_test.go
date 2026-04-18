package event_handlers_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingChannelArchiveTestingSuite struct {
	sut *event_handlers.ChannelArchiveHandler

	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester                    *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenHandlingChannelArchiveBeforeEach(t *testing.T) *WhenHandlingChannelArchiveTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := event_handlers.NewChannelArchiveHandler(zap.NewNop(), mockSlackConversationRepository, mockBrokerRequester)

	return &WhenHandlingChannelArchiveTestingSuite{
		sut:                                 sut,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:                    mockBrokerRequester,
	}
}

func TestWhenHandlingChannelArchive(t *testing.T) {
	t.Parallel()

	t.Run("Given the channel has multiple linked broker conversations", func(t *testing.T) {
		t.Parallel()

		convId1 := uuid.New()
		convId2 := uuid.New()

		initSuite := func(suite *WhenHandlingChannelArchiveTestingSuite, channelId string) {
			suite.mockSlackConversationRepository.EXPECT().
				ListByChannelId(channelId).
				Return([]*models.SlackConversation{
					{SlackChannelId: channelId, SlackThreadTs: "1.0", BrokerConversationID: convId1},
					{SlackChannelId: channelId, SlackThreadTs: "2.0", BrokerConversationID: convId2},
				}, nil)
			suite.mockBrokerRequester.EXPECT().ForceConcludeConversation(convId1).Return(nil)
			suite.mockBrokerRequester.EXPECT().ForceConcludeConversation(convId2).Return(nil)
		}

		t.Run("Should force-conclude every broker conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelArchiveBeforeEach(t)
			req := &requests.ChannelArchiveRequest{}
			req.Event.Channel = "C123"

			initSuite(suite, req.Event.Channel)

			err := suite.sut.Handle(req)
			assert.NoError(t, err)
		})
	})

	t.Run("Given the lookup fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelArchiveTestingSuite, channelId string) {
			suite.mockSlackConversationRepository.EXPECT().ListByChannelId(channelId).Return(nil, errors.New("boom"))
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelArchiveBeforeEach(t)
			req := &requests.ChannelArchiveRequest{}
			req.Event.Channel = "C123"

			initSuite(suite, req.Event.Channel)

			err := suite.sut.Handle(req)
			assert.Error(t, err)
		})
	})
}
