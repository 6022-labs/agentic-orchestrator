package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
)

type DatabaseConnectionMigrationRepositoryInterface interface {
	Save(entity *models.DatabaseConnectionMigration) (*models.DatabaseConnectionMigration, error)
	GetAll() ([]*models.DatabaseConnectionMigration, error)
}
