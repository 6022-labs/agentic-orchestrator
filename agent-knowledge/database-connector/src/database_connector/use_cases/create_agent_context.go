package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateAgentContextInterface interface {
	Create(request *requests.CreateAgentContextRequest) error
}

type CreateAgentContext struct {
	agentDbClientFactory          factories_db_clients.AgentDatabaseClientFactoryInterface
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface
}

func NewCreateAgentContext(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface,
) *CreateAgentContext {
	return &CreateAgentContext{
		agentDbClientFactory:          agentDbClientFactory,
		agentContextRepositoryFactory: agentContextRepositoryFactory,
	}
}

func (c *CreateAgentContext) Create(request *requests.CreateAgentContextRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agentDbClient, agentContextRepository, err := c.generateRepositories(request.AgentId)
	if err != nil {
		return err
	}
	defer agentDbClient.Close()

	existingAgentContext, err := agentContextRepository.GetByKey(request.Key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get agent context by key: %w", err)
	}
	if existingAgentContext != nil {
		return fmt.Errorf("agent context with key %s already exists", request.Key)
	}

	_, err = agentContextRepository.Save(&models.AgentContext{
		Key: request.Key,
	})
	if err != nil {
		return fmt.Errorf("failed to create agent context: %w", err)
	}

	return nil
}

func (g *CreateAgentContext) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentContextRepositoryInterface, error) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentContextRepository := g.agentContextRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentContextRepository, nil
}
