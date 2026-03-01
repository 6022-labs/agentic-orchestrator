package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/services/interfaces"
	"github.com/google/uuid"
)

type AgentLLMSecretStorageInterface interface {
	Get(id uuid.UUID) (*string, error)
	Store(id uuid.UUID, credentials string) error
	Delete(id uuid.UUID) error
}

type AgentLLMSecretStorage struct {
	vaultSecretStorage interfaces.VaultSecretStorageInterface
}

func NewAgentLLMSecretStorage(vaultSecretStorage interfaces.VaultSecretStorageInterface) *AgentLLMSecretStorage {
	return &AgentLLMSecretStorage{
		vaultSecretStorage: vaultSecretStorage,
	}
}

func (s *AgentLLMSecretStorage) Get(id uuid.UUID) (*string, error) {
	path := s.getAgentLLMSecretPath(id)

	credentials, err := s.vaultSecretStorage.ReadKey(path, "value")
	if err != nil {
		return nil, err
	}

	return &credentials, nil
}

func (s *AgentLLMSecretStorage) Store(id uuid.UUID, secret string) error {
	path := s.getAgentLLMSecretPath(id)
	return s.vaultSecretStorage.StoreKey(path, "value", secret)
}

func (s *AgentLLMSecretStorage) Delete(id uuid.UUID) error {
	path := s.getAgentLLMSecretPath(id)
	return s.vaultSecretStorage.Delete(path)
}

func (s *AgentLLMSecretStorage) getAgentLLMSecretPath(id uuid.UUID) string {
	return fmt.Sprintf(models.LLMSecretVaultPathTemplate, id.String())
}
