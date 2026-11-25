package services_test

import (
	"context"

	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/metrics/interfaces"
	http_node_models "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_http_remote_agent_mocks/clients_mocks"
	remote_agent_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_http_remote_agent_mocks/metrics_mocks/interfaces_mocks"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRequestingRemoteAgentCompletionsTestingSuite struct {
	mockRemoteAgentMetrics *remote_agent_metrics_mocks.MockRemoteAgentMetricsInterface
	sut                    *services.RemoteAgentCompletionsRequester

	mockRemoteAgentClient     *clients_mocks.MockRemoteAgentClientInterface
	x402PaymentPayloadSigner  *interfaces_mocks.MockX402PaymentPayloadSignerInterface
	x402PaymentChallengeCache *interfaces_mocks.MockX402PaymentChallengeCacheInterface
}

func WhenRequestingRemoteAgentCompletionsBeforeEach(t *testing.T) *WhenRequestingRemoteAgentCompletionsTestingSuite {
	mockController := gomock.NewController(t)

	mockRemoteAgentMetrics := remote_agent_metrics_mocks.NewMockRemoteAgentMetricsInterface(mockController)

	mockRemoteAgentClient := clients_mocks.NewMockRemoteAgentClientInterface(mockController)
	x402PaymentPayloadSigner := interfaces_mocks.NewMockX402PaymentPayloadSignerInterface(mockController)
	x402PaymentChallengeCache := interfaces_mocks.NewMockX402PaymentChallengeCacheInterface(mockController)

	sut := services.NewRemoteAgentCompletionsRequester(
		zap.NewNop(),
		mockRemoteAgentMetrics,
		mockRemoteAgentClient,
		x402PaymentPayloadSigner,
		x402PaymentChallengeCache,
	)

	return &WhenRequestingRemoteAgentCompletionsTestingSuite{
		mockRemoteAgentMetrics: mockRemoteAgentMetrics,
		sut:                    sut,

		mockRemoteAgentClient:     mockRemoteAgentClient,
		x402PaymentPayloadSigner:  x402PaymentPayloadSigner,
		x402PaymentChallengeCache: x402PaymentChallengeCache,
	}
}

func (suite *WhenRequestingRemoteAgentCompletionsTestingSuite) noCachedChallenge() {
	suite.x402PaymentChallengeCache.EXPECT().Read(gomock.Any(), "agent.eth").Return(nil, false)
}

func turnForTest() *models.RemoteAgentInput {
	return &models.RemoteAgentInput{
		Instructions: "instructions",
		Input: []*models.RemoteAgentInputMessage{
			{Role: enums.RemoteAgentInputRoleUser, Content: "hello"},
			{Role: enums.RemoteAgentInputRoleAssistant, Content: "hi"},
		},
		Metadata: map[string]string{"ensDomain": "agent.eth"},
	}
}

func completedResponseForTest() *http_node_models.ExternalResponse {
	return &http_node_models.ExternalResponse{
		Status: "completed",
		Output: []*http_node_models.ExternalResponseOutput{
			{
				Type:    "message",
				Role:    "assistant",
				Content: []*http_node_models.ExternalResponseOutputContent{{Type: "output_text", Text: "answer"}},
			},
		},
	}
}

func unavailableForTest(status int, retryAfter time.Duration) error {
	return cb_errors.NewRemoteAgentUnavailableError("agent.eth", status, retryAfter, "")
}

func encodeForTest(t *testing.T, value any) string {
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(encoded)
}

func TestWhenRequestingRemoteAgentCompletions(t *testing.T) {
	t.Parallel()

	t.Run("Given a responses endpoint", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite, turn *models.RemoteAgentInput) {
			suite.noCachedChallenge()
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				DoAndReturn(func(_ context.Context, _, _ string, request *http_node_models.ExternalCreateResponseRequest, _ *http_node_models.ExternalRequestHeaders) (*http_node_models.ExternalResponse, error) {
					assert.Equal(t, turn.Instructions, request.Instructions)
					assert.Len(t, request.Input, 2)
					assert.Equal(t, "message", request.Input[0].Type)
					assert.Equal(t, "user", request.Input[0].Role)
					assert.Equal(t, "assistant", request.Input[1].Role)
					assert.Nil(t, request.Conversation)
					return &http_node_models.ExternalResponse{
						Status: "completed",
						Output: []*http_node_models.ExternalResponseOutput{
							{Type: "reasoning"},
							{
								Type: "message",
								Role: "assistant",
								Content: []*http_node_models.ExternalResponseOutputContent{
									{Type: "output_text", Text: "answer"},
								},
							},
						},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should map the turn and extract message output texts", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			turn := turnForTest()
			initSuite(suite, turn)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())

			completions, err := suite.sut.RequestCompletions(context.Background(), card, turn)

			assert.NoError(t, err)
			assert.Equal(t, []string{"answer"}, completions)
		})
	})

	t.Run("Given the turn carries a memory session", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		t.Run("Should send it as the conversation ref", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			turn := turnForTest()
			turn.RemoteConversationId = "3f2b7c41-9a1e-4d0b-8c33-2f6a1b9d5c91"

			suite.noCachedChallenge()
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				DoAndReturn(func(_ context.Context, _, _ string, request *http_node_models.ExternalCreateResponseRequest, _ *http_node_models.ExternalRequestHeaders) (*http_node_models.ExternalResponse, error) {
					assert.NotNil(t, request.Conversation)
					assert.Equal(t, turn.RemoteConversationId, request.Conversation.Id)
					return &http_node_models.ExternalResponse{
						Status: "completed",
						Output: []*http_node_models.ExternalResponseOutput{
							{
								Type: "message",
								Role: "assistant",
								Content: []*http_node_models.ExternalResponseOutputContent{
									{Type: "output_text", Text: "answer"},
								},
							},
						},
					}, nil
				}).
				Times(1)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())

			completions, err := suite.sut.RequestCompletions(context.Background(), card, turn)

			assert.NoError(t, err)
			assert.Equal(t, []string{"answer"}, completions)
		})
	})

	t.Run("Given the response carries an error payload", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				Return(
					&http_node_models.ExternalResponse{
						Error: &http_node_models.ExternalResponseError{Code: "server_error", Message: "boom"},
					},
					nil,
				).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(suite)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))

			_, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			assert.Error(t, err)
		})
	})

	t.Run("Given the response carries no message output", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				Return(&http_node_models.ExternalResponse{Status: "completed"}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(suite)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))

			_, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			assert.Error(t, err)
		})
	})

	t.Run("Given a chat completions endpoint", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{ChatCompletions: "https://node.example/api/chat"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite, turn *models.RemoteAgentInput) {
			suite.noCachedChallenge()
			suite.mockRemoteAgentClient.EXPECT().
				CreateChatCompletion(gomock.Any(), card.Endpoints.ChatCompletions, "agent.eth", gomock.Any(), gomock.Nil()).
				DoAndReturn(func(_ context.Context, _, _ string, request *http_node_models.ExternalChatCompletionRequest, _ *http_node_models.ExternalRequestHeaders) (*http_node_models.ExternalChatCompletionResponse, error) {
					assert.Len(t, request.Messages, 3)
					assert.Equal(t, "system", request.Messages[0].Role)
					assert.Equal(t, turn.Instructions, request.Messages[0].Content)
					return &http_node_models.ExternalChatCompletionResponse{
						Choices: []*http_node_models.ExternalChatCompletionChoice{
							{
								Message: &http_node_models.ExternalChatCompletionMessage{
									Role:    "assistant",
									Content: "chat answer",
								},
							},
						},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should prepend instructions as a system message and extract choices", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			turn := turnForTest()
			initSuite(suite, turn)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())

			completions, err := suite.sut.RequestCompletions(context.Background(), card, turn)

			assert.NoError(t, err)
			assert.Equal(t, []string{"chat answer"}, completions)
		})
	})

	t.Run("Given the turn carries a memory session and the card only advertises chat completions", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{ChatCompletions: "https://node.example/api/chat"},
		}

		t.Run("Should carry the session in metadata, the only field that dialect has", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			turn := turnForTest()
			turn.RemoteConversationId = "3f2b7c41-9a1e-4d0b-8c33-2f6a1b9d5c91"
			turn.Metadata["memorySessionId"] = turn.RemoteConversationId

			suite.noCachedChallenge()
			suite.mockRemoteAgentClient.EXPECT().
				CreateChatCompletion(gomock.Any(), card.Endpoints.ChatCompletions, "agent.eth", gomock.Any(), gomock.Nil()).
				DoAndReturn(func(_ context.Context, _, _ string, request *http_node_models.ExternalChatCompletionRequest, _ *http_node_models.ExternalRequestHeaders) (*http_node_models.ExternalChatCompletionResponse, error) {
					assert.Equal(t, turn.RemoteConversationId, request.Metadata["memorySessionId"])
					return &http_node_models.ExternalChatCompletionResponse{
						Choices: []*http_node_models.ExternalChatCompletionChoice{
							{
								Message: &http_node_models.ExternalChatCompletionMessage{
									Role:    "assistant",
									Content: "chat answer",
								},
							},
						},
					}, nil
				}).
				Times(1)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())

			completions, err := suite.sut.RequestCompletions(context.Background(), card, turn)

			assert.NoError(t, err)
			assert.Equal(t, []string{"chat answer"}, completions)
		})
	})

	t.Run("Given the node answers 402 then accepts the signed free call", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(t *testing.T, suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{
				X402Version: models.X402Version,
				Accepts:     []models.X402PaymentOption{{Scheme: models.X402SchemeExact, Network: "eip155:80002"}},
			}
			paymentPayload := &models.X402PaymentPayload{
				X402Version: models.X402Version,
				Scheme:      models.X402SchemeExact,
			}
			expectedSignatureHeader := encodeForTest(t, paymentPayload)

			suite.x402PaymentChallengeCache.EXPECT().Read(gomock.Any(), "agent.eth").Return(nil, false)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				Return(nil, cb_errors.NewPaymentRequiredError("agent.eth", encodeForTest(t, challenge))).
				Times(1)
			suite.x402PaymentChallengeCache.EXPECT().Store(gomock.Any(), "agent.eth", challenge)
			suite.x402PaymentPayloadSigner.EXPECT().Sign(challenge, "0").Return(paymentPayload, nil)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(),
					card.Endpoints.Responses,
					"agent.eth",
					gomock.Any(),
					&http_node_models.ExternalRequestHeaders{PaymentSignature: expectedSignatureHeader},
				).
				Return(completedResponseForTest(), nil).
				Times(1)
		}

		t.Run("Should cache the challenge and retry once signed", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(t, suite)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeSettled, "eip155:80002")
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())

			completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			assert.NoError(t, err)
			assert.Equal(t, []string{"answer"}, completions)
		})
	})

	t.Run("Given cached payment challenge exist", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(t *testing.T, suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{X402Version: models.X402Version}
			paymentPayload := &models.X402PaymentPayload{
				X402Version: models.X402Version,
				Scheme:      models.X402SchemeExact,
			}
			expectedSignatureHeader := encodeForTest(t, paymentPayload)

			suite.x402PaymentChallengeCache.EXPECT().Read(gomock.Any(), "agent.eth").Return(challenge, true)
			suite.x402PaymentPayloadSigner.EXPECT().Sign(challenge, "0").Return(paymentPayload, nil)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(),
					card.Endpoints.Responses,
					"agent.eth",
					gomock.Any(),
					&http_node_models.ExternalRequestHeaders{PaymentSignature: expectedSignatureHeader},
				).
				Return(completedResponseForTest(), nil).
				Times(1)
		}

		t.Run("Should pre-sign and call once", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(t, suite)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())
			// The pre-signed call was accepted, so it settled on the first attempt.
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeSettled, gomock.Any())

			completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			assert.NoError(t, err)
			assert.Equal(t, []string{"answer"}, completions)
		})
	})

	t.Run("Given the node rejects the pre-signed call with identical terms", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(t *testing.T, suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{
				X402Version: models.X402Version,
				Accepts:     []models.X402PaymentOption{{Scheme: models.X402SchemeExact, Network: "eip155:80002"}},
			}
			paymentPayload := &models.X402PaymentPayload{
				X402Version: models.X402Version,
				Scheme:      models.X402SchemeExact,
			}
			expectedSignatureHeader := encodeForTest(t, paymentPayload)

			suite.x402PaymentChallengeCache.EXPECT().Read(gomock.Any(), "agent.eth").Return(challenge, true)
			suite.x402PaymentPayloadSigner.EXPECT().Sign(challenge, "0").Return(paymentPayload, nil)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(),
					card.Endpoints.Responses,
					"agent.eth",
					gomock.Any(),
					&http_node_models.ExternalRequestHeaders{PaymentSignature: expectedSignatureHeader},
				).
				Return(nil, cb_errors.NewPaymentRequiredError("agent.eth", encodeForTest(t, challenge))).
				Times(1)
			suite.x402PaymentChallengeCache.EXPECT().Store(gomock.Any(), "agent.eth", challenge)
		}

		t.Run("Should return the PaymentRequiredError without a doomed retry", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(t, suite)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeSkipped, "eip155:80002")
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))

			_, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			var paymentRequired *cb_errors.PaymentRequiredError
			assert.ErrorAs(t, err, &paymentRequired)
		})
	})

	t.Run("Given the signer cannot sign after a 402", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(t *testing.T, suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{X402Version: models.X402Version}

			suite.x402PaymentChallengeCache.EXPECT().Read(gomock.Any(), "agent.eth").Return(nil, false)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				Return(nil, cb_errors.NewPaymentRequiredError("agent.eth", encodeForTest(t, challenge))).
				Times(1)
			suite.x402PaymentChallengeCache.EXPECT().Store(gomock.Any(), "agent.eth", challenge)
			suite.x402PaymentPayloadSigner.EXPECT().Sign(challenge, "0").Return(nil, assert.AnError)
		}

		t.Run("Should return the PaymentRequiredError without retrying", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(t, suite)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeSkipped, gomock.Any())
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))

			_, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			var paymentRequired *cb_errors.PaymentRequiredError
			assert.ErrorAs(t, err, &paymentRequired)
		})
	})

	t.Run("Given the node answers 402 twice", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(t *testing.T, suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{X402Version: models.X402Version}
			paymentPayload := &models.X402PaymentPayload{X402Version: models.X402Version}
			expectedSignatureHeader := encodeForTest(t, paymentPayload)

			suite.x402PaymentChallengeCache.EXPECT().Read(gomock.Any(), "agent.eth").Return(nil, false)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				Return(nil, cb_errors.NewPaymentRequiredError("agent.eth", encodeForTest(t, challenge))).
				Times(1)
			suite.x402PaymentChallengeCache.EXPECT().Store(gomock.Any(), "agent.eth", challenge)
			suite.x402PaymentPayloadSigner.EXPECT().Sign(challenge, "0").Return(paymentPayload, nil)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(),
					card.Endpoints.Responses,
					"agent.eth",
					gomock.Any(),
					&http_node_models.ExternalRequestHeaders{PaymentSignature: expectedSignatureHeader},
				).
				Return(nil, cb_errors.NewPaymentRequiredError("agent.eth", "retry-header")).
				Times(1)
		}

		t.Run("Should return the retry PaymentRequiredError", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(t, suite)

			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeRejected, gomock.Any())
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))

			_, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			var paymentRequired *cb_errors.PaymentRequiredError
			assert.ErrorAs(t, err, &paymentRequired)
			assert.Equal(t, "retry-header", paymentRequired.PaymentRequirements)
		})
	})

	t.Run("Given the node is rate limiting then answers", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, 3*time.Second)),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(completedResponseForTest(), nil),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())
		}

		t.Run("Should retry after the delay the node asked for", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				assert.NoError(t, err)
				assert.Equal(t, []string{"answer"}, completions)
				assert.Equal(t, 3*time.Second, time.Since(start))
			})
		})
	})

	t.Run("Given the node asks to wait exactly ten seconds then answers", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, 10*time.Second)),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(completedResponseForTest(), nil),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())
		}

		t.Run("Should retry after the delay the node asked for", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				assert.NoError(t, err)
				assert.Equal(t, []string{"answer"}, completions)
				assert.Equal(t, 10*time.Second, time.Since(start))
			})
		})
	})

	t.Run("Given the node is rate limiting without a delay then answers", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, 0)),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(completedResponseForTest(), nil),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())
		}

		t.Run("Should retry after one second", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				assert.NoError(t, err)
				assert.Equal(t, []string{"answer"}, completions)
				assert.Equal(t, time.Second, time.Since(start))
			})
		})
	})

	t.Run("Given the node keeps rate limiting", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, time.Second)),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, time.Second)),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, time.Second)),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
		}

		t.Run("Should return the unavailable error after three attempts", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				var unavailable *cb_errors.RemoteAgentUnavailableError
				assert.ErrorAs(t, err, &unavailable)
				assert.Nil(t, completions)
				assert.Equal(t, 2*time.Second, time.Since(start))
			})
		})
	})

	t.Run("Given the node asks to wait longer than ten seconds", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, 10*time.Second+time.Millisecond)),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
		}

		t.Run("Should return the unavailable error without waiting", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				var unavailable *cb_errors.RemoteAgentUnavailableError
				assert.ErrorAs(t, err, &unavailable)
				assert.Nil(t, completions)
				assert.Equal(t, time.Duration(0), time.Since(start))
			})
		})
	})

	t.Run("Given the node answers 502", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusBadGateway, 0)),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
		}

		t.Run("Should return the unavailable error without retrying", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				var unavailable *cb_errors.RemoteAgentUnavailableError
				assert.ErrorAs(t, err, &unavailable)
				assert.Nil(t, completions)
				assert.Equal(t, time.Duration(0), time.Since(start))
			})
		})
	})

	t.Run("Given the node answers 503", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusServiceUnavailable, 0)),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
		}

		t.Run("Should return the unavailable error without retrying", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				var unavailable *cb_errors.RemoteAgentUnavailableError
				assert.ErrorAs(t, err, &unavailable)
				assert.Nil(t, completions)
				assert.Equal(t, time.Duration(0), time.Since(start))
			})
		})
	})

	t.Run("Given the node answers 504", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, unavailableForTest(http.StatusGatewayTimeout, 0)),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
		}

		t.Run("Should return the unavailable error without retrying", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				start := time.Now()

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				var unavailable *cb_errors.RemoteAgentUnavailableError
				assert.ErrorAs(t, err, &unavailable)
				assert.Nil(t, completions)
				assert.Equal(t, time.Duration(0), time.Since(start))
			})
		})
	})

	t.Run("Given the node answers 502 to a pre-signed call", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{X402Version: models.X402Version}

			suite.x402PaymentChallengeCache.EXPECT().Read(gomock.Any(), "agent.eth").Return(challenge, true)
			suite.x402PaymentPayloadSigner.EXPECT().Sign(challenge, "0").Return(&models.X402PaymentPayload{}, nil)
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Not(gomock.Nil())).
				Return(nil, unavailableForTest(http.StatusBadGateway, 0))
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeSkipped, gomock.Any())
		}

		t.Run("Should record the payment as skipped", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			var unavailable *cb_errors.RemoteAgentUnavailableError
			assert.ErrorAs(t, err, &unavailable)
		})
	})

	t.Run("Given the node rate limits the signed call after a 402", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(t *testing.T, suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{X402Version: models.X402Version}

			suite.noCachedChallenge()
			suite.x402PaymentChallengeCache.EXPECT().Store(gomock.Any(), "agent.eth", gomock.Any())
			suite.x402PaymentPayloadSigner.EXPECT().Sign(gomock.Any(), "0").Return(&models.X402PaymentPayload{}, nil)
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, cb_errors.NewPaymentRequiredError("agent.eth", encodeForTest(t, challenge))),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Not(gomock.Nil())).
					Return(nil, unavailableForTest(http.StatusTooManyRequests, time.Second)),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Not(gomock.Nil())).
					Return(completedResponseForTest(), nil),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil())
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeSettled, gomock.Any())
		}

		t.Run("Should retry the signed call", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(t, suite)

				completions, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

				assert.NoError(t, err)
				assert.Equal(t, []string{"answer"}, completions)
			})
		})
	})

	t.Run("Given the node answers 502 to the signed call after a 402", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(t *testing.T, suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			challenge := &models.X402PaymentChallenge{X402Version: models.X402Version}

			suite.noCachedChallenge()
			suite.x402PaymentChallengeCache.EXPECT().Store(gomock.Any(), "agent.eth", gomock.Any())
			suite.x402PaymentPayloadSigner.EXPECT().Sign(gomock.Any(), "0").Return(&models.X402PaymentPayload{}, nil)
			gomock.InOrder(
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
					Return(nil, cb_errors.NewPaymentRequiredError("agent.eth", encodeForTest(t, challenge))),
				suite.mockRemoteAgentClient.EXPECT().
					CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Not(gomock.Nil())).
					Return(nil, unavailableForTest(http.StatusBadGateway, 0)),
			)
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordPaymentSettlement(gomock.Any(), metrics_interfaces.PaymentOutcomeSkipped, gomock.Any())
		}

		t.Run("Should record the payment as skipped", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
			initSuite(t, suite)

			_, err := suite.sut.RequestCompletions(context.Background(), card, turnForTest())

			var unavailable *cb_errors.RemoteAgentUnavailableError
			assert.ErrorAs(t, err, &unavailable)
		})
	})

	t.Run("Given the broker stops while waiting to retry", func(t *testing.T) {
		t.Parallel()

		card := &models.RemoteAgentCard{
			Protocol:  "6022",
			Endpoints: models.RemoteAgentCardEndpoints{Responses: "https://node.example/api/responses"},
		}

		initSuite := func(suite *WhenRequestingRemoteAgentCompletionsTestingSuite) {
			suite.noCachedChallenge()
			suite.mockRemoteAgentClient.EXPECT().
				CreateResponse(gomock.Any(), card.Endpoints.Responses, "agent.eth", gomock.Any(), gomock.Nil()).
				Return(nil, unavailableForTest(http.StatusTooManyRequests, 9*time.Second))
			suite.mockRemoteAgentMetrics.EXPECT().
				RecordInvocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(gomock.Nil()))
		}

		t.Run("Should stop waiting at once", func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				suite := WhenRequestingRemoteAgentCompletionsBeforeEach(t)
				initSuite(suite)
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				start := time.Now()

				_, err := suite.sut.RequestCompletions(ctx, card, turnForTest())

				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, 2*time.Second, time.Since(start))
			})
		})
	})
}
