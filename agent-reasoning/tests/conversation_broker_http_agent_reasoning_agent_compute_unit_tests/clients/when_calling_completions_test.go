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
}

func WhenCallingCompletionsBeforeEach(handler http.HandlerFunc) *WhenCallingCompletionsTestingSuite {
	server := httptest.NewServer(handler)
	sut := clients.NewAgentReasoningAgentComputeClient(clients.NewAgentReasoningAgentComputeClientParams{
		Logger:                             zap.NewNop(),
		HttpClient:                         http.DefaultClient,
		AgentReasoningAgentComputeSettings: &settings.AgentReasoningAgentComputeSettings{BaseUrl: server.URL},
	})
	return &WhenCallingCompletionsTestingSuite{sut: sut, server: server}
}

func WhenCallingCompletionsAfterEach(suite *WhenCallingCompletionsTestingSuite) {
	suite.server.Close()
}

func TestWhenCallingCompletions(t *testing.T) {
	t.Parallel()

	t.Run("Given the upstream returns 200", func(t *testing.T) {
		t.Parallel()

		t.Run("Should POST to /agent/{addr}/{tokenId}/compute and decode the response", func(t *testing.T) {
			t.Parallel()

			var capturedPath string
			var capturedBody []byte
			suite := WhenCallingCompletionsBeforeEach(func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.Path
				capturedBody, _ = io.ReadAll(r.Body)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(&models.ExternalAgentComputeResponses{
					IsFinal:   true,
					Responses: []*string{lo.ToPtr("ok")},
				})
			})
			defer WhenCallingCompletionsAfterEach(suite)

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

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingCompletionsBeforeEach(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})
			defer WhenCallingCompletionsAfterEach(suite)

			_, err := suite.sut.Completions("0xabc", uint64(1), &models.ExternalAgentComputeRequest{})

			assert.Error(t, err)
		})
	})
}