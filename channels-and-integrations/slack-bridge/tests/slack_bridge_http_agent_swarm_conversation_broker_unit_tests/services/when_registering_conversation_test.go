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

type WhenRegisteringConversationTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenRegisteringConversationBeforeEach(t *testing.T) *WhenRegisteringConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenRegisteringConversationTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenRegisteringConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker creates the conversation", func(t *testing.T) {
		t.Parallel()

		bridgeId := uuid.New()
		participantIds := []uuid.UUID{uuid.New(), uuid.New()}

		initSuite := func(suite *WhenRegisteringConversationTestingSuite) uuid.UUID {
			conversationId := uuid.New()

			suite.mockClient.EXPECT().
				PostConversationsWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				DoAndReturn(func(_ any, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostConversationsResponse, error) {
					var request clients.RequestsRegisterConversationRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, bridgeId.String(), *request.BridgeId)
					assert.Equal(t, "demo", *request.Description)
					assert.Equal(t, []string{participantIds[0].String(), participantIds[1].String()}, *request.ParticipantIds)

					return &clients.PostConversationsResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusCreated},
						JSON201:      &clients.ResponsesConversationResponse{Id: lo.ToPtr(conversationId.String())},
					}, nil
				}).
				Times(1)

			return conversationId
		}

		t.Run("Should return the conversation id", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringConversationBeforeEach(t)
			conversationId := initSuite(suite)

			got, err := suite.sut.RegisterConversation(bridgeId, lo.ToPtr("demo"), participantIds)

			assert.NoError(t, err)
			assert.Equal(t, conversationId, got)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringConversationTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostConversationsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringConversationBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.RegisterConversation(uuid.New(), nil, nil)

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringConversationTestingSuite) {
			suite.mockClient.EXPECT().
				PostConversationsWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringConversationBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.RegisterConversation(uuid.New(), nil, nil)

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})
}
