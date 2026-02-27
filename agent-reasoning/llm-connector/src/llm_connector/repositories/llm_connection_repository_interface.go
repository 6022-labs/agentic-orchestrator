package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models/enums"
	"github.com/google/uuid"
)

type LLMConnectionRepositoryInterface interface {
	Save(entity *models.LLMConnection) (*models.LLMConnection, error)
	GetByAgentId(agentId uuid.UUID) (*models.LLMConnection, error)
	GetSetupStatusByAgentId(agentId uuid.UUID) (*enums.LLMConnectionSetupStatus, error)
	DeleteByAgentId(agentId uuid.UUID) error
}
