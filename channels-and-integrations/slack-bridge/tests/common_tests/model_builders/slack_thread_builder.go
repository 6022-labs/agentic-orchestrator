package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/google/uuid"
)

type SlackThreadBuilder struct {
	id            string
	channelId     string
	isConcluded   bool
	unconcludedAt *time.Time
}

func NewSlackThreadBuilder() *SlackThreadBuilder {
	id := uuid.New().String()
	channelId := uuid.New().String()
	isConcluded := false

	return &SlackThreadBuilder{
		id:            id,
		channelId:     channelId,
		isConcluded:   isConcluded,
		unconcludedAt: nil,
	}
}

func (b *SlackThreadBuilder) WithId(id string) *SlackThreadBuilder {
	b.id = id
	return b
}

func (b *SlackThreadBuilder) WithChannelId(channelId string) *SlackThreadBuilder {
	b.channelId = channelId
	return b
}

func (b *SlackThreadBuilder) WithIsConcluded(isConcluded bool) *SlackThreadBuilder {
	b.isConcluded = isConcluded
	return b
}

func (b *SlackThreadBuilder) WithUnconcludedAt(unconcludedAt *time.Time) *SlackThreadBuilder {
	b.unconcludedAt = unconcludedAt
	return b
}

func (b *SlackThreadBuilder) Build() *models.SlackThread {
	return &models.SlackThread{
		Id:            b.id,
		ChannelId:     b.channelId,
		IsConcluded:   b.isConcluded,
		UnconcludedAt: b.unconcludedAt,
	}
}
