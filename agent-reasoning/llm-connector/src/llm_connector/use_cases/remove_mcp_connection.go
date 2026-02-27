package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/use_cases/requests"
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
	if err != nil && err != gorm.ErrRecordNotFound {
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
