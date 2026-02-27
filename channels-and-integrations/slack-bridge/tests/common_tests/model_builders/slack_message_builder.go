package model_builders

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type SlackMessageBuilder struct {
	id        string
	content   string
	threadId  string
	channelId string

	isFromUser  bool
	isFromAgent bool

	authorUserId            *string
	authorCollectionAddress *string
	authorCollectionTokenId *uint64
}

func NewSlackMessageBuilder() *SlackMessageBuilder {
	id := uuid.New().String()
	content := faker.Sentence()
	threadId := uuid.New().String()
	channelId := uuid.New().String()

	isFromUser := false
	isFromAgent := false

	return &SlackMessageBuilder{
		id:        id,
		content:   content,
		threadId:  threadId,
		channelId: channelId,

		isFromUser:  isFromUser,
		isFromAgent: isFromAgent,

		authorUserId:            nil,
		authorCollectionAddress: nil,
		authorCollectionTokenId: nil,
	}
}

func (b *SlackMessageBuilder) WithId(id string) *SlackMessageBuilder {
	b.id = id
	return b
}

func (b *SlackMessageBuilder) WithContent(content string) *SlackMessageBuilder {
	b.content = content
	return b
}

func (b *SlackMessageBuilder) WithThreadId(threadId string) *SlackMessageBuilder {
	b.threadId = threadId
	return b
}

func (b *SlackMessageBuilder) WithChannelId(channelId string) *SlackMessageBuilder {
	b.channelId = channelId
	return b
}

func (b *SlackMessageBuilder) WithIsFromUser(isFromUser bool) *SlackMessageBuilder {
	b.isFromUser = isFromUser
	return b
}

func (b *SlackMessageBuilder) WithIsFromAgent(isFromAgent bool) *SlackMessageBuilder {
	b.isFromAgent = isFromAgent
	return b
}

func (b *SlackMessageBuilder) WithAuthorUserId(authorUserId string) *SlackMessageBuilder {
	b.authorUserId = &authorUserId
	return b
}

func (b *SlackMessageBuilder) WithAuthorCollectionAddress(authorCollectionAddress string) *SlackMessageBuilder {
	b.authorCollectionAddress = &authorCollectionAddress
	return b
}

func (b *SlackMessageBuilder) WithAuthorCollectionTokenId(authorCollectionTokenId uint64) *SlackMessageBuilder {
	b.authorCollectionTokenId = &authorCollectionTokenId
	return b
}

func (b *SlackMessageBuilder) Build() *models.SlackMessage {
	return &models.SlackMessage{
		Id:        b.id,
		Content:   b.content,
		ThreadId:  b.threadId,
		ChannelId: b.channelId,

		IsFromUser:  b.isFromUser,
		IsFromAgent: b.isFromAgent,

		AuthorUserId:            b.authorUserId,
		AuthorCollectionAddress: b.authorCollectionAddress,
		AuthorCollectionTokenId: b.authorCollectionTokenId,
	}
}
