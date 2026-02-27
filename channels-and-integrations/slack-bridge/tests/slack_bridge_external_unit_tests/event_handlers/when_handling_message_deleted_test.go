package event_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingMessageDeletedTestingSuite struct {
	sut *event_handlers.MessageDeletedHandler

	mockSlackThreadConclude *interfaces_mocks.MockSlackThreadConcludeInterface
}

func WhenHandlingMessageDeletedBeforeEach(t *testing.T) *WhenHandlingMessageDeletedTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackThreadConclude := interfaces_mocks.NewMockSlackThreadConcludeInterface(mockController)

	sut := event_handlers.NewMessageDeletedHandler(
		zap.NewNop(),
		mockSlackThreadConclude,
	)

	return &WhenHandlingMessageDeletedTestingSuite{
		sut:                     sut,
		mockSlackThreadConclude: mockSlackThreadConclude,
	}
}

func TestWhenHandlingMessageDeleted(t *testing.T) {
	t.Parallel()

	t.Run("Given ConcludeIfExists returns an error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageDeletedTestingSuite) {
			suite.mockSlackThreadConclude.
				EXPECT().
				ConcludeIfExists("1234567890.123456", "C1234567890").
				Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageDeletedBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.MessageDeletedRequest{
				Event: requests.MessageDeletedEvent{
					PreviousMessage: requests.MessageDeletedEventPreviousMessage{
						TS: "1234567890.123456",
					},
					Channel: "C1234567890",
				},
			})

			assert.Error(t, err)
		})
	})

	t.Run("Given ConcludeIfExists does not return an error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageDeletedTestingSuite) {
			suite.mockSlackThreadConclude.
				EXPECT().
				ConcludeIfExists("1234567890.123456", "C1234567890").
				Return(nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageDeletedBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Handle(&requests.MessageDeletedRequest{
				Event: requests.MessageDeletedEvent{
					PreviousMessage: requests.MessageDeletedEventPreviousMessage{
						TS: "1234567890.123456",
					},
					Channel: "C1234567890",
				},
			})

			assert.NoError(t, err)
		})
	})
}
