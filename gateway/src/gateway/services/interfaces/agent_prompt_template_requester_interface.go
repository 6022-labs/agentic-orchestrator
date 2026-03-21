package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/google/uuid"
)

type AgentPromptTemplateRequesterInterface interface {
	SaveAgentPromptTemplate(agentId uuid.UUID, key enums.AgentPromptTemplateKey, content string) error
	GetAgentPromptTemplateByAgentIdAndKey(agentId uuid.UUID, key enums.AgentPromptTemplateKey) (*string, error)
}
