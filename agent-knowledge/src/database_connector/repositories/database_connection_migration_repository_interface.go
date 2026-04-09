package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
)

type DatabaseConnectionMigrationRepositoryInterface interface {
	Save(entity *models.DatabaseConnectionMigration) (*models.DatabaseConnectionMigration, error)
	GetAll() ([]*models.DatabaseConnectionMigration, error)
}
