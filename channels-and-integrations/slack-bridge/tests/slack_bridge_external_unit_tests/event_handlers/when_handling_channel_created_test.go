package event_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
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

	slackSettings := &settings.SlackSettings{
		AcceptedChannelPrefix: "ai",
	}

	mockSlackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)

	sut := event_handlers.NewChannelCreatedHandler(
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
				Sync(gomock.Any()).
				Times(0)
		}

		t.Run("Should not sync the channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelCreatedBeforeEach(t)

			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelCreatedRequest{
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
				Sync(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelCreatedBeforeEach(t)

			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelCreatedRequest{
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
				Sync(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelCreatedBeforeEach(t)

			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelCreatedRequest{
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
