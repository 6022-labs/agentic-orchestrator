package services

import (
	"slices"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases/requests"
	"github.com/google/uuid"
)

// LLMChatCompletionUseCase defines what we need from the LLM chat completion use case
// without importing the use_cases package (avoids import cycle).
type LLMChatCompletionUseCase interface {
	Complete(request *requests.LLMConnectionChatCompletionRequest) (*string, error)
}

type InternalLLMConnectionChatCompletionRequester struct {
	llmConnectionChatCompletion LLMChatCompletionUseCase
}

func NewInternalLLMConnectionChatCompletionRequester(
	llmConnectionChatCompletion LLMChatCompletionUseCase,
) *InternalLLMConnectionChatCompletionRequester {
	return &InternalLLMConnectionChatCompletionRequester{
		llmConnectionChatCompletion: llmConnectionChatCompletion,
	}
}

func (r *InternalLLMConnectionChatCompletionRequester) ChatCompletion(
	agentId uuid.UUID,
	messages []*models.AgentContextMessage,
) (*string, error) {
	slices.SortFunc(messages, func(a, b *models.AgentContextMessage) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		} else if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}
		return 0
	})

	requestMessages := make([]*models.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		requestMessages[i] = &models.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	request := &requests.LLMConnectionChatCompletionRequest{
		AgentId:  agentId,
		Messages: requestMessages,
	}

	return r.llmConnectionChatCompletion.Complete(request)
}
