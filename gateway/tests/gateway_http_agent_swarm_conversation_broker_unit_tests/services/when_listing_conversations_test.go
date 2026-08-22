package services_test

import (
	"context"

	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingConversationsTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenListingConversationsBeforeEach(t *testing.T) *WhenListingConversationsTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenListingConversationsTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenListingConversations(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker returns conversations", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationsTestingSuite) uuid.UUID {
			conversationId := uuid.New()

			suite.mockClient.EXPECT().GetConversationsWithResponse(gomock.Any(), &clients.GetConversationsParams{
				Page:            lo.ToPtr(1),
				ElementsPerPage: lo.ToPtr(20),
			}).Return(&clients.GetConversationsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &clients.ResponsesListConversationsResponse{
					Page:            lo.ToPtr(1),
					ElementsPerPage: lo.ToPtr(20),
					TotalElements:   lo.ToPtr(1),
					Conversations: &[]clients.ResponsesConversationResponse{{
						Id:        lo.ToPtr(conversationId.String()),
						Status:    lo.ToPtr(clients.EnumsConversationStatus("ongoing")),
						CreatedAt: lo.ToPtr("2026-07-04T10:00:00Z"),
						UpdatedAt: lo.ToPtr("2026-07-04T11:00:00Z"),
					}},
				},
			}, nil).Times(1)

			return conversationId
		}

		t.Run("Should map the generated response to the use case response", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			conversationId := initSuite(suite)

			result, err := suite.sut.ListConversations(context.Background(), lo.ToPtr(uint(1)), lo.ToPtr(uint(20)))

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, uint(1), result.Page)
				assert.Equal(t, uint(20), result.ElementsPerPage)
				assert.Equal(t, uint(1), result.TotalElements)
				if assert.Len(t, result.Conversations, 1) {
					conversation := result.Conversations[0]
					assert.Equal(t, conversationId, conversation.ID)
					assert.Equal(t, "ongoing", conversation.Status)
					assert.Equal(t, time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC), conversation.CreatedAt)
				}
			}
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationsTestingSuite) {
			suite.mockClient.EXPECT().
				GetConversationsWithResponse(gomock.Any(), gomock.Any()).
				Return(&clients.GetConversationsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ListConversations(context.Background(), nil, nil)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationsTestingSuite) {
			suite.mockClient.EXPECT().
				GetConversationsWithResponse(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationsBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.ListConversations(context.Background(), nil, nil)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
