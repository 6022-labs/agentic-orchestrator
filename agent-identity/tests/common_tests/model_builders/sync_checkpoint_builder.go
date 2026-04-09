package model_builders

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"github.com/google/uuid"
)

type SyncCheckpointBuilder struct {
	id    uuid.UUID
	block uint64
}

func NewSyncCheckpointBuilder() *SyncCheckpointBuilder {
	return &SyncCheckpointBuilder{
		id:    uuid.New(),
		block: 0,
	}
}

func (b *SyncCheckpointBuilder) WithID(id uuid.UUID) *SyncCheckpointBuilder {
	b.id = id
	return b
}

func (b *SyncCheckpointBuilder) WithBlock(block uint64) *SyncCheckpointBuilder {
	b.block = block
	return b
}

func (b *SyncCheckpointBuilder) Build() *models.SyncCheckpoint {
	return &models.SyncCheckpoint{
		ID:    b.id,
		Block: b.block,
	}
}
