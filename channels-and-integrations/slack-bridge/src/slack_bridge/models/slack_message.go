package models

import (
	"time"

	"github.com/google/uuid"
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
