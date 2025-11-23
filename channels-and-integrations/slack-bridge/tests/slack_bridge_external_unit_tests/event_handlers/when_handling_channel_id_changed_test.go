package event_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
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

	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)

	sut := event_handlers.NewChannelIdChangedHandler(
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
				UpdateId(gomock.Any(), gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelIdChangedBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelIdChangedRequest{
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
				UpdateId(gomock.Any(), gomock.Any()).
				Return(nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelIdChangedBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelIdChangedRequest{
				Event: requests.ChannelIdChangedEvent{
					OldChannelId: "channel_id",
					NewChannelId: "new_channel_id",
				},
			})

			assert.NoError(t, err)
		})
	})
}
