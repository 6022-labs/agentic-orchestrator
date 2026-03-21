package models

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/google/uuid"
)

type ExternalSaveAgentPromptTemplateRequest struct {
	AgentId uuid.UUID                    `json:"agentId"`
	Content string                       `json:"content"`
	Key     enums.AgentPromptTemplateKey `json:"key"`
}
