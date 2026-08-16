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
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenIngestingMessageTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenIngestingMessageBeforeEach(t *testing.T) *WhenIngestingMessageTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenIngestingMessageTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenIngestingMessage(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker accepts the message", func(t *testing.T) {
		t.Parallel()

		conversationId := uuid.New()
		participantId := uuid.New()

		initSuite := func(suite *WhenIngestingMessageTestingSuite) uuid.UUID {
			messageId := uuid.New()

			suite.mockClient.EXPECT().
				PostConversationsIdMessagesWithBodyWithResponse(
					gomock.Any(),
					conversationId.String(),
					"application/json",
					gomock.Any(),
				).
				DoAndReturn(func(_ any, _ string, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostConversationsIdMessagesResponse, error) {
					var request clients.RequestsIngestMessageRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, participantId.String(), *request.ParticipantId)
					assert.Equal(t, "hello", *request.Content)

					return &clients.PostConversationsIdMessagesResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusCreated},
						JSON201:      &clients.ResponsesConversationMessageResponse{Id: lo.ToPtr(messageId.String())},
					}, nil
				}).
				Times(1)

			return messageId
		}

		t.Run("Should return the message id", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			messageId := initSuite(suite)

			got, err := suite.sut.IngestMessage(conversationId, participantId, "hello")

			assert.NoError(t, err)
			assert.Equal(t, messageId, got)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenIngestingMessageTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdMessagesWithBodyWithResponse(gomock.Any(), gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostConversationsIdMessagesResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
					Body:         []byte("not found"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.IngestMessage(uuid.New(), uuid.New(), "hello")

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenIngestingMessageTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsIdMessagesWithBodyWithResponse(gomock.Any(), gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.IngestMessage(uuid.New(), uuid.New(), "hello")

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})
}
