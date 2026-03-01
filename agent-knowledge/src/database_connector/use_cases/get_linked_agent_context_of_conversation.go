package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetLinkedAgentContextOfConversationInterface interface {
	Get(request *requests.GetLinkedAgentContextOfConversationRequest) (*responses.GetLinkedAgentContextOfConversationResponse, error)
}

type GetLinkedAgentContextOfConversation struct {
	agentDbClientFactory                             factories_db_clients.AgentDatabaseClientFactoryInterface
	agentContextRepositoryFactory                    factories_repositories.AgentContextRepositoryFactoryInterface
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface
}

func NewGetLinkedAgentContextOfConversation(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface,
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface,
) *GetLinkedAgentContextOfConversation {
	return &GetLinkedAgentContextOfConversation{
		agentDbClientFactory:                             agentDbClientFactory,
		agentContextRepositoryFactory:                    agentContextRepositoryFactory,
		agentContextConversationMappingRepositoryFactory: agentContextConversationMappingRepositoryFactory,
	}
}

func (g *GetLinkedAgentContextOfConversation) Get(request *requests.GetLinkedAgentContextOfConversationRequest) (*responses.GetLinkedAgentContextOfConversationResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	agentDbClient, agentContextRepository, agentContextConversationMappingRepository, err := g.generateRepositories(request.AgentId)
	if err != nil {
		return nil, err
	}
	defer agentDbClient.Close()

	conversationMapping, err := agentContextConversationMappingRepository.GetByConversationId(request.ConversationId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if conversationMapping == nil {
		return &responses.GetLinkedAgentContextOfConversationResponse{
			Status: responses.GetLinkedAgentContextOfConversationResponseStatusNotFound,
		}, nil
	}

	agentContext, err := agentContextRepository.GetById(conversationMapping.AgentContextId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get agent context: %w", err)
	}
	if agentContext == nil {
		return nil, fmt.Errorf("agent context not found")
	}

	return &responses.GetLinkedAgentContextOfConversationResponse{
		Id:     agentContext.ID,
		Key:    agentContext.Key,
		Status: responses.GetLinkedAgentContextOfConversationResponseStatusFound,
	}, nil
}

func (g *GetLinkedAgentContextOfConversation) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentContextRepositoryInterface, repositories.AgentContextConversationMappingRepositoryInterface, error) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentContextRepository := g.agentContextRepositoryFactory.Create(agentDbClient)
	agentContextConversationMappingRepository := g.agentContextConversationMappingRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentContextRepository, agentContextConversationMappingRepository, nil
}
