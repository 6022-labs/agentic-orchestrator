package model_builders

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/go-faker/faker/v4"
)

type SlackAgentBuilder struct {
	collectionAddress string
	collectionTokenId uint64
	name              string
	ensDomain         string
	role              enums.AgentRole
	ownerWallet       string
	iconImageCid      *string
	isProcessing      bool
	lastProcessedAt   *time.Time
}

func NewSlackAgentBuilder() *SlackAgentBuilder {
	collectionAddress := common_tests.RandomEthAddress().Hex()
	collectionTokenId := uint64(0)
	role := enums.AgentRoleExpert
	name := faker.Name()
	ensDomain := fmt.Sprintf("%s.6022.eth", name)
	ownerWallet := common_tests.RandomEthAddress().Hex()
	iconImageCid := faker.Word()

	return &SlackAgentBuilder{
		collectionAddress: collectionAddress,
		collectionTokenId: collectionTokenId,
		role:              role,
		name:              name,
		ensDomain:         ensDomain,
		isProcessing:      false,
		ownerWallet:       ownerWallet,
		iconImageCid:      &iconImageCid,
	}
}

func (b *SlackAgentBuilder) WithCollectionAddress(collectionAddress string) *SlackAgentBuilder {
	b.collectionAddress = collectionAddress
	return b
}

func (b *SlackAgentBuilder) WithCollectionTokenId(collectionTokenId uint64) *SlackAgentBuilder {
	b.collectionTokenId = collectionTokenId
	return b
}

func (b *SlackAgentBuilder) WithRole(role enums.AgentRole) *SlackAgentBuilder {
	b.role = role
	return b
}

func (b *SlackAgentBuilder) WithName(name string) *SlackAgentBuilder {
	b.name = name
	return b
}

func (b *SlackAgentBuilder) WithEnsDomain(ensDomain string) *SlackAgentBuilder {
	b.ensDomain = ensDomain
	return b
}

func (b *SlackAgentBuilder) WithOwnerWallet(ownerWallet string) *SlackAgentBuilder {
	b.ownerWallet = ownerWallet
	return b
}

func (b *SlackAgentBuilder) WithIconImageCid(iconImageCid *string) *SlackAgentBuilder {
	b.iconImageCid = iconImageCid
	return b
}

func (b *SlackAgentBuilder) WithIsProcessing(isProcessing bool) *SlackAgentBuilder {
	b.isProcessing = isProcessing
	return b
}

func (b *SlackAgentBuilder) WithLastProcessedAt(lastProcessedAt *time.Time) *SlackAgentBuilder {
	b.lastProcessedAt = lastProcessedAt
	return b
}

func (b *SlackAgentBuilder) Build() *models.SlackAgent {
	return &models.SlackAgent{
		CollectionAddress: b.collectionAddress,
		CollectionTokenId: b.collectionTokenId,
		Role:              b.role,
		Name:              b.name,
		EnsDomain:         b.ensDomain,
		OwnerWallet:       b.ownerWallet,
		IconImageCid:      b.iconImageCid,
		IsProcessing:      b.isProcessing,
		LastProcessedAt:   b.lastProcessedAt,
	}
}
