package interfaces

import "github.com/google/uuid"

type AgentPrivateKeyStorageInterface interface {
	Get(id uuid.UUID) (*string, error)
	Store(id uuid.UUID, privateKey string) error
	Delete(id uuid.UUID) error
}
