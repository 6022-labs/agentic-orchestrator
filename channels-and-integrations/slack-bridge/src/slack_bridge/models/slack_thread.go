package models

import (
	"strconv"
	"time"
)

type SlackThread struct {
	Id            string     `gorm:"primaryKey;type:string;not null" json:"id"`
	IsConcluded   bool       `gorm:"type:boolean;not null" json:"is_concluded"`
	UnconcludedAt *time.Time `gorm:"type:timestamp" json:"unconcluded_at,omitempty"`

	ChannelId string        `gorm:"primaryKey" json:"channel_id"`
	Channel   *SlackChannel `gorm:"foreignKey:ChannelId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *SlackThread) CreatedAt() time.Time {
	// Convert the ID to a float64 to extract the timestamp
	timestampFloat, err := strconv.ParseFloat(t.Id, 64)
	if err != nil {
		return time.Time{} // Return zero time on error
	}
	seconds := int64(timestampFloat)
	nanos := int64((timestampFloat - float64(seconds)) * 1e9)
	return time.Unix(seconds, nanos)
}
