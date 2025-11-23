package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type DatabaseConnectionRepository struct {
	client db_clients.OrchestratorDatabaseClientInterface
}

func NewDatabaseConnectionRepository(client db_clients.OrchestratorDatabaseClientInterface) *DatabaseConnectionRepository {
	return &DatabaseConnectionRepository{
		client: client,
	}
}

func (r *DatabaseConnectionRepository) Save(entity *models.DatabaseConnection) (*models.DatabaseConnection, error) {
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *DatabaseConnectionRepository) GetByAgentId(agentId uuid.UUID) (*models.DatabaseConnection, error) {
	var databaseConnection models.DatabaseConnection
	if err := r.client.GetDB().Where("agent_id = ?", agentId).First(&databaseConnection).Error; err != nil {
		return nil, err
	}

	return &databaseConnection, nil
}

func (r *DatabaseConnectionRepository) GetSetupStatusByAgentId(agentId uuid.UUID) (*enums.DatabaseConnectionSetupStatus, error) {
	var setupStatus enums.DatabaseConnectionSetupStatus
	if err := r.client.GetDB().
		Model(&models.DatabaseConnection{}).
		Select("setup_status").
		Where("agent_id = ?", agentId).
		Limit(1).
		Scan(&setupStatus).
		Error; err != nil {
		return nil, err
	}

	if !enums.GetAllDatabaseConnectionSetupStatuses().Contains(setupStatus) {
		return nil, nil
	}

	return setupStatus.ToPointer(), nil
}

func (r *DatabaseConnectionRepository) GetMigratedMigrations(id uuid.UUID) ([]string, error) {
	var migrations []string

	if err := r.client.GetDB().
		Model(&models.DatabaseConnectionMigration{}).
		Select("name").
		Joins("JOIN database_connection_migration_connections dc_mc ON dc_mc.database_connection_migration_id = database_connection_migrations.id").
		Where("dc_mc.database_connection_id = ?", id).
		Scan(&migrations).
		Error; err != nil {
		return nil, err
	}

	return migrations, nil
}

func (r *DatabaseConnectionRepository) DeleteByAgentId(agentId uuid.UUID) error {
	var databaseConnection models.DatabaseConnection
	if err := r.client.GetDB().Where("agent_id = ?", agentId).Delete(&databaseConnection).Error; err != nil {
		return err
	}

	return nil
}
