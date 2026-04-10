package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
)

type MCPConnectionHeaderAuthValueStorageInterface interface {
	Get(mcpConnection *models.MCPConnection) (map[string]string, error)
	Store(mcpConnection *models.MCPConnection, headerAuth map[string]string) error
	Delete(mcpConnection *models.MCPConnection) error
}

type MCPConnectionHeaderAuthValueStorage struct {
	vaultSecretStorage interfaces.VaultSecretStorageInterface
}

func NewMCPConnectionHeaderAuthValueStorage(vaultSecretStorage interfaces.VaultSecretStorageInterface) *MCPConnectionHeaderAuthValueStorage {
	return &MCPConnectionHeaderAuthValueStorage{
		vaultSecretStorage: vaultSecretStorage,
	}
}

func (storage *MCPConnectionHeaderAuthValueStorage) Get(mcpConnection *models.MCPConnection) (map[string]string, error) {
	if len(mcpConnection.HeaderAuthKeys) == 0 {
		return make(map[string]string), nil
	}

	secretPath := storage.getMCPConnectionHeaderAuthPath(mcpConnection)
	return storage.vaultSecretStorage.Read(secretPath)
}

func (storage *MCPConnectionHeaderAuthValueStorage) Store(mcpConnection *models.MCPConnection, headerAuth map[string]string) error {
	secretPath := storage.getMCPConnectionHeaderAuthPath(mcpConnection)
	return storage.vaultSecretStorage.Store(secretPath, headerAuth)
}

func (s *MCPConnectionHeaderAuthValueStorage) Delete(mcpConnection *models.MCPConnection) error {
	path := s.getMCPConnectionHeaderAuthPath(mcpConnection)
	return s.vaultSecretStorage.Delete(path)
}

func (s *MCPConnectionHeaderAuthValueStorage) getMCPConnectionHeaderAuthPath(mcpConnection *models.MCPConnection) string {
	return fmt.Sprintf(models.MCPAuthHeaderSecretVaultPathTemplate, mcpConnection.AgentID, mcpConnection.Name)
}
