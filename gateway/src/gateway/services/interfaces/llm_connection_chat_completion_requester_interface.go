package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models"
	"github.com/google/uuid"
)

type LLMConnectionChatCompletionRequesterInterface interface {
	ChatCompletion(agentId uuid.UUID, messages []*models.AgentContextMessage) (*string, error)
}
