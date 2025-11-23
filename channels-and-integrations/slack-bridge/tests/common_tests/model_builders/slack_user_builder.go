package model_builders

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/go-faker/faker/v4"
)

type SlackUserBuilder struct {
	id       string
	name     string
	realName string

	walletAddress *string

	pairedHumanAgentCollectionAddress *string
	pairedHumanAgentCollectionTokenId *uint64
}

func NewSlackUserBuilder() *SlackUserBuilder {
	id := "default-slack-user-id"
	name := faker.Name()
	realName := faker.Name()

	return &SlackUserBuilder{
		id:       id,
		name:     name,
		realName: realName,

		walletAddress: nil,

		pairedHumanAgentCollectionAddress: nil,
		pairedHumanAgentCollectionTokenId: nil,
	}
}

func (b *SlackUserBuilder) WithId(id string) *SlackUserBuilder {
	b.id = id
	return b
}

func (b *SlackUserBuilder) WithName(name string) *SlackUserBuilder {
	b.name = name
	return b
}

func (b *SlackUserBuilder) WithRealName(realName string) *SlackUserBuilder {
	b.realName = realName
	return b
}

func (b *SlackUserBuilder) WithWalletAddress(walletAddress *string) *SlackUserBuilder {
	b.walletAddress = walletAddress
	return b
}

func (b *SlackUserBuilder) WithPairedHumanAgentCollectionAddress(address *string) *SlackUserBuilder {
	b.pairedHumanAgentCollectionAddress = address
	return b
}

func (b *SlackUserBuilder) WithPairedHumanAgentCollectionTokenId(tokenId *uint64) *SlackUserBuilder {
	b.pairedHumanAgentCollectionTokenId = tokenId
	return b
}

func (b *SlackUserBuilder) Build() *models.SlackUser {
	return &models.SlackUser{
		Id:       b.id,
		Name:     b.name,
		RealName: b.realName,

		WalletAddress: b.walletAddress,

		PairedHumanAgentCollectionAddress: b.pairedHumanAgentCollectionAddress,
		PairedHumanAgentCollectionTokenId: b.pairedHumanAgentCollectionTokenId,
	}
}
