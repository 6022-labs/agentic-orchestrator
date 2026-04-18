package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SlackMessage is the dedup guard + bridge for Slack message ↔ broker message
// id. Composite PK keyed by (channel, ts) since ts alone is only unique within
// a channel/workspace.
type SlackMessage struct {
	SlackChannelId  string    `gorm:"primaryKey;type:text;not null" json:"slackChannelId"`
	SlackMessageTs  string    `gorm:"primaryKey;type:text;not null" json:"slackMessageTs"`
	BrokerMessageID uuid.UUID `gorm:"type:uuid;not null" json:"brokerMessageId"`
	CreatedAt       time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
}

func (entity *SlackMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	return
}
