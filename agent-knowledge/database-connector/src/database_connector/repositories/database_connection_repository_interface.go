package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type DatabaseConnectionRepositoryInterface interface {
	Save(entity *models.DatabaseConnection) (*models.DatabaseConnection, error)
	GetByAgentId(agentId uuid.UUID) (*models.DatabaseConnection, error)
	GetSetupStatusByAgentId(agentId uuid.UUID) (*enums.DatabaseConnectionSetupStatus, error)
	GetMigratedMigrations(id uuid.UUID) ([]string, error)
	DeleteByAgentId(agentId uuid.UUID) error
}
