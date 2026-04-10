package use_cases

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"gorm.io/gorm"
)

type RemoveMCPConnectionInterface interface {
	Remove(request *requests.IdRequest) error
}

type RemoveMCPConnection struct {
	mcpConnectionRepository             repositories.MCPConnectionRepositoryInterface
	mcpConnectionHeaderAuthValueStorage services.MCPConnectionHeaderAuthValueStorageInterface
}

func NewRemoveMCPConnection(
	mcpConnectionRepository repositories.MCPConnectionRepositoryInterface,
	mcpConnectionHeaderAuthValueStorage services.MCPConnectionHeaderAuthValueStorageInterface,
) *RemoveMCPConnection {
	return &RemoveMCPConnection{
		mcpConnectionRepository:             mcpConnectionRepository,
		mcpConnectionHeaderAuthValueStorage: mcpConnectionHeaderAuthValueStorage,
	}
}

func (d *RemoveMCPConnection) Remove(request *requests.IdRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	mcpConnection, err := d.mcpConnectionRepository.GetById(request.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to find MCPConnection: %w", err)
	}
	if mcpConnection == nil {
		return fmt.Errorf("MCPConnection with ID %s not found", request.ID)
	}

	if len(mcpConnection.HeaderAuthKeys) > 0 {
		err = d.mcpConnectionHeaderAuthValueStorage.Delete(mcpConnection)
		if err != nil {
			return fmt.Errorf("failed to delete secret for MCPConnection: %w", err)
		}
	}

	err = d.mcpConnectionRepository.Delete(mcpConnection)
	if err != nil {
		return fmt.Errorf("failed to delete MCPConnection: %w", err)
	}

	return nil
}
