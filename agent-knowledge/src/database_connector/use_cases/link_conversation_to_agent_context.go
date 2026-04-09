package use_cases

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LinkConversationToAgentContextInterface interface {
	Link(request *requests.LinkConversationToAgentContextRequest) error
}

type LinkConversationToAgentContext struct {
	agentDbClientFactory                             factories_db_clients.AgentDatabaseClientFactoryInterface
	agentContextRepositoryFactory                    factories_repositories.AgentContextRepositoryFactoryInterface
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface
}

func NewLinkConversationToAgentContext(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface,
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface,
) *LinkConversationToAgentContext {
	return &LinkConversationToAgentContext{
		agentDbClientFactory:                             agentDbClientFactory,
		agentContextRepositoryFactory:                    agentContextRepositoryFactory,
		agentContextConversationMappingRepositoryFactory: agentContextConversationMappingRepositoryFactory,
	}
}

func (l *LinkConversationToAgentContext) Link(request *requests.LinkConversationToAgentContextRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agentDbClient, agentContextRepository, agentContextConversationMappingRepository, err := l.generateRepositories(request.AgentId)
	if err != nil {
		return fmt.Errorf("failed to generate repositories: %w", err)
	}
	defer agentDbClient.Close()

	// Check if the conversation mapping already exists
	existingMapping, err := agentContextConversationMappingRepository.GetByConversationId(request.ConversationId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get conversation mapping: %w", err)
	}
	if existingMapping != nil {
		return nil // Mapping already exists
	}

	agentContext, err := agentContextRepository.GetByKey(request.Key)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get agent context by key: %w", err)
	}
	if agentContext == nil {
		return fmt.Errorf("agent context not found for key: %s", request.Key)
	}

	newMapping := &models.AgentContextConversationMapping{
		ConversationId: request.ConversationId,
		AgentContextId: agentContext.ID,
	}

	_, err = agentContextConversationMappingRepository.Save(newMapping)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkConversationToAgentContext) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentContextRepositoryInterface, repositories.AgentContextConversationMappingRepositoryInterface, error) {
	agentDbClient, err := l.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentContextRepository := l.agentContextRepositoryFactory.Create(agentDbClient)
	agentContextConversationMappingRepository := l.agentContextConversationMappingRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentContextRepository, agentContextConversationMappingRepository, nil
}
