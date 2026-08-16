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
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingConversationMessagesTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenGettingConversationMessagesBeforeEach(t *testing.T) *WhenGettingConversationMessagesTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenGettingConversationMessagesTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenGettingConversationMessages(t *testing.T) {
	t.Parallel()

	t.Run("Given the conversation has messages", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationMessagesTestingSuite, conversationId uuid.UUID) uuid.UUID {
			messageId := uuid.New()

			suite.mockClient.EXPECT().
				GetConversationsIdMessagesWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdMessagesResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200: &[]clients.ResponsesConversationMessageResponse{{
						Id:             lo.ToPtr(messageId.String()),
						ConversationId: lo.ToPtr(conversationId.String()),
						Content:        lo.ToPtr("hello"),
					}},
				}, nil).
				Times(1)

			return messageId
		}

		t.Run("Should return the mapped messages", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)
			conversationId := uuid.New()
			messageId := initSuite(suite, conversationId)

			result, err := suite.sut.GetConversationMessages(conversationId)

			assert.NoError(t, err)
			if assert.Len(t, result, 1) {
				assert.Equal(t, messageId, result[0].ID)
				assert.Equal(t, conversationId, result[0].ConversationId)
				assert.Equal(t, "hello", result[0].Content)
			}
		})
	})

	t.Run("Given the conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationMessagesTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdMessagesWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdMessagesResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
				}, nil).
				Times(1)
		}

		t.Run("Should return a NotFoundError", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.GetConversationMessages(conversationId)

			var target *errors.NotFoundError
			assert.ErrorAs(t, err, &target)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationMessagesTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdMessagesWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdMessagesResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.GetConversationMessages(conversationId)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationMessagesTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdMessagesWithResponse(gomock.Any(), conversationId.String()).
				Return(nil, stderrors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.GetConversationMessages(conversationId)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
