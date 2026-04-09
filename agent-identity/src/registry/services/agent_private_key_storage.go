package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	common_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services/interfaces"
	"github.com/google/uuid"
)

type AgentPrivateKeyStorage struct {
	vaultSecretStorage common_interfaces.VaultSecretStorageInterface
}

func NewAgentPrivateKeyStorage(vaultSecretStorage common_interfaces.VaultSecretStorageInterface) *AgentPrivateKeyStorage {
	return &AgentPrivateKeyStorage{vaultSecretStorage: vaultSecretStorage}
}

func (s *AgentPrivateKeyStorage) Get(id uuid.UUID) (*string, error) {
	path := s.getAgentPrivateKeyPath(id)
	privateKey, err := s.vaultSecretStorage.Read(path)
	if err != nil {
		return nil, err
	}
	return &privateKey, nil
}

func (s *AgentPrivateKeyStorage) Store(id uuid.UUID, privateKey string) error {
	return s.vaultSecretStorage.Store(s.getAgentPrivateKeyPath(id), privateKey)
}

func (s *AgentPrivateKeyStorage) Delete(id uuid.UUID) error {
	return s.vaultSecretStorage.Delete(s.getAgentPrivateKeyPath(id))
}

func (s *AgentPrivateKeyStorage) getAgentPrivateKeyPath(id uuid.UUID) string {
	return fmt.Sprintf(models.AgentPrivateKeyVaultPathTemplate, id.String())
}
