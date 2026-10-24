package services_test

import (
	"context"

	stderrors "errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingConversationTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenGettingConversationBeforeEach(t *testing.T) *WhenGettingConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenGettingConversationTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenGettingConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given the conversation exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200: &clients.ResponsesConversationResponse{
						Id:     lo.ToPtr(conversationId.String()),
						Status: lo.ToPtr(clients.EnumsConversationStatus("concluded")),
					},
				}, nil).
				Times(1)
		}

		t.Run("Should return the mapped conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.GetConversation(context.Background(), conversationId)

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, conversationId, result.ID)
				assert.Equal(t, "concluded", result.Status)
			}
		})
	})

	t.Run("Given the conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
				}, nil).
				Times(1)
		}

		t.Run("Should return a NotFoundError", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.GetConversation(context.Background(), conversationId)

			var target *errors.NotFoundError
			assert.ErrorAs(t, err, &target)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.GetConversation(context.Background(), conversationId)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdWithResponse(gomock.Any(), conversationId.String()).
				Return(nil, stderrors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.GetConversation(context.Background(), conversationId)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
