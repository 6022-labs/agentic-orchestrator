package clients_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/settings"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCallingCompletionsTestingSuite struct {
	sut    *clients.AgentReasoningAgentComputeClient
	server *httptest.Server

	handler http.HandlerFunc
}

func WhenCallingCompletionsBeforeEach() *WhenCallingCompletionsTestingSuite {
	suite := &WhenCallingCompletionsTestingSuite{}

	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if suite.handler == nil {
			http.NotFound(w, r)
			return
		}
		suite.handler(w, r)
	}))

	suite.sut = clients.NewAgentReasoningAgentComputeClient(clients.NewAgentReasoningAgentComputeClientParams{
		Logger:                             zap.NewNop(),
		HttpClient:                         http.DefaultClient,
		AgentReasoningAgentComputeSettings: &settings.AgentReasoningAgentComputeSettings{BaseUrl: suite.server.URL},
	})

	return suite
}

func WhenCallingCompletionsAfterEach(suite *WhenCallingCompletionsTestingSuite) {
	suite.server.Close()
}

func TestWhenCallingCompletions(t *testing.T) {
	t.Parallel()

	t.Run("Given the upstream returns 200", func(t *testing.T) {
		t.Parallel()

		var capturedPath string
		var capturedBody []byte

		initSuite := func(suite *WhenCallingCompletionsTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.Path
				capturedBody, _ = io.ReadAll(r.Body)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(&models.ExternalAgentComputeResponses{
					IsFinal:   true,
					Responses: []*string{lo.ToPtr("ok")},
				})
			}
		}

		t.Run("Should POST to /agent/{addr}/{tokenId}/compute and decode the response", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingCompletionsBeforeEach()
			defer WhenCallingCompletionsAfterEach(suite)
			initSuite(suite)

			result, err := suite.sut.Completions("0xabc", uint64(7), &models.ExternalAgentComputeRequest{
				Platform: "slack",
			})

			assert.NoError(t, err)
			assert.True(t, result.IsFinal)
			assert.Equal(t, "/agent/0xabc/7/compute", capturedPath)
			assert.Contains(t, string(capturedBody), `"platform":"slack"`)
		})
	})

	t.Run("Given the upstream returns non-200", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingCompletionsTestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingCompletionsBeforeEach()
			defer WhenCallingCompletionsAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.Completions("0xabc", uint64(1), &models.ExternalAgentComputeRequest{})

			assert.Error(t, err)
		})
	})
}
