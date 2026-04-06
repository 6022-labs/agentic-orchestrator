package models

import "github.com/google/uuid"

type SyncCheckpoint struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Block uint64
}

func (SyncCheckpoint) TableName() string {
	return "sync_checkpoints"
}
