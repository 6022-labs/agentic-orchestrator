package clients_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCallingPostCompletionTestingSuite struct {
	sut    *clients.BridgeCallbackClient
	server *httptest.Server
}

func WhenCallingPostCompletionBeforeEach(handler http.HandlerFunc) *WhenCallingPostCompletionTestingSuite {
	server := httptest.NewServer(handler)

	sut := clients.NewBridgeCallbackClient(clients.NewBridgeCallbackClientParams{
		Logger:     zap.NewNop(),
		HttpClient: http.DefaultClient,
	})

	return &WhenCallingPostCompletionTestingSuite{
		sut:    sut,
		server: server,
	}
}

func WhenCallingPostCompletionAfterEach(suite *WhenCallingPostCompletionTestingSuite) {
	suite.server.Close()
}

func TestWhenCallingPostCompletion(t *testing.T) {
	t.Parallel()

	t.Run("Given the bridge returns 2xx", func(t *testing.T) {
		t.Parallel()

		t.Run("Should POST the payload to the callback URL", func(t *testing.T) {
			t.Parallel()

			var capturedMethod string
			var capturedPath string
			var capturedBody []byte
			suite := WhenCallingPostCompletionBeforeEach(func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				capturedPath = r.URL.Path
				capturedBody, _ = io.ReadAll(r.Body)
				w.WriteHeader(http.StatusNoContent)
			})
			defer WhenCallingPostCompletionAfterEach(suite)

			body := &models.ExternalCompletionCallbackRequest{
				ConversationId: uuid.New(),
				ParticipantId:  uuid.New(),
				Content:        "hello",
			}

			err := suite.sut.PostCompletion(suite.server.URL+"/internal/broker/callback", body)

			assert.NoError(t, err)
			assert.Equal(t, http.MethodPost, capturedMethod)
			assert.Equal(t, "/internal/broker/callback", capturedPath)

			var decoded models.ExternalCompletionCallbackRequest
			assert.NoError(t, json.Unmarshal(capturedBody, &decoded))
			assert.Equal(t, body.ConversationId, decoded.ConversationId)
			assert.Equal(t, body.ParticipantId, decoded.ParticipantId)
			assert.Equal(t, body.Content, decoded.Content)
		})
	})

	t.Run("Given the bridge returns non-2xx", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingPostCompletionBeforeEach(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})
			defer WhenCallingPostCompletionAfterEach(suite)

			err := suite.sut.PostCompletion(suite.server.URL, &models.ExternalCompletionCallbackRequest{})

			assert.Error(t, err)
		})
	})
}
