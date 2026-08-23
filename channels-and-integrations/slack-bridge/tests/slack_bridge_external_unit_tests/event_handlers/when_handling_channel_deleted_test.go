package event_handlers_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	channel_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	channel_traces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/traces_mocks/interfaces_mocks"
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

	mockChannelMetrics := channel_metrics_mocks.NewMockChannelMetricsInterface(mockController)
	mockChannelMetrics.EXPECT().
		RecordEvent(gomock.Any(), "channel.deleted", gomock.Any(), gomock.Any()).
		Times(1)

	// Spans are not asserted here.
	mockChannelTracer := channel_traces_mocks.NewMockChannelTracerInterface(mockController)
	mockSpan := channel_traces_mocks.NewMockSpan(mockController)
	mockSpan.EXPECT().End().AnyTimes()
	mockSpan.EXPECT().Fail(gomock.Any()).AnyTimes()
	mockChannelTracer.EXPECT().
		StartEventHandling(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ string) (context.Context, traces_interfaces.Span) {
			return ctx, mockSpan
		}).
		AnyTimes()

	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := event_handlers.NewChannelDeletedHandler(
		mockChannelMetrics,
		mockChannelTracer,
		zap.NewNop(),
		mockSlackConversationRepository,
		mockBrokerRequester,
	)

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
				ListByChannelId(gomock.Any(), channelId).
				Return([]*models.SlackConversation{
					{SlackChannelId: channelId, SlackThreadTs: "1.0", BrokerConversationID: convId},
				}, nil)
			suite.mockBrokerRequester.EXPECT().ForceConcludeConversation(gomock.Any(), convId).Return(nil)
		}

		t.Run("Should force-conclude every broker conversation in the channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelDeletedBeforeEach(t)
			req := &requests.ChannelDeletedRequest{}
			req.Event.Channel = "C123"

			initSuite(suite, req.Event.Channel)

			err := suite.sut.Handle(context.Background(), req)
			assert.NoError(t, err)
		})
	})
}
