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

type WhenHandlingMessageChangedTestingSuite struct {
	sut *event_handlers.MessageChangedHandler

	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester             *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenHandlingMessageChangedBeforeEach(t *testing.T) *WhenHandlingMessageChangedTestingSuite {
	mockController := gomock.NewController(t)

	mockChannelMetrics := channel_metrics_mocks.NewMockChannelMetricsInterface(mockController)
	mockChannelMetrics.EXPECT().
		RecordEvent(gomock.Any(), "message.updated", gomock.Any(), gomock.Any()).
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

	sut := event_handlers.NewMessageChangedHandler(
		mockChannelMetrics,
		mockChannelTracer,
		zap.NewNop(),
		mockSlackConversationRepository,
		mockBrokerRequester,
	)

	return &WhenHandlingMessageChangedTestingSuite{
		sut:                             sut,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:             mockBrokerRequester,
	}
}

func newMessageChangedRequest(channel, prevTs string, hidden bool) *requests.MessageChangedRequest {
	req := &requests.MessageChangedRequest{}
	req.Event.Channel = channel
	req.Event.PreviousMessage.TS = prevTs
	req.Event.Message.Hidden = hidden
	return req
}

func TestWhenHandlingMessageChanged(t *testing.T) {
	t.Parallel()

	t.Run("Given the changed message is not hidden", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be a no-op", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageChangedBeforeEach(t)
			req := newMessageChangedRequest("C123", "1.0", false)

			err := suite.sut.Handle(context.Background(), req)
			assert.NoError(t, err)
		})
	})

	t.Run("Given the hidden message has a linked broker conversation", func(t *testing.T) {
		t.Parallel()

		convId := uuid.New()

		initSuite := func(suite *WhenHandlingMessageChangedTestingSuite, channel, ts string) {
			suite.mockSlackConversationRepository.EXPECT().
				GetBySlackIds(gomock.Any(), channel, ts).
				Return(&models.SlackConversation{BrokerConversationID: convId}, nil)
			suite.mockBrokerRequester.EXPECT().ForceConcludeConversation(gomock.Any(), convId).Return(nil)
		}

		t.Run("Should force-conclude the broker conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageChangedBeforeEach(t)
			req := newMessageChangedRequest("C123", "1.0", true)

			initSuite(suite, req.Event.Channel, req.Event.PreviousMessage.TS)

			err := suite.sut.Handle(context.Background(), req)
			assert.NoError(t, err)
		})
	})

	t.Run("Given the hidden message has no linked broker conversation", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageChangedTestingSuite, channel, ts string) {
			suite.mockSlackConversationRepository.EXPECT().
				GetBySlackIds(gomock.Any(), channel, ts).
				Return(nil, nil)
		}

		t.Run("Should be a no-op", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageChangedBeforeEach(t)
			req := newMessageChangedRequest("C123", "1.0", true)

			initSuite(suite, req.Event.Channel, req.Event.PreviousMessage.TS)

			err := suite.sut.Handle(context.Background(), req)
			assert.NoError(t, err)
		})
	})
}
