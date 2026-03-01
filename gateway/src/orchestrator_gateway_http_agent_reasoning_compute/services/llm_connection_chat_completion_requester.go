package services

import (
	"fmt"
	"sort"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_reasoning_compute/clients"
	agent_compute_models "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_reasoning_compute/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type LLMConnectionCompletionRequester struct {
	client clients.AgentReasoningComputeClientInterface
}

func NewLLMConnectionCompletionRequester(
	client clients.AgentReasoningComputeClientInterface,
) *LLMConnectionCompletionRequester {
	return &LLMConnectionCompletionRequester{
		client: client,
	}
}

func (requester *LLMConnectionCompletionRequester) ChatCompletion(
	agentId uuid.UUID,
	messages []*models.AgentContextMessage,
) (*string, error) {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(messages[j].CreatedAt)
	})

	requestMessages := lo.Map(messages, func(message *models.AgentContextMessage, _ int) *agent_compute_models.ExternalChatCompletionMessage {
		return &agent_compute_models.ExternalChatCompletionMessage{
			Content: message.Content,
			Role:    message.Role.String(),
		}
	})

	request := &agent_compute_models.ExternalChatCompletionRequest{
		AgentId:  agentId,
		Messages: requestMessages,
	}

	response, err := requester.client.ChatCompletion(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat completion: %w", err)
	}

	return response, nil
}
