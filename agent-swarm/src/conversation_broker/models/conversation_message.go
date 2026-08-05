package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConversationMessage is a single message in a Conversation, always attributed
// to a Participant. Author kind (AI agent vs human) is derived from the
// Participant's role on read.
type ConversationMessage struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`

	ConversationID uuid.UUID `gorm:"type:uuid;not null;index:idx_messages_conv_time,priority:1" json:"conversationId"`
	ParticipantID  uuid.UUID `gorm:"type:uuid;not null" json:"participantId"`

	Content string `gorm:"type:text;not null" json:"content"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;index:idx_messages_conv_time,priority:2" json:"createdAt"`

	Conversation *Conversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Participant  *Participant  `gorm:"foreignKey:ParticipantID;references:ID" json:"-"`
}

func (entity *ConversationMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	return
}
