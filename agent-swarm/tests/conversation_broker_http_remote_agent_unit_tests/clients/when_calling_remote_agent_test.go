package clients_test

import (
	"context"

	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func rateLimitedHandler(retryAfter string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if retryAfter != "" {
			w.Header().Set("Retry-After", retryAfter)
		}
		w.WriteHeader(http.StatusTooManyRequests)
	}
}

func (suite *WhenCallingRemoteAgentTestingSuite) retryAfter(t *testing.T) time.Duration {
	_, err := suite.sut.CreateResponse(context.Background(),
		suite.server.URL,
		"agent.eth",
		&models.ExternalCreateResponseRequest{},
		nil,
	)

	var unavailable *cb_errors.RemoteAgentUnavailableError
	require.ErrorAs(t, err, &unavailable)
	return unavailable.RetryAfter
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
				_, _ = w.Write(
					[]byte(
						`{"protocol":"6022","protocolVersion":"0.1","name":"claudine","endpoints":{"responses":"https://node.example/ap` +
							`i/responses"},"paymentMethods":["x402"]}`,
					),
				)
			}
		}

		t.Run("Should GET /.well-known/6022 and decode the card", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			card, raw, err := suite.sut.GetWellKnownCard(context.Background(), suite.server.URL+"/")

			assert.NoError(t, err)
			assert.Equal(t, "/.well-known/6022", capturedPath)
			assert.Equal(
				t,
				`{"protocol":"6022","protocolVersion":"0.1","name":"claudine","endpoints":{"responses":"https://node.example/api/`+
					`responses"},"paymentMethods":["x402"]}`,
				string(raw),
			)
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

			_, _, err := suite.sut.GetWellKnownCard(context.Background(), suite.server.URL)

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
				_, _ = w.Write(
					[]byte(
						`{"id":"resp_1","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_` +
							`text","text":"hi"}]}]}`,
					),
				)
			}
		}

		t.Run("Should POST the request without a payment header and decode the response", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			response, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL+"/api/responses",
				"agent.eth",
				&models.ExternalCreateResponseRequest{
					Instructions: "be nice",
					Input: []*models.ExternalResponseInputMessage{
						{Type: "message", Role: "user", Content: "hello"},
					},
				},
				nil,
			)

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

			_, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalCreateResponseRequest{},
				&models.ExternalRequestHeaders{PaymentSignature: "signed-payload"},
			)

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

			_, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalCreateResponseRequest{},
				nil,
			)

			var paymentRequired *cb_errors.PaymentRequiredError
			assert.ErrorAs(t, err, &paymentRequired)
			assert.Equal(t, "agent.eth", paymentRequired.EnsDomain)
			assert.Equal(t, "base64payload", paymentRequired.PaymentRequirements)
			assert.Equal(t, 1, requestCount)
		})
	})

	t.Run("Given the node answers 500", func(t *testing.T) {
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

			_, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalCreateResponseRequest{},
				nil,
			)

			assert.Error(t, err)
			var paymentRequired *cb_errors.PaymentRequiredError
			assert.False(t, errors.As(err, &paymentRequired))
			var unavailable *cb_errors.RemoteAgentUnavailableError
			assert.False(t, errors.As(err, &unavailable))
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

			response, err := suite.sut.CreateChatCompletion(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalChatCompletionRequest{},
				nil,
			)

			assert.NoError(t, err)
			assert.Len(t, response.Choices, 1)
			assert.Equal(t, "hey", response.Choices[0].Message.Content)
		})
	})

	t.Run("Given the node answers 429", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte("busy"))
			}
		}

		t.Run("Should return a RemoteAgentUnavailableError carrying the status and body", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalCreateResponseRequest{},
				nil,
			)

			var unavailable *cb_errors.RemoteAgentUnavailableError
			assert.ErrorAs(t, err, &unavailable)
			assert.Equal(t, "agent.eth", unavailable.EnsDomain)
			assert.Equal(t, http.StatusTooManyRequests, unavailable.StatusCode)
			assert.Equal(t, "busy", unavailable.Body)
		})
	})

	t.Run("Given the node answers 502", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte("busy"))
			}
		}

		t.Run("Should return a RemoteAgentUnavailableError carrying the status and body", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalCreateResponseRequest{},
				nil,
			)

			var unavailable *cb_errors.RemoteAgentUnavailableError
			assert.ErrorAs(t, err, &unavailable)
			assert.Equal(t, "agent.eth", unavailable.EnsDomain)
			assert.Equal(t, http.StatusBadGateway, unavailable.StatusCode)
			assert.Equal(t, "busy", unavailable.Body)
		})
	})

	t.Run("Given the node answers 503", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("busy"))
			}
		}

		t.Run("Should return a RemoteAgentUnavailableError carrying the status and body", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalCreateResponseRequest{},
				nil,
			)

			var unavailable *cb_errors.RemoteAgentUnavailableError
			assert.ErrorAs(t, err, &unavailable)
			assert.Equal(t, "agent.eth", unavailable.EnsDomain)
			assert.Equal(t, http.StatusServiceUnavailable, unavailable.StatusCode)
			assert.Equal(t, "busy", unavailable.Body)
		})
	})

	t.Run("Given the node answers 504", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusGatewayTimeout)
				_, _ = w.Write([]byte("busy"))
			}
		}

		t.Run("Should return a RemoteAgentUnavailableError carrying the status and body", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.CreateResponse(context.Background(),
				suite.server.URL,
				"agent.eth",
				&models.ExternalCreateResponseRequest{},
				nil,
			)

			var unavailable *cb_errors.RemoteAgentUnavailableError
			assert.ErrorAs(t, err, &unavailable)
			assert.Equal(t, "agent.eth", unavailable.EnsDomain)
			assert.Equal(t, http.StatusGatewayTimeout, unavailable.StatusCode)
			assert.Equal(t, "busy", unavailable.Body)
		})
	})

	t.Run("Given the node answers 429 with a Retry-After in seconds", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("7")
		}

		t.Run("Should carry that delay", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Equal(t, 7*time.Second, suite.retryAfter(t))
		})
	})

	t.Run("Given the node answers 429 with no Retry-After", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("")
		}

		t.Run("Should carry no delay", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Equal(t, time.Duration(0), suite.retryAfter(t))
		})
	})

	t.Run("Given the node answers 429 with a negative Retry-After", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("-1")
		}

		t.Run("Should carry no delay", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Equal(t, time.Duration(0), suite.retryAfter(t))
		})
	})

	t.Run("Given the node answers 429 with a Retry-After that is not a delay", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("soon")
		}

		t.Run("Should carry no delay", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Equal(t, time.Duration(0), suite.retryAfter(t))
		})
	})

	t.Run("Given the node answers 429 with a Retry-After date in the past", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("Wed, 25 Nov 2020 13:00:00 GMT")
		}

		t.Run("Should carry no delay", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Equal(t, time.Duration(0), suite.retryAfter(t))
		})
	})

	t.Run("Given the node answers 429 with a Retry-After of a day", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("86400")
		}

		t.Run("Should carry a delay of hours", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Greater(t, suite.retryAfter(t), time.Hour)
		})
	})

	t.Run("Given the node answers 429 with a Retry-After larger than 64 bits", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("99999999999999999999999")
		}

		t.Run("Should carry a delay of hours", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Greater(t, suite.retryAfter(t), time.Hour)
		})
	})

	t.Run("Given the node answers 429 with a Retry-After date far ahead", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCallingRemoteAgentTestingSuite) {
			suite.handler = rateLimitedHandler("Wed, 25 Nov 2099 13:00:00 GMT")
		}

		t.Run("Should carry a delay of hours", func(t *testing.T) {
			t.Parallel()

			suite := WhenCallingRemoteAgentBeforeEach()
			defer WhenCallingRemoteAgentAfterEach(suite)
			initSuite(suite)

			assert.Greater(t, suite.retryAfter(t), time.Hour)
		})
	})
}
