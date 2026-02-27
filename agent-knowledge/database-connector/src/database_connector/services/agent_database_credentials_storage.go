package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/google/uuid"
)

type AgentDatabaseCredentialsStorageInterface interface {
	Get(id uuid.UUID) (*string, error)
	Store(id uuid.UUID, credentials string) error
	Delete(id uuid.UUID) error
}

type AgentDatabaseCredentialsStorage struct {
	vaultSecretStorage VaultSecretStorageInterface
}

func NewAgentDatabaseCredentialsStorage(vaultSecretStorage VaultSecretStorageInterface) *AgentDatabaseCredentialsStorage {
	return &AgentDatabaseCredentialsStorage{
		vaultSecretStorage: vaultSecretStorage,
	}
}

func (s *AgentDatabaseCredentialsStorage) Get(id uuid.UUID) (*string, error) {
	path := s.getAgentDatabaseCredentialsPath(id)

	credentials, err := s.vaultSecretStorage.Read(path)
	if err != nil {
		return nil, err
	}

	return &credentials, nil
}

func (s *AgentDatabaseCredentialsStorage) Store(id uuid.UUID, credentials string) error {
	path := s.getAgentDatabaseCredentialsPath(id)

	err := s.vaultSecretStorage.Store(path, credentials)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentDatabaseCredentialsStorage) Delete(id uuid.UUID) error {
	path := s.getAgentDatabaseCredentialsPath(id)

	err := s.vaultSecretStorage.Delete(path)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentDatabaseCredentialsStorage) getAgentDatabaseCredentialsPath(id uuid.UUID) string {
	return fmt.Sprintf(models.DatabaseCredentialsVaultPathTemplate, id.String())
}
