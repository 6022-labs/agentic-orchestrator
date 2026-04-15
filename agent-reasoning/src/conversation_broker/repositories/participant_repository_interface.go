package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ParticipantRepositoryInterface interface {
	Save(entity *models.Participant) (*models.Participant, error)
	GetById(id uuid.UUID) (*models.Participant, error)
	GetByCollectionAddressAndCollectionTokenId(collectionAddress string, collectionTokenId string) (*models.Participant, error)
	GetNextMintedParticipantToProcess(processingIntervalSeconds uint64) (*models.Participant, error)
}
