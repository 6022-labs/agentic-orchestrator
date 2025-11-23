package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
)

type DatabaseConnectionMigrationRepository struct {
	client db_clients.OrchestratorDatabaseClientInterface
}

func NewDatabaseConnectionMigrationRepository(client db_clients.OrchestratorDatabaseClientInterface) *DatabaseConnectionMigrationRepository {
	return &DatabaseConnectionMigrationRepository{
		client: client,
	}
}

func (r *DatabaseConnectionMigrationRepository) Save(entity *models.DatabaseConnectionMigration) (*models.DatabaseConnectionMigration, error) {
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *DatabaseConnectionMigrationRepository) GetAll() ([]*models.DatabaseConnectionMigration, error) {
	var entities []*models.DatabaseConnectionMigration
	if err := r.client.GetDB().Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}
