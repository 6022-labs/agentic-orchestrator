package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/google/uuid"
)

type AgentPrivateKeyStorage struct {
	vaultSecretStorage interfaces.VaultSecretStorageInterface
}

func NewAgentPrivateKeyStorage(vaultSecretStorage interfaces.VaultSecretStorageInterface) *AgentPrivateKeyStorage {
	return &AgentPrivateKeyStorage{
		vaultSecretStorage: vaultSecretStorage,
	}
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
	path := s.getAgentPrivateKeyPath(id)

	err := s.vaultSecretStorage.Store(path, privateKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentPrivateKeyStorage) Delete(id uuid.UUID) error {
	path := s.getAgentPrivateKeyPath(id)

	err := s.vaultSecretStorage.Delete(path)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentPrivateKeyStorage) getAgentPrivateKeyPath(id uuid.UUID) string {
	return fmt.Sprintf(models.AgentPrivateKeyVaultPathTemplate, id.String())
}
