package event_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingMessageChangedTestingSuite struct {
	sut *event_handlers.MessageChangedHandler

	mockSlackThreadConclude *interfaces_mocks.MockSlackThreadConcludeInterface
}

func WhenHandlingMessageChangedBeforeEach(t *testing.T) *WhenHandlingMessageChangedTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackThreadConclude := interfaces_mocks.NewMockSlackThreadConcludeInterface(mockController)

	sut := event_handlers.NewMessageChangedHandler(
		zap.NewNop(),
		mockSlackThreadConclude,
	)

	return &WhenHandlingMessageChangedTestingSuite{
		sut:                     sut,
		mockSlackThreadConclude: mockSlackThreadConclude,
	}
}

func TestWhenHandlingMessageChanged(t *testing.T) {
	t.Parallel()

	t.Run("Given message is hidden", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error from ConcludeIfExists", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageChangedBeforeEach(t)
			suite.mockSlackThreadConclude.
				EXPECT().
				ConcludeIfExists("1234567890.123456", "C1234567890").
				Return(assert.AnError)

			err := suite.sut.Handle(&requests.MessageChangedRequest{
				Event: requests.MessageChangedEvent{
					PreviousMessage: requests.MessageChangedEventPreviousMessage{
						TS: "1234567890.123456",
					},
					Message: requests.MessageChangedEventMessage{
						Hidden: true,
					},
					Channel: "C1234567890",
				},
			})

			assert.Error(t, err)
		})

		t.Run("Should not return error if ConcludeIfExists succeeds", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageChangedBeforeEach(t)
			suite.mockSlackThreadConclude.
				EXPECT().
				ConcludeIfExists("1234567890.123456", "C1234567890").
				Return(nil)

			err := suite.sut.Handle(&requests.MessageChangedRequest{
				Event: requests.MessageChangedEvent{
					PreviousMessage: requests.MessageChangedEventPreviousMessage{
						TS: "1234567890.123456",
					},
					Message: requests.MessageChangedEventMessage{
						Hidden: true,
					},
					Channel: "C1234567890",
				},
			})

			assert.NoError(t, err)
		})
	})

	t.Run("Given message is not hidden", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not call ConcludeIfExists and return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageChangedBeforeEach(t)
			// No EXPECT() call, ConcludeIfExists should not be called

			err := suite.sut.Handle(&requests.MessageChangedRequest{
				Event: requests.MessageChangedEvent{
					PreviousMessage: requests.MessageChangedEventPreviousMessage{
						TS: "1234567890.123456",
					},
					Message: requests.MessageChangedEventMessage{
						Hidden: false,
					},
					Channel: "C1234567890",
				},
			})

			assert.NoError(t, err)
		})
	})
}
