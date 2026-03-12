package models

type SyncCheckpoint struct {
	Block uint64 `gorm:"not null"`
}

func (SyncCheckpoint) TableName() string {
	return "sync_checkpoints"
}
