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

type WhenDetachingParticipantTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenDetachingParticipantBeforeEach(t *testing.T) *WhenDetachingParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenDetachingParticipantTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenDetachingParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker detaches the participant", func(t *testing.T) {
		t.Parallel()

		conversationId := uuid.New()
		participantId := uuid.New()

		initSuite := func(suite *WhenDetachingParticipantTestingSuite) {
			suite.mockClient.EXPECT().
				DeleteConversationsIdParticipantsParticipantIdWithResponse(
					gomock.Any(),
					conversationId.String(),
					participantId.String(),
				).
				Return(&clients.DeleteConversationsIdParticipantsParticipantIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNoContent},
				}, nil).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenDetachingParticipantBeforeEach(t)
			initSuite(suite)

			assert.NoError(t, suite.sut.DetachParticipant(conversationId, participantId))
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenDetachingParticipantTestingSuite) {
			suite.mockClient.EXPECT().
				DeleteConversationsIdParticipantsParticipantIdWithResponse(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&clients.DeleteConversationsIdParticipantsParticipantIdResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
					Body:         []byte("not found"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDetachingParticipantBeforeEach(t)
			initSuite(suite)

			assert.Error(t, suite.sut.DetachParticipant(uuid.New(), uuid.New()))
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenDetachingParticipantTestingSuite) {
			suite.mockClient.EXPECT().
				DeleteConversationsIdParticipantsParticipantIdWithResponse(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDetachingParticipantBeforeEach(t)
			initSuite(suite)

			assert.Error(t, suite.sut.DetachParticipant(uuid.New(), uuid.New()))
		})
	})
}
