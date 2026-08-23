package event_handlers_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	channel_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/metrics_mocks/interfaces_mocks"
	channel_traces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/traces_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingChannelCreatedTestingSuite struct {
	sut *event_handlers.ChannelCreatedHandler

	slackSettings *settings.SlackSettings

	mockSlackChannelSynchronizer *interfaces_mocks.MockSlackChannelSynchronizerInterface
}

func WhenHandlingChannelCreatedBeforeEach(t *testing.T) *WhenHandlingChannelCreatedTestingSuite {
	mockController := gomock.NewController(t)

	mockChannelMetrics := channel_metrics_mocks.NewMockChannelMetricsInterface(mockController)
	mockChannelMetrics.EXPECT().
		RecordEvent(gomock.Any(), "channel.created", gomock.Any(), gomock.Any()).
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

	slackSettings := &settings.SlackSettings{
		AcceptedChannelPrefix: "ai",
	}

	mockSlackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)

	sut := event_handlers.NewChannelCreatedHandler(
		mockChannelMetrics,
		mockChannelTracer,
		zap.NewNop(),
		slackSettings,
		mockSlackChannelSynchronizer,
	)

	return &WhenHandlingChannelCreatedTestingSuite{
		sut: sut,

		slackSettings: slackSettings,

		mockSlackChannelSynchronizer: mockSlackChannelSynchronizer,
	}
}

func TestWhenHandlingChannelCreated(t *testing.T) {
	t.Parallel()

	t.Run("Given channel is not prefixed with accepted prefix", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelCreatedTestingSuite) {
			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), gomock.Any()).
				Times(0)
		}

		t.Run("Should not sync the channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelCreatedBeforeEach(t)

			initSuite(suite)

			err := suite.sut.Handle(context.Background(), &requests.ChannelCreatedRequest{
				Event: requests.ChannelCreatedEvent{
					Channel: requests.ChannelCreatedEventChannel{
						Id:   "channel_id",
						Name: "general", // Not prefixed with "ai"
					},
				},
			})

			assert.NoError(t, err)
		})
	})

	t.Run("Given an error occurs while syncing the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelCreatedTestingSuite) {
			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelCreatedBeforeEach(t)

			initSuite(suite)

			err := suite.sut.Handle(context.Background(), &requests.ChannelCreatedRequest{
				Event: requests.ChannelCreatedEvent{
					Channel: requests.ChannelCreatedEventChannel{
						Id:   "channel_id",
						Name: "ai-agent-channel",
					},
				},
			})

			assert.Error(t, err)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelCreatedTestingSuite) {
			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelCreatedBeforeEach(t)

			initSuite(suite)

			err := suite.sut.Handle(context.Background(), &requests.ChannelCreatedRequest{
				Event: requests.ChannelCreatedEvent{
					Channel: requests.ChannelCreatedEventChannel{
						Id:   "channel_id",
						Name: "ai-agent-channel",
					},
				},
			})

			assert.NoError(t, err)
		})
	})
}
