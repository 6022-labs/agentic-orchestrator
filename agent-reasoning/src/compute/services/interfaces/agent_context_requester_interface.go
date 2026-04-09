package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/google/uuid"
)

type AgentContextRequesterInterface interface {
	CreateAgentContext(agentId uuid.UUID, key string) error
	GetAgentContextByKey(agentId uuid.UUID, key string) (*models.AgentContext, *bool, error)
	GetAgentContexts(agentId uuid.UUID, search, sortBy, sortDirection *string) ([]*models.AgentContext, error)
}
