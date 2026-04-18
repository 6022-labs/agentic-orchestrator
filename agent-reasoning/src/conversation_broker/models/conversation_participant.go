package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConversationParticipant attaches a Participant (agent or user) to a
// Conversation. Composite primary key.
type ConversationParticipant struct {
	ConversationID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"conversationId"`
	ParticipantID  uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"participantId"`

	JoinedAt time.Time `gorm:"type:timestamptz;not null" json:"joinedAt"`

	Conversation *Conversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Participant  *Participant  `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (entity *ConversationParticipant) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.JoinedAt.IsZero() {
		entity.JoinedAt = time.Now()
	}
	return
}
