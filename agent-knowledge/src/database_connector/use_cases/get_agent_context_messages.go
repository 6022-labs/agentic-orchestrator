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
)

type GetAgentContextMessagesInterface interface {
	Get(request *requests.GetAgentContextMessagesRequest) ([]*responses.GetAgentContextMessageResponse, error)
}

type GetAgentContextMessages struct {
	agentDbClientFactory                 factories_db_clients.AgentDatabaseClientFactoryInterface
	agentContextMessageRepositoryFactory factories_repositories.AgentContextMessageRepositoryFactoryInterface
}

func NewGetAgentContextMessages(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentContextMessageRepositoryFactory factories_repositories.AgentContextMessageRepositoryFactoryInterface,
) *GetAgentContextMessages {
	return &GetAgentContextMessages{
		agentDbClientFactory:                 agentDbClientFactory,
		agentContextMessageRepositoryFactory: agentContextMessageRepositoryFactory,
	}
}

func (g *GetAgentContextMessages) Get(request *requests.GetAgentContextMessagesRequest) ([]*responses.GetAgentContextMessageResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	agentDbClient, agentContextMessageRepository, err := g.generateRepositories(request.AgentId)
	if err != nil {
		return nil, err
	}
	defer agentDbClient.Close()

	messages, err := agentContextMessageRepository.GetByAgentContextKey(request.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	response := make([]*responses.GetAgentContextMessageResponse, len(messages))

	for i, message := range messages {
		response[i] = &responses.GetAgentContextMessageResponse{
			Hash:           message.Hash,
			Role:           message.Role,
			Content:        message.Content,
			Author:         message.Author,
			Platform:       message.Platform,
			CreatedAt:      message.CreatedAt,
			ConversationId: message.ConversationId,
			AgentContextId: message.AgentContextId,
		}
	}

	return response, nil
}

func (g *GetAgentContextMessages) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentContextMessageRepositoryInterface, error) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentContextMessageRepository := g.agentContextMessageRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentContextMessageRepository, nil
}
