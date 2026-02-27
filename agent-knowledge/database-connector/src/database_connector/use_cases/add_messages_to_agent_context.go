package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddMessagesToAgentContextInterface interface {
	Add(requests *requests.AddMessagesToAgentContextRequest) error
}

type AddMessagesToAgentContext struct {
	agentDbClientFactory                             factories_db_clients.AgentDatabaseClientFactoryInterface
	agentContextRepositoryFactory                    factories_repositories.AgentContextRepositoryFactoryInterface
	agentContextMessageRepositoryFactory             factories_repositories.AgentContextMessageRepositoryFactoryInterface
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface
}

func NewAddMessagesToAgentContext(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface,
	agentContextMessageRepositoryFactory factories_repositories.AgentContextMessageRepositoryFactoryInterface,
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface,
) *AddMessagesToAgentContext {
	return &AddMessagesToAgentContext{
		agentDbClientFactory:                             agentDbClientFactory,
		agentContextRepositoryFactory:                    agentContextRepositoryFactory,
		agentContextMessageRepositoryFactory:             agentContextMessageRepositoryFactory,
		agentContextConversationMappingRepositoryFactory: agentContextConversationMappingRepositoryFactory,
	}
}

func (a *AddMessagesToAgentContext) Add(request *requests.AddMessagesToAgentContextRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agentDbClient, agentContextRepository, agentContextMessageRepository, agentContextConversationMappingRepository, err := a.generateRepositories(request.AgentId)
	if err != nil {
		return err
	}
	defer agentDbClient.Close()

	var agentContextId *uuid.UUID

	if request.ConversationId != nil {
		conversationMapping, err := agentContextConversationMappingRepository.GetByConversationId(*request.ConversationId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if conversationMapping == nil {
			return fmt.Errorf("conversation mapping not found")
		}

		agentContextId = &conversationMapping.AgentContextId
	}

	if request.ContextKey != nil && agentContextId == nil {
		agentContext, err := agentContextRepository.GetByKey(*request.ContextKey)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if agentContext == nil {
			return fmt.Errorf("agent context not found")
		}

		agentContextId = &agentContext.ID
	}

	if agentContextId == nil {
		return fmt.Errorf("agent context not found")
	}

	alreadyExistingMessages, err := agentContextMessageRepository.GetByAgentContextId(*agentContextId)
	if err != nil {
		return fmt.Errorf("failed to get existing messages: %w", err)
	}

	nonExistingMessages := make([]*models.AgentContextMessage, 0)
	for _, message := range request.Messages {
		exists := false

		potentialNewAgentContextMessage := message.ToAgentContextMessage(*agentContextId)
		err = potentialNewAgentContextMessage.ComputeHash()
		if err != nil {
			return fmt.Errorf("failed to compute hash: %w", err)
		}

		for _, existingMessage := range alreadyExistingMessages {
			if existingMessage.Hash == potentialNewAgentContextMessage.Hash {
				exists = true
				break
			}
		}

		if !exists {
			nonExistingMessages = append(nonExistingMessages, potentialNewAgentContextMessage)
		}
	}

	err = agentContextMessageRepository.SaveMultiple(nonExistingMessages)
	if err != nil {
		return fmt.Errorf("failed to save messages: %w", err)
	}

	return nil
}

func (a *AddMessagesToAgentContext) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentContextRepositoryInterface, repositories.AgentContextMessageRepositoryInterface, repositories.AgentContextConversationMappingRepositoryInterface, error) {
	dbClient, err := a.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create agent database client: %w", err)
	}

	agentContextRepository := a.agentContextRepositoryFactory.Create(dbClient)
	agentContextMessageRepository := a.agentContextMessageRepositoryFactory.Create(dbClient)
	agentContextConversationMappingRepository := a.agentContextConversationMappingRepositoryFactory.Create(dbClient)

	return dbClient, agentContextRepository, agentContextMessageRepository, agentContextConversationMappingRepository, nil
}
