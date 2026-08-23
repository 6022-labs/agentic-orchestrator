package services_test

import (
	"context"

	"encoding/base64"
	"encoding/json"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
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
}
