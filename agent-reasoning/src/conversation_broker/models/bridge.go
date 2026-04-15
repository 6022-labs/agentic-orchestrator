package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bridge represents a connected bridge instance (e.g. slack-bridge-prod).
// The broker posts completion callbacks to bridges.callback_url.
type Bridge struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`

	Kind        enums.BridgeKind `gorm:"type:text;not null;index:idx_bridges_kind" json:"kind"`
	CallbackUrl string           `gorm:"type:text;not null" json:"callbackUrl"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updatedAt"`
}

func (entity *Bridge) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	return
}

func (entity *Bridge) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
