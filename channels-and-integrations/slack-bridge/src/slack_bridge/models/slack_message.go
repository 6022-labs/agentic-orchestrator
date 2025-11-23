package models

import (
	"strconv"
	"time"
)

type SlackMessage struct {
	Id          string `gorm:"primaryKey;type:string;not null" json:"id"`
	IsFromUser  bool   `gorm:"type:boolean;not null" json:"is_from_user"`
	IsFromAgent bool   `gorm:"type:boolean;not null" json:"is_from_agent"`
	Content     string `gorm:"type:string;not null" json:"content"`

	ThreadId  string        `gorm:"type:string;not null" json:"thread_id"`
	ChannelId string        `gorm:"type:string;not null" json:"channel_id"`
	Thread    *SlackThread  `gorm:"foreignKey:ThreadId,ChannelId;references:Id,ChannelId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Channel   *SlackChannel `gorm:"foreignKey:ChannelId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// With user paring to human agents, a message can have both an author_user_id and an author_collection_address/author_collection_token_id
	AuthorUserId *string    `gorm:"type:string;index" json:"author_user_id"` // Filled if the author is a Slack user
	AuthorUser   *SlackUser `gorm:"foreignKey:AuthorUserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	AuthorCollectionAddress *string     `gorm:"type:string" json:"author_collection_address"`  // Filled if the author is an agent
	AuthorCollectionTokenId *uint64     `gorm:"type:bigint" json:"author_collection_token_id"` // Filled if the author is an agent
	AuthorAgent             *SlackAgent `gorm:"foreignKey:AuthorCollectionAddress,AuthorCollectionTokenId;references:CollectionAddress,CollectionTokenId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (m *SlackMessage) CreatedAt() time.Time {
	// Convert the ID to a float64 to extract the timestamp
	timestampFloat, err := strconv.ParseFloat(m.Id, 64)
	if err != nil {
		return time.Time{} // Return zero time on error
	}
	seconds := int64(timestampFloat)
	nanos := int64((timestampFloat - float64(seconds)) * 1e9)
	return time.Unix(seconds, nanos)
}
