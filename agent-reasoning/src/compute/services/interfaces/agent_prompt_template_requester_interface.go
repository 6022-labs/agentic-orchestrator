package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/google/uuid"
)

type AgentPromptTemplateRequesterInterface interface {
	GetAgentPromptTemplateByAgentIdAndKey(agentId uuid.UUID, key enums.AgentPromptTemplateKey) (*string, error)
}
