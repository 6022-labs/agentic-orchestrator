package event_handlers_test

import (
	"context"

	"testing"

	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	channel_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	channel_traces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/traces_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingChannelIdChangedTestingSuite struct {
	sut *event_handlers.ChannelIdChangedHandler

	mockSlackChannelRepository *repositories_mocks.MockSlackChannelRepositoryInterface
}

func WhenHandlingChannelIdChangedBeforeEach(t *testing.T) *WhenHandlingChannelIdChangedTestingSuite {
	mockController := gomock.NewController(t)

	mockChannelMetrics := channel_metrics_mocks.NewMockChannelMetricsInterface(mockController)
	mockChannelMetrics.EXPECT().
		RecordEvent(gomock.Any(), "channel.id_changed", gomock.Any(), gomock.Any()).
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

	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)

	sut := event_handlers.NewChannelIdChangedHandler(
		mockChannelMetrics,
		mockChannelTracer,
		zap.NewNop(),
		mockSlackChannelRepository,
	)

	return &WhenHandlingChannelIdChangedTestingSuite{
		sut: sut,

		mockSlackChannelRepository: mockSlackChannelRepository,
	}
}

func TestWhenHandlingChannelIdChanged(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while updating channel id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelIdChangedTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				UpdateId(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelIdChangedBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(context.Background(), &requests.ChannelIdChangedRequest{
				Event: requests.ChannelIdChangedEvent{
					OldChannelId: "channel_id",
					NewChannelId: "new_channel_id",
				},
			})

			assert.Error(t, err)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelIdChangedTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				UpdateId(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelIdChangedBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(context.Background(), &requests.ChannelIdChangedRequest{
				Event: requests.ChannelIdChangedEvent{
					OldChannelId: "channel_id",
					NewChannelId: "new_channel_id",
				},
			})

			assert.NoError(t, err)
		})
	})
}
