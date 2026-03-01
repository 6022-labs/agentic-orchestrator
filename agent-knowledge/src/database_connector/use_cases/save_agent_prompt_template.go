package use_cases

import (
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

type SaveAgentPromptTemplateInterface interface {
	Save(request *requests.SaveAgentPromptTemplateRequest) error
}

type SaveAgentPromptTemplate struct {
	agentDbClientFactory                 factories_db_clients.AgentDatabaseClientFactoryInterface
	agentPromptTemplateRepositoryFactory factories_repositories.AgentPromptTemplateRepositoryFactoryInterface
}

func NewSaveAgentPromptTemplate(
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	agentPromptTemplateRepositoryFactory factories_repositories.AgentPromptTemplateRepositoryFactoryInterface,
) *SaveAgentPromptTemplate {
	return &SaveAgentPromptTemplate{
		agentDbClientFactory:                 agentDbClientFactory,
		agentPromptTemplateRepositoryFactory: agentPromptTemplateRepositoryFactory,
	}
}

func (u *SaveAgentPromptTemplate) Save(request *requests.SaveAgentPromptTemplateRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agentDbClient, agentPromptTemplateRepository, err := u.generateRepositories(request.AgentId)
	if err != nil {
		return err
	}
	defer agentDbClient.Close()

	agentPromptTemplate, err := agentPromptTemplateRepository.GetByKey(request.Key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	if agentPromptTemplate == nil {
		agentPromptTemplate = &models.AgentPromptTemplate{
			Key: request.Key,
		}
	}

	agentPromptTemplate.Content = request.Content

	_, err = agentPromptTemplateRepository.Save(agentPromptTemplate)
	if err != nil {
		return fmt.Errorf("failed to save agent prompt template: %w", err)
	}

	return nil
}

func (g *SaveAgentPromptTemplate) generateRepositories(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, repositories.AgentPromptTemplateRepositoryInterface, error) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentPromptTemplateRepository := g.agentPromptTemplateRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentPromptTemplateRepository, nil
}
