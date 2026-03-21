package models

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/google/uuid"
)

type ExternalGetAgentPromptTemplateRequest struct {
	AgentId     uuid.UUID                    `json:"agentId"`
	TemplateKey enums.AgentPromptTemplateKey `json:"templateKey"`
}
