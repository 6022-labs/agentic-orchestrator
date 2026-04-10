package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/samber/lo"
)

type CreateMCPConnectionInterface interface {
	Create(request *requests.CreateMCPConnectionRequest) error
}

type CreateMCPConnection struct {
	mcpConnectionRepository             repositories.MCPConnectionRepositoryInterface
	mcpConnectionHeaderAuthValueStorage services.MCPConnectionHeaderAuthValueStorageInterface
}

func NewCreateMCPConnection(
	mcpConnectionRepository repositories.MCPConnectionRepositoryInterface,
	mcpConnectionHeaderAuthValueStorage services.MCPConnectionHeaderAuthValueStorageInterface,
) *CreateMCPConnection {
	return &CreateMCPConnection{
		mcpConnectionRepository:             mcpConnectionRepository,
		mcpConnectionHeaderAuthValueStorage: mcpConnectionHeaderAuthValueStorage,
	}
}

func (c *CreateMCPConnection) Create(request *requests.CreateMCPConnectionRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	nameAlreadyExists, err := c.mcpConnectionRepository.ExistsByAgentIdAndName(request.AgentId, request.Name)
	if err != nil || nameAlreadyExists == nil {
		return fmt.Errorf("failed to check if MCPConnection already exists by name: %w", err)
	}
	if *nameAlreadyExists {
		return errors.NewValidationError("name", "MCPConnection name already exists for this agent")
	}

	urlAlreadyExists, err := c.mcpConnectionRepository.ExistsByAgentIdAndUrl(request.AgentId, request.Url)
	if err != nil || urlAlreadyExists == nil {
		return fmt.Errorf("failed to check if MCPConnection already exists by url: %w", err)
	}
	if *urlAlreadyExists {
		return errors.NewValidationError("url", "MCPConnection URL already exists for this agent")
	}

	mcpConnection := &models.MCPConnection{
		Url:     request.Url,
		Name:    request.Name,
		AgentID: request.AgentId,
		HeaderAuthKeys: lo.MapToSlice(request.HeaderAuth, func(key string, _ string) models.MCPConnectionHeaderAuthKey {
			return models.MCPConnectionHeaderAuthKey(key)
		}),
	}

	if len(request.HeaderAuth) > 0 {
		err = c.mcpConnectionHeaderAuthValueStorage.Store(mcpConnection, request.HeaderAuth)
		if err != nil {
			return fmt.Errorf("failed to store header auth values for MCPConnection: %w", err)
		}
	}

	_, err = c.mcpConnectionRepository.Save(mcpConnection)
	if err != nil {
		return fmt.Errorf("failed to create MCPConnection: %w", err)
	}

	return nil
}
