package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/responses"
	"github.com/google/uuid"
)

type GetAgentContextsInterface interface {
	Get(request *requests.GetAgentContextsRequest) ([]*responses.GetAgentContextResponse, error)
}

type GetAgentContexts struct {
	agentDbClientFactory          factories_db_clients.AgentDatabaseClientFactoryInterface
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface
}

func NewGetAgentContexts(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface,
) *GetAgentContexts {
	return &GetAgentContexts{
		agentDbClientFactory:          agentDbClientFactory,
		agentContextRepositoryFactory: agentContextRepositoryFactory,
	}
}

func (g *GetAgentContexts) Get(request *requests.GetAgentContextsRequest) ([]*responses.GetAgentContextResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	request.Sanitize()

	agentDbClient, agentContextRepository, err := g.generateRepositories(request.AgentId)
	if err != nil {
		return nil, err
	}
	defer agentDbClient.Close()

	agentContexts, err := agentContextRepository.Search(repositories.AgentContextSearchParams{
		Search:        request.Search,
		SortBy:        request.SortBy.ToDatabaseField(),
		SortDirection: request.SortDirection.ToDatabaseDirection(),
	})
	if err != nil {
		return nil, err
	}

	return responses.GetAgentContextResponsesFromAgentContexts(agentContexts), nil
}

func (g *GetAgentContexts) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentContextRepositoryInterface, error) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentContextRepository := g.agentContextRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentContextRepository, nil
}
