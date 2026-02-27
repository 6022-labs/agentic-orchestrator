package event_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingChannelLeftTestingSuite struct {
	sut *event_handlers.ChannelLeftHandler

	mockSlackThreadRepository *repositories_mocks.MockSlackThreadRepositoryInterface
}

func WhenHandlingChannelLeftBeforeEach(t *testing.T) *WhenHandlingChannelLeftTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackThreadRepository := repositories_mocks.NewMockSlackThreadRepositoryInterface(mockController)

	sut := event_handlers.NewChannelLeftHandler(
		zap.NewNop(),
		mockSlackThreadRepository,
	)

	return &WhenHandlingChannelLeftTestingSuite{
		sut:                       sut,
		mockSlackThreadRepository: mockSlackThreadRepository,
	}
}

func TestWhenHandlingChannelLeft(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while concluding threads by channel id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelLeftTestingSuite) {
			suite.mockSlackThreadRepository.
				EXPECT().
				ConcludeThreadsByChannelId(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelLeftBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelLeftRequest{
				Event: requests.ChannelLeftEvent{
					Channel: "channel_id",
				},
			})

			assert.Error(t, err)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelLeftTestingSuite) {
			suite.mockSlackThreadRepository.
				EXPECT().
				ConcludeThreadsByChannelId(gomock.Any()).
				Return(nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelLeftBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelLeftRequest{
				Event: requests.ChannelLeftEvent{
					Channel: "channel_id",
				},
			})

			assert.NoError(t, err)
		})
	})
}
