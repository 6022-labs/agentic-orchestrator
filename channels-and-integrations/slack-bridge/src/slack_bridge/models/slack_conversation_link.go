package models

import (
	"time"

	"github.com/google/uuid"
)

// SlackConversationLink maps a Slack (channel, thread) pair to the broker
// conversation id it was registered as. Composite PK keyed by the Slack pair
// so ingest-side lookups are O(1); unique index on broker id so callback-side
// lookups from broker → Slack are also O(1).
type SlackConversationLink struct {
	SlackChannelId       string    `gorm:"primaryKey;type:text;not null" json:"slackChannelId"`
	SlackThreadTs        string    `gorm:"primaryKey;type:text;not null" json:"slackThreadTs"`
	BrokerConversationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_slack_conversation_link_broker" json:"brokerConversationId"`
	CreatedAt            time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
}
