package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	factories_db_clients "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/factories/db_clients"
	factories_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/factories/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
)

type SetupAgentDatabaseInterface interface {
	Setup(request *requests.SetupAgentDatabaseRequest) error
}

type SetupAgentDatabase struct {
	logger                                           *zap.Logger
	agentDbClientFactory                             factories_db_clients.AgentDatabaseClientFactoryInterface
	databaseConnectionRepository                     repositories.DatabaseConnectionRepositoryInterface
	agentContextRepositoryFactory                    factories_repositories.AgentContextRepositoryFactoryInterface
	agentPromptTemplateRepositoryFactory             factories_repositories.AgentPromptTemplateRepositoryFactoryInterface
	agentContextMessageRepositoryFactory             factories_repositories.AgentContextMessageRepositoryFactoryInterface
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface
}

func NewSetupAgentDatabase(
	logger *zap.Logger,
	agentDbClientFactory factories_db_clients.AgentDatabaseClientFactoryInterface,
	databaseConnectionRepository repositories.DatabaseConnectionRepositoryInterface,
	agentContextRepositoryFactory factories_repositories.AgentContextRepositoryFactoryInterface,
	agentPromptTemplateRepositoryFactory factories_repositories.AgentPromptTemplateRepositoryFactoryInterface,
	agentContextMessageRepositoryFactory factories_repositories.AgentContextMessageRepositoryFactoryInterface,
	agentContextConversationMappingRepositoryFactory factories_repositories.AgentContextConversationMappingRepositoryFactoryInterface,
) *SetupAgentDatabase {
	return &SetupAgentDatabase{
		logger:                                           logger,
		agentDbClientFactory:                             agentDbClientFactory,
		databaseConnectionRepository:                     databaseConnectionRepository,
		agentContextRepositoryFactory:                    agentContextRepositoryFactory,
		agentPromptTemplateRepositoryFactory:             agentPromptTemplateRepositoryFactory,
		agentContextMessageRepositoryFactory:             agentContextMessageRepositoryFactory,
		agentContextConversationMappingRepositoryFactory: agentContextConversationMappingRepositoryFactory,
	}
}

func (u *SetupAgentDatabase) Setup(request *requests.SetupAgentDatabaseRequest) (err error) {
	err = request.Validate()
	if err != nil {
		return
	}

	databaseConnection, err := u.databaseConnectionRepository.GetByAgentId(request.AgentId)
	if err != nil && err != gorm.ErrRecordNotFound {
		err = fmt.Errorf("failed to get database connection for agent %s: %w", request.AgentId, err)
		return
	}
	if databaseConnection == nil {
		err = fmt.Errorf("no database connection found for agent %s", request.AgentId)
		return
	}
	if databaseConnection.SetupStatus != enums.DatabaseConnectionSetupStatusConfigured {
		err = fmt.Errorf("database connection for agent %s is not in a state to run migrations", request.AgentId)
		return
	}

	databaseConnection.SetupStatus = enums.DatabaseConnectionSetupStatusEstablishing
	databaseConnection, err = u.databaseConnectionRepository.Save(databaseConnection)
	if err != nil {
		err = fmt.Errorf("failed to update database connection setup status: %w", err)
		return
	}

	finalStatus := enums.DatabaseConnectionSetupStatusConfigured
	defer func() {
		databaseConnection.SetupStatus = finalStatus
		_, deferredErr := u.databaseConnectionRepository.Save(databaseConnection)
		if deferredErr != nil {
			err = fmt.Errorf("failed to update database connection setup status to final status: %w", deferredErr)
		}
	}()

	agentDbClient, err := u.agentDbClientFactory.CreateFromDatabaseConnection(databaseConnection)
	if err != nil {
		err = fmt.Errorf("failed to create agent db client: %w", err)
		return
	}
	defer agentDbClient.Close()

	err = agentDbClient.Migrate()
	if err != nil {
		return
	}

	if request.CopyFromAgentId != nil {
		err = u.copy(*request.CopyFromAgentId, request.AgentId)
		if err != nil {
			err = fmt.Errorf("failed to copy data from agent %s to agent %s: %w", *request.CopyFromAgentId, request.AgentId, err)
			return
		}
	}

	finalStatus = enums.DatabaseConnectionSetupStatusEstablished

	// Return deferred error if any
	return
}

func (u *SetupAgentDatabase) copy(fromAgentId, toAgentId uuid.UUID) error {
	fromAgentDbClient, fromAgentContextRepository, fromAgentPromptTemplateRepository, fromAgentContextMessageRepository, fromAgentContextConversationMappingRepository, err := u.generateRepositories(fromAgentId)
	if err != nil {
		return err
	}
	defer fromAgentDbClient.Close()

	toAgentDbClient, toAgentContextRepository, toAgentPromptTemplateRepository, toAgentContextMessageRepository, toAgentContextConversationMappingRepository, err := u.generateRepositories(toAgentId)
	if err != nil {
		return err
	}
	defer toAgentDbClient.Close()

	fromAgentContextMessages, err := fromAgentContextMessageRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get from agent context messages: %w", err)
	}

	fromAgentPromptTemplates, err := fromAgentPromptTemplateRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get from agent prompt templates: %w", err)
	}

	fromAgentContextConversationMappings, err := fromAgentContextConversationMappingRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get from agent context conversation mappings: %w", err)
	}

	fromAgentContexts, err := fromAgentContextRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get from agent contexts: %w", err)
	}

	tx, err := toAgentDbClient.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	_, err = toAgentContextRepository.CreateMultipleUsingTransaction(tx, fromAgentContexts)
	if err != nil {
		toAgentDbClient.RollbackTransaction(tx)
		return fmt.Errorf("failed to save to agent contexts: %w", err)
	}

	_, err = toAgentPromptTemplateRepository.CreateMultipleUsingTransaction(tx, fromAgentPromptTemplates)
	if err != nil {
		toAgentDbClient.RollbackTransaction(tx)
		return fmt.Errorf("failed to save to agent prompt templates: %w", err)
	}

	_, err = toAgentContextConversationMappingRepository.CreateMultipleUsingTransaction(tx, fromAgentContextConversationMappings)
	if err != nil {
		toAgentDbClient.RollbackTransaction(tx)
		return fmt.Errorf("failed to save to agent context conversation mappings: %w", err)
	}

	_, err = toAgentContextMessageRepository.CreateMultipleUsingTransaction(tx, fromAgentContextMessages)
	if err != nil {
		toAgentDbClient.RollbackTransaction(tx)
		return fmt.Errorf("failed to save to agent context messages: %w", err)
	}

	err = toAgentDbClient.CommitTransaction(tx)
	if err != nil {
		toAgentDbClient.RollbackTransaction(tx)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (g *SetupAgentDatabase) generateRepositories(agentId uuid.UUID) (
	db_clients.AgentDatabaseClientInterface,
	repositories.AgentContextRepositoryInterface,
	repositories.AgentPromptTemplateRepositoryInterface,
	repositories.AgentContextMessageRepositoryInterface,
	repositories.AgentContextConversationMappingRepositoryInterface,
	error,
) {
	agentDbClient, err := g.agentDbClientFactory.Create(agentId)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to create agent db client: %w", err)
	}

	agentContextRepository := g.agentContextRepositoryFactory.Create(agentDbClient)
	agentPromptTemplateRepository := g.agentPromptTemplateRepositoryFactory.Create(agentDbClient)
	agentContextMessageRepository := g.agentContextMessageRepositoryFactory.Create(agentDbClient)
	agentContextConversationMappingRepository := g.agentContextConversationMappingRepositoryFactory.Create(agentDbClient)

	return agentDbClient, agentContextRepository, agentPromptTemplateRepository, agentContextMessageRepository, agentContextConversationMappingRepository, nil
}
