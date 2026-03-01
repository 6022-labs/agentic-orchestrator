package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAgentContextByKeyInterface interface {
	Get(request *requests.GetAgentContextByKeyRequest) (*responses.GetAgentContextResponse, error)
}

type GetAgentContextByKey struct {
	agentDbClientFactory          factories_db_clients.AgentDatabaseClientFactoryInterface
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface
}

func NewGetAgentContextByKey(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface,
) *GetAgentContextByKey {
	return &GetAgentContextByKey{
		agentDbClientFactory:          agentDbClientFactory,
		agentContextRepositoryFactory: agentContextRepositoryFactory,
	}
}

func (g *GetAgentContextByKey) Get(request *requests.GetAgentContextByKeyRequest) (*responses.GetAgentContextResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	agentDbClient, agentContextRepository, err := g.generateRepositories(request.AgentId)
	if err != nil {
		return nil, err
	}
	defer agentDbClient.Close()

	agentContext, err := agentContextRepository.GetByKey(request.Key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if agentContext == nil {
		return nil, nil
	}

	return responses.GetAgentContextResponseFromAgentContext(agentContext), nil
}

func (g *GetAgentContextByKey) generateRepositories(agentId uuid.UUID) (dbClient db_clients.AgentDatabaseClientInterface, agentContextRepository repositories.AgentContextRepositoryInterface, err error) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, err
	}

	agentContextRepository = g.agentContextRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentContextRepository, nil
}
