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

type WhenListingConversationParticipantsTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenListingConversationParticipantsBeforeEach(t *testing.T) *WhenListingConversationParticipantsTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenListingConversationParticipantsTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenListingConversationParticipants(t *testing.T) {
	t.Parallel()

	t.Run("Given the conversation has participants", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationParticipantsTestingSuite, conversationId uuid.UUID) uuid.UUID {
			participantId := uuid.New()

			suite.mockClient.EXPECT().
				GetConversationsIdParticipantsWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdParticipantsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200: &[]clients.ResponsesParticipantResponse{{
						Id:                lo.ToPtr(participantId.String()),
						Role:              lo.ToPtr(clients.EnumsParticipantRole("expert")),
						CollectionAddress: lo.ToPtr("0xabc"),
						CollectionTokenId: lo.ToPtr(7),
						EnsDomain:         lo.ToPtr("ada.experts.80002.6022.eth"),
					}},
				}, nil).
				Times(1)

			return participantId
		}

		t.Run("Should return the mapped participants", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)
			conversationId := uuid.New()
			participantId := initSuite(suite, conversationId)

			result, err := suite.sut.ListConversationParticipants(context.Background(), conversationId)

			assert.NoError(t, err)
			if assert.Len(t, result, 1) {
				participant := result[0]
				assert.Equal(t, participantId, participant.ID)
				if assert.NotNil(t, participant.Role) {
					assert.Equal(t, "expert", *participant.Role)
				}
				if assert.NotNil(t, participant.CollectionTokenId) {
					assert.Equal(t, uint64(7), *participant.CollectionTokenId)
				}
			}
		})
	})

	t.Run("Given the conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationParticipantsTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdParticipantsWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdParticipantsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
				}, nil).
				Times(1)
		}

		t.Run("Should return a NotFoundError", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.ListConversationParticipants(context.Background(), conversationId)

			var target *errors.NotFoundError
			assert.ErrorAs(t, err, &target)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationParticipantsTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdParticipantsWithResponse(gomock.Any(), conversationId.String()).
				Return(&clients.GetConversationsIdParticipantsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.ListConversationParticipants(context.Background(), conversationId)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationParticipantsTestingSuite, conversationId uuid.UUID) {
			suite.mockClient.EXPECT().
				GetConversationsIdParticipantsWithResponse(gomock.Any(), conversationId.String()).
				Return(nil, stderrors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			result, err := suite.sut.ListConversationParticipants(context.Background(), conversationId)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
