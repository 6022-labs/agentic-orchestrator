package clients_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCallingRemoteAgentTestingSuite struct {
	sut    *clients.RemoteAgentClient
	server *httptest.Server

	handler http.HandlerFunc
}

func WhenCallingRemoteAgentBeforeEach() *WhenCallingRemoteAgentTestingSuite {
	suite := &WhenCallingRemoteAgentTestingSuite{}

	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if suite.handler == nil {
			http.NotFound(w, r)
			return
		}
		suite.handler(w, r)
	}))

	suite.sut = clients.NewRemoteAgentClient(clients.NewRemoteAgentClientParams{
		Logger:              zap.NewNop(),
		HttpClient:          http.DefaultClient,
		WellKnownHttpClient: http.DefaultClient,
	})

	return suite
}

func WhenCallingRemoteAgentAfterEach(suite *WhenCallingRemoteAgentTestingSuite) {
	suite.server.Close()
}

func TestWhenCallingRemoteAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given the node serves a valid 6022 card", func(t *testing.T) {
		t.Parallel()

		var capturedPath string

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"protocol":"6022","protocolVersion":"0.1","name":"claudine","endpoints":{"responses":"https://node.example/api/responses"},"paymentMethods":["x402"]}`))
			}
		}

		t.Run("Should GET /.well-known/6022 and decode the card", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			card, err := suite.sut.GetWellKnownCard(suite.server.URL + "/")

			assert.NoError(t, err)
			assert.Equal(t, "/.well-known/6022", capturedPath)
			assert.Equal(t, "claudine", card.Name)
			assert.Equal(t, "https://node.example/api/responses", card.Endpoints.Responses)
			assert.Equal(t, []string{"x402"}, card.PaymentMethods)
		})
	})

	t.Run("Given the card carries a foreign protocol", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"protocol":"other"}`))
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.GetWellKnownCard(suite.server.URL)

			assert.Error(t, err)
		})
	})

	t.Run("Given the node answers the responses call", func(t *testing.T) {
		t.Parallel()

		var capturedMethod string
		var capturedBody []byte
		var capturedSignatureHeader string

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				capturedBody, _ = io.ReadAll(r.Body)
				capturedSignatureHeader = r.Header.Get("PAYMENT-SIGNATURE")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"resp_1","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hi"}]}]}`))
			}
		}

		t.Run("Should POST the request without a payment header and decode the response", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			response, err := suite.sut.CreateResponse(suite.server.URL+"/api/responses", "agent.eth", &models.ExternalCreateResponseRequest{
				Instructions: "be nice",
				Input:        []*models.ExternalResponseInputMessage{{Type: "message", Role: "user", Content: "hello"}},
			}, nil)

			assert.NoError(t, err)
			assert.Equal(t, http.MethodPost, capturedMethod)
			assert.Equal(t, "completed", response.Status)
			assert.Empty(t, capturedSignatureHeader)

			var decoded map[string]any
			assert.NoError(t, json.Unmarshal(capturedBody, &decoded))
			assert.Equal(t, "be nice", decoded["instructions"])
		})
	})

	t.Run("Given a payment signature is provided", func(t *testing.T) {
		t.Parallel()

		var capturedSignatureHeader string

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedSignatureHeader = r.Header.Get("PAYMENT-SIGNATURE")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"resp_1","status":"completed"}`))
			}
		}

		t.Run("Should send it as the PAYMENT-SIGNATURE header", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.CreateResponse(suite.server.URL, "agent.eth", &models.ExternalCreateResponseRequest{}, &models.ExternalRequestHeaders{PaymentSignature: "signed-payload"})

			assert.NoError(t, err)
			assert.Equal(t, "signed-payload", capturedSignatureHeader)
		})
	})

	t.Run("Given the node answers 402 with a PAYMENT-REQUIRED header", func(t *testing.T) {
		t.Parallel()

		var requestCount int

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.Header().Set(clients.PaymentRequiredHeader, "base64payload")
				w.WriteHeader(http.StatusPaymentRequired)
			}
		}

		t.Run("Should return a PaymentRequiredError carrying the header without retrying", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.CreateResponse(suite.server.URL, "agent.eth", &models.ExternalCreateResponseRequest{}, nil)

			var paymentRequired *cb_errors.PaymentRequiredError
			assert.ErrorAs(t, err, &paymentRequired)
			assert.Equal(t, "agent.eth", paymentRequired.EnsDomain)
			assert.Equal(t, "base64payload", paymentRequired.PaymentRequirements)
			assert.Equal(t, 1, requestCount)
		})
	})

	t.Run("Given the node answers 5xx", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		t.Run("Should return a plain error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.CreateResponse(suite.server.URL, "agent.eth", &models.ExternalCreateResponseRequest{}, nil)

			assert.Error(t, err)
			var paymentRequired *cb_errors.PaymentRequiredError
			assert.False(t, errors.As(err, &paymentRequired))
		})
	})

	t.Run("Given the node answers the chat completions call", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"hey"}}]}`))
			}
		}

		t.Run("Should decode the choices", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			response, err := suite.sut.CreateChatCompletion(suite.server.URL, "agent.eth", &models.ExternalChatCompletionRequest{}, nil)

			assert.NoError(t, err)
			assert.Len(t, response.Choices, 1)
			assert.Equal(t, "hey", response.Choices[0].Message.Content)
		})
	})
}
