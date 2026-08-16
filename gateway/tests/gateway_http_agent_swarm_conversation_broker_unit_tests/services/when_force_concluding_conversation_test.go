package services_test

import (
	stderrors "errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_swarm_conversation_broker_mocks/clients_mocks"
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

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				PostConversationsIdForceConcludeWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.PostConversationsIdForceConcludeResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNoContent},
				}, nil).
				Times(1)
		}

		t.Run("Should return no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			err := suite.sut.ForceConcludeConversation(conversationId)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				PostConversationsIdForceConcludeWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.PostConversationsIdForceConcludeResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
				}, nil).
				Times(1)
		}

		t.Run("Should return a NotFoundError", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			err := suite.sut.ForceConcludeConversation(conversationId)

			var target *errors.NotFoundError
			assert.ErrorAs(t, err, &target)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				PostConversationsIdForceConcludeWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.PostConversationsIdForceConcludeResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			err := suite.sut.ForceConcludeConversation(conversationId)

			assert.Error(t, err)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForceConcludingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				PostConversationsIdForceConcludeWithResponse(gomock.Any(), conversationId.String()).
				Return(nil, stderrors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForceConcludingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			err := suite.sut.ForceConcludeConversation(conversationId)

			assert.Error(t, err)
		})
	})
}
