package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/use_cases/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAgentPromptTemplateInterface interface {
	Get(request *requests.GetAgentPromptTemplateRequest) (*responses.GetAgentPromptTemplateResponse, error)
}

type GetAgentPromptTemplate struct {
	agentDatabaseCredentialsStorage      services.AgentDatabaseCredentialsStorageInterface
	databaseConnectionRepository         repositories.DatabaseConnectionRepositoryInterface
	agentDbClientFactory                 factories_db_clients.AgentDatabaseClientFactoryInterface
	agentPromptTemplateRepositoryFactory factories_repositories.AgentPromptTemplateRepositoryFactoryInterface
}

func NewGetAgentPromptTemplate(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	databaseConnectionRepository repositories.DatabaseConnectionRepositoryInterface,
	agentDatabaseCredentialsStorage services.AgentDatabaseCredentialsStorageInterface,
	agentPromptTemplateRepositoryFactory factories_repositories.AgentPromptTemplateRepositoryFactoryInterface,
) *GetAgentPromptTemplate {
	return &GetAgentPromptTemplate{
		agentDbClientFactory:                 agentDbClientFactory,
		databaseConnectionRepository:         databaseConnectionRepository,
		agentDatabaseCredentialsStorage:      agentDatabaseCredentialsStorage,
		agentPromptTemplateRepositoryFactory: agentPromptTemplateRepositoryFactory,
	}
}

func (u *GetAgentPromptTemplate) Get(request *requests.GetAgentPromptTemplateRequest) (*responses.GetAgentPromptTemplateResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	agentDbClient, agentPromptTemplateRepository, err := u.generateRepositories(request.AgentId)
	if err != nil {
		return nil, err
	}
	defer agentDbClient.Close()

	agentPromptTemplate, err := agentPromptTemplateRepository.GetByKey(request.TemplateKey)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	if agentPromptTemplate == nil {
		return &responses.GetAgentPromptTemplateResponse{
			PromptTemplate: "",
		}, nil
	}

	return &responses.GetAgentPromptTemplateResponse{
		PromptTemplate: agentPromptTemplate.Content,
	}, nil
}

func (g *GetAgentPromptTemplate) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentPromptTemplateRepositoryInterface, error) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentPromptTemplateRepository := g.agentPromptTemplateRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentPromptTemplateRepository, nil
}
