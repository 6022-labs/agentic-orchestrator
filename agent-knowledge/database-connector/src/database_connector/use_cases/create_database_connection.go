package use_cases

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateDatabaseConnectionInterface interface {
	Create(request *requests.CreateDatabaseConnectionRequest) (*responses.CreateDatabaseConnectionResponse, error)
}

type CreateDatabaseConnection struct {
	agentDatabaseCredentialsStorage   services.AgentDatabaseCredentialsStorageInterface
	databaseConnectionRepository      repositories.DatabaseConnectionRepositoryInterface
	agentDatabaseCredentialsGenerator services.AgentDatabaseCredentialsGeneratorInterface
}

func NewCreateDatabaseConnection(
	databaseConnectionRepository repositories.DatabaseConnectionRepositoryInterface,
	agentDatabaseCredentialsStorage services.AgentDatabaseCredentialsStorageInterface,
	agentDatabaseCredentialsGenerator services.AgentDatabaseCredentialsGeneratorInterface,
) *CreateDatabaseConnection {
	return &CreateDatabaseConnection{
		databaseConnectionRepository:      databaseConnectionRepository,
		agentDatabaseCredentialsStorage:   agentDatabaseCredentialsStorage,
		agentDatabaseCredentialsGenerator: agentDatabaseCredentialsGenerator,
	}
}

func (c *CreateDatabaseConnection) Create(request *requests.CreateDatabaseConnectionRequest) (*responses.CreateDatabaseConnectionResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	existingDatabaseConnection, err := c.databaseConnectionRepository.GetByAgentId(request.AgentID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error getting existing database connection: %w", err)
	}
	if existingDatabaseConnection != nil {
		return nil, fmt.Errorf("database connection already exists for agent %s", request.AgentID)
	}

	agentDatabaseUsername := strings.ReplaceAll(request.AgentID.String(), "-", "")

	databaseConnection := &models.DatabaseConnection{
		ID:          uuid.New(),
		Host:        request.Host,
		Port:        request.Port,
		Driver:      request.Driver,
		AgentID:     request.AgentID,
		Username:    agentDatabaseUsername,
		SetupStatus: enums.DatabaseConnectionSetupStatusConfigured,
	}

	createdDatabaseConnection, err := c.databaseConnectionRepository.Save(databaseConnection)
	if err != nil {
		return nil, fmt.Errorf("error saving database connection: %w", err)
	}

	// Generate client database user for the agent
	agentDatabasePassword := c.agentDatabaseCredentialsGenerator.Generate(databaseConnection.ID.String())

	// Store agent database credentials in vault
	err = c.agentDatabaseCredentialsStorage.Store(createdDatabaseConnection.AgentID, agentDatabasePassword)
	if err != nil {
		return nil, fmt.Errorf("error storing agent database credentials: %w", err)
	}

	return &responses.CreateDatabaseConnectionResponse{
		DatabaseName:     databaseConnection.AgentID.String(),
		DatabaseUsername: agentDatabaseUsername,
		DatabasePassword: agentDatabasePassword,
	}, nil
}
