package models

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/google/uuid"
)

type ExternalGetAgentPromptTemplateRequest struct {
	AgentId     uuid.UUID                    `json:"agentId"`
	TemplateKey enums.AgentPromptTemplateKey `json:"templateKey"`
}
