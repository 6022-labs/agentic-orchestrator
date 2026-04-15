package services_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/http_channels_and_integrations_bridge_mocks/clients_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenPostingCompletionCallbackTestingSuite struct {
	sut *services.BridgeCallbackRequester

	mockClient *clients_mocks.MockBridgeCallbackClientInterface
}

func WhenPostingCompletionCallbackBeforeEach(t *testing.T) *WhenPostingCompletionCallbackTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockBridgeCallbackClientInterface(mockController)

	sut := services.NewBridgeCallbackRequester(mockClient)

	return &WhenPostingCompletionCallbackTestingSuite{
		sut: sut,

		mockClient: mockClient,
	}
}

func TestWhenPostingCompletionCallback(t *testing.T) {
	t.Parallel()

	t.Run("Given the underlying client fails", func(t *testing.T) {
		t.Parallel()

		bridge := model_builders.NewBridgeBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		agent := model_builders.NewParticipantBuilder().Build()

		initSuite := func(suite *WhenPostingCompletionCallbackTestingSuite) {
			suite.mockClient.EXPECT().
				PostCompletion(bridge.CallbackUrl, gomock.Any()).
				Return(errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingCompletionCallbackBeforeEach(t)
			initSuite(suite)

			err := suite.sut.PostCompletion(bridge, conversation, agent, "hello")

			assert.Error(t, err)
		})
	})

	t.Run("Given the call succeeds", func(t *testing.T) {
		t.Parallel()

		bridge := model_builders.NewBridgeBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		agent := model_builders.NewParticipantBuilder().Build()
		content := "completed response"

		initSuite := func(suite *WhenPostingCompletionCallbackTestingSuite) {
			suite.mockClient.EXPECT().
				PostCompletion(bridge.CallbackUrl, gomock.Any()).
				DoAndReturn(func(_ string, body *models.ExternalCompletionCallbackRequest) error {
					assert.Equal(t, conversation.ID, body.ConversationId)
					assert.Equal(t, agent.ID, body.ParticipantId)
					assert.Equal(t, content, body.Content)
					return nil
				}).
				Times(1)
		}

		t.Run("Should send the bridge.callback_url with the mapped payload", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingCompletionCallbackBeforeEach(t)
			initSuite(suite)

			err := suite.sut.PostCompletion(bridge, conversation, agent, content)

			assert.NoError(t, err)
		})
	})
}
