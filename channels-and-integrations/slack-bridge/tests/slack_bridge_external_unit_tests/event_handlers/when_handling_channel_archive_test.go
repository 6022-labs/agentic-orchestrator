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

type WhenHandlingChannelArchiveTestingSuite struct {
	sut *event_handlers.ChannelArchiveHandler

	mockSlackThreadRepository *repositories_mocks.MockSlackThreadRepositoryInterface
}

func WhenHandlingChannelArchiveBeforeEach(t *testing.T) *WhenHandlingChannelArchiveTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackThreadRepository := repositories_mocks.NewMockSlackThreadRepositoryInterface(mockController)

	sut := event_handlers.NewChannelArchiveHandler(
		zap.NewNop(),
		mockSlackThreadRepository,
	)

	return &WhenHandlingChannelArchiveTestingSuite{
		sut: sut,

		mockSlackThreadRepository: mockSlackThreadRepository,
	}
}

func TestWhenHandlingChannelArchive(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while concluding threads by channel id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelArchiveTestingSuite) {
			suite.mockSlackThreadRepository.
				EXPECT().
				ConcludeThreadsByChannelId(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelArchiveBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelArchiveRequest{
				Event: requests.ChannelArchiveEvent{
					Channel: "channel_id",
				},
			})

			assert.Error(t, err)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingChannelArchiveTestingSuite) {
			suite.mockSlackThreadRepository.
				EXPECT().
				ConcludeThreadsByChannelId(gomock.Any()).
				Return(nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingChannelArchiveBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.ChannelArchiveRequest{
				Event: requests.ChannelArchiveEvent{
					Channel: "channel_id",
				},
			})

			assert.NoError(t, err)
		})
	})
}
