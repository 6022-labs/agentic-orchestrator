package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/clients"
	http_node_models "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/models"
	"go.uber.org/zap"
)

type RemoteAgentCompletionsRequester struct {
	logger                    *zap.Logger
	client                    clients.RemoteAgentClientInterface
	x402PaymentPayloadSigner  service_interfaces.X402PaymentPayloadSignerInterface
	x402PaymentChallengeCache service_interfaces.X402PaymentChallengeCacheInterface
}

func NewRemoteAgentCompletionsRequester(
	logger *zap.Logger,
	client clients.RemoteAgentClientInterface,
	x402PaymentPayloadSigner service_interfaces.X402PaymentPayloadSignerInterface,
	x402PaymentChallengeCache service_interfaces.X402PaymentChallengeCacheInterface,
) *RemoteAgentCompletionsRequester {
	return &RemoteAgentCompletionsRequester{
		logger:                    logger,
		client:                    client,
		x402PaymentPayloadSigner:  x402PaymentPayloadSigner,
		x402PaymentChallengeCache: x402PaymentChallengeCache,
	}
}

func (r *RemoteAgentCompletionsRequester) RequestCompletions(
	card *models.RemoteAgentCard,
	input *models.RemoteAgentInput,
) ([]string, error) {
	ensDomain := input.Metadata["ensDomain"]

	if card.Endpoints.Responses != "" {
		return r.requestWithFreeCallRetry(ensDomain, func(paymentSignature string) ([]string, error) {
			return r.requestResponsesTurn(card.Endpoints.Responses, ensDomain, input, paymentSignature)
		})
	}
	if card.Endpoints.ChatCompletions != "" {
		return r.requestWithFreeCallRetry(ensDomain, func(paymentSignature string) ([]string, error) {
			return r.requestChatCompletionTurn(card.Endpoints.ChatCompletions, ensDomain, input, paymentSignature)
		})
	}
	return nil, fmt.Errorf("card advertises no endpoint")
}

// requestWithFreeCallRetry pre-signs from cache, then after a 402 caches the fresh challenge and retries once
// signed — unless the node just rejected these exact terms, where a retry would fail identically.
func (r *RemoteAgentCompletionsRequester) requestWithFreeCallRetry(
	ensDomain string,
	request func(paymentSignature string) ([]string, error),
) ([]string, error) {
	cachedChallenge, _ := r.x402PaymentChallengeCache.Read(ensDomain)

	completions, err := request(r.freeCallHeader(ensDomain, cachedChallenge))

	var paymentRequired *cb_errors.PaymentRequiredError
	if err == nil || !errors.As(err, &paymentRequired) {
		return completions, err
	}

	challenge, decodeErr := r.decodePaymentChallengeHeader(paymentRequired.PaymentRequirements)
	if decodeErr != nil {
		r.logger.Warn("free-call x402 retry skipped", zap.String("ensDomain", ensDomain), zap.Error(decodeErr))
		return nil, err
	}
	r.x402PaymentChallengeCache.Store(ensDomain, challenge)

	if reflect.DeepEqual(challenge, cachedChallenge) {
		return nil, err
	}

	signatureHeader, signErr := r.signFreeCallHeader(challenge)
	if signErr != nil {
		r.logger.Warn("free-call x402 retry skipped", zap.String("ensDomain", ensDomain), zap.Error(signErr))
		return nil, err
	}

	return request(signatureHeader)
}

func (r *RemoteAgentCompletionsRequester) freeCallHeader(
	ensDomain string,
	challenge *models.X402PaymentChallenge,
) string {
	if challenge == nil {
		return ""
	}

	signatureHeader, err := r.signFreeCallHeader(challenge)
	if err != nil {
		r.logger.Warn("pre-signed x402 free call skipped", zap.String("ensDomain", ensDomain), zap.Error(err))
		return ""
	}
	return signatureHeader
}

func (r *RemoteAgentCompletionsRequester) decodePaymentChallengeHeader(
	paymentRequiredHeader string,
) (*models.X402PaymentChallenge, error) {
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(paymentRequiredHeader))
	if err != nil {
		return nil, fmt.Errorf("PAYMENT-REQUIRED header is not valid base64: %w", err)
	}

	var paymentChallenge models.X402PaymentChallenge
	if err := json.Unmarshal(decoded, &paymentChallenge); err != nil {
		return nil, fmt.Errorf("PAYMENT-REQUIRED header is not valid JSON: %w", err)
	}
	return &paymentChallenge, nil
}

func (r *RemoteAgentCompletionsRequester) requestHeaders(
	paymentSignature string,
) *http_node_models.ExternalRequestHeaders {
	if paymentSignature == "" {
		return nil
	}
	return &http_node_models.ExternalRequestHeaders{PaymentSignature: paymentSignature}
}

func (r *RemoteAgentCompletionsRequester) signFreeCallHeader(
	paymentChallenge *models.X402PaymentChallenge,
) (string, error) {
	paymentPayload, err := r.x402PaymentPayloadSigner.Sign(paymentChallenge, "0")
	if err != nil {
		return "", err
	}

	encoded, err := json.Marshal(paymentPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payment payload: %w", err)
	}
	return base64.StdEncoding.EncodeToString(encoded), nil
}

func (r *RemoteAgentCompletionsRequester) requestResponsesTurn(
	endpointUrl, ensDomain string,
	turn *models.RemoteAgentInput,
	paymentSignature string,
) ([]string, error) {
	input := make([]*http_node_models.ExternalResponseInputMessage, len(turn.Input))
	for i, message := range turn.Input {
		input[i] = &http_node_models.ExternalResponseInputMessage{
			Type:    "message",
			Role:    message.Role.String(),
			Content: message.Content,
		}
	}

	response, err := r.client.CreateResponse(endpointUrl, ensDomain, &http_node_models.ExternalCreateResponseRequest{
		Instructions: turn.Instructions,
		Input:        input,
		Metadata:     turn.Metadata,
	}, r.requestHeaders(paymentSignature))
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("agent node returned error %s: %s", response.Error.Code, response.Error.Message)
	}

	completions := make([]string, 0, len(response.Output))
	for _, output := range response.Output {
		if output.Type != "message" {
			continue
		}
		for _, content := range output.Content {
			if content.Text != "" {
				completions = append(completions, content.Text)
			}
		}
	}

	if len(completions) == 0 {
		return nil, fmt.Errorf("agent node returned no message output (status %s)", response.Status)
	}

	return completions, nil
}

func (r *RemoteAgentCompletionsRequester) requestChatCompletionTurn(
	endpointUrl, ensDomain string,
	turn *models.RemoteAgentInput,
	paymentSignature string,
) ([]string, error) {
	messages := make([]*http_node_models.ExternalChatCompletionMessage, 0, len(turn.Input)+1)
	if turn.Instructions != "" {
		messages = append(messages, &http_node_models.ExternalChatCompletionMessage{
			Role:    "system",
			Content: turn.Instructions,
		})
	}
	for _, message := range turn.Input {
		messages = append(messages, &http_node_models.ExternalChatCompletionMessage{
			Role:    message.Role.String(),
			Content: message.Content,
		})
	}

	response, err := r.client.CreateChatCompletion(
		endpointUrl,
		ensDomain,
		&http_node_models.ExternalChatCompletionRequest{
			Messages: messages,
			Metadata: turn.Metadata,
		},
		r.requestHeaders(paymentSignature),
	)
	if err != nil {
		return nil, err
	}

	completions := make([]string, 0, len(response.Choices))
	for _, choice := range response.Choices {
		if choice.Message != nil && choice.Message.Content != "" {
			completions = append(completions, choice.Message.Content)
		}
	}

	if len(completions) == 0 {
		return nil, fmt.Errorf("agent node returned no chat completion choice")
	}

	return completions, nil
}
