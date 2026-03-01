package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/google/uuid"
)

type MCPConnectionSearchParams struct {
	AgentId          *uuid.UUID
	Search           *string
	Offset           *uint
	Limit            *uint
	OrderBy          *string
	OrderByDirection *string
}

type MCPConnectionRepositoryInterface interface {
	Save(mcpConnection *models.MCPConnection) (*models.MCPConnection, error)
	Delete(mcpConnection *models.MCPConnection) error
	DeleteById(id uuid.UUID) error
	GetById(id uuid.UUID) (*models.MCPConnection, error)
	GetByAgentId(agentId uuid.UUID) ([]*models.MCPConnection, error)
	Search(params MCPConnectionSearchParams) ([]*models.MCPConnection, error)
	ExistsByAgentIdAndUrl(agentId uuid.UUID, url string) (*bool, error)
	ExistsByAgentIdAndName(agentId uuid.UUID, name string) (*bool, error)
}
