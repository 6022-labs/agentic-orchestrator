package clients_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCallingPostCompletionTestingSuite struct {
	sut    *clients.BridgeCallbackClient
	server *httptest.Server

	handler http.HandlerFunc
}

func WhenCallingPostCompletionBeforeEach() *WhenCallingPostCompletionTestingSuite {
	suite := &WhenCallingPostCompletionTestingSuite{}

	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if suite.handler == nil {
			http.NotFound(w, r)
			return
		}
		suite.handler(w, r)
	}))

	suite.sut = clients.NewBridgeCallbackClient(clients.NewBridgeCallbackClientParams{
		Logger:     zap.NewNop(),
		HttpClient: http.DefaultClient,
	})

	return suite
}

func WhenCallingPostCompletionAfterEach(suite *WhenCallingPostCompletionTestingSuite) {
	suite.server.Close()
}

func TestWhenCallingPostCompletion(t *testing.T) {
	t.Parallel()

	t.Run("Given the bridge returns 2xx", func(t *testing.T) {
		t.Parallel()

		var capturedMethod string
		var capturedPath string
		var capturedBody []byte

		initSuite := func(suite *WhenCallingPostCompletionTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				capturedPath = r.URL.Path
				capturedBody, _ = io.ReadAll(r.Body)
				w.WriteHeader(http.StatusNoContent)
			}
		}

		t.Run("Should POST the payload to the callback URL", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingPostCompletionBeforeEach()
			defer WhenCallingPostCompletionAfterEach(suite)
			initSuite(suite)

			body := &models.ExternalCompletionCallbackRequest{
				ConversationId:    uuid.New(),
				CollectionAddress: "0xabc",
				CollectionTokenId: uint64(42),
				Content:           "hello",
			}

			err := suite.sut.PostCompletion(suite.server.URL, body)

			assert.NoError(t, err)
			assert.Equal(t, http.MethodPost, capturedMethod)
			assert.Equal(t, clients.BridgeCallbackEndpoint, capturedPath)

			var decoded models.ExternalCompletionCallbackRequest
			assert.NoError(t, json.Unmarshal(capturedBody, &decoded))
			assert.Equal(t, body.ConversationId, decoded.ConversationId)
			assert.Equal(t, body.CollectionAddress, decoded.CollectionAddress)
			assert.Equal(t, body.CollectionTokenId, decoded.CollectionTokenId)
			assert.Equal(t, body.Content, decoded.Content)
		})
	})

	t.Run("Given the bridge returns non-2xx", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingPostCompletionTestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingPostCompletionBeforeEach()
			defer WhenCallingPostCompletionAfterEach(suite)
			initSuite(suite)

			err := suite.sut.PostCompletion(suite.server.URL, &models.ExternalCompletionCallbackRequest{})

			assert.Error(t, err)
		})
	})
}
