package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/google/uuid"
)

type LLMConnectionChatCompletionRequest struct {
	AgentId  uuid.UUID                       `json:"agentId"`
	Messages []*models.ChatCompletionMessage `json:"messages"`
}

func (r *LLMConnectionChatCompletionRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "is required")
	}

	if len(r.Messages) == 0 {
		return errors.NewValidationError("messages", "at least one message is required")
	}

	for _, message := range r.Messages {
		if message == nil {
			return errors.NewValidationError("messages", "contains nil message")
		}
		if len(strings.TrimSpace(message.Role)) == 0 {
			return errors.NewValidationError("messages.role", "is required")
		}
		if len(strings.TrimSpace(message.Content)) == 0 {
			return errors.NewValidationError("messages.content", "is required")
		}
	}

	return nil
}
