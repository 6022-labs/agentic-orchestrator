package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Conversation is a generic, bridge-agnostic conversation. Bridge-side
// identifiers (Slack channel_id + thread_ts etc.) are kept by the originating
// bridge in its own linking table; the broker only holds an internal UUID.
type Conversation struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`

	BridgeID    uuid.UUID `gorm:"type:uuid;not null;index" json:"bridgeId"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`

	Status        enums.ConversationStatus `gorm:"type:text;not null;default:'active';index:idx_conversations_status_updated,priority:1" json:"status"`
	ConcludedAt   *time.Time               `gorm:"type:timestamptz" json:"concludedAt,omitempty"`
	UnconcludedAt *time.Time               `gorm:"type:timestamptz" json:"unconcludedAt,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;index:idx_conversations_status_updated,priority:2" json:"updatedAt"`

	Bridge *Bridge `gorm:"foreignKey:BridgeID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (entity *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	return
}

func (entity *Conversation) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
