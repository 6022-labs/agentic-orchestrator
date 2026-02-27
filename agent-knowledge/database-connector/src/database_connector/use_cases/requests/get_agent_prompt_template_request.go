package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type GetAgentPromptTemplateRequest struct {
	AgentId     uuid.UUID                    `json:"agentId"`
	TemplateKey enums.AgentPromptTemplateKey `json:"templateKey"`
}

func (r *GetAgentPromptTemplateRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "Agent ID cannot be empty")
	}

	if !enums.GetAllAgentPromptTemplateKeys().Contains(r.TemplateKey) {
		return errors.NewValidationError("templateKey", "Invalid template key, must be one of "+enums.GetAllAgentPromptTemplateKeys().String())
	}

	return nil
}
