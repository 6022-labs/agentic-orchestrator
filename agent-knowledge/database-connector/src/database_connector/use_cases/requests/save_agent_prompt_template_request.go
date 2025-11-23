package requests

import (
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type SaveAgentPromptTemplateRequest struct {
	AgentId uuid.UUID                    `json:"agentId"`
	Key     enums.AgentPromptTemplateKey `json:"key"`
	Content string                       `json:"content"`
}

func (r *SaveAgentPromptTemplateRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "Agent ID cannot be empty")
	}

	if len(strings.TrimSpace(r.Content)) == 0 {
		return errors.NewValidationError("content", "Content cannot be empty")
	}

	if !enums.GetAllAgentPromptTemplateKeys().Contains(r.Key) {
		return errors.NewValidationError("key", "Invalid value, key must be one of: "+enums.GetAllAgentPromptTemplateKeys().String())
	}

	return nil
}
