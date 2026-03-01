package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
)

type RemoveDatabaseConnectionInterface interface {
	Remove(request *requests.AgentIdRequest) error
}

type RemoveDatabaseConnection struct {
	agentDatabaseCredentialsStorage services.AgentDatabaseCredentialsStorageInterface
	databaseConnectionRepository    repositories.DatabaseConnectionRepositoryInterface
}

func NewRemoveDatabaseConnection(
	databaseConnectionRepository repositories.DatabaseConnectionRepositoryInterface,
	agentDatabaseCredentialsStorage services.AgentDatabaseCredentialsStorageInterface,
) *RemoveDatabaseConnection {
	return &RemoveDatabaseConnection{
		databaseConnectionRepository:    databaseConnectionRepository,
		agentDatabaseCredentialsStorage: agentDatabaseCredentialsStorage,
	}
}

func (r *RemoveDatabaseConnection) Remove(request *requests.AgentIdRequest) error {
	err := request.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}

	err = r.agentDatabaseCredentialsStorage.Delete(request.AgentId)
	if err != nil {
		return fmt.Errorf("failed to delete agent database credentials: %w", err)
	}

	err = r.databaseConnectionRepository.DeleteByAgentId(request.AgentId)
	if err != nil {
		return fmt.Errorf("failed to delete database connection: %w", err)
	}

	return nil
}
