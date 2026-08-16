package services_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenAttachingParticipantTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenAttachingParticipantBeforeEach(t *testing.T) *WhenAttachingParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenAttachingParticipantTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenAttachingParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker attaches the participant", func(t *testing.T) {
		t.Parallel()

		conversationId := uuid.New()
		participantId := uuid.New()

		initSuite := func(suite *WhenAttachingParticipantTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdParticipantsWithBodyWithResponse(
					gomock.Any(),
					conversationId.String(),
					"application/json",
					gomock.Any(),
				).
				DoAndReturn(func(_ any, _ string, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostConversationsIdParticipantsResponse, error) {
					var request clients.RequestsConversationParticipantRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, participantId.String(), *request.ParticipantId)

					return &clients.PostConversationsIdParticipantsResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusCreated},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			initSuite(suite)

			assert.NoError(t, suite.sut.AttachParticipant(conversationId, participantId))
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAttachingParticipantTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdParticipantsWithBodyWithResponse(gomock.Any(), gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostConversationsIdParticipantsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
					Body:         []byte("not found"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			initSuite(suite)

			assert.Error(t, suite.sut.AttachParticipant(uuid.New(), uuid.New()))
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAttachingParticipantTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdParticipantsWithBodyWithResponse(gomock.Any(), gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			initSuite(suite)

			assert.Error(t, suite.sut.AttachParticipant(uuid.New(), uuid.New()))
		})
	})
}
