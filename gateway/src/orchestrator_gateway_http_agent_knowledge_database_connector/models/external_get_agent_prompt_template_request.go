package models

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/google/uuid"
)

type ExternalGetAgentPromptTemplateRequest struct {
	AgentId     uuid.UUID                    `json:"agentId"`
	TemplateKey enums.AgentPromptTemplateKey `json:"templateKey"`
}
