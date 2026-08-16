package services_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenForceConcludingConversationTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenForceConcludingConversationBeforeEach(t *testing.T) *WhenForceConcludingConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenForceConcludingConversationTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenForceConcludingConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker concludes the conversation", func(t *testing.T) {
		t.Parallel()

		conversationId := uuid.New()

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdForceConcludeWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.PostConversationsIdForceConcludeResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNoContent},
				}, nil).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			initSuite(suite)

			assert.NoError(t, suite.sut.ForceConcludeConversation(conversationId))
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdForceConcludeWithResponse(gomock.Any(), gomock.Any()).
				Return(&clients.PostConversationsIdForceConcludeResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
					Body:         []byte("not found"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			initSuite(suite)

			assert.Error(t, suite.sut.ForceConcludeConversation(uuid.New()))
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdForceConcludeWithResponse(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			initSuite(suite)

			assert.Error(t, suite.sut.ForceConcludeConversation(uuid.New()))
		})
	})
}
