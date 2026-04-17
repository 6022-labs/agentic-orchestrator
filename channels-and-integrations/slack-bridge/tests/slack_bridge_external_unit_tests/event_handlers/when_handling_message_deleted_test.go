package event_handlers_test

import (
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
	"gorm.io/gorm"
)

type WhenHandlingMessageDeletedTestingSuite struct {
	sut *event_handlers.MessageDeletedHandler

	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester                    *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenHandlingMessageDeletedBeforeEach(t *testing.T) *WhenHandlingMessageDeletedTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := event_handlers.NewMessageDeletedHandler(zap.NewNop(), mockSlackConversationRepository, mockBrokerRequester)

	return &WhenHandlingMessageDeletedTestingSuite{
		sut:                                 sut,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:                    mockBrokerRequester,
	}
}

func newMessageDeletedRequest(channel, prevTs string) *requests.MessageDeletedRequest {
	req := &requests.MessageDeletedRequest{}
	req.Event.Channel = channel
	req.Event.PreviousMessage.TS = prevTs
	return req
}

func TestWhenHandlingMessageDeleted(t *testing.T) {
	t.Parallel()

	t.Run("Given the deleted message has a linked broker conversation", func(t *testing.T) {
		t.Parallel()

		convId := uuid.New()

		initSuite := func(suite *WhenHandlingMessageDeletedTestingSuite, channel, ts string) {
			suite.mockSlackConversationRepository.EXPECT().
				GetBySlackIds(channel, ts).
				Return(&models.SlackConversation{BrokerConversationID: convId}, nil)
			suite.mockBrokerRequester.EXPECT().ForceConcludeConversation(convId).Return(nil)
		}

		t.Run("Should force-conclude the broker conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageDeletedBeforeEach(t)
			req := newMessageDeletedRequest("C123", "1.0")

			initSuite(suite, req.Event.Channel, req.Event.PreviousMessage.TS)

			err := suite.sut.Handle(req)
			assert.NoError(t, err)
		})
	})

	t.Run("Given the deleted message has no linked broker conversation", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageDeletedTestingSuite, channel, ts string) {
			suite.mockSlackConversationRepository.EXPECT().
				GetBySlackIds(channel, ts).
				Return(nil, gorm.ErrRecordNotFound)
		}

		t.Run("Should be a no-op", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageDeletedBeforeEach(t)
			req := newMessageDeletedRequest("C123", "1.0")

			initSuite(suite, req.Event.Channel, req.Event.PreviousMessage.TS)

			err := suite.sut.Handle(req)
			assert.NoError(t, err)
		})
	})
}
