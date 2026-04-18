package clients_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/settings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCallingConversationBrokerTestingSuite struct {
	sut    *clients.ConversationBrokerClient
	server *httptest.Server

	handler http.HandlerFunc
}

func WhenCallingConversationBrokerBeforeEach() *WhenCallingConversationBrokerTestingSuite {
	suite := &WhenCallingConversationBrokerTestingSuite{}

	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if suite.handler == nil {
			http.NotFound(w, r)
			return
		}
		suite.handler(w, r)
	}))

	suite.sut = clients.NewConversationBrokerClient(clients.NewConversationBrokerClientParams{
		Logger:     zap.NewNop(),
		HttpClient: http.DefaultClient,
		ConversationBrokerSettings: &settings.ConversationBrokerSettings{BaseUrl: suite.server.URL},
	})

	return suite
}

func WhenCallingConversationBrokerAfterEach(suite *WhenCallingConversationBrokerTestingSuite) {
	suite.server.Close()
}

func TestWhenCallingConversationBroker(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker accepts a RegisterBridge", func(t *testing.T) {
		t.Parallel()

		var capturedPath string
		var capturedMethod string
		var capturedBody []byte
		assignedId := uuid.New()

		initSuite := func(suite *WhenCallingConversationBrokerTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				capturedPath = r.URL.Path
				capturedBody, _ = io.ReadAll(r.Body)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(&models.ExternalBridgeResponse{
					ID:          assignedId,
					Kind:        "slack",
					BaseUrl: "http://slack.local/cb",
				})
			}
		}

		t.Run("Should POST /bridges and decode the assigned id", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingConversationBrokerBeforeEach()
			defer WhenCallingConversationBrokerAfterEach(suite)
			initSuite(suite)

			resp, err := suite.sut.RegisterBridge(&models.ExternalRegisterBridgeRequest{
				Kind:        "slack",
				BaseUrl: "http://slack.local/cb",
			})

			assert.NoError(t, err)
			assert.Equal(t, http.MethodPost, capturedMethod)
			assert.Equal(t, "/bridges", capturedPath)
			assert.Contains(t, string(capturedBody), `"kind":"slack"`)
			assert.Equal(t, assignedId, resp.ID)
		})
	})

	t.Run("Given the broker accepts an IngestMessage", func(t *testing.T) {
		t.Parallel()

		var capturedPath string
		conversationId := uuid.New()
		participantId := uuid.New()
		messageId := uuid.New()

		initSuite := func(suite *WhenCallingConversationBrokerTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(&models.ExternalConversationMessageResponse{
					ID:             messageId,
					ConversationID: conversationId,
					ParticipantID:  participantId,
					Content:        "hi",
				})
			}
		}

		t.Run("Should POST /conversations/{id}/messages", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingConversationBrokerBeforeEach()
			defer WhenCallingConversationBrokerAfterEach(suite)
			initSuite(suite)

			resp, err := suite.sut.IngestMessage(conversationId, &models.ExternalIngestMessageRequest{
				ParticipantID: participantId,
				Content:       "hi",
			})

			assert.NoError(t, err)
			assert.Equal(t, "/conversations/"+conversationId.String()+"/messages", capturedPath)
			assert.Equal(t, messageId, resp.ID)
		})
	})

	t.Run("Given the broker returns 404", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingConversationBrokerTestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingConversationBrokerBeforeEach()
			defer WhenCallingConversationBrokerAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.GetConversation(uuid.New())

			assert.Error(t, err)
		})
	})
}
