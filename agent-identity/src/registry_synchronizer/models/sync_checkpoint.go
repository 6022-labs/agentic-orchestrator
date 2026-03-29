package models

type SyncCheckpoint struct {
	Block uint64 `gorm:"primaryKey"`
}

func (SyncCheckpoint) TableName() string {
	return "sync_checkpoints"
}
